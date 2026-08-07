//go:build e2e

package tract_e2e_test

// Covers the tract basic scenario end-to-end against real Postgres/CouchDB (via
// tests/docker-compose.yaml), driving it through the real TractsAPI gRPC surface (via
// tests/harness's bufconn server + auth interceptor) rather than calling the service layer
// directly. Kept in its own package (rather than tests/e2e) so it doesn't depend on that
// package's other files compiling.
//
// Run: docker compose -f tests/docker-compose.yaml up -d
//      go test -tags e2e ./tests/tract_e2e/...

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/domain"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/tract"
	"github.com/ruf-dev/artel/internal/service/v1/tract/script"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/tracts_api"
	"github.com/ruf-dev/artel/tests/harness"
)

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

	authClient  pb.AuthAPIClient
	tractClient pb.TractsAPIClient
}

func TestTractE2E(t *testing.T) {
	suite.Run(t, new(TractE2ESuite))
}

func (s *TractE2ESuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	cfg := config.EnvironmentConfig{}

	var credsEncrypted bool
	s.repos, s.svcs, credsEncrypted = harness.BuildServices(s.T(), s.db, cfg)

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

	s.couchURL = harness.CouchURL(s.T())
	s.couchUser, s.couchPass = harness.CouchCreds(s.T())

	s.couchInstance = harness.GetCouchInstance(s.T(), ctx, s.svcs, s.couchURL)

	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, credsEncrypted, s.svcs.DockerHost, s.svcs.SetupWizard,
		false,
	)
	tractsImpl := tracts_api.New(ctx, s.svcs.TractService(), false)

	conn := harness.NewBufconnServer(s.T(), s.svcs, authImpl.Register, tractsImpl.Register)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.tractClient = pb.NewTractsAPIClient(conn)
}

func (s *TractE2ESuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err)
	}
}

// registerAndLogin registers+logs in a fresh user via the real AuthAPI RPCs, stamps the
// username stand-in row (see e2e_test.go's comment on why — the real auth interceptor reads
// UserName off the persisted row), and returns an authenticated context plus the user's uuid.
func (s *TractE2ESuite) registerAndLogin(slug string) (context.Context, uuid.UUID) {
	email := slug + "@test.local"
	password := "test-password-" + slug

	registerReq := &pb.Register_Request{
		Email:    email,
		Password: password,
	}
	registerResp, err := s.authClient.Register(context.Background(), registerReq)
	s.Require().NoError(err)

	userUuid, err := uuid.Parse(registerResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), userUuid)
	})

	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, slug, registerResp.Id)
	s.Require().NoError(err, "set username stand-in for vault creation")

	passwordCreds := &pb.PasswordCredentials{Email: email, Password: password}
	loginMethod := &pb.Login_Request_Password{Password: passwordCreds}
	loginReq := &pb.Login_Request{Method: loginMethod}
	loginResp, err := s.authClient.Login(context.Background(), loginReq)
	s.Require().NoError(err)
	s.Require().NotEmpty(loginResp.Token)

	authedCtx := harness.AuthedContext(context.Background(), loginResp.Token)

	return authedCtx, userUuid
}

// TestTractBasicScenario mirrors the UI's create -> add step -> save -> run -> view run flow:
// create a tract with one builtin action step, create and link a manual trigger, run it, then
// read the run back and assert the step completed.
func (s *TractE2ESuite) TestTractBasicScenario() {
	slug := harness.Slug(s.T())

	authedCtx, userUuid := s.registerAndLogin(slug)

	// The list_files action step below is a builtin vault tool that resolves against *a* real
	// vault (by membership, then by its CouchDB database) — bypass the real CouchDB
	// provisioning (Vault.CreateVault) by pointing a vault row directly at our own throwaway
	// couch instance and creating the database it names.
	vaultDbName := slug + "_vault_db"
	err := putCouchDatabase(s.couchURL, s.couchUser, s.couchPass, vaultDbName)
	s.Require().NoError(err, "create couch database for vault")
	s.T().Cleanup(func() {
		_ = deleteCouchDatabase(s.couchURL, s.couchUser, s.couchPass, vaultDbName)
	})

	couchInstanceUuid, err := uuid.Parse(s.couchInstance)
	s.Require().NoError(err)
	testVault, err := s.repos.Vaults().
		Upsert(context.Background(), userUuid, couchInstanceUuid, slug+"_vault", vaultDbName, "active", "")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Vaults().Delete(context.Background(), testVault.Uuid)
	})

	// ExecuteBuiltinToolForUser resolves the vault via membership, not vault ownership —
	// without this row the run fails with "no vault available to run a builtin tool".
	err = s.repos.VaultMembers().Add(context.Background(), testVault.Uuid, userUuid, artel_q.VaultRoleOwner)
	s.Require().NoError(err)

	// 1. Create a tract with a single builtin action step, via the real CreateTract RPC.
	actionStep := &pb.ActionStep{
		Mcp:    "artel",
		Tool:   "list_files",
		Params: map[string]string{},
	}
	step := &pb.TractStep{
		Id:   "list_files",
		Name: "list_files",
		Kind: &pb.TractStep_Action{Action: actionStep},
	}
	definition := &pb.TractDefinition{Steps: []*pb.TractStep{step}}
	createTractReq := &pb.CreateTract_Request{
		Name:       "e2e basic tract",
		Definition: definition,
	}
	createTractResp, err := s.tractClient.CreateTract(authedCtx, createTractReq)
	s.Require().NoError(err)
	s.Empty(createTractResp.Warnings)

	tractUuid := createTractResp.Tract.Uuid
	s.T().Cleanup(func() {
		deleteTractReq := &pb.DeleteTract_Request{Uuid: tractUuid}
		_, _ = s.tractClient.DeleteTract(authedCtx, deleteTractReq)
	})

	// 2. Create and link a manual trigger — the same "Add trigger" flow the Trigger inspector
	// panel drives.
	createTriggerReq := &pb.CreateTrigger_Request{
		Name:          "e2e manual trigger",
		Kind:          "manual",
		Source:        "generic",
		Config:        "{}",
		PayloadSchema: "{}",
	}
	createTriggerResp, err := s.tractClient.CreateTrigger(authedCtx, createTriggerReq)
	s.Require().NoError(err)

	triggerUuid := createTriggerResp.Trigger.Uuid
	s.T().Cleanup(func() {
		deleteTriggerReq := &pb.DeleteTrigger_Request{Uuid: triggerUuid}
		_, _ = s.tractClient.DeleteTrigger(authedCtx, deleteTriggerReq)
	})

	linkTriggerReq := &pb.LinkTrigger_Request{TriggerUuid: triggerUuid, TractUuid: tractUuid}
	_, err = s.tractClient.LinkTrigger(authedCtx, linkTriggerReq)
	s.Require().NoError(err)

	// 3. Run it — the same call the canvas builder's Run button makes. RunTract fires the
	// engine asynchronously against the server-lifecycle context and returns immediately with
	// an empty response (see tracts_api.RunTract's doc comment), unlike calling
	// tract.Service.StartRun directly, which returned the finished run synchronously — so poll
	// ListRuns until the run leaves "running", same pattern gitlab_trigger_e2e uses for its
	// webhook-triggered runs.
	runTractReq := &pb.RunTract_Request{TractUuid: tractUuid, Params: "{}"}
	_, err = s.tractClient.RunTract(authedCtx, runTractReq)
	s.Require().NoError(err)

	var runItem *pb.TractRunItem
	s.Require().Eventually(func() bool {
		listRunsReq := &pb.ListRuns_Request{TractUuid: tractUuid, Limit: 1}
		listRunsResp, listErr := s.tractClient.ListRuns(authedCtx, listRunsReq)
		if listErr != nil || len(listRunsResp.Runs) == 0 {
			return false
		}

		runItem = listRunsResp.Runs[0]

		return runItem.Status != string(domain.TractRunRunning)
	}, 5*time.Second, 50*time.Millisecond, "tract run did not finish")

	s.Equal(string(domain.TractRunDone), runItem.Status)

	// 4. Read the run back — the same call the log panel makes — and assert the step recorded.
	getRunReq := &pb.GetRun_Request{RunUuid: runItem.Uuid}
	getRunResp, err := s.tractClient.GetRun(authedCtx, getRunReq)
	s.Require().NoError(err)
	s.Equal(string(domain.TractRunDone), getRunResp.Run.Status)
	s.Require().Len(getRunResp.Steps, 1)
	s.Equal("list_files", getRunResp.Steps[0].StepId)
	s.Equal(string(domain.TractRunStepDone), getRunResp.Steps[0].Status)
}
