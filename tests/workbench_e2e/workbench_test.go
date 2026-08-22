//go:build e2e

// Package workbench_e2e_test covers the Workbench feature end to end: a real per-(vault, user)
// Docker container/volume lifecycle against test-dockerd (tests/docker-compose.yaml), driven
// through the real VaultsAPI gRPC surface where an RPC exists
// (CreateWorkbench/StartWorkbench/StopWorkbench/DeleteWorkbench) and through the service layer
// directly where it doesn't (GetWorkbench is service-only by design — see
// internal/service/v1/workbench/workbench.go). Every vault member gets their own workbench;
// direct service-layer calls thread the calling user's identity via user_context.WithUserContext
// (see workbenchCtx) the same way the gRPC auth interceptor injects it for RPC calls.
//
// Kept in its own package (rather than tests/e2e) so it doesn't depend on that package's other
// files compiling — mirrors tests/tract_e2e and tests/vault.
//
// Run: docker compose -f tests/docker-compose.yaml up -d --wait
//
//	go test -tags "e2e e2e_bootstrap" -count=1 ./tests/bootstrap/... -run TestEnvSetup
//	go test -tags e2e ./tests/workbench_e2e/...
package workbench_e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	dockerclient "github.com/docker/docker/client"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/clients/workbenchdocker"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/service/v1/workbench"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/vaults_api"
	"github.com/ruf-dev/artel/internal/utils"
	"github.com/ruf-dev/artel/tests/harness"
)

// workbenchCtx returns a context.Background() carrying an injected user_context.UserContext for
// userUuid — needed for tests that call workbench.Service methods directly (GetWorkbench has no
// RPC), since those methods resolve the caller from ctx the same way a gRPC handler's ctx would
// have it injected by the auth interceptor (see internal/service/v1/workbench/workbench.go).
func workbenchCtx(userUuid uuid.UUID) context.Context {
	uc := user_context.UserContext{UserUuid: userUuid}

	return user_context.WithUserContext(context.Background(), uc)
}

// WorkbenchSuite exercises the full workbench lifecycle (create -> start -> stop -> get -> delete)
// against a real Docker daemon (test-dockerd), asserting both artel's own DB-backed state and the
// actual container/volume state on the daemon at each step.
type WorkbenchSuite struct {
	suite.Suite
	db    *sql.DB
	repos *repopg.Repos
	svcs  *svcv1.Services

	dockerHostURL string
	dockerSDK     *dockerclient.Client

	// dockerClient talks to the same dockerHostURL through the typed workbenchdocker.Client
	// wrapper (the exact production WriteFilesToVolume/ReadFilesFromVolume codepath
	// workbench.Service.materializeVault/syncWorkbenchToVault themselves use), so the vault-sync
	// tests can simulate an in-container edit and read the workspace back without going through
	// the service layer twice — dockerSDK above stays for the raw container/volume state
	// assertions the rest of the suite already does.
	dockerClient *workbenchdocker.Client

	authClient  pb.AuthAPIClient
	vaultClient pb.VaultsAPIClient

	// currentVaultUuid/currentUserUuid track the vault and its creator from the most recent
	// setupUser call within a single test, purely so TearDownTest has something to best-effort
	// clean up — not exercised by assertions themselves.
	currentVaultUuid uuid.UUID
	currentUserUuid  uuid.UUID
}

func TestWorkbench(t *testing.T) {
	suite.Run(t, new(WorkbenchSuite))
}

func (s *WorkbenchSuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	// Workbench doesn't gate on subscriptions/plan checks anywhere in its own service code (see
	// internal/service/v1/workbench/workbench.go — no SubscriptionService dependency at all), so
	// unlike QuotaSuite the default config (SubscriptionsEnabled=false, routing through
	// FreeService) is fine: FreeService.CheckActive always returns nil, so the auth interceptor's
	// per-RPC subscription check passes without a subscription row ever being created.
	cfg := config.EnvironmentConfig{}

	var credsEncrypted bool
	s.repos, s.svcs, credsEncrypted = harness.BuildServices(s.T(), s.db, cfg)

	// svcv1.New deliberately leaves Services.Workbench nil — the real app wires it up in
	// internal/app/custom.go, after ExternalConnections exists, rather than inside svcv1.New
	// itself (see that struct's field doc comment). Mirror that exact construction here since
	// harness.BuildServices can't do it generically for every suite.
	workbenchSvc := workbench.New(
		s.repos.Workbenches(), s.repos.Vaults(), s.repos.VaultMembers(), s.repos.DockerHosts(),
		s.svcs.ExternalConnections, s.svcs.Notes, nil,
	)
	s.svcs.Workbench = workbenchSvc

	s.dockerHostURL = harness.DockerHostURL(s.T())
	// Confirms bootstrap already registered this host — fails fast with a clear message
	// otherwise, same contract as GetCouchInstance/GetS3Instance.
	harness.GetDockerHost(s.T(), ctx, s.svcs, s.dockerHostURL)

	dockerSDK, err := dockerclient.NewClientWithOpts(
		dockerclient.WithHost(s.dockerHostURL),
		dockerclient.WithAPIVersionNegotiation(),
	)
	s.Require().NoError(err, "create docker sdk client for direct container/volume assertions")
	s.dockerSDK = dockerSDK

	dockerClient, err := workbenchdocker.New(s.dockerHostURL, workbenchdocker.TLSConfig{})
	s.Require().NoError(err, "create typed workbenchdocker client for volume file assertions")
	s.dockerClient = dockerClient

	s.startGrpcServer(credsEncrypted)
}

// startGrpcServer builds a real *grpc.Server chained with the production auth interceptor,
// registers the AuthAPI and VaultsAPI implementations onto it, and serves it over an in-memory
// bufconn listener — mirrors QuotaSuite.startGrpcServer.
func (s *WorkbenchSuite) startGrpcServer(credsEncrypted bool) {
	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, credsEncrypted, s.svcs.DockerHost, s.svcs.SetupWizard,
	)
	workbenchTerminalShellHandler := vaults_api.NewWorkbenchTerminalShellHandler(
		s.svcs.Auth, s.repos.VaultMembers(), s.svcs.Workbench,
	)
	vaultsImpl := vaults_api.NewVaultsImpl(s.svcs.Vault, s.svcs.Workbench, workbenchTerminalShellHandler)

	conn := harness.NewBufconnServer(s.T(), s.svcs, authImpl.Register, vaultsImpl.Register)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.vaultClient = pb.NewVaultsAPIClient(conn)
}

func (s *WorkbenchSuite) TearDownSuite() {
	if s.dockerSDK != nil {
		_ = s.dockerSDK.Close()
	}

	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

// TearDownTest is a best-effort DeleteWorkbench cleanup in case a test fails partway through the
// lifecycle — test-dockerd persists across `make test-e2e` invocations (same as
// test-couchdb/test-postgres), so a leftover container/volume from a failed run would otherwise
// leak into the next one.
func (s *WorkbenchSuite) TearDownTest() {
	if s.currentVaultUuid == uuid.Nil {
		return
	}

	_ = s.svcs.Workbench.DeleteWorkbench(workbenchCtx(s.currentUserUuid), s.currentVaultUuid)
	s.currentVaultUuid = uuid.Nil
	s.currentUserUuid = uuid.Nil
}

type workbenchTestUser struct {
	userUuid  uuid.UUID
	authedCtx context.Context
	vaultUuid uuid.UUID
}

// setupUser registers a user via the real AuthAPI RPC, logs in, and creates a vault via
// VaultsAPI — the minimum a caller needs before driving the workbench lifecycle. Mirrors
// QuotaSuite.setupUser, minus the MCP-key issuance this suite doesn't need.
//
// suffix disambiguates multiple setupUser calls within the same test (harness.Slug is derived
// purely from t.Name(), so two calls in one test would otherwise collide on the same email and
// the second Register would fail with AlreadyExists) — same pattern as
// tests/e2e/postgres_test.go's setupUser.
func (s *WorkbenchSuite) setupUser(suffix string) workbenchTestUser {
	ctx := context.Background()

	slug := harness.Slug(s.T()) + "_" + suffix
	email := slug + "@test.local"
	password := "test-password-workbench"

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
	// registration never populates domain.User.Username — see harness.GetCouchInstance's callers
	// for the same workaround (tests/e2e/quota_test.go).
	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, slug, userUuid)
	s.Require().NoError(err, "set username stand-in for vault creation")

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

	s.currentVaultUuid = vaultUuid
	s.currentUserUuid = userUuid

	testUser := workbenchTestUser{
		userUuid:  userUuid,
		authedCtx: authedCtx,
		vaultUuid: vaultUuid,
	}

	return testUser
}

// seedAnthropicConnection writes a dummy (non-functional) Anthropic external_connections row
// directly via the repo layer, bypassing externalconnections.Service.AddAnthropicConnection's
// live validation against the Anthropic API — the workbench's own auth-mode logic
// (startWithApiKey) only reads this row back (GetAnthropicApiKey), it never validates it, and
// deploy/workbench/entrypoint.sh keeps the container alive regardless of whether the in-container
// `claude` process can actually authenticate with it. Mirrors
// externalconnections.Service.AddAnthropicConnection's construction of the stored row.
func (s *WorkbenchSuite) seedAnthropicConnection(userUuid uuid.UUID) {
	ctx := context.Background()

	creds := domain.AnthropicKeyCredentials{ApiKey: "sk-ant-e2e-dummy"}

	credJSON, err := json.Marshal(creds)
	s.Require().NoError(err, "marshal dummy anthropic credentials")

	conn := domain.ExternalConnection{
		UserUuid:        userUuid,
		Provider:        domain.ProviderAnthropic,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: credJSON,
	}

	_, err = s.repos.ExternalConnections().Upsert(ctx, conn)
	s.Require().NoError(err, "seed dummy anthropic connection")
}

// containerState returns the live container state for containerID from test-dockerd directly,
// independent of what artel's own workbenches row claims — this is what actually confirms
// artel's Docker orchestration worked, not just that the DB row was written.
func (s *WorkbenchSuite) containerState(containerID string) (running bool, status string) {
	inspect, err := s.dockerSDK.ContainerInspect(context.Background(), containerID)
	s.Require().NoError(err, "inspect workbench container on test-dockerd")
	s.Require().NotNil(inspect.State)

	return inspect.State.Running, inspect.State.Status
}

// volumeExists reports whether a named volume exists on test-dockerd.
func (s *WorkbenchSuite) volumeExists(volumeName string) bool {
	_, err := s.dockerSDK.VolumeInspect(context.Background(), volumeName)

	return err == nil
}

// containerExists reports whether a container exists (in any state) on test-dockerd.
func (s *WorkbenchSuite) containerExists(containerID string) bool {
	_, err := s.dockerSDK.ContainerInspect(context.Background(), containerID)

	return err == nil
}

// TestFullLifecycle drives CreateWorkbench -> StartWorkbench(api_key) -> StopWorkbench ->
// GetWorkbench -> DeleteWorkbench, asserting both artel's reported state and the real
// container/volume state on test-dockerd at every step.
func (s *WorkbenchSuite) TestFullLifecycle() {
	ctx := context.Background()

	user := s.setupUser("user")
	s.seedAnthropicConnection(user.userUuid)

	// Step 1: CreateWorkbench via the real VaultsAPI RPC.
	createReq := &pb.CreateWorkbench_Request{VaultId: user.vaultUuid.String()}

	createResp, err := s.vaultClient.CreateWorkbench(user.authedCtx, createReq)
	s.Require().NoError(err)
	s.Equal(string(domain.WorkbenchStatusCreated), createResp.Status)

	wb, err := s.svcs.Workbench.GetWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().NoError(err)
	s.Equal(domain.WorkbenchStatusCreated, wb.Status)
	s.NotEmpty(wb.ContainerId, "workbench row should carry the created container id")
	s.NotEmpty(wb.VolumeName)

	s.True(s.containerExists(wb.ContainerId), "container should exist on test-dockerd after CreateWorkbench")
	s.True(s.volumeExists(wb.VolumeName), "volume should exist on test-dockerd after CreateWorkbench")

	running, _ := s.containerState(wb.ContainerId)
	s.False(running, "container should not be running yet — CreateWorkbench only creates it")

	// Confirm the network attachment matches what workbenchdocker.CreateContainer wires up.
	inspect, err := s.dockerSDK.ContainerInspect(ctx, wb.ContainerId)
	s.Require().NoError(err)
	s.Require().NotNil(inspect.NetworkSettings)
	_, attached := inspect.NetworkSettings.Networks["workbench-net"]
	s.True(attached, "container should be attached to workbench-net")

	// Step 2: StartWorkbench via the real VaultsAPI RPC, api_key auth mode.
	startReq := &pb.StartWorkbench_Request{
		VaultId:  user.vaultUuid.String(),
		AuthMode: string(domain.WorkbenchAuthModeAPIKey),
	}

	startResp, err := s.vaultClient.StartWorkbench(user.authedCtx, startReq)
	s.Require().NoError(err)
	s.Equal(string(domain.WorkbenchStatusRunning), startResp.Status)
	s.Equal(string(domain.WorkbenchAuthModeAPIKey), startResp.AuthMode)

	running, _ = s.containerState(wb.ContainerId)
	s.True(running, "container should actually be running on test-dockerd after StartWorkbench")

	// The terminal reverse-proxy (internal/transport/vaults_api/workbench_terminal.go) dials
	// ResolveTerminalTarget's address directly from Artel's own process — not from inside
	// test-dockerd. test-dockerd is itself a docker:dind *container* (tests/docker-compose.yaml),
	// so this is the regression check for workbenchdocker.ContainerAddress returning the
	// container's IP on its inner, dind-private workbench-net (unreachable from here, this test
	// process being the same vantage point as Artel's own) instead of the published-port address
	// it actually resolves to — which only works because test-dockerd's compose service publishes
	// a matching fixed host port range for the reverse proxy to dial.
	target, err := s.svcs.Workbench.ResolveTerminalTarget(ctx, user.vaultUuid, user.userUuid)
	s.Require().NoError(err)

	httpClient := http.Client{Timeout: 2 * time.Second}

	// Polled rather than a single request: entrypoint.sh's tmux-session-then-ttyd startup takes a
	// beat after `docker start` returns before ttyd is actually accepting connections — nothing to
	// do with the reachability this test exists to check.
	var terminalResp *http.Response
	s.Require().Eventually(func() bool {
		resp, getErr := httpClient.Get(target)
		if getErr != nil {
			return false
		}

		terminalResp = resp

		return true
	}, 5*time.Second, 200*time.Millisecond,
		"terminal target %q did not become reachable from outside test-dockerd, same as from Artel's own process", target)

	defer utils.CloseWithLog(terminalResp.Body, "workbench terminal probe response body")

	s.Equal(http.StatusOK, terminalResp.StatusCode, "ttyd should serve its index page over the resolved terminal target")

	// Step 3: StopWorkbench via the real VaultsAPI RPC.
	stopReq := &pb.StopWorkbench_Request{VaultId: user.vaultUuid.String()}

	stopResp, err := s.vaultClient.StopWorkbench(user.authedCtx, stopReq)
	s.Require().NoError(err)
	s.Equal(string(domain.WorkbenchStatusStopped), stopResp.Status)

	running, status := s.containerState(wb.ContainerId)
	s.False(running, "container should be stopped on test-dockerd after StopWorkbench")
	s.Equal("exited", status)
	s.True(s.containerExists(wb.ContainerId), "container should still exist (not removed) after StopWorkbench")
	s.True(s.volumeExists(wb.VolumeName), "volume should still exist after StopWorkbench")

	// Step 4: GetWorkbench via the service layer directly (no gRPC RPC exists).
	wb, err = s.svcs.Workbench.GetWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().NoError(err)
	s.Equal(domain.WorkbenchStatusStopped, wb.Status)

	// Step 5: DeleteWorkbench via the service layer directly. (TestDeleteWorkbench_RPC below
	// covers the real DeleteWorkbench RPC instead, so both call paths get e2e coverage.)
	err = s.svcs.Workbench.DeleteWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().NoError(err)

	s.False(s.containerExists(wb.ContainerId), "container should be gone from test-dockerd after DeleteWorkbench")
	s.False(s.volumeExists(wb.VolumeName), "volume should be gone from test-dockerd after DeleteWorkbench")

	_, err = s.svcs.Workbench.GetWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().Error(err, "workbench row should be gone after DeleteWorkbench")
}

// TestCreateWorkbench_IdempotentRetry confirms a second CreateWorkbench call for the same vault
// returns the existing row rather than provisioning a duplicate container/volume — see
// workbench.Service.CreateWorkbench's doc comment ("Safely re-invocable").
func (s *WorkbenchSuite) TestCreateWorkbench_IdempotentRetry() {
	user := s.setupUser("user")

	createReq := &pb.CreateWorkbench_Request{VaultId: user.vaultUuid.String()}

	firstResp, err := s.vaultClient.CreateWorkbench(user.authedCtx, createReq)
	s.Require().NoError(err)

	secondResp, err := s.vaultClient.CreateWorkbench(user.authedCtx, createReq)
	s.Require().NoError(err)

	s.Equal(firstResp.Status, secondResp.Status)

	wb, err := s.svcs.Workbench.GetWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().NoError(err)

	err = s.svcs.Workbench.DeleteWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().NoError(err)

	s.False(s.containerExists(wb.ContainerId))
}

// TestCreateWorkbench_TwoMembersGetDistinctWorkbenches confirms two different members of the
// same vault each calling the real CreateWorkbench RPC get their own, separate workbench row —
// distinct containers/volumes, keyed by (vault, user) rather than just vault — the core behavior
// migrations/072_workbench_per_user.sql introduces.
func (s *WorkbenchSuite) TestCreateWorkbench_TwoMembersGetDistinctWorkbenches() {
	owner := s.setupUser("owner")

	// setupUser always creates its own vault too — unused here, but its registered cleanup still
	// best-effort tears it down at test end, same as any other setupUser call.
	member := s.setupUser("member")

	err := s.repos.VaultMembers().Add(context.Background(), owner.vaultUuid, member.userUuid, artel_q.VaultRoleReader)
	s.Require().NoError(err, "add second user as a member of the owner's vault")

	createReq := &pb.CreateWorkbench_Request{VaultId: owner.vaultUuid.String()}

	ownerResp, err := s.vaultClient.CreateWorkbench(owner.authedCtx, createReq)
	s.Require().NoError(err)

	memberResp, err := s.vaultClient.CreateWorkbench(member.authedCtx, createReq)
	s.Require().NoError(err)

	ownerWb, err := s.svcs.Workbench.GetWorkbench(workbenchCtx(owner.userUuid), owner.vaultUuid)
	s.Require().NoError(err)

	memberWb, err := s.svcs.Workbench.GetWorkbench(workbenchCtx(member.userUuid), owner.vaultUuid)
	s.Require().NoError(err)

	s.Equal(string(domain.WorkbenchStatusCreated), ownerResp.Status)
	s.Equal(string(domain.WorkbenchStatusCreated), memberResp.Status)

	s.NotEqual(ownerWb.Uuid, memberWb.Uuid, "each member should get their own workbench row")
	s.NotEqual(ownerWb.ContainerId, memberWb.ContainerId, "each member's container must be distinct")
	s.NotEqual(ownerWb.VolumeName, memberWb.VolumeName, "each member's volume must be distinct")

	s.True(s.containerExists(ownerWb.ContainerId))
	s.True(s.containerExists(memberWb.ContainerId))

	// TearDownTest only best-effort cleans up the LAST setupUser call's (vault, user) pair
	// (member's), since currentVaultUuid/currentUserUuid get overwritten by the second setupUser
	// call above — clean up both workbenches explicitly here rather than leaking the owner's
	// container/volume on test-dockerd.
	err = s.svcs.Workbench.DeleteWorkbench(workbenchCtx(member.userUuid), owner.vaultUuid)
	s.Require().NoError(err, "clean up the member's workbench")

	err = s.svcs.Workbench.DeleteWorkbench(workbenchCtx(owner.userUuid), owner.vaultUuid)
	s.Require().NoError(err, "clean up the owner's workbench")
}

// TestDeleteWorkbench_RPC drives the real DeleteWorkbench RPC end to end, confirming it tears
// down the container/volume on test-dockerd and the row is gone afterward — TestFullLifecycle
// exercises DeleteWorkbench via the service layer directly, this covers the RPC path.
func (s *WorkbenchSuite) TestDeleteWorkbench_RPC() {
	user := s.setupUser("user")

	createReq := &pb.CreateWorkbench_Request{VaultId: user.vaultUuid.String()}

	_, err := s.vaultClient.CreateWorkbench(user.authedCtx, createReq)
	s.Require().NoError(err)

	wb, err := s.svcs.Workbench.GetWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().NoError(err)
	s.Require().NotEmpty(wb.ContainerId)

	deleteReq := &pb.DeleteWorkbench_Request{VaultId: user.vaultUuid.String()}

	deleteResp, err := s.vaultClient.DeleteWorkbench(user.authedCtx, deleteReq)
	s.Require().NoError(err)
	s.Equal(string(domain.WorkbenchStatusRemoved), deleteResp.Status)

	s.False(s.containerExists(wb.ContainerId), "container should be gone from test-dockerd after DeleteWorkbench RPC")
	s.False(s.volumeExists(wb.VolumeName), "volume should be gone from test-dockerd after DeleteWorkbench RPC")

	_, err = s.svcs.Workbench.GetWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().Error(err, "workbench row should be gone after DeleteWorkbench RPC")
}

// createAndStartWorkbench creates and starts user's workbench in api_key auth mode via the real
// VaultsAPI RPCs — mirrors TestFullLifecycle's own Create+Start sequence, so the vault-sync tests
// below exercise the same call path the rest of the suite already covers rather than inventing a
// second one (subscription_login). Callers must have already seeded an Anthropic connection for
// user.userUuid (seedAnthropicConnection) so startWithApiKey's key resolution succeeds. Returns
// the running workbench's domain row so callers can read its ContainerId.
func (s *WorkbenchSuite) createAndStartWorkbench(user workbenchTestUser) domain.Workbench {
	createReq := &pb.CreateWorkbench_Request{VaultId: user.vaultUuid.String()}

	_, err := s.vaultClient.CreateWorkbench(user.authedCtx, createReq)
	s.Require().NoError(err, "create workbench")

	startReq := &pb.StartWorkbench_Request{
		VaultId:  user.vaultUuid.String(),
		AuthMode: string(domain.WorkbenchAuthModeAPIKey),
	}

	startResp, err := s.vaultClient.StartWorkbench(user.authedCtx, startReq)
	s.Require().NoError(err, "start workbench")
	s.Require().Equal(string(domain.WorkbenchStatusRunning), startResp.Status)

	wb, err := s.svcs.Workbench.GetWorkbench(workbenchCtx(user.userUuid), user.vaultUuid)
	s.Require().NoError(err, "get workbench")
	s.Require().NotEmpty(wb.ContainerId)

	return wb
}

// TestMaterializeOnStart_VaultNotesLandInContainerVolume seeds a note into the vault before the
// workbench exists, then starts the workbench and confirms the seeded note actually landed in the
// container's real Docker volume — exercises workbench.Service.materializeVault's
// ExportFolder->unzip->WriteFilesToVolume path end to end.
func (s *WorkbenchSuite) TestMaterializeOnStart_VaultNotesLandInContainerVolume() {
	ctx := context.Background()
	user := s.setupUser("materialize")
	s.seedAnthropicConnection(user.userUuid)

	notePath := "hello.md"
	noteContent := "# Hello\nseeded before the workbench starts"

	err := s.svcs.Notes.SaveNote(workbenchCtx(user.userUuid), user.vaultUuid, notePath, noteContent)
	s.Require().NoError(err, "seed a note before workbench creation")

	wb := s.createAndStartWorkbench(user)

	files, err := s.dockerClient.ReadFilesFromVolume(ctx, wb.ContainerId)
	s.Require().NoError(err, "read files from the workbench's container volume")

	content, ok := files[notePath]
	s.Require().True(ok, "the seeded note should have been materialized into the workbench volume")
	s.Equal(noteContent, string(content))
}

// TestStopWorkbench_PushesContainerEditBackToVault_NoConflict edits a note's file only inside the
// running workbench's container volume (simulating `claude` editing it), stops the workbench, and
// confirms the edit landed back in the vault's note with no reported conflict — exercises
// workbench.Service.syncWorkbenchToVault's "safe to overwrite" branch (the vault-side copy is
// unchanged since materialization).
func (s *WorkbenchSuite) TestStopWorkbench_PushesContainerEditBackToVault_NoConflict() {
	ctx := context.Background()
	user := s.setupUser("pushback")
	s.seedAnthropicConnection(user.userUuid)

	notePath := "pushback.md"
	original := "# Original\noriginal vault content, unchanged after materialization"

	err := s.svcs.Notes.SaveNote(workbenchCtx(user.userUuid), user.vaultUuid, notePath, original)
	s.Require().NoError(err)

	wb := s.createAndStartWorkbench(user)

	edited := "# Edited\nedited inside the workbench container"
	editedFiles := map[string][]byte{notePath: []byte(edited)}

	err = s.dockerClient.WriteFilesToVolume(ctx, wb.ContainerId, editedFiles)
	s.Require().NoError(err, "simulate an in-container edit via the same tar-write mechanism materializeVault uses")

	stopReq := &pb.StopWorkbench_Request{VaultId: user.vaultUuid.String()}

	stopResp, err := s.vaultClient.StopWorkbench(user.authedCtx, stopReq)
	s.Require().NoError(err, "stop workbench")
	s.Empty(stopResp.SyncConflicts, "an untouched-since-materialization note should sync back without conflict")

	note, err := s.svcs.Notes.GetNote(workbenchCtx(user.userUuid), user.vaultUuid, notePath)
	s.Require().NoError(err)
	s.Equal(edited, note.Content, "the vault's note should now hold the container-side edit")
}

// TestStopWorkbench_ConcurrentEditConflict_VaultSideNotClobbered edits the same note both inside
// the workbench's container volume AND via notes.Service.SaveNote (simulating a concurrent edit
// through the Obsidian app) with genuinely different content in each place, then stops the
// workbench and confirms: the path is reported as a conflict, and the vault's note still holds the
// vault-side edit — it must NOT be clobbered by the container's edit. Exercises
// workbench.Service.syncWorkbenchToVault's conflict-detection branch (the vault-side mtime moved
// past the materialization snapshot, and the two contents genuinely diverged).
func (s *WorkbenchSuite) TestStopWorkbench_ConcurrentEditConflict_VaultSideNotClobbered() {
	ctx := context.Background()
	user := s.setupUser("conflict")
	s.seedAnthropicConnection(user.userUuid)

	notePath := "conflict.md"
	original := "# Original\noriginal vault content"

	err := s.svcs.Notes.SaveNote(workbenchCtx(user.userUuid), user.vaultUuid, notePath, original)
	s.Require().NoError(err)

	wb := s.createAndStartWorkbench(user)

	containerEdit := "# Container edit\nedited inside the workbench, never reaches the vault"
	err = s.dockerClient.WriteFilesToVolume(ctx, wb.ContainerId, map[string][]byte{notePath: []byte(containerEdit)})
	s.Require().NoError(err, "simulate an in-container edit")

	vaultEdit := "# Vault edit\nedited concurrently via the Obsidian app while the workbench was running"

	err = s.svcs.Notes.SaveNote(workbenchCtx(user.userUuid), user.vaultUuid, notePath, vaultEdit)
	s.Require().NoError(err, "simulate a concurrent vault-side edit")

	stopReq := &pb.StopWorkbench_Request{VaultId: user.vaultUuid.String()}

	stopResp, err := s.vaultClient.StopWorkbench(user.authedCtx, stopReq)
	s.Require().NoError(err, "stop workbench")
	s.Contains(stopResp.SyncConflicts, notePath, "diverging concurrent edits should be reported as a conflict")

	note, err := s.svcs.Notes.GetNote(workbenchCtx(user.userUuid), user.vaultUuid, notePath)
	s.Require().NoError(err)
	s.Equal(vaultEdit, note.Content, "the vault-side edit must not be clobbered by the container's conflicting edit")
}
