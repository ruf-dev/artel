package workbench

import (
	"context"
	"database/sql"
	"testing"

	"github.com/containerd/errdefs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/clients/workbenchdocker"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// fakeWorkbenches is a hand-written repository.Workbenches for exercising Service's branching
// logic without a live Postgres.
type fakeWorkbenches struct {
	getByVaultID         func(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error)
	create               func(ctx context.Context, vaultID, userID uuid.UUID, volumeName string, dockerHostID uuid.UUID) (domain.Workbench, error)
	markContainerCreated func(ctx context.Context, vaultID uuid.UUID, containerID string) error
	markConfiguring      func(ctx context.Context, vaultID uuid.UUID) error
	markRunning          func(ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode) error
	markStopped          func(ctx context.Context, vaultID uuid.UUID) error
	delete               func(ctx context.Context, vaultID uuid.UUID) error
}

func (f *fakeWorkbenches) Create(
	ctx context.Context,
	vaultID, userID uuid.UUID,
	volumeName string,
	dockerHostID uuid.UUID,
) (domain.Workbench, error) {
	if f.create == nil {
		return domain.Workbench{}, nil
	}

	return f.create(ctx, vaultID, userID, volumeName, dockerHostID)
}

func (f *fakeWorkbenches) GetByVaultID(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	return f.getByVaultID(ctx, vaultID)
}

func (f *fakeWorkbenches) MarkContainerCreated(ctx context.Context, vaultID uuid.UUID, containerID string) error {
	if f.markContainerCreated == nil {
		return nil
	}

	return f.markContainerCreated(ctx, vaultID, containerID)
}

func (f *fakeWorkbenches) MarkConfiguring(ctx context.Context, vaultID uuid.UUID) error {
	if f.markConfiguring == nil {
		return nil
	}

	return f.markConfiguring(ctx, vaultID)
}

func (f *fakeWorkbenches) MarkRunning(ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode) error {
	if f.markRunning == nil {
		return nil
	}

	return f.markRunning(ctx, vaultID, authMode)
}

func (f *fakeWorkbenches) MarkStopped(ctx context.Context, vaultID uuid.UUID) error {
	if f.markStopped == nil {
		return nil
	}

	return f.markStopped(ctx, vaultID)
}

func (f *fakeWorkbenches) MarkRemoved(context.Context, uuid.UUID) error {
	return nil
}

func (f *fakeWorkbenches) Delete(ctx context.Context, vaultID uuid.UUID) error {
	if f.delete == nil {
		return nil
	}

	return f.delete(ctx, vaultID)
}

func (f *fakeWorkbenches) WithTx(sqldb.DB) repository.Workbenches {
	return f
}

// fakeVaults is a hand-written repository.Vaults exposing only GetByID's behavior as
// configurable — the rest of the interface is unused by Service.
type fakeVaults struct {
	getByID func(ctx context.Context, id uuid.UUID) (domain.Vault, error)
}

func (f *fakeVaults) Upsert(context.Context, uuid.UUID, uuid.UUID, string, string, string, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) GetByID(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
	return f.getByID(ctx, id)
}

func (f *fakeVaults) GetByNameAndUser(context.Context, uuid.UUID, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) UpdateStatus(context.Context, uuid.UUID, string) error {
	return nil
}

func (f *fakeVaults) SetLiveSyncPassphrase(context.Context, uuid.UUID, string) error {
	return nil
}

func (f *fakeVaults) ListByMembership(context.Context, uuid.UUID) ([]domain.Vault, error) {
	return nil, nil
}

func (f *fakeVaults) Delete(context.Context, uuid.UUID) error {
	return nil
}

func (f *fakeVaults) LinkS3Bucket(context.Context, uuid.UUID, uuid.UUID, string) error {
	return nil
}

func (f *fakeVaults) UnlinkS3Bucket(context.Context, uuid.UUID) error {
	return nil
}

func (f *fakeVaults) SetUseCouchDBForBinaries(context.Context, uuid.UUID, bool) error {
	return nil
}

func (f *fakeVaults) WithTx(sqldb.DB) repository.Vaults {
	return f
}

// fakeDockerHosts is a hand-written repository.DockerHosts for exercising Service's per-workbench
// docker host resolution without a live Postgres. Only Get/PickLeastLoaded are exercised by
// workbench.Service today; the rest are stubbed trivially since these tests never drive them.
type fakeDockerHosts struct {
	register        func(ctx context.Context, url, caCert, clientCert, clientKey string) (uuid.UUID, error)
	get             func(ctx context.Context, id uuid.UUID) (domain.DockerHost, error)
	getWithCreds    func(ctx context.Context, id uuid.UUID) (domain.DockerHost, error)
	list            func(ctx context.Context) ([]domain.DockerHost, error)
	update          func(ctx context.Context, id uuid.UUID, url string, caCert, clientCert, clientKey *string) error
	delete          func(ctx context.Context, id uuid.UUID) error
	exists          func(ctx context.Context) (bool, error)
	pickLeastLoaded func(ctx context.Context) (domain.DockerHost, error)
}

func (f *fakeDockerHosts) Register(ctx context.Context, url, caCert, clientCert, clientKey string) (uuid.UUID, error) {
	if f.register == nil {
		return uuid.New(), nil
	}

	return f.register(ctx, url, caCert, clientCert, clientKey)
}

func (f *fakeDockerHosts) Get(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
	if f.get == nil {
		return domain.DockerHost{Uuid: id}, nil
	}

	return f.get(ctx, id)
}

// GetWithCreds falls back to Get's behavior/fixture when getWithCreds isn't set — most tests
// here don't care about TLS material, only that resolveClient resolves the right host, so they
// only wire up get/pickLeastLoaded and expect GetWithCreds to resolve the same way.
func (f *fakeDockerHosts) GetWithCreds(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
	if f.getWithCreds == nil {
		return f.Get(ctx, id)
	}

	return f.getWithCreds(ctx, id)
}

func (f *fakeDockerHosts) List(ctx context.Context) ([]domain.DockerHost, error) {
	if f.list == nil {
		return nil, nil
	}

	return f.list(ctx)
}

func (f *fakeDockerHosts) Update(ctx context.Context, id uuid.UUID, url string, caCert, clientCert, clientKey *string) error {
	if f.update == nil {
		return nil
	}

	return f.update(ctx, id, url, caCert, clientCert, clientKey)
}

func (f *fakeDockerHosts) Delete(ctx context.Context, id uuid.UUID) error {
	if f.delete == nil {
		return nil
	}

	return f.delete(ctx, id)
}

func (f *fakeDockerHosts) Exists(ctx context.Context) (bool, error) {
	if f.exists == nil {
		return false, nil
	}

	return f.exists(ctx)
}

func (f *fakeDockerHosts) PickLeastLoaded(ctx context.Context) (domain.DockerHost, error) {
	if f.pickLeastLoaded == nil {
		return domain.DockerHost{}, nil
	}

	return f.pickLeastLoaded(ctx)
}

func (f *fakeDockerHosts) WithTx(sqldb.DB) repository.DockerHosts {
	return f
}

// defaultDockerHostID/defaultDockerHost back newFakeDockerHosts's default single-host pool —
// shared read-only across tests, so a workbench fixture that needs *some* valid DockerHostUuid
// (most of the docker-op tests below) can just point at this one without every test minting its
// own host.
var defaultDockerHostID = uuid.New()

var defaultDockerHost = domain.DockerHost{Uuid: defaultDockerHostID, Url: "docker://default-host"}

// newFakeDockerHosts returns a fakeDockerHosts backed by a single always-available host
// (defaultDockerHost) — the common case for tests that don't care which host handles a
// workbench, just that resolution succeeds.
func newFakeDockerHosts() *fakeDockerHosts {
	return &fakeDockerHosts{
		get: func(context.Context, uuid.UUID) (domain.DockerHost, error) {
			return defaultDockerHost, nil
		},
		pickLeastLoaded: func(context.Context) (domain.DockerHost, error) {
			return defaultDockerHost, nil
		},
	}
}

// fakeDocker is a hand-written dockerClient. Every method defaults to a no-op success; tests
// override only the fields they need to exercise.
type fakeDocker struct {
	createVolume    func(ctx context.Context, name string) error
	removeVolume    func(ctx context.Context, name string) error
	createContainer func(ctx context.Context, opts workbenchdocker.CreateOpts) (string, error)
	startContainer  func(ctx context.Context, containerID string, env map[string]string) error
	stopContainer   func(ctx context.Context, containerID string) error
	removeContainer func(ctx context.Context, containerID string) error
	capturePane     func(ctx context.Context, containerID string) (string, error)
	sendKeys        func(ctx context.Context, containerID string, keys string) error
}

func newFakeDocker() *fakeDocker {
	return &fakeDocker{
		createVolume: func(context.Context, string) error { return nil },
		removeVolume: func(context.Context, string) error { return nil },
		createContainer: func(context.Context, workbenchdocker.CreateOpts) (string, error) {
			return "container-1", nil
		},
		startContainer:  func(context.Context, string, map[string]string) error { return nil },
		stopContainer:   func(context.Context, string) error { return nil },
		removeContainer: func(context.Context, string) error { return nil },
		capturePane:     func(context.Context, string) (string, error) { return "", nil },
		sendKeys:        func(context.Context, string, string) error { return nil },
	}
}

func (f *fakeDocker) CreateVolume(ctx context.Context, name string) error {
	return f.createVolume(ctx, name)
}

func (f *fakeDocker) RemoveVolume(ctx context.Context, name string) error {
	return f.removeVolume(ctx, name)
}

func (f *fakeDocker) CreateContainer(ctx context.Context, opts workbenchdocker.CreateOpts) (string, error) {
	return f.createContainer(ctx, opts)
}

func (f *fakeDocker) StartContainer(ctx context.Context, containerID string, env map[string]string) error {
	return f.startContainer(ctx, containerID, env)
}

func (f *fakeDocker) StopContainer(ctx context.Context, containerID string) error {
	return f.stopContainer(ctx, containerID)
}

func (f *fakeDocker) RemoveContainer(ctx context.Context, containerID string) error {
	return f.removeContainer(ctx, containerID)
}

func (f *fakeDocker) CapturePane(ctx context.Context, containerID string) (string, error) {
	return f.capturePane(ctx, containerID)
}

func (f *fakeDocker) SendKeys(ctx context.Context, containerID string, keys string) error {
	return f.sendKeys(ctx, containerID, keys)
}

// fakeDockerClientFactory adapts a single fakeDocker into a Service.newDockerClient func,
// ignoring the host/tlsCfg arguments — the common case for tests where only one docker host is
// in play.
func fakeDockerClientFactory(d *fakeDocker) func(string, workbenchdocker.TLSConfig) (dockerClient, error) {
	return func(string, workbenchdocker.TLSConfig) (dockerClient, error) {
		return d, nil
	}
}

// fakeExternalConnections is a hand-written externalConnectionService for exercising
// StartWorkbench's api_key resolution without a live external-connections service.
type fakeExternalConnections struct {
	getAnthropicApiKey func(ctx context.Context, userUuid uuid.UUID) (string, error)
}

func (f *fakeExternalConnections) GetAnthropicApiKey(ctx context.Context, userUuid uuid.UUID) (string, error) {
	return f.getAnthropicApiKey(ctx, userUuid)
}

func newFakeExternalConnections() *fakeExternalConnections {
	return &fakeExternalConnections{
		getAnthropicApiKey: func(context.Context, uuid.UUID) (string, error) {
			return "sk-ant-fake", nil
		},
	}
}

var errBoom = rerrors.New("boom")

func TestCreateWorkbench_Idempotent_ReturnsExisting(t *testing.T) {
	vaultID := uuid.New()
	existing := domain.Workbench{Uuid: uuid.New(), VaultUuid: vaultID, VolumeName: "workbench-existing"}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return existing, nil
		},
		create: func(context.Context, uuid.UUID, uuid.UUID, string, uuid.UUID) (domain.Workbench, error) {
			t.Fatal("Create should not be called when a workbench already exists")
			return domain.Workbench{}, nil
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			t.Fatal("vaults.GetByID should not be called when a workbench already exists")
			return domain.Vault{}, nil
		},
	}
	docker := newFakeDocker()
	docker.createVolume = func(context.Context, string) error {
		t.Fatal("docker.CreateVolume should not be called when a workbench already exists")
		return nil
	}

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	got, err := svc.CreateWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, existing, got)
}

func TestCreateWorkbench_HappyPath(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()
	vault := domain.Vault{Uuid: vaultID, UserUuid: userID}

	var gotCreateVolumeName string
	var gotCreateOpts workbenchdocker.CreateOpts
	var gotCreateVaultID, gotCreateUserID, gotCreateDockerHostID uuid.UUID
	var gotCreateVolumeParam string
	var gotMarkContainerCreatedVaultID uuid.UUID
	var gotMarkContainerCreatedContainerID string
	var gotNewDockerClientHost string

	configuring := domain.Workbench{Uuid: uuid.New(), VaultUuid: vaultID, UserUuid: userID, Status: domain.WorkbenchStatusConfiguring}
	final := domain.Workbench{
		Uuid: uuid.New(), VaultUuid: vaultID, UserUuid: userID, Status: domain.WorkbenchStatusCreated,
		ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID,
	}

	getByVaultIDCallCount := 0
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			getByVaultIDCallCount++
			if getByVaultIDCallCount == 1 {
				return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "error getting workbench by vault id")
			}

			return final, nil
		},
		create: func(_ context.Context, vaultID, userID uuid.UUID, volumeName string, dockerHostID uuid.UUID) (domain.Workbench, error) {
			gotCreateVaultID = vaultID
			gotCreateUserID = userID
			gotCreateVolumeParam = volumeName
			gotCreateDockerHostID = dockerHostID

			return configuring, nil
		},
		markContainerCreated: func(_ context.Context, vaultID uuid.UUID, containerID string) error {
			gotMarkContainerCreatedVaultID = vaultID
			gotMarkContainerCreatedContainerID = containerID

			return nil
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return vault, nil
		},
	}

	docker := newFakeDocker()
	docker.createVolume = func(_ context.Context, name string) error {
		gotCreateVolumeName = name
		return nil
	}
	docker.createContainer = func(_ context.Context, opts workbenchdocker.CreateOpts) (string, error) {
		gotCreateOpts = opts
		return "container-1", nil
	}

	newDockerClient := func(host string, _ workbenchdocker.TLSConfig) (dockerClient, error) {
		gotNewDockerClientHost = host
		return docker, nil
	}

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), newDockerClient)

	got, err := svc.CreateWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, final, got)

	expectedVolumeName := "workbench-" + vaultID.String()
	require.Equal(t, expectedVolumeName, gotCreateVolumeName)
	require.Equal(t, expectedVolumeName, gotCreateOpts.Name)
	require.Equal(t, expectedVolumeName, gotCreateOpts.VolumeName)
	require.Equal(t, vaultID, gotCreateVaultID)
	require.Equal(t, userID, gotCreateUserID)
	require.Equal(t, expectedVolumeName, gotCreateVolumeParam)
	require.Equal(t, vaultID, gotMarkContainerCreatedVaultID)
	require.Equal(t, "container-1", gotMarkContainerCreatedContainerID)
	require.Equal(t, 2, getByVaultIDCallCount)
	require.Equal(t, defaultDockerHostID, gotCreateDockerHostID, "the least-loaded host's uuid must be threaded into repo.Create")
	require.Equal(t, defaultDockerHost.Url, gotNewDockerClientHost, "CreateVolume/CreateContainer's client must be built from the picked host's url")
}

func TestCreateWorkbench_NoDockerHostsAvailable_ReturnsClearError(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
		create: func(context.Context, uuid.UUID, uuid.UUID, string, uuid.UUID) (domain.Workbench, error) {
			t.Fatal("Create should not be called when no docker hosts are registered")
			return domain.Workbench{}, nil
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, nil
		},
	}
	dockerHosts := &fakeDockerHosts{
		// A real DockerHosts.PickLeastLoaded runs its sqlc :one query through
		// pg_err.UnwrapPgErr (see internal/repository/pg/repos/dockerhosts/dockerhosts.go),
		// which turns sql.ErrNoRows into user_errors.NotFound rather than surfacing
		// sql.ErrNoRows itself — so that's what an empty pool actually looks like to the
		// service, and what's simulated here.
		pickLeastLoaded: func(context.Context) (domain.DockerHost, error) {
			return domain.DockerHost{}, user_errors.NotFound
		},
	}

	svc := New(workbenches, vaults, dockerHosts, newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.ErrorIs(t, err, user_errors.NoDockerHostsAvailable)
}

func TestCreateWorkbench_GetByVaultIDUnexpectedError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}
	vaults := &fakeVaults{}
	docker := newFakeDocker()

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_VaultLookupError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, errBoom
		},
	}
	docker := newFakeDocker()

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_RepoCreateError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
		create: func(context.Context, uuid.UUID, uuid.UUID, string, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, nil
		},
	}
	docker := newFakeDocker()
	docker.createVolume = func(context.Context, string) error {
		t.Fatal("docker.CreateVolume should not be called when the repo insert fails")
		return nil
	}

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_CreateVolumeError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
		markContainerCreated: func(context.Context, uuid.UUID, string) error {
			t.Fatal("MarkContainerCreated should not be called when CreateVolume fails")
			return nil
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, nil
		},
	}
	docker := newFakeDocker()
	docker.createVolume = func(context.Context, string) error {
		return errBoom
	}
	docker.createContainer = func(context.Context, workbenchdocker.CreateOpts) (string, error) {
		t.Fatal("CreateContainer should not be called when CreateVolume fails")
		return "", nil
	}

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_CreateContainerError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
		markContainerCreated: func(context.Context, uuid.UUID, string) error {
			t.Fatal("MarkContainerCreated should not be called when CreateContainer fails")
			return nil
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, nil
		},
	}
	docker := newFakeDocker()
	docker.createContainer = func(context.Context, workbenchdocker.CreateOpts) (string, error) {
		return "", errBoom
	}

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_MarkContainerCreatedError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
		markContainerCreated: func(context.Context, uuid.UUID, string) error {
			return errBoom
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, nil
		},
	}
	docker := newFakeDocker()

	svc := New(workbenches, vaults, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestGetWorkbench_Passthrough(t *testing.T) {
	vaultID := uuid.New()
	want := domain.Workbench{Uuid: uuid.New(), VaultUuid: vaultID}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(_ context.Context, id uuid.UUID) (domain.Workbench, error) {
			require.Equal(t, vaultID, id)
			return want, nil
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	got, err := svc.GetWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestGetWorkbench_Error_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	_, err := svc.GetWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestStartWorkbench_UnknownMode_ReturnsNotImplemented(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			t.Fatal("GetByVaultID should not be called for an unimplemented auth mode")
			return domain.Workbench{}, nil
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	_, err := svc.StartWorkbench(context.Background(), uuid.New(), domain.WorkbenchAuthMode("bogus"))
	require.ErrorIs(t, err, user_errors.WorkbenchAuthModeNotImplemented)
}

func TestStartWorkbench_SubscriptionLogin_HappyPath(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}
	running := domain.Workbench{
		VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID,
		Status: domain.WorkbenchStatusRunning, AuthMode: domain.WorkbenchAuthModeSubscriptionLogin,
	}

	var gotStartContainerID string
	var gotEnv map[string]string
	var gotMarkConfiguringVaultID uuid.UUID
	var gotMarkRunningVaultID uuid.UUID
	var gotMarkRunningAuthMode domain.WorkbenchAuthMode
	markConfiguringCalled := false

	callCount := 0
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			callCount++
			if callCount == 1 {
				return wb, nil
			}

			return running, nil
		},
		markConfiguring: func(_ context.Context, vaultID uuid.UUID) error {
			markConfiguringCalled = true
			gotMarkConfiguringVaultID = vaultID
			return nil
		},
		markRunning: func(_ context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode) error {
			require.True(t, markConfiguringCalled, "MarkConfiguring must be called before MarkRunning")
			gotMarkRunningVaultID = vaultID
			gotMarkRunningAuthMode = authMode
			return nil
		},
	}

	docker := newFakeDocker()
	docker.startContainer = func(_ context.Context, containerID string, env map[string]string) error {
		require.True(t, markConfiguringCalled, "MarkConfiguring must be called before StartContainer")
		gotStartContainerID = containerID
		gotEnv = env
		return nil
	}

	externalConnections := newFakeExternalConnections()
	externalConnections.getAnthropicApiKey = func(context.Context, uuid.UUID) (string, error) {
		t.Fatal("GetAnthropicApiKey should not be called for subscription_login mode")
		return "", nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), externalConnections, fakeDockerClientFactory(docker))

	got, err := svc.StartWorkbench(context.Background(), vaultID, domain.WorkbenchAuthModeSubscriptionLogin)
	require.NoError(t, err)
	require.Equal(t, running, got)
	require.Equal(t, "container-1", gotStartContainerID)
	require.Nil(t, gotEnv)
	require.Equal(t, vaultID, gotMarkConfiguringVaultID)
	require.Equal(t, vaultID, gotMarkRunningVaultID)
	require.Equal(t, domain.WorkbenchAuthModeSubscriptionLogin, gotMarkRunningAuthMode)
}

func TestStartWorkbench_SubscriptionLogin_StartContainerError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
		markRunning: func(context.Context, uuid.UUID, domain.WorkbenchAuthMode) error {
			t.Fatal("MarkRunning should not be called when StartContainer fails")
			return nil
		},
	}

	docker := newFakeDocker()
	docker.startContainer = func(context.Context, string, map[string]string) error {
		return errBoom
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.StartWorkbench(context.Background(), uuid.New(), domain.WorkbenchAuthModeSubscriptionLogin)
	require.Error(t, err)
}

func TestStartWorkbench_SubscriptionLogin_MarkConfiguringError_ShortCircuits(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
		markConfiguring: func(context.Context, uuid.UUID) error {
			return errBoom
		},
	}

	docker := newFakeDocker()
	docker.startContainer = func(context.Context, string, map[string]string) error {
		t.Fatal("StartContainer should not be called when MarkConfiguring fails")
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.StartWorkbench(context.Background(), uuid.New(), domain.WorkbenchAuthModeSubscriptionLogin)
	require.Error(t, err)
}

func TestStartWorkbench_ApiKey_HappyPath(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, UserUuid: userID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}
	running := domain.Workbench{
		VaultUuid: vaultID, UserUuid: userID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID,
		Status: domain.WorkbenchStatusRunning,
	}

	var gotAnthropicUserID uuid.UUID
	var gotStartContainerID string
	var gotEnv map[string]string
	var gotMarkConfiguringVaultID uuid.UUID
	var gotMarkRunningVaultID uuid.UUID
	var gotMarkRunningAuthMode domain.WorkbenchAuthMode
	markConfiguringCalled := false

	callCount := 0
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			callCount++
			if callCount == 1 {
				return wb, nil
			}

			return running, nil
		},
		markConfiguring: func(_ context.Context, vaultID uuid.UUID) error {
			markConfiguringCalled = true
			gotMarkConfiguringVaultID = vaultID
			return nil
		},
		markRunning: func(_ context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode) error {
			require.True(t, markConfiguringCalled, "MarkConfiguring must be called before MarkRunning")
			gotMarkRunningVaultID = vaultID
			gotMarkRunningAuthMode = authMode
			return nil
		},
	}

	docker := newFakeDocker()
	docker.startContainer = func(_ context.Context, containerID string, env map[string]string) error {
		require.True(t, markConfiguringCalled, "MarkConfiguring must be called before StartContainer")
		gotStartContainerID = containerID
		gotEnv = env
		return nil
	}

	externalConnections := newFakeExternalConnections()
	externalConnections.getAnthropicApiKey = func(_ context.Context, userUuid uuid.UUID) (string, error) {
		gotAnthropicUserID = userUuid
		return "sk-ant-test", nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), externalConnections, fakeDockerClientFactory(docker))

	got, err := svc.StartWorkbench(context.Background(), vaultID, domain.WorkbenchAuthModeAPIKey)
	require.NoError(t, err)
	require.Equal(t, running, got)
	require.Equal(t, userID, gotAnthropicUserID)
	require.Equal(t, "container-1", gotStartContainerID)
	require.Equal(t, map[string]string{"ANTHROPIC_API_KEY": "sk-ant-test"}, gotEnv)
	require.Equal(t, vaultID, gotMarkConfiguringVaultID)
	require.Equal(t, vaultID, gotMarkRunningVaultID)
	require.Equal(t, domain.WorkbenchAuthModeAPIKey, gotMarkRunningAuthMode)
}

func TestStartWorkbench_GetByVaultIDError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	_, err := svc.StartWorkbench(context.Background(), uuid.New(), domain.WorkbenchAuthModeAPIKey)
	require.Error(t, err)
}

func TestStartWorkbench_NoAnthropicConnection_ReturnsClearError(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	externalConnections := newFakeExternalConnections()
	externalConnections.getAnthropicApiKey = func(context.Context, uuid.UUID) (string, error) {
		return "", user_errors.LlmKeyRequired
	}

	docker := newFakeDocker()
	docker.startContainer = func(context.Context, string, map[string]string) error {
		t.Fatal("StartContainer should not be called when no anthropic connection is found")
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), externalConnections, fakeDockerClientFactory(docker))

	_, err := svc.StartWorkbench(context.Background(), uuid.New(), domain.WorkbenchAuthModeAPIKey)
	require.ErrorIs(t, err, user_errors.WorkbenchMissingAnthropicConnection)
}

func TestStartWorkbench_StartContainerError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
		markRunning: func(context.Context, uuid.UUID, domain.WorkbenchAuthMode) error {
			t.Fatal("MarkRunning should not be called when StartContainer fails")
			return nil
		},
	}

	docker := newFakeDocker()
	docker.startContainer = func(context.Context, string, map[string]string) error {
		return errBoom
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.StartWorkbench(context.Background(), uuid.New(), domain.WorkbenchAuthModeAPIKey)
	require.Error(t, err)
}

func TestStartWorkbench_ApiKey_MarkConfiguringError_ShortCircuits(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
		markConfiguring: func(context.Context, uuid.UUID) error {
			return errBoom
		},
	}

	docker := newFakeDocker()
	docker.startContainer = func(context.Context, string, map[string]string) error {
		t.Fatal("StartContainer should not be called when MarkConfiguring fails")
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.StartWorkbench(context.Background(), uuid.New(), domain.WorkbenchAuthModeAPIKey)
	require.Error(t, err)
}

func TestStopWorkbench_HappyPath(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}

	var stoppedID string
	var markedStoppedVaultID uuid.UUID

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
		markStopped: func(_ context.Context, vaultID uuid.UUID) error {
			markedStoppedVaultID = vaultID
			return nil
		},
	}

	docker := newFakeDocker()
	docker.stopContainer = func(_ context.Context, id string) error {
		stoppedID = id
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.StopWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, wb.ContainerId, stoppedID)
	require.Equal(t, vaultID, markedStoppedVaultID)
}

func TestStopWorkbench_GetByVaultIDError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	err := svc.StopWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestStopWorkbench_StopContainerError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
		markStopped: func(context.Context, uuid.UUID) error {
			t.Fatal("MarkStopped should not be called when StopContainer fails")
			return nil
		},
	}

	docker := newFakeDocker()
	docker.stopContainer = func(context.Context, string) error {
		return errBoom
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.StopWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

// TestStopWorkbench_ResolvesClientFromWorkbenchsOwnDockerHost pins two workbenches to two
// distinguishable docker hosts and asserts each one's StopContainer call lands on the docker
// client built from *its own* DockerHostUuid, not a single shared client — i.e. resolveClient
// genuinely branches per-workbench rather than always resolving the same host.
func TestStopWorkbench_ResolvesClientFromWorkbenchsOwnDockerHost(t *testing.T) {
	hostAID := uuid.New()
	hostBID := uuid.New()
	hostA := domain.DockerHost{Uuid: hostAID, Url: "docker://host-a"}
	hostB := domain.DockerHost{Uuid: hostBID, Url: "docker://host-b"}

	dockerHosts := &fakeDockerHosts{
		get: func(_ context.Context, id uuid.UUID) (domain.DockerHost, error) {
			switch id {
			case hostAID:
				return hostA, nil
			case hostBID:
				return hostB, nil
			default:
				t.Fatalf("unexpected docker host id %s", id)
				return domain.DockerHost{}, nil
			}
		},
	}

	dockerA := newFakeDocker()
	dockerB := newFakeDocker()

	var stoppedOnA, stoppedOnB bool
	dockerA.stopContainer = func(context.Context, string) error {
		stoppedOnA = true
		return nil
	}
	dockerB.stopContainer = func(context.Context, string) error {
		stoppedOnB = true
		return nil
	}

	newDockerClient := func(host string, _ workbenchdocker.TLSConfig) (dockerClient, error) {
		switch host {
		case hostA.Url:
			return dockerA, nil
		case hostB.Url:
			return dockerB, nil
		default:
			t.Fatalf("unexpected host url %q", host)
			return nil, nil
		}
	}

	vaultIDOnA := uuid.New()
	wbOnA := domain.Workbench{VaultUuid: vaultIDOnA, ContainerId: "container-a", DockerHostUuid: &hostAID}
	workbenchesOnA := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wbOnA, nil
		},
	}

	svcOnA := New(workbenchesOnA, &fakeVaults{}, dockerHosts, newFakeExternalConnections(), newDockerClient)

	err := svcOnA.StopWorkbench(context.Background(), vaultIDOnA)
	require.NoError(t, err)
	require.True(t, stoppedOnA, "StopContainer must be called on the client built from the workbench's own (host A) docker host")
	require.False(t, stoppedOnB, "the host B client must not be touched by a workbench pinned to host A")

	vaultIDOnB := uuid.New()
	wbOnB := domain.Workbench{VaultUuid: vaultIDOnB, ContainerId: "container-b", DockerHostUuid: &hostBID}
	workbenchesOnB := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wbOnB, nil
		},
	}

	svcOnB := New(workbenchesOnB, &fakeVaults{}, dockerHosts, newFakeExternalConnections(), newDockerClient)

	err = svcOnB.StopWorkbench(context.Background(), vaultIDOnB)
	require.NoError(t, err)
	require.True(t, stoppedOnB, "StopContainer must be called on the client built from the workbench's own (host B) docker host")
}

func TestDeleteWorkbench_NotFound_IsNoop(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
	}
	docker := newFakeDocker()
	docker.stopContainer = func(context.Context, string) error {
		t.Fatal("StopContainer should not be called when the workbench row is already gone")
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.DeleteWorkbench(context.Background(), uuid.New())
	require.NoError(t, err)
}

func TestDeleteWorkbench_HappyPath_StopsRemovesAndDeletes(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", VolumeName: "workbench-vol", DockerHostUuid: &defaultDockerHostID}

	var stoppedID, removedContainerID, removedVolumeName string
	var deletedVaultID uuid.UUID

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
		delete: func(_ context.Context, id uuid.UUID) error {
			deletedVaultID = id
			return nil
		},
	}

	docker := newFakeDocker()
	docker.stopContainer = func(_ context.Context, id string) error {
		stoppedID = id
		return nil
	}
	docker.removeContainer = func(_ context.Context, id string) error {
		removedContainerID = id
		return nil
	}
	docker.removeVolume = func(_ context.Context, name string) error {
		removedVolumeName = name
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, wb.ContainerId, stoppedID)
	require.Equal(t, wb.ContainerId, removedContainerID)
	require.Equal(t, wb.VolumeName, removedVolumeName)
	require.Equal(t, vaultID, deletedVaultID)
}

func TestDeleteWorkbench_NoContainerYet_SkipsContainerCalls(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "", VolumeName: "workbench-vol", DockerHostUuid: &defaultDockerHostID}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
	}

	docker := newFakeDocker()
	docker.stopContainer = func(context.Context, string) error {
		t.Fatal("StopContainer should not be called when ContainerId is empty")
		return nil
	}
	docker.removeContainer = func(context.Context, string) error {
		t.Fatal("RemoveContainer should not be called when ContainerId is empty")
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
}

func TestDeleteWorkbench_StopContainerNotFound_ContinuesTeardown(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", VolumeName: "workbench-vol", DockerHostUuid: &defaultDockerHostID}

	var removedContainerID string
	var deleted bool

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
		delete: func(context.Context, uuid.UUID) error {
			deleted = true
			return nil
		},
	}

	docker := newFakeDocker()
	docker.stopContainer = func(context.Context, string) error {
		return errdefs.ErrNotFound
	}
	docker.removeContainer = func(_ context.Context, id string) error {
		removedContainerID = id
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, wb.ContainerId, removedContainerID)
	require.True(t, deleted)
}

func TestDeleteWorkbench_StopContainerGenericError_Propagates(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", VolumeName: "workbench-vol", DockerHostUuid: &defaultDockerHostID}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
		delete: func(context.Context, uuid.UUID) error {
			t.Fatal("Delete should not be called when StopContainer fails with a non-NotFound error")
			return nil
		},
	}

	docker := newFakeDocker()
	docker.stopContainer = func(context.Context, string) error {
		return errBoom
	}
	docker.removeContainer = func(context.Context, string) error {
		t.Fatal("RemoveContainer should not be called when StopContainer fails")
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.Error(t, err)
}

func TestDeleteWorkbench_RemoveVolumeNotFound_ContinuesToDelete(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, VolumeName: "workbench-vol", DockerHostUuid: &defaultDockerHostID}

	var deleted bool

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
		delete: func(context.Context, uuid.UUID) error {
			deleted = true
			return nil
		},
	}

	docker := newFakeDocker()
	docker.removeVolume = func(context.Context, string) error {
		return errdefs.ErrNotFound
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestDeleteWorkbench_RepoDeleteError_Propagates(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, VolumeName: "workbench-vol", DockerHostUuid: &defaultDockerHostID}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
		delete: func(context.Context, uuid.UUID) error {
			return errBoom
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.Error(t, err)
}

func TestGetLoginPrompt_UrlPresent(t *testing.T) {
	vaultID := uuid.New()

	pane := "Opening browser to sign in...\n" +
		"Browser didn't open? Use the url below to sign in (c to copy)\n" +
		"\n" +
		"https://platform.claude.com/oauth/authorize?code=true&state=abc123\n" +
		"\n" +
		"Paste code here if prompted >"

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	docker := newFakeDocker()
	docker.capturePane = func(_ context.Context, containerID string) (string, error) {
		require.Equal(t, "container-1", containerID)
		return pane, nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	got, err := svc.GetLoginPrompt(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, domain.WorkbenchLoginStateURLPresent, got.State)
	require.Equal(t, "https://platform.claude.com/oauth/authorize?code=true&state=abc123", got.URL)
}

func TestGetLoginPrompt_OAuthErrorPresent(t *testing.T) {
	vaultID := uuid.New()

	pane := "https://platform.claude.com/oauth/authorize?code=true\n" +
		"\n" +
		"OAuth error: Invalid code. Please make sure the full code was copied\n" +
		"Paste code here if prompted >"

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	docker := newFakeDocker()
	docker.capturePane = func(context.Context, string) (string, error) {
		return pane, nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	got, err := svc.GetLoginPrompt(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, domain.WorkbenchLoginStateError, got.State)
	require.Equal(t, "OAuth error: Invalid code. Please make sure the full code was copied", got.ErrorMessage)
}

func TestGetLoginPrompt_NeitherPresent_Authorized(t *testing.T) {
	vaultID := uuid.New()

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	docker := newFakeDocker()
	docker.capturePane = func(context.Context, string) (string, error) {
		return "│ > hello, what can you help me with today?\n", nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	got, err := svc.GetLoginPrompt(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, domain.WorkbenchLoginStateAuthorized, got.State)
}

func TestGetLoginPrompt_NeitherPresent_MenuStillShowing_Pending(t *testing.T) {
	vaultID := uuid.New()

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	docker := newFakeDocker()
	docker.capturePane = func(context.Context, string) (string, error) {
		pane := "Select login method:\n" +
			"> Claude account with subscription\n" +
			"  Anthropic Console account\n" +
			"  3rd-party platform"

		return pane, nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	got, err := svc.GetLoginPrompt(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, domain.WorkbenchLoginStatePending, got.State)
}

func TestGetLoginPrompt_GetByVaultIDError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	_, err := svc.GetLoginPrompt(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestGetLoginPrompt_CapturePaneError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	docker := newFakeDocker()
	docker.capturePane = func(context.Context, string) (string, error) {
		return "", errBoom
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	_, err := svc.GetLoginPrompt(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestSubmitLoginCode_HappyPath(t *testing.T) {
	vaultID := uuid.New()

	var gotContainerID, gotKeys string

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	docker := newFakeDocker()
	docker.sendKeys = func(_ context.Context, containerID string, keys string) error {
		gotContainerID = containerID
		gotKeys = keys
		return nil
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.SubmitLoginCode(context.Background(), vaultID, "the-pasted-code")
	require.NoError(t, err)
	require.Equal(t, "container-1", gotContainerID)
	require.Equal(t, "the-pasted-code", gotKeys)
}

func TestSubmitLoginCode_GetByVaultIDError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(newFakeDocker()))

	err := svc.SubmitLoginCode(context.Background(), uuid.New(), "code")
	require.Error(t, err)
}

func TestSubmitLoginCode_SendKeysError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{ContainerId: "container-1", DockerHostUuid: &defaultDockerHostID}, nil
		},
	}

	docker := newFakeDocker()
	docker.sendKeys = func(context.Context, string, string) error {
		return errBoom
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDockerHosts(), newFakeExternalConnections(), fakeDockerClientFactory(docker))

	err := svc.SubmitLoginCode(context.Background(), uuid.New(), "code")
	require.Error(t, err)
}
