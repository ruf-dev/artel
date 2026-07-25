//go:build e2e

package tract_e2e_test

// Covers the tract basic scenario end-to-end against real Postgres/CouchDB (via
// tests/docker-compose.yaml), calling the service layer directly the same way
// internal/app/custom.go wires it — no mocks. Kept in its own package (rather than
// tests/e2e) so it doesn't depend on that package's other files compiling.
//
// Run: docker compose -f tests/docker-compose.yaml up -d
//      go test -tags e2e ./tests/tract_e2e/...

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/tract"
	"github.com/ruf-dev/artel/internal/service/v1/tract/script"
	"github.com/stretchr/testify/suite"
)

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return def
}

func randomEmail() string {
	return fmt.Sprintf("tract_e2e_%08x@test.local", rand.Uint32())
}

// putCouchDatabase creates dbName on the test CouchDB instance — the vault row this test
// creates bypasses the real Vault.CreateVault provisioning, so the database itself needs to
// exist for the list_files action step to query it successfully.
func putCouchDatabase(couchURL, user, pass, dbName string) error {
	req, err := http.NewRequest(http.MethodPut, couchURL+"/"+dbName, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, pass)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusPreconditionFailed {
		return fmt.Errorf("unexpected status creating couch database: %d", resp.StatusCode)
	}
	return nil
}

func deleteCouchDatabase(couchURL, user, pass, dbName string) error {
	req, err := http.NewRequest(http.MethodDelete, couchURL+"/"+dbName, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, pass)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	return nil
}

type TractE2ESuite struct {
	suite.Suite
	db            *sql.DB
	repos         *repopg.Repos
	svcs          *svcv1.Services
	couchInstance string
	couchURL      string
	couchUser     string
	couchPass     string
}

func TestTractE2E(t *testing.T) {
	suite.Run(t, new(TractE2ESuite))
}

func (s *TractE2ESuite) SetupSuite() {
	ctx := context.Background()

	pgDSN := envOrDefault("PG_DSN", "postgres://artel:artel_db@localhost:15434/artel_db?sslmode=disable")
	db, err := sql.Open("postgres", pgDSN)
	s.Require().NoError(err, "open postgres")

	err = db.Ping()
	s.Require().NoError(err, "ping postgres — is tests/docker-compose.yaml up?")
	s.db = db

	goose.SetLogger(goose.NopLogger())

	err = goose.SetDialect("postgres")
	s.Require().NoError(err)

	err = goose.Up(db, "../../migrations")
	s.Require().NoError(err, "run migrations")

	encKey := make([]byte, 32)
	encryptor, err := cryptoutil.NewAESEncryptor(encKey)
	s.Require().NoError(err, "create AES encryptor")
	s.repos = repopg.New(db, encryptor)

	cfg := config.EnvironmentConfig{}
	svcs, err := svcv1.New(s.repos, cfg)
	s.Require().NoError(err, "init services")
	s.svcs = svcs

	// Tract is normally wired in internal/app/custom.go (not svcv1.New) because its
	// ToolExecutor composes the already-built Mcp/Mom services — mirror that wiring exactly
	// since this test drives the service layer directly.
	tractToolExecutor := tract.NewToolExecutor(s.svcs.McpService(), s.svcs.MomService())
	tractLlmExecutor := tract.NewLlmExecutor(s.repos.ExternalConnections())
	scriptEngines := script.NewRegistry(script.NewJavaScriptEngine())
	s.svcs.Tract = tract.New(
		s.repos.Tracts(),
		s.repos.TractTemplates(),
		s.repos.Triggers(),
		s.repos.TriggerPresets(),
		s.repos.ExternalConnections(),
		s.repos.McpDefinitions(),
		tractToolExecutor,
		s.svcs.SubscriptionService(),
		scriptEngines,
		tractLlmExecutor,
	)
	s.svcs.McpService().SetTractService(ctx, s.svcs.TractService())

	s.couchURL = envOrDefault("COUCH_URL", "http://localhost:15985")
	s.couchUser = envOrDefault("COUCH_USER", "admin")
	s.couchPass = envOrDefault("COUCH_PASS", "admin")

	s.couchInstance, err = s.svcs.CouchInstance.RegisterCouchInstance(ctx, s.couchURL, s.couchUser, s.couchPass)
	s.Require().NoError(err, "register couch instance")
}

func (s *TractE2ESuite) TearDownSuite() {
	ctx := context.Background()
	if s.couchInstance != "" {
		err := s.svcs.CouchInstance.DeleteCouchInstance(ctx, s.couchInstance)
		s.NoError(err)
	}
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err)
	}
}

// TestTractBasicScenario mirrors the UI's create -> add step -> save -> run -> view run flow:
// create a tract with one builtin action step, create and link a manual trigger, run it, then
// read the run back and assert the step completed.
func (s *TractE2ESuite) TestTractBasicScenario() {
	ctx := context.Background()

	email := randomEmail()
	user, err := s.svcs.Auth.Register(ctx, email, "test-password-e2e")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), user.Uuid)
	})

	uc := user_context.UserContext{UserUuid: user.Uuid}
	userCtx := user_context.WithUserContext(ctx, uc)

	// The list_files action step below is a builtin vault tool that resolves against *a* real
	// vault (by membership, then by its CouchDB database) — bypass the real CouchDB
	// provisioning (Vault.CreateVault) by pointing a vault row directly at our own throwaway
	// couch instance and creating the database it names.
	const vaultDbName = "e2e_tract_vault_db"
	err = putCouchDatabase(s.couchURL, s.couchUser, s.couchPass, vaultDbName)
	s.Require().NoError(err, "create couch database for vault")
	s.T().Cleanup(func() {
		_ = deleteCouchDatabase(s.couchURL, s.couchUser, s.couchPass, vaultDbName)
	})

	couchInstanceUuid, err := uuid.Parse(s.couchInstance)
	s.Require().NoError(err)
	testVault, err := s.repos.Vaults().
		Upsert(ctx, user.Uuid, couchInstanceUuid, "e2e_tract_vault", vaultDbName, "active", "")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Vaults().Delete(context.Background(), testVault.Uuid)
	})

	// ExecuteBuiltinToolForUser resolves the vault via membership, not vault ownership —
	// without this row the run fails with "no vault available to run a builtin tool".
	err = s.repos.VaultMembers().Add(ctx, testVault.Uuid, user.Uuid, artel_q.VaultRoleOwner)
	s.Require().NoError(err)

	// 1. Create a tract with a single builtin action step.
	step := domain.TractStep{
		Id:     "list_files",
		Name:   "list_files",
		Type:   "action",
		Mcp:    "artel",
		Tool:   "list_files",
		Params: map[string]string{},
	}
	definition := domain.TractDefinition{Steps: []domain.TractStep{step}}
	tr, warnings, err := s.svcs.Tract.CreateTract(userCtx, "e2e basic tract", "", definition)
	s.Require().NoError(err)
	s.Empty(warnings)
	s.T().Cleanup(func() {
		_ = s.svcs.Tract.DeleteTract(context.Background(), tr.Uuid)
	})

	// 2. Create and link a manual trigger — the same "Add trigger" flow the Trigger inspector
	// panel drives.
	payloadSchema := domain.ToolSchema{Properties: map[string]domain.ToolProperty{}}
	trigger, _, err := s.svcs.Tract.CreateTrigger(
		userCtx,
		"e2e manual trigger",
		"manual",
		"generic",
		json.RawMessage(`{}`),
		payloadSchema,
	)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Tract.DeleteTrigger(context.Background(), trigger.Uuid)
	})

	err = s.svcs.Tract.LinkTrigger(userCtx, trigger.Uuid, tr.Uuid, nil)
	s.Require().NoError(err)

	// 3. Run it — the same call the canvas builder's Run button makes — and assert it completes.
	run, err := s.svcs.Tract.StartRun(userCtx, tr, json.RawMessage(`{}`), tract.StartedByManual, trigger.Uuid)
	s.Require().NoError(err)
	s.Equal(domain.TractRunDone, run.Status)

	// 4. Read the run back — the same call the log panel makes — and assert the step recorded.
	gotRun, steps, err := s.svcs.Tract.GetRun(userCtx, run.Uuid)
	s.Require().NoError(err)
	s.Equal(domain.TractRunDone, gotRun.Status)
	s.Require().Len(steps, 1)
	s.Equal("list_files", steps[0].StepId)
	s.Equal(domain.TractRunStepDone, steps[0].Status)
}
