//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/transport/mcp_api"
	"github.com/ruf-dev/artel/migrations"

	"github.com/stretchr/testify/suite"
)

// QuotaSuite exercises PaidService.CheckStorageQuota end to end: a real CouchDB and a real
// S3-compatible store (MinIO) sit behind the notes/MCP write path, and a user's plan quota is
// overridden down to zero so the very next write is rejected before anything is read from or
// written to either backend.
type QuotaSuite struct {
	suite.Suite
	db              *sql.DB
	repos           *repopg.Repos
	svcs            *svcv1.Services
	mcpHdlr         *mcp_api.McpHandler
	couchInstanceID string
	s3InstanceID    string
}

func TestQuota(t *testing.T) {
	suite.Run(t, new(QuotaSuite))
}

func (s *QuotaSuite) SetupSuite() {
	ctx := context.Background()

	pgDSN := envOrDefault("PG_DSN", "postgres://artel:artel_db@localhost:15434/artel_db?sslmode=disable")
	db, err := sql.Open("postgres", pgDSN)
	s.Require().NoError(err, "open postgres")

	err = db.Ping()
	s.Require().NoError(err, "ping postgres — is the container running?")
	s.db = db

	err = migrations.ApplyMigration(db)
	s.Require().NoError(err, "run migrations")

	encKey := make([]byte, 32)
	encryptor, err := cryptoutil.NewAESEncryptor(encKey)
	s.Require().NoError(err, "create AES encryptor")
	s.repos = repopg.New(db, encryptor)

	// SubscriptionsEnabled must be true here: with it false (the E2ESuite default) every
	// quota check runs through FreeService, which always passes and never touches CouchDB/S3.
	cfg := config.EnvironmentConfig{}
	cfg.SubscriptionsEnabled = true

	s.svcs, err = svcv1.New(s.repos, cfg)
	s.Require().NoError(err, "init services")

	err = s.repos.SystemSettings().CompleteSetup(ctx)
	s.Require().NoError(err, "complete setup so Register/Login aren't gated in tests")

	err = s.repos.SystemSettings().UpdateRegistrationMode(ctx, domain.RegistrationModeSelfRegister)
	s.Require().NoError(err, "enable self-registration so tests can Register directly")

	couchURL := envOrDefault("COUCH_URL", "http://localhost:15985")
	couchUser := envOrDefault("COUCH_USER", "admin")
	couchPass := envOrDefault("COUCH_PASS", "admin")

	s.couchInstanceID, err = s.svcs.CouchInstance.RegisterCouchInstance(ctx, couchURL, couchUser, couchPass)
	s.Require().NoError(err, "register couch instance")

	s3Endpoint := envOrDefault("S3_ENDPOINT", "localhost:19000")

	s.s3InstanceID, err = s.svcs.S3Instance.RegisterS3Instance(
		ctx, s3Endpoint, "us-east-1", "minioadmin", "minioadmin", false, true,
	)
	s.Require().NoError(err, "register s3 instance — is the MinIO container running?")

	s.mcpHdlr = mcp_api.NewMcpHandler(s.svcs.Mcp, s.svcs.Mom)
}

func (s *QuotaSuite) TearDownSuite() {
	ctx := context.Background()

	if s.couchInstanceID != "" {
		err := s.svcs.CouchInstance.DeleteCouchInstance(ctx, s.couchInstanceID)
		s.NoError(err, "delete couch instance")
	}

	if s.s3InstanceID != "" {
		err := s.svcs.S3Instance.DeleteS3Instance(ctx, s.s3InstanceID)
		s.NoError(err, "delete s3 instance")
	}

	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

type quotaTestUser struct {
	userUuid  uuid.UUID
	userCtx   context.Context
	vaultUuid uuid.UUID
	rawToken  string
}

// setupUser registers a user, creates a vault for them, and issues an MCP key scoped to it —
// the minimum a caller needs to attempt a write.
func (s *QuotaSuite) setupUser(namePrefix string) quotaTestUser {
	ctx := context.Background()

	email := randomEmail()
	password := "test-password-quota"

	user, err := s.svcs.Auth.Register(ctx, email, password)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), user.Uuid)
	})

	// UserName must be non-empty: vault creation derives the CouchDB database name from it, and
	// an empty UserName produces a name starting with "-", which CouchDB rejects. Email/password
	// registration never populates domain.User.Username (defaults to "" — see
	// migrations/007_telegram_auth.sql), so a local-part stand-in is used here instead.
	uc := user_context.UserContext{UserUuid: user.Uuid, UserName: strings.SplitN(email, "@", 2)[0]}
	userCtx := user_context.WithUserContext(ctx, uc)

	vault, err := s.svcs.Vault.CreateVault(userCtx, namePrefix+"_vault")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vault.Uuid)
	})

	rawToken, key, err := s.svcs.Mcp.CreateKey(userCtx, vault.Uuid, namePrefix+"_key")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Mcp.RevokeKey(context.Background(), key.Uuid)
	})

	testUser := quotaTestUser{
		userUuid:  user.Uuid,
		userCtx:   userCtx,
		vaultUuid: vault.Uuid,
		rawToken:  rawToken,
	}

	return testUser
}

// overrideQuota activates the user's subscription and pins their effective CouchDB/S3 quotas —
// a nil bound leaves that backend on the plan's default quota (5MB couch / 20MB s3 from the
// seeded "basic" plan) so a test can exhaust one backend without also tripping the other.
func (s *QuotaSuite) overrideQuota(userUuid uuid.UUID, couchQuotaBytes, s3QuotaBytes *int64) {
	ctx := context.Background()
	subsRepo := s.repos.Subscriptions()

	_, err := subsRepo.Upsert(ctx, userUuid, true)
	s.Require().NoError(err, "activate subscription")

	sub, err := subsRepo.GetByUser(ctx, userUuid)
	s.Require().NoError(err, "get subscription")

	sub.CouchQuotaOverrideBytes = couchQuotaBytes
	sub.S3QuotaOverrideBytes = s3QuotaBytes

	_, err = subsRepo.UpsertOverrides(ctx, sub)
	s.Require().NoError(err, "override subscription quota")
}

type mcpRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// callWriteFile drives the same MCP HTTP path a real Obsidian/agent client uses (mirrors
// TestUserSessionVaultMCPWriteNotesRead in e2e_test.go), returning the RPC-level error (nil on
// success) rather than an HTTP-transport error — quota rejection is a well-formed JSON-RPC
// error response, not an HTTP failure.
func (s *QuotaSuite) callWriteFile(rawToken, path, content string) *mcpRpcError {
	writeArgs := map[string]any{
		"path":    path,
		"content": content,
	}
	callParams := map[string]any{
		"name":      "write_file",
		"arguments": writeArgs,
	}
	body := mcpCall("tools/call", callParams)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+rawToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.mcpHdlr.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code)

	var rpcResp struct {
		Error *mcpRpcError `json:"error"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &rpcResp)
	s.Require().NoError(err)

	return rpcResp.Error
}

// TestCouchQuotaExceeded_RejectsWriteBeforeUpload overrides a user's CouchDB quota down to zero
// bytes — already-at-quota — and confirms a subsequent note write is rejected with a
// CouchStorageQuotaExceeded error before CouchDB is touched, not written and then rolled back.
func (s *QuotaSuite) TestCouchQuotaExceeded_RejectsWriteBeforeUpload() {
	user := s.setupUser("couch_quota")

	// Sanity check: with the default plan quota (5MB), a small note writes fine.
	sanityPath := "quota/before.md"
	rpcErr := s.callWriteFile(user.rawToken, sanityPath, "# Before quota override\nstill under quota")
	s.Require().Nil(rpcErr, "sanity write should succeed under the default plan quota")

	zero := int64(0)
	s.overrideQuota(user.userUuid, &zero, nil)

	blockedPath := "quota/after.md"
	rpcErr = s.callWriteFile(user.rawToken, blockedPath, "# Should never land\nquota is exhausted")
	s.Require().NotNil(rpcErr, "write should be rejected once couch quota is exhausted")
	s.Contains(rpcErr.Details, "couchdb storage quota exceeded", "error should explain the couch quota was exceeded")

	notes, err := s.svcs.Notes.ListNotes(user.userCtx, user.vaultUuid)
	s.Require().NoError(err)

	var paths []string
	for _, n := range notes {
		paths = append(paths, n.Path)
	}

	s.Contains(paths, sanityPath, "the pre-override note should still be present")
	s.NotContains(paths, blockedPath, "the rejected note must never have been written to CouchDB")
}

// TestS3QuotaExceeded_RejectsUploadBeforeWrite links a vault to a real S3 (MinIO) bucket,
// overrides the user's S3 quota down to zero, and confirms a subsequent binary upload is
// rejected before any object is put into the bucket.
func (s *QuotaSuite) TestS3QuotaExceeded_RejectsUploadBeforeWrite() {
	user := s.setupUser("s3_quota")

	s3InstanceUuid, err := uuid.Parse(s.s3InstanceID)
	s.Require().NoError(err)

	bucketName := fmt.Sprintf("quota-test-%08x", rand.Uint32())

	err = s.svcs.Vault.LinkS3Bucket(user.userCtx, user.vaultUuid, s3InstanceUuid, bucketName)
	s.Require().NoError(err, "link s3 bucket to vault")

	// Vaults default to storing binaries in CouchDB (use_couchdb_for_binaries defaults true —
	// see migrations/041_vault_binary_storage.sql), even once an S3 bucket is linked; without
	// this, non-markdown writes would land in CouchDB instead of the linked bucket.
	err = s.svcs.Vault.SetUseCouchDBForBinaries(user.userCtx, user.vaultUuid, false)
	s.Require().NoError(err, "route binary storage to the linked s3 bucket")

	content := base64.StdEncoding.EncodeToString([]byte("just a few binary bytes"))

	// Sanity check: with the default plan quota (20MB), a small binary file uploads fine.
	sanityPath := "quota/before.bin"
	rpcErr := s.callWriteFile(user.rawToken, sanityPath, content)
	s.Require().Nil(rpcErr, "sanity upload should succeed under the default plan quota")

	zero := int64(0)
	s.overrideQuota(user.userUuid, nil, &zero)

	blockedPath := "quota/after.bin"
	rpcErr = s.callWriteFile(user.rawToken, blockedPath, content)
	s.Require().NotNil(rpcErr, "upload should be rejected once s3 quota is exhausted")
	s.Contains(rpcErr.Details, "s3 storage quota exceeded", "error should explain the s3 quota was exceeded")

	bucketCfg := s3client.Config{
		Endpoint:  envOrDefault("S3_ENDPOINT", "localhost:19000"),
		Region:    "us-east-1",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		UseSSL:    false,
		PathStyle: true,
	}

	bucketClient, err := s3client.New(bucketCfg, bucketName)
	s.Require().NoError(err)

	objects, err := bucketClient.List(context.Background())
	s.Require().NoError(err)

	var objectPaths []string
	for _, o := range objects {
		objectPaths = append(objectPaths, o.Path)
	}

	s.Contains(objectPaths, sanityPath, "the pre-override file should still be present")
	s.NotContains(objectPaths, blockedPath, "the rejected file must never have been written to S3")
}
