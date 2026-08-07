//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/config"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_keys_api"
	"github.com/ruf-dev/artel/internal/transport/notes_api"
	"github.com/ruf-dev/artel/internal/transport/vaults_api"
	"github.com/ruf-dev/artel/tests/harness"
)

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mcpCall(method string, params any) string {
	p, err := json.Marshal(params)
	if err != nil {
		panic(fmt.Sprintf("mcpCall: marshal params: %v", err))
	}
	return fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":%q,"params":%s}`, method, p)
}

type E2ESuite struct {
	suite.Suite
	db              *sql.DB
	repos           *repopg.Repos
	svcs            *svcv1.Services
	mcpHdlr         *mcp_api.McpHandler
	couchInstanceID string

	authClient   pb.AuthAPIClient
	vaultClient  pb.VaultsAPIClient
	mcpKeyClient pb.McpKeysAPIClient
	notesClient  pb.NotesAPIClient
}

func TestE2E(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(E2ESuite))
}

func (s *E2ESuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	cfg := config.EnvironmentConfig{}

	var credsEncrypted bool
	s.repos, s.svcs, credsEncrypted = harness.BuildServices(s.T(), s.db, cfg)

	couchURL := harness.CouchURL(s.T())

	s.couchInstanceID = harness.GetCouchInstance(s.T(), ctx, s.svcs, couchURL)

	s.mcpHdlr = mcp_api.NewMcpHandler(s.svcs.Mcp, s.svcs.Mom)

	s.startGrpcServer(credsEncrypted)
}
func (s *E2ESuite) startGrpcServer(credsEncrypted bool) {
	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, credsEncrypted, s.svcs.DockerHost, s.svcs.SetupWizard,
		false,
	)
	vaultsImpl := vaults_api.NewVaultsImpl(s.svcs.Vault, s.svcs.Workbench, false)
	mcpKeysImpl := mcp_keys_api.NewMcpKeysImpl(s.svcs.Mcp, s.svcs.Mom, false)
	notesImpl := notes_api.NewNotesImpl(s.svcs.Notes, false)

	conn := harness.NewBufconnServer(
		s.T(), s.svcs, authImpl.Register, vaultsImpl.Register, mcpKeysImpl.Register, notesImpl.Register,
	)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.vaultClient = pb.NewVaultsAPIClient(conn)
	s.mcpKeyClient = pb.NewMcpKeysAPIClient(conn)
	s.notesClient = pb.NewNotesAPIClient(conn)
}

func (s *E2ESuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

func (s *E2ESuite) TestUserSessionVaultMCPWriteNotesRead() {
	slug := harness.Slug(s.T())
	email := slug + "@test.local"
	password := "test-password-" + slug

	// 1. Register user via the real AuthAPI RPC
	registerReq := &pb.Register_Request{
		Email:    email,
		Password: password,
	}
	registerResp, err := s.authClient.Register(context.Background(), registerReq)
	s.Require().NoError(err)

	userUuid := registerResp.Id
	s.T().Cleanup(func() {
		id, parseErr := uuid.Parse(userUuid)
		if parseErr == nil {
			_ = s.repos.Users().Delete(context.Background(), id)
		}
	})

	// Email/password registration never populates domain.User.Username (defaults to "" — see
	// migrations/007_telegram_auth.sql). Vault creation derives the CouchDB database name from
	// it, and an empty UserName produces a name starting with "-", which CouchDB rejects. This
	// suite now authenticates through the real gRPC auth interceptor (which reads UserName off
	// the persisted user row, unlike a hand-built user_context.UserContext), so the stand-in has
	// to be written to the row itself.
	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, slug, userUuid)
	s.Require().NoError(err, "set username stand-in for vault creation")

	// 2. Login via the real AuthAPI RPC → session token
	passwordCreds := &pb.PasswordCredentials{Email: email, Password: password}
	loginMethod := &pb.Login_Request_Password{Password: passwordCreds}
	loginReq := &pb.Login_Request{Method: loginMethod}
	loginResp, err := s.authClient.Login(context.Background(), loginReq)
	s.Require().NoError(err)
	s.Require().NotEmpty(loginResp.Token)

	// 3. Carry the session token as outgoing gRPC metadata, same as a real client authenticating
	// against middleware.GrpcAuthInterceptor.
	authedCtx := harness.AuthedContext(context.Background(), loginResp.Token)

	// 4. Create vault via the real VaultsAPI RPC
	createVaultReq := &pb.CreateVault_Request{Name: slug + "_vault"}
	createVaultResp, err := s.vaultClient.CreateVault(authedCtx, createVaultReq)
	s.Require().NoError(err)

	vaultUuid := createVaultResp.Id
	s.T().Cleanup(func() {
		id, parseErr := uuid.Parse(vaultUuid)
		if parseErr == nil {
			_ = s.svcs.Vault.DeleteVault(context.Background(), id)
		}
	})

	// 5. Create MCP key via the real McpKeysAPI RPC
	createKeyReq := &pb.CreateMcpKey_Request{
		VaultId: vaultUuid,
		Name:    slug + "_key",
	}
	createKeyResp, err := s.mcpKeyClient.CreateMcpKey(authedCtx, createKeyReq)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		id, parseErr := uuid.Parse(createKeyResp.Key.Id)
		if parseErr == nil {
			_ = s.svcs.Mcp.RevokeKey(context.Background(), id)
		}
	})

	// 6. Write a note via the real MCP JSON-RPC-over-HTTP handler — this surface isn't gRPC (it's
	// the tool-call protocol AI clients speak), so httptest against the real handler already is
	// the real transport for it, not a shortcut.
	writeArgs := map[string]any{
		"path":    "e2e/hello.md",
		"content": "# Hello\nThis is a test note from the e2e suite.",
	}
	callParams := map[string]any{
		"name":      "write_file",
		"arguments": writeArgs,
	}
	body := mcpCall("tools/call", callParams)
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+createKeyResp.RawToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.mcpHdlr.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code)

	var rpcResp struct {
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &rpcResp)
	s.Require().NoError(err)
	s.Require().Nil(rpcResp.Error, "mcp write_file returned rpc error: %+v", rpcResp.Error)

	// 7. Read notes back via the real NotesAPI RPC
	listNotesReq := &pb.ListNotes_Request{VaultId: vaultUuid}
	listResp, err := s.notesClient.ListNotes(authedCtx, listNotesReq)
	s.Require().NoError(err)

	found := false
	for _, n := range listResp.Notes {
		if n.Path == "e2e/hello.md" {
			found = true
			break
		}
	}
	s.Require().True(found, "note 'e2e/hello.md' not found in ListNotes result")
}
