//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/transport/mcp_api"
	"github.com/stretchr/testify/suite"
)

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func randomEmail() string {
	return fmt.Sprintf("e2e_%08x@test.local", rand.Uint32())
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
}

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ESuite))
}

func (s *E2ESuite) SetupSuite() {
	ctx := context.Background()

	pgDSN := envOrDefault("PG_DSN", "postgres://artel:artel_db@localhost:15434/artel_db?sslmode=disable")
	db, err := sql.Open("postgres", pgDSN)
	s.Require().NoError(err, "open postgres")

	err = db.Ping()
	s.Require().NoError(err, "ping postgres — is the container running?")
	s.db = db

	goose.SetLogger(goose.NopLogger())

	err = goose.SetDialect("postgres")
	s.Require().NoError(err)

	err = goose.Up(db, "../../migrations")
	s.Require().NoError(err, "run migrations")

	encKey := make([]byte, 32)
	s.repos = repopg.New(db, encKey)

	cfg := config.EnvironmentConfig{}
	s.svcs, err = svcv1.New(s.repos, cfg)
	s.Require().NoError(err, "init services")

	couchURL := envOrDefault("COUCH_URL", "http://localhost:15985")
	couchUser := envOrDefault("COUCH_USER", "admin")
	couchPass := envOrDefault("COUCH_PASS", "admin")

	s.couchInstanceID, err = s.svcs.CouchInstance.RegisterCouchInstance(ctx, couchURL, couchUser, couchPass)
	s.Require().NoError(err, "register couch instance")

	s.mcpHdlr = mcp_api.NewMcpHandler(s.svcs.Mcp, s.svcs.Email, s.svcs.Mom)
}

func (s *E2ESuite) TearDownSuite() {
	ctx := context.Background()
	if s.couchInstanceID != "" {
		err := s.svcs.CouchInstance.DeleteCouchInstance(ctx, s.couchInstanceID)
		s.NoError(err, "delete couch instance")
	}
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

func (s *E2ESuite) TestUserSessionVaultMCPWriteNotesRead() {
	ctx := context.Background()

	// 1. Register user
	email := randomEmail()
	password := "test-password-e2e"
	user, err := s.svcs.Auth.Register(ctx, email, password)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), user.Uuid)
	})

	// 2. Login → session
	session, err := s.svcs.Auth.Login(ctx, email, password)
	s.Require().NoError(err)
	s.Require().NotEmpty(session.Token)

	// 3. Build user context (replaces the gRPC auth interceptor for direct service calls)
	uc := user_context.UserContext{UserUuid: user.Uuid}
	userCtx := user_context.WithUserContext(ctx, uc)

	// 4. Create vault
	vault, err := s.svcs.Vault.CreateVault(userCtx, "e2e_test_vault")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vault.Uuid)
	})

	// 5. Enable HasNotes (default permissions have it off)
	_, err = s.repos.UserPermissions().Upsert(ctx, user.Uuid, false, false, false, true)
	s.Require().NoError(err)

	// 6. Create MCP key for this vault
	rawToken, key, err := s.svcs.Mcp.CreateKey(userCtx, vault.Uuid, "e2e-key")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Mcp.RevokeKey(context.Background(), key.Uuid)
	})

	// 7. Write a note via MCP HTTP handler
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
	req.Header.Set("Authorization", "Bearer "+rawToken)
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

	// 8. Read notes back via notes service API
	notes, err := s.svcs.Notes.ListNotes(userCtx, vault.Uuid)
	s.Require().NoError(err)

	found := false
	for _, n := range notes {
		if n.Path == "e2e/hello.md" {
			found = true
			break
		}
	}
	s.Require().True(found, "note 'e2e/hello.md' not found in ListNotes result")
}
