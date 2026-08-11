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
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_keys_api"
	"github.com/ruf-dev/artel/internal/transport/notes_api"
	"github.com/ruf-dev/artel/internal/transport/vaults_api"
	"github.com/ruf-dev/artel/tests/harness"

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

	authClient   pb.AuthAPIClient
	vaultClient  pb.VaultsAPIClient
	mcpKeyClient pb.McpKeysAPIClient
	notesClient  pb.NotesAPIClient
}

func TestQuota(t *testing.T) {
	suite.Run(t, new(QuotaSuite))
}

func (s *QuotaSuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	// SubscriptionsEnabled must be true here: with it false (the E2ESuite default) every
	// quota check runs through FreeService, which always passes and never touches CouchDB/S3.
	cfg := config.EnvironmentConfig{}
	cfg.SubscriptionsEnabled = true

	var credsEncrypted bool
	s.repos, s.svcs, credsEncrypted = harness.BuildServices(s.T(), s.db, cfg)

	couchURL := harness.CouchURL(s.T())

	s.couchInstanceID = harness.GetCouchInstance(s.T(), ctx, s.svcs, couchURL)

	s3Endpoint := harness.S3Endpoint(s.T())

	s.s3InstanceID = harness.GetS3Instance(s.T(), ctx, s.svcs, s3Endpoint)

	s.mcpHdlr = mcp_api.NewMcpHandler(s.svcs.Mcp, s.svcs.Mom)

	s.startGrpcServer(credsEncrypted)
}

// startGrpcServer builds a real *grpc.Server chained with the production auth interceptor,
// registers the AuthAPI, VaultsAPI, McpKeysAPI and NotesAPI implementations onto it, and serves it
// over an in-memory bufconn listener so the suite's RPCs travel through the real transport + auth
// stack without binding a TCP port.
func (s *QuotaSuite) startGrpcServer(credsEncrypted bool) {
	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, credsEncrypted, s.svcs.DockerHost, s.svcs.SetupWizard,
	)
	vaultsImpl := vaults_api.NewVaultsImpl(s.svcs.Vault, s.svcs.Workbench)
	mcpKeysImpl := mcp_keys_api.NewMcpKeysImpl(s.svcs.Mcp, s.svcs.Mom)
	notesImpl := notes_api.NewNotesImpl(s.svcs.Notes)

	conn := harness.NewBufconnServer(
		s.T(), s.svcs, authImpl.Register, vaultsImpl.Register, mcpKeysImpl.Register, notesImpl.Register,
	)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.vaultClient = pb.NewVaultsAPIClient(conn)
	s.mcpKeyClient = pb.NewMcpKeysAPIClient(conn)
	s.notesClient = pb.NewNotesAPIClient(conn)
}

func (s *QuotaSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

type quotaTestUser struct {
	userUuid  uuid.UUID
	userCtx   context.Context
	authedCtx context.Context
	vaultUuid uuid.UUID
	rawToken  string
}

// setupUser registers a user via the real AuthAPI RPC, creates a vault for them via VaultsAPI, and
// issues an MCP key scoped to it via McpKeysAPI — the minimum a caller needs to attempt a write.
func (s *QuotaSuite) setupUser() quotaTestUser {
	ctx := context.Background()

	slug := harness.Slug(s.T())
	email := slug + "@test.local"
	// A fixed password, not slug-derived: bcrypt rejects passwords over 72 bytes, and slug embeds
	// the full (possibly long, subtest-qualified) t.Name() — uniqueness only matters for email.
	password := "test-password-quota"

	registerReq := &pb.Register_Request{Email: email, Password: password}

	registerResp, err := s.authClient.Register(ctx, registerReq)
	s.Require().NoError(err)

	userUuid, err := uuid.Parse(registerResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), userUuid)
	})

	// UserName must be non-empty: vault creation derives the CouchDB database name from it, and
	// an empty UserName produces a name starting with "-", which CouchDB rejects. Email/password
	// registration never populates domain.User.Username (defaults to "" — see
	// migrations/007_telegram_auth.sql), so the slug is written directly onto the persisted user
	// row, since this suite now authenticates through the real gRPC auth interceptor (which reads
	// UserName off that row, not off a hand-built context).
	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, slug, userUuid)
	s.Require().NoError(err, "set username stand-in for vault creation")

	// With SubscriptionsEnabled=true, middleware.authWithSession calls
	// subscriptionService.CheckActive on every authenticated RPC — so the subscription has to be
	// activated before CreateVault/CreateMcpKey below, not only once overrideQuota runs later in
	// the test. overrideQuota itself re-activates (Upsert is idempotent) once it pins a quota
	// bound.
	_, err = s.repos.Subscriptions().Upsert(ctx, userUuid, true)
	s.Require().NoError(err, "activate subscription so authenticated RPCs pass CheckActive")

	passwordCreds := &pb.PasswordCredentials{Email: email, Password: password}
	loginMethod := &pb.Login_Request_Password{Password: passwordCreds}
	loginReq := &pb.Login_Request{Method: loginMethod}

	loginResp, err := s.authClient.Login(ctx, loginReq)
	s.Require().NoError(err)
	s.Require().NotEmpty(loginResp.Token)

	authedCtx := harness.AuthedContext(ctx, loginResp.Token)

	createVaultReq := &pb.CreateVault_Request{Name: slug + "_vault"}

	createVaultResp, err := s.vaultClient.CreateVault(authedCtx, createVaultReq)
	s.Require().NoError(err)

	vaultUuid, err := uuid.Parse(createVaultResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vaultUuid)
	})

	createKeyReq := &pb.CreateMcpKey_Request{
		VaultId: createVaultResp.Id,
		Name:    slug + "_key",
	}

	createKeyResp, err := s.mcpKeyClient.CreateMcpKey(authedCtx, createKeyReq)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		keyUuid, parseErr := uuid.Parse(createKeyResp.Key.Id)
		if parseErr == nil {
			_ = s.svcs.Mcp.RevokeKey(context.Background(), keyUuid)
		}
	})

	// LinkS3Bucket/SetUseCouchDBForBinaries (used by TestS3QuotaExceeded_RejectsUploadBeforeWrite)
	// are called directly against the service layer rather than through a gRPC client, and
	// SetUseCouchDBForBinaries reads the caller's identity via user_context.GetUserContext — a
	// hand-built context is still needed alongside the gRPC-authenticated one above.
	uc := user_context.UserContext{UserUuid: userUuid, UserName: slug}
	userCtx := user_context.WithUserContext(ctx, uc)

	testUser := quotaTestUser{
		userUuid:  userUuid,
		userCtx:   userCtx,
		authedCtx: authedCtx,
		vaultUuid: vaultUuid,
		rawToken:  createKeyResp.RawToken,
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
	user := s.setupUser()

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

	listNotesReq := &pb.ListNotes_Request{VaultId: user.vaultUuid.String()}

	listResp, err := s.notesClient.ListNotes(user.authedCtx, listNotesReq)
	s.Require().NoError(err)

	var paths []string
	for _, n := range listResp.Notes {
		paths = append(paths, n.Path)
	}

	s.Contains(paths, sanityPath, "the pre-override note should still be present")
	s.NotContains(paths, blockedPath, "the rejected note must never have been written to CouchDB")
}

// TestS3QuotaExceeded_RejectsUploadBeforeWrite links a vault to a real S3 (MinIO) bucket,
// overrides the user's S3 quota down to zero, and confirms a subsequent binary upload is
// rejected before any object is put into the bucket.
func (s *QuotaSuite) TestS3QuotaExceeded_RejectsUploadBeforeWrite() {
	user := s.setupUser()

	s3InstanceUuid, err := uuid.Parse(s.s3InstanceID)
	s.Require().NoError(err)

	bucketName := fmt.Sprintf("quota-test-%08x", rand.Uint32())

	err = s.svcs.Vault.LinkS3Bucket(user.userCtx, user.vaultUuid, &s3InstanceUuid, bucketName)
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
