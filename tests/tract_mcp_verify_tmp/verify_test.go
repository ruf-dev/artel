//go:build tractverify

package tract_mcp_verify_tmp

// THROWAWAY verification test — not part of the deliverable. Exercises the new MCP tract
// builtin tools (list_tract_actions, create_tract, run_tract, get_tract_runs) end-to-end
// against the real dev Postgres/CouchDB. Deleted after use.

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/tract"
	"github.com/ruf-dev/artel/internal/transport/mcp_api"
	"github.com/stretchr/testify/suite"
)

func randomEmail() string {
	return fmt.Sprintf("tractverify_%08x@test.local", rand.Uint32())
}

func mcpCall(method string, params any) string {
	p, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":%q,"params":%s}`, method, p)
}

type TractVerifySuite struct {
	suite.Suite
	db              *sql.DB
	repos           *repopg.Repos
	svcs            *svcv1.Services
	mcpHdlr         *mcp_api.McpHandler
	couchInstanceID string
}

func TestTractVerify(t *testing.T) {
	suite.Run(t, new(TractVerifySuite))
}

func (s *TractVerifySuite) SetupSuite() {
	ctx := context.Background()

	pgDSN := "postgres://artel:artel@localhost:15433/artel_db?sslmode=disable"
	db, err := sql.Open("postgres", pgDSN)
	s.Require().NoError(err)
	err = db.Ping()
	s.Require().NoError(err, "ping postgres")
	s.db = db

	encKey := make([]byte, 32)
	s.repos = repopg.New(db, encKey)

	cfg := config.EnvironmentConfig{}
	s.svcs, err = svcv1.New(s.repos, cfg)
	s.Require().NoError(err)

	// Mirror custom.go's Tract wiring exactly.
	tractToolExecutor := tract.NewToolExecutor(s.svcs.McpService(), s.svcs.MomService())
	s.svcs.Tract = tract.New(
		s.repos.Tracts(),
		s.repos.Triggers(),
		s.repos.ExternalConnections(),
		s.repos.McpDefinitions(),
		tractToolExecutor,
	)
	s.svcs.McpService().SetTractService(ctx, s.svcs.TractService())

	// Use 127.0.0.1 rather than localhost: a pre-existing (unrelated, real) couch_instances row
	// already occupies the "http://localhost:15984" URL and must not be touched or reused —
	// its credentials are encrypted with a different key than this throwaway test's zero key.
	couchURL := "http://127.0.0.1:15984"
	s.couchInstanceID, err = s.svcs.CouchInstance.RegisterCouchInstance(ctx, couchURL, "admin", "admin")
	s.Require().NoError(err, "register couch instance")

	s.mcpHdlr = mcp_api.NewMcpHandler(s.svcs.Mcp, s.svcs.Mom)
}

func (s *TractVerifySuite) TearDownSuite() {
	ctx := context.Background()
	if s.couchInstanceID != "" {
		err := s.svcs.CouchInstance.DeleteCouchInstance(ctx, s.couchInstanceID)
		s.NoError(err, "delete couch instance")
	}
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err)
	}
}

func (s *TractVerifySuite) TestTractMcpTools() {
	ctx := context.Background()

	email := randomEmail()
	user, err := s.svcs.Auth.Register(ctx, email, "test-password-tractverify")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), user.Uuid)
	})

	uc := user_context.UserContext{UserUuid: user.Uuid}
	userCtx := user_context.WithUserContext(ctx, uc)

	// Bypass Vault.CreateVault: it provisions a real CouchDB database/user over HTTP, which this
	// test doesn't need (the tract MCP tools under test never touch CouchDB). Insert the vault
	// row directly instead, pointed at our own throwaway couch instance.
	couchInstanceUuid, err := uuid.Parse(s.couchInstanceID)
	s.Require().NoError(err)
	vault, err := s.repos.Vaults().
		Upsert(ctx, user.Uuid, couchInstanceUuid, "tractverify_vault", "tractverify_vault_db", "active", "")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Vaults().Delete(context.Background(), vault.Uuid)
	})

	rawToken, key, err := s.svcs.Mcp.CreateKey(userCtx, vault.Uuid, "tractverify-key")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Mcp.RevokeKey(context.Background(), key.Uuid)
	})

	doCall := func(body string) map[string]any {
		req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+rawToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		s.mcpHdlr.ServeHTTP(w, req)
		s.Require().Equal(http.StatusOK, w.Code, "body: %s", w.Body.String())

		var resp map[string]any
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		s.Require().NoError(err, "unmarshal response: %s", w.Body.String())
		s.Require().Nil(resp["error"], "rpc error in response: %s", w.Body.String())
		return resp
	}

	// 1. tools/list must include the new tract tools.
	listResp := doCall(mcpCall("tools/list", map[string]any{}))
	result := listResp["result"].(map[string]any)
	toolsRaw := result["tools"].([]any)
	names := map[string]bool{}
	for _, t := range toolsRaw {
		tm := t.(map[string]any)
		names[tm["name"].(string)] = true
	}
	for _, want := range []string{"list_tract_actions", "create_tract", "create_trigger", "link_trigger", "list_tracts", "get_tract", "get_tract_runs", "run_tract"} {
		s.Require().True(names[want], "tools/list missing %s", want)
	}

	// 2. list_tract_actions
	actionsResp := doCall(mcpCall("tools/call", map[string]any{
		"name":      "list_tract_actions",
		"arguments": map[string]any{},
	}))
	s.T().Logf("list_tract_actions result: %+v", actionsResp["result"])

	// 3. create_tract with a minimal definition + webhook_trigger
	definition := map[string]any{
		"root": map[string]any{
			"type":   "action",
			"tool":   "list_tract_actions",
			"params": map[string]any{},
		},
	}
	createArgs := map[string]any{
		"name":       "tractverify tract",
		"definition": definition,
		"webhook_trigger": map[string]any{
			"name":   "tractverify trigger",
			"source": "manual",
		},
	}
	createResp := doCall(mcpCall("tools/call", map[string]any{
		"name":      "create_tract",
		"arguments": createArgs,
	}))
	s.T().Logf("create_tract result: %+v", createResp["result"])

	createResult := extractToolResult(s.T(), createResp)
	tractUuid, ok := createResult["tract_uuid"].(string)
	s.Require().True(ok, "create_tract result missing tract_uuid: %+v", createResult)
	s.Require().NotEmpty(tractUuid)

	// 4. run_tract
	runResp := doCall(mcpCall("tools/call", map[string]any{
		"name": "run_tract",
		"arguments": map[string]any{
			"tract_uuid": tractUuid,
		},
	}))
	s.T().Logf("run_tract result: %+v", runResp["result"])

	// 5. get_tract_runs
	runsResp := doCall(mcpCall("tools/call", map[string]any{
		"name": "get_tract_runs",
		"arguments": map[string]any{
			"tract_uuid": tractUuid,
		},
	}))
	s.T().Logf("get_tract_runs result: %+v", runsResp["result"])
}

func extractToolResult(t *testing.T, resp map[string]any) map[string]any {
	result := resp["result"].(map[string]any)
	content := result["content"].([]any)
	first := content[0].(map[string]any)
	text, ok := first["text"].(string)
	if !ok {
		t.Fatalf("tool result content[0] has no text field: %+v", first)
	}
	var out map[string]any
	err := json.Unmarshal([]byte(text), &out)
	if err != nil {
		t.Fatalf("unmarshal tool result text: %v (%s)", err, text)
	}
	return out
}
