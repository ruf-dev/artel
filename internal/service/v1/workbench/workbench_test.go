package workbench

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/clients/workbenchdocker"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// withUser returns a context.Background() carrying an injected user_context.UserContext for
// userUuid — mirrors internal/service/v1/vault/postgres_test.go's withUser, the idiom this
// package's Service methods expect (they resolve the caller via user_context.GetUserContext,
// same as vault.Service.requireVaultMember) rather than taking an explicit userID parameter.
func withUser(userUuid uuid.UUID) context.Context {
	uc := user_context.UserContext{UserUuid: userUuid, UserName: "tester"}

	return user_context.WithUserContext(context.Background(), uc)
}

// workbenchKey identifies a fakeWorkbenchesRepo row, mirroring the (vault_id, user_id) unique
// key migrations/072_workbench_per_user.sql put on the real table.
type workbenchKey struct {
	vaultID, userID uuid.UUID
}

// fakeWorkbenchesRepo is a hand-rolled, in-memory-map-backed fake of repository.Workbenches —
// real create/get/update/delete semantics keyed by (vaultID, userID), so tests can assert on
// actual row identity/count rather than just call arguments. Any *Func field, when set,
// overrides the default map-backed behavior for that method — used by tests that need to inject
// a specific error or a fixed row regardless of what's in the map.
type fakeWorkbenchesRepo struct {
	rows map[workbenchKey]domain.Workbench

	getByVaultAndUserFunc func(ctx context.Context, vaultID, userID uuid.UUID) (domain.Workbench, error)
	listByVaultIDFunc     func(ctx context.Context, vaultID uuid.UUID) ([]domain.Workbench, error)
}

func newFakeWorkbenchesRepo() *fakeWorkbenchesRepo {
	return &fakeWorkbenchesRepo{
		rows: make(map[workbenchKey]domain.Workbench),
	}
}

func (f *fakeWorkbenchesRepo) Create(
	ctx context.Context, vaultID, userID uuid.UUID, volumeName string, dockerHostID uuid.UUID,
) (domain.Workbench, error) {
	wb := domain.Workbench{
		Uuid:           uuid.New(),
		VaultUuid:      vaultID,
		UserUuid:       userID,
		Status:         domain.WorkbenchStatusConfiguring,
		VolumeName:     volumeName,
		DockerHostUuid: &dockerHostID,
	}

	f.rows[workbenchKey{vaultID, userID}] = wb

	return wb, nil
}

func (f *fakeWorkbenchesRepo) GetByVaultAndUser(ctx context.Context, vaultID, userID uuid.UUID) (domain.Workbench, error) {
	if f.getByVaultAndUserFunc != nil {
		return f.getByVaultAndUserFunc(ctx, vaultID, userID)
	}

	wb, ok := f.rows[workbenchKey{vaultID, userID}]
	if !ok {
		return domain.Workbench{}, sql.ErrNoRows
	}

	return wb, nil
}

func (f *fakeWorkbenchesRepo) ListByVaultID(ctx context.Context, vaultID uuid.UUID) ([]domain.Workbench, error) {
	if f.listByVaultIDFunc != nil {
		return f.listByVaultIDFunc(ctx, vaultID)
	}

	var out []domain.Workbench

	for key, wb := range f.rows {
		if key.vaultID == vaultID {
			out = append(out, wb)
		}
	}

	return out, nil
}

func (f *fakeWorkbenchesRepo) MarkContainerCreated(ctx context.Context, vaultID, userID uuid.UUID, containerID string) error {
	key := workbenchKey{vaultID, userID}

	wb, ok := f.rows[key]
	if !ok {
		return sql.ErrNoRows
	}

	wb.Status = domain.WorkbenchStatusCreated
	wb.ContainerId = containerID
	f.rows[key] = wb

	return nil
}

func (f *fakeWorkbenchesRepo) MarkConfiguring(ctx context.Context, vaultID, userID uuid.UUID) error {
	key := workbenchKey{vaultID, userID}

	wb, ok := f.rows[key]
	if !ok {
		return sql.ErrNoRows
	}

	wb.Status = domain.WorkbenchStatusConfiguring
	f.rows[key] = wb

	return nil
}

func (f *fakeWorkbenchesRepo) MarkRunning(ctx context.Context, vaultID, userID uuid.UUID, authMode domain.WorkbenchAuthMode) error {
	key := workbenchKey{vaultID, userID}

	wb, ok := f.rows[key]
	if !ok {
		return sql.ErrNoRows
	}

	wb.Status = domain.WorkbenchStatusRunning
	wb.AuthMode = authMode
	f.rows[key] = wb

	return nil
}

func (f *fakeWorkbenchesRepo) MarkStopped(ctx context.Context, vaultID, userID uuid.UUID) error {
	key := workbenchKey{vaultID, userID}

	wb, ok := f.rows[key]
	if !ok {
		return sql.ErrNoRows
	}

	wb.Status = domain.WorkbenchStatusStopped
	f.rows[key] = wb

	return nil
}

func (f *fakeWorkbenchesRepo) MarkRemoved(ctx context.Context, vaultID, userID uuid.UUID) error {
	key := workbenchKey{vaultID, userID}

	wb, ok := f.rows[key]
	if !ok {
		return sql.ErrNoRows
	}

	wb.Status = domain.WorkbenchStatusRemoved
	f.rows[key] = wb

	return nil
}

func (f *fakeWorkbenchesRepo) Delete(ctx context.Context, vaultID, userID uuid.UUID) error {
	delete(f.rows, workbenchKey{vaultID, userID})

	return nil
}

func (f *fakeWorkbenchesRepo) WithTx(tx postgres.DB) repository.Workbenches {
	panic("not implemented")
}

// fakeVaultMembersRepo is a hand-rolled fake of repository.VaultMembers — only Get is exercised
// by this package's tests (CreateWorkbench's membership check), the rest panic if called
// unexpectedly.
type fakeVaultMembersRepo struct {
	getFunc func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error)
}

func (f *fakeVaultMembersRepo) Add(ctx context.Context, vaultID, userID uuid.UUID, role artel_q.VaultRole) error {
	panic("not implemented")
}

func (f *fakeVaultMembersRepo) Remove(ctx context.Context, vaultID, userID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeVaultMembersRepo) Get(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
	return f.getFunc(ctx, vaultID, userID)
}

func (f *fakeVaultMembersRepo) ListByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMember, error) {
	panic("not implemented")
}

func (f *fakeVaultMembersRepo) ListByVaultWithUsers(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMemberInfo, error) {
	panic("not implemented")
}

func (f *fakeVaultMembersRepo) WithTx(tx postgres.DB) repository.VaultMembers {
	panic("not implemented")
}

// alwaysVaultMember is a fakeVaultMembersRepo that accepts every membership lookup.
func alwaysVaultMember() *fakeVaultMembersRepo {
	return &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{VaultUuid: vaultID, UserUuid: userID}, nil
		},
	}
}

// fakeDockerHostsRepo is a hand-rolled fake of repository.DockerHosts — only the methods
// exercised by a given test set a func field, the rest panic if called unexpectedly.
type fakeDockerHostsRepo struct {
	getWithCredsFunc    func(ctx context.Context, id uuid.UUID) (domain.DockerHost, error)
	pickLeastLoadedFunc func(ctx context.Context) (domain.DockerHost, error)
}

func (f *fakeDockerHostsRepo) Register(ctx context.Context, url, caCert, clientCert, clientKey string) (uuid.UUID, error) {
	panic("not implemented")
}

func (f *fakeDockerHostsRepo) Get(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
	panic("not implemented")
}

func (f *fakeDockerHostsRepo) GetWithCreds(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
	return f.getWithCredsFunc(ctx, id)
}

func (f *fakeDockerHostsRepo) List(ctx context.Context) ([]domain.DockerHost, error) {
	panic("not implemented")
}

func (f *fakeDockerHostsRepo) Update(ctx context.Context, id uuid.UUID, url string, caCert, clientCert, clientKey *string) error {
	panic("not implemented")
}

func (f *fakeDockerHostsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeDockerHostsRepo) Exists(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (f *fakeDockerHostsRepo) PickLeastLoaded(ctx context.Context) (domain.DockerHost, error) {
	return f.pickLeastLoadedFunc(ctx)
}

func (f *fakeDockerHostsRepo) WithTx(tx postgres.DB) repository.DockerHosts {
	panic("not implemented")
}

// fakeDockerClient is a hand-rolled fake of the package-private dockerClient interface — only
// the methods exercised by a given test set a func field, the rest panic if called unexpectedly.
type fakeDockerClient struct {
	containerAddressFunc func(ctx context.Context, containerID string) (string, error)
	createVolumeFunc     func(ctx context.Context, name string) error
	createContainerFunc  func(ctx context.Context, opts workbenchdocker.CreateOpts) (string, error)
	startContainerFunc   func(ctx context.Context, containerID string, env map[string]string) error
	stopContainerFunc    func(ctx context.Context, containerID string) error
	removeContainerFunc  func(ctx context.Context, containerID string) error
	removeVolumeFunc     func(ctx context.Context, name string) error

	listTmuxWindowsFunc  func(ctx context.Context, containerID string) ([]domain.TerminalTab, error)
	newTmuxWindowFunc    func(ctx context.Context, containerID string) (domain.TerminalTab, error)
	selectTmuxWindowFunc func(ctx context.Context, containerID, windowID string) error
	killTmuxWindowFunc   func(ctx context.Context, containerID, windowID string) error

	containerAddressCalls []string

	// listTmuxWindowsCalls/newTmuxWindowCalls/selectTmuxWindowCalls/killTmuxWindowCalls count
	// invocations of the corresponding method — used by tests to assert a docker call was (or was
	// deliberately not) made, e.g. WorkbenchNotRunning short-circuiting before ever reaching the
	// docker client, or CloseTerminalTab refusing to kill a workbench's last remaining window.
	listTmuxWindowsCalls  int
	newTmuxWindowCalls    int
	selectTmuxWindowCalls int
	killTmuxWindowCalls   int
}

func (f *fakeDockerClient) CreateVolume(ctx context.Context, name string) error {
	if f.createVolumeFunc != nil {
		return f.createVolumeFunc(ctx, name)
	}

	return nil
}

func (f *fakeDockerClient) RemoveVolume(ctx context.Context, name string) error {
	if f.removeVolumeFunc != nil {
		return f.removeVolumeFunc(ctx, name)
	}

	return nil
}

func (f *fakeDockerClient) CreateContainer(ctx context.Context, opts workbenchdocker.CreateOpts) (string, error) {
	if f.createContainerFunc != nil {
		return f.createContainerFunc(ctx, opts)
	}

	return "container-" + opts.Name, nil
}

func (f *fakeDockerClient) StartContainer(ctx context.Context, containerID string, env map[string]string) error {
	if f.startContainerFunc != nil {
		return f.startContainerFunc(ctx, containerID, env)
	}

	return nil
}

func (f *fakeDockerClient) StopContainer(ctx context.Context, containerID string) error {
	if f.stopContainerFunc != nil {
		return f.stopContainerFunc(ctx, containerID)
	}

	return nil
}

func (f *fakeDockerClient) RemoveContainer(ctx context.Context, containerID string) error {
	if f.removeContainerFunc != nil {
		return f.removeContainerFunc(ctx, containerID)
	}

	return nil
}

func (f *fakeDockerClient) ContainerAddress(ctx context.Context, containerID string) (string, error) {
	f.containerAddressCalls = append(f.containerAddressCalls, containerID)

	return f.containerAddressFunc(ctx, containerID)
}

func (f *fakeDockerClient) ListTmuxWindows(ctx context.Context, containerID string) ([]domain.TerminalTab, error) {
	f.listTmuxWindowsCalls++

	if f.listTmuxWindowsFunc != nil {
		return f.listTmuxWindowsFunc(ctx, containerID)
	}

	return nil, nil
}

func (f *fakeDockerClient) NewTmuxWindow(ctx context.Context, containerID string) (domain.TerminalTab, error) {
	f.newTmuxWindowCalls++

	if f.newTmuxWindowFunc != nil {
		return f.newTmuxWindowFunc(ctx, containerID)
	}

	return domain.TerminalTab{}, nil
}

func (f *fakeDockerClient) SelectTmuxWindow(ctx context.Context, containerID, windowID string) error {
	f.selectTmuxWindowCalls++

	if f.selectTmuxWindowFunc != nil {
		return f.selectTmuxWindowFunc(ctx, containerID, windowID)
	}

	return nil
}

func (f *fakeDockerClient) KillTmuxWindow(ctx context.Context, containerID, windowID string) error {
	f.killTmuxWindowCalls++

	if f.killTmuxWindowFunc != nil {
		return f.killTmuxWindowFunc(ctx, containerID, windowID)
	}

	return nil
}

// testContainerId is the container id newTestService's fake workbenches repo reports, asserted
// on to confirm the right container was addressed.
const testContainerId = "container-1"

// newTestService builds a Service wired to fakes: workbenchesRepo returns a workbench row in
// status with a non-nil DockerHostUuid (required for resolveClient to proceed past the
// WorkbenchMissingDockerHost early return), dockerHostsRepo returns a bare docker host, and
// newDockerClient always returns client.
func newTestService(t *testing.T, status domain.WorkbenchStatus, client *fakeDockerClient) *Service {
	t.Helper()

	hostUuid := uuid.New()

	workbenchesRepo := &fakeWorkbenchesRepo{
		getByVaultAndUserFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (domain.Workbench, error) {
			wb := domain.Workbench{
				Uuid:           uuid.New(),
				VaultUuid:      vaultID,
				UserUuid:       userID,
				Status:         status,
				ContainerId:    testContainerId,
				DockerHostUuid: &hostUuid,
			}

			return wb, nil
		},
	}

	dockerHostsRepo := &fakeDockerHostsRepo{
		getWithCredsFunc: func(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
			return domain.DockerHost{Uuid: id}, nil
		},
	}

	newDockerClient := func(host string, tlsCfg workbenchdocker.TLSConfig) (dockerClient, error) {
		return client, nil
	}

	svc := &Service{
		workbenchesRepo: workbenchesRepo,
		dockerHostsRepo: dockerHostsRepo,
		newDockerClient: newDockerClient,
	}

	return svc
}

func TestResolveTerminalTarget_Running(t *testing.T) {
	client := &fakeDockerClient{
		containerAddressFunc: func(ctx context.Context, containerID string) (string, error) {
			return "172.18.0.7:7681", nil
		},
	}

	svc := newTestService(t, domain.WorkbenchStatusRunning, client)

	target, err := svc.ResolveTerminalTarget(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target != "http://172.18.0.7:7681" {
		t.Fatalf("target = %q, want %q", target, "http://172.18.0.7:7681")
	}

	if len(client.containerAddressCalls) != 1 || client.containerAddressCalls[0] != testContainerId {
		t.Fatalf("ContainerAddress calls = %v, want exactly [%s]", client.containerAddressCalls, testContainerId)
	}
}

func TestResolveTerminalTarget_NotRunning(t *testing.T) {
	statuses := []domain.WorkbenchStatus{
		domain.WorkbenchStatusCreated,
		domain.WorkbenchStatusConfiguring,
		domain.WorkbenchStatusStopped,
		domain.WorkbenchStatusRemoved,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			client := &fakeDockerClient{}

			svc := newTestService(t, status, client)

			_, err := svc.ResolveTerminalTarget(context.Background(), uuid.New(), uuid.New())
			if !errors.Is(err, user_errors.WorkbenchNotRunning) {
				t.Fatalf("expected user_errors.WorkbenchNotRunning, got %v", err)
			}

			if len(client.containerAddressCalls) != 0 {
				t.Fatalf("ContainerAddress called %d times for a non-running workbench, want 0", len(client.containerAddressCalls))
			}
		})
	}
}

func TestResolveTerminalTarget_DockerClientError(t *testing.T) {
	wantErr := errors.New("boom")

	client := &fakeDockerClient{
		containerAddressFunc: func(ctx context.Context, containerID string) (string, error) {
			return "", wantErr
		},
	}

	svc := newTestService(t, domain.WorkbenchStatusRunning, client)

	_, err := svc.ResolveTerminalTarget(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped ContainerAddress error, got %v", err)
	}
}

// newCreateWorkbenchService builds a Service wired to fakes sufficient to drive CreateWorkbench
// end to end: an always-accepting vault membership check, a single-host docker pool, and a
// fakeDockerClient whose CreateVolume/CreateContainer/StartContainer all default to success.
func newCreateWorkbenchService(t *testing.T, workbenchesRepo *fakeWorkbenchesRepo) *Service {
	t.Helper()

	hostUuid := uuid.New()
	client := &fakeDockerClient{}

	dockerHostsRepo := &fakeDockerHostsRepo{
		pickLeastLoadedFunc: func(ctx context.Context) (domain.DockerHost, error) {
			return domain.DockerHost{Uuid: hostUuid}, nil
		},
		getWithCredsFunc: func(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
			return domain.DockerHost{Uuid: id}, nil
		},
	}

	newDockerClient := func(host string, tlsCfg workbenchdocker.TLSConfig) (dockerClient, error) {
		return client, nil
	}

	svc := &Service{
		workbenchesRepo: workbenchesRepo,
		vaultMembers:    alwaysVaultMember(),
		dockerHostsRepo: dockerHostsRepo,
		newDockerClient: newDockerClient,
	}

	return svc
}

// TestCreateWorkbench_TwoUsersGetDistinctWorkbenches confirms two different members of the same
// vault each calling CreateWorkbench get their own, separate workbench row — the core new
// per-(vault, user) behavior migrations/072_workbench_per_user.sql introduces, replacing the
// prior one-workbench-per-vault model.
func TestCreateWorkbench_TwoUsersGetDistinctWorkbenches(t *testing.T) {
	workbenchesRepo := newFakeWorkbenchesRepo()
	svc := newCreateWorkbenchService(t, workbenchesRepo)

	vaultID := uuid.New()
	user1 := uuid.New()
	user2 := uuid.New()

	wb1, err := svc.CreateWorkbench(withUser(user1), vaultID)
	if err != nil {
		t.Fatalf("unexpected error creating workbench for user1: %v", err)
	}

	wb2, err := svc.CreateWorkbench(withUser(user2), vaultID)
	if err != nil {
		t.Fatalf("unexpected error creating workbench for user2: %v", err)
	}

	if wb1.Uuid == wb2.Uuid {
		t.Fatalf("expected distinct workbench rows, got the same uuid %s for both users", wb1.Uuid)
	}

	if wb1.UserUuid != user1 {
		t.Fatalf("wb1.UserUuid = %s, want %s", wb1.UserUuid, user1)
	}

	if wb2.UserUuid != user2 {
		t.Fatalf("wb2.UserUuid = %s, want %s", wb2.UserUuid, user2)
	}

	if wb1.VolumeName == wb2.VolumeName {
		t.Fatalf("expected distinct volume names to avoid Docker resource collisions, got %q for both", wb1.VolumeName)
	}

	wantVolume1 := fmt.Sprintf("workbench-%s-%s", vaultID, user1)
	wantVolume2 := fmt.Sprintf("workbench-%s-%s", vaultID, user2)

	if wb1.VolumeName != wantVolume1 {
		t.Fatalf("wb1.VolumeName = %q, want %q", wb1.VolumeName, wantVolume1)
	}

	if wb2.VolumeName != wantVolume2 {
		t.Fatalf("wb2.VolumeName = %q, want %q", wb2.VolumeName, wantVolume2)
	}

	if len(workbenchesRepo.rows) != 2 {
		t.Fatalf("expected 2 rows stored, got %d", len(workbenchesRepo.rows))
	}
}

// TestCreateWorkbench_IdempotentSameUser confirms a retry by the same user for the same vault
// returns the existing row rather than creating a second one.
func TestCreateWorkbench_IdempotentSameUser(t *testing.T) {
	workbenchesRepo := newFakeWorkbenchesRepo()
	svc := newCreateWorkbenchService(t, workbenchesRepo)

	vaultID := uuid.New()
	user1 := uuid.New()

	first, err := svc.CreateWorkbench(withUser(user1), vaultID)
	if err != nil {
		t.Fatalf("unexpected error on first create: %v", err)
	}

	second, err := svc.CreateWorkbench(withUser(user1), vaultID)
	if err != nil {
		t.Fatalf("unexpected error on second create: %v", err)
	}

	if first.Uuid != second.Uuid {
		t.Fatalf("expected the same workbench row on retry, got %s then %s", first.Uuid, second.Uuid)
	}

	if len(workbenchesRepo.rows) != 1 {
		t.Fatalf("expected exactly 1 row stored, got %d", len(workbenchesRepo.rows))
	}
}

// TestCreateWorkbench_RequiresVaultMembership confirms a caller who isn't a member of vaultID is
// rejected with user_errors.WorkbenchRequiresVaultMembership before any workbench row is
// created.
func TestCreateWorkbench_RequiresVaultMembership(t *testing.T) {
	workbenchesRepo := newFakeWorkbenchesRepo()
	svc := newCreateWorkbenchService(t, workbenchesRepo)

	svc.vaultMembers = &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, user_errors.NotFound
		},
	}

	_, err := svc.CreateWorkbench(withUser(uuid.New()), uuid.New())
	if !errors.Is(err, user_errors.WorkbenchRequiresVaultMembership) {
		t.Fatalf("expected user_errors.WorkbenchRequiresVaultMembership, got %v", err)
	}

	if len(workbenchesRepo.rows) != 0 {
		t.Fatalf("expected no workbench row created, got %d", len(workbenchesRepo.rows))
	}
}

// TestCreateWorkbench_Unauthenticated confirms a ctx with no injected user_context.UserContext
// is rejected rather than reaching the membership check or Docker calls.
func TestCreateWorkbench_Unauthenticated(t *testing.T) {
	workbenchesRepo := newFakeWorkbenchesRepo()
	svc := newCreateWorkbenchService(t, workbenchesRepo)

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	if !errors.Is(err, user_errors.Unauthenticated) {
		t.Fatalf("expected user_errors.Unauthenticated, got %v", err)
	}
}

// newDeletableWorkbenchRow seeds workbenchesRepo with a row for (vaultID, userID) in a state
// deleteWorkbenchRow can tear down (has a container id and a resolvable docker host).
func newDeletableWorkbenchRow(workbenchesRepo *fakeWorkbenchesRepo, vaultID, userID uuid.UUID, hostUuid uuid.UUID) {
	wb := domain.Workbench{
		Uuid:           uuid.New(),
		VaultUuid:      vaultID,
		UserUuid:       userID,
		Status:         domain.WorkbenchStatusStopped,
		ContainerId:    "container-" + userID.String(),
		VolumeName:     "workbench-" + vaultID.String() + "-" + userID.String(),
		DockerHostUuid: &hostUuid,
	}

	workbenchesRepo.rows[workbenchKey{vaultID, userID}] = wb
}

// TestDeleteWorkbenchesForVault_MultipleRows confirms every member's workbench row for vaultID
// is torn down (container/volume removed via the fake docker client, DB row deleted) — the
// behavior internal/transport/vaults_api/delete.go's DeleteVault now relies on instead of
// tearing down a single workbench.
func TestDeleteWorkbenchesForVault_MultipleRows(t *testing.T) {
	workbenchesRepo := newFakeWorkbenchesRepo()

	hostUuid := uuid.New()
	vaultID := uuid.New()
	user1 := uuid.New()
	user2 := uuid.New()

	newDeletableWorkbenchRow(workbenchesRepo, vaultID, user1, hostUuid)
	newDeletableWorkbenchRow(workbenchesRepo, vaultID, user2, hostUuid)

	var removedVolumes []string

	client := &fakeDockerClient{
		removeVolumeFunc: func(ctx context.Context, name string) error {
			removedVolumes = append(removedVolumes, name)

			return nil
		},
	}

	dockerHostsRepo := &fakeDockerHostsRepo{
		getWithCredsFunc: func(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
			return domain.DockerHost{Uuid: id}, nil
		},
	}

	newDockerClient := func(host string, tlsCfg workbenchdocker.TLSConfig) (dockerClient, error) {
		return client, nil
	}

	svc := &Service{
		workbenchesRepo: workbenchesRepo,
		dockerHostsRepo: dockerHostsRepo,
		newDockerClient: newDockerClient,
	}

	err := svc.DeleteWorkbenchesForVault(context.Background(), vaultID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(workbenchesRepo.rows) != 0 {
		t.Fatalf("expected all rows deleted, got %d remaining", len(workbenchesRepo.rows))
	}

	if len(removedVolumes) != 2 {
		t.Fatalf("expected 2 volumes removed, got %d (%v)", len(removedVolumes), removedVolumes)
	}
}

// TestDeleteWorkbenchesForVault_PartialFailureContinues confirms a Docker failure tearing down
// one member's row doesn't stop the rest from being torn down — the whole point of joining
// errors rather than failing fast, per DeleteWorkbenchesForVault's doc comment.
func TestDeleteWorkbenchesForVault_PartialFailureContinues(t *testing.T) {
	workbenchesRepo := newFakeWorkbenchesRepo()

	hostUuid := uuid.New()
	vaultID := uuid.New()
	failingUser := uuid.New()
	okUser := uuid.New()

	newDeletableWorkbenchRow(workbenchesRepo, vaultID, failingUser, hostUuid)
	newDeletableWorkbenchRow(workbenchesRepo, vaultID, okUser, hostUuid)

	boom := errors.New("boom")

	client := &fakeDockerClient{
		removeVolumeFunc: func(ctx context.Context, name string) error {
			if name == "workbench-"+vaultID.String()+"-"+failingUser.String() {
				return boom
			}

			return nil
		},
	}

	dockerHostsRepo := &fakeDockerHostsRepo{
		getWithCredsFunc: func(ctx context.Context, id uuid.UUID) (domain.DockerHost, error) {
			return domain.DockerHost{Uuid: id}, nil
		},
	}

	newDockerClient := func(host string, tlsCfg workbenchdocker.TLSConfig) (dockerClient, error) {
		return client, nil
	}

	svc := &Service{
		workbenchesRepo: workbenchesRepo,
		dockerHostsRepo: dockerHostsRepo,
		newDockerClient: newDockerClient,
	}

	err := svc.DeleteWorkbenchesForVault(context.Background(), vaultID)
	if err == nil {
		t.Fatal("expected a combined error from the failing row's teardown")
	}

	if !errors.Is(err, boom) {
		t.Fatalf("expected the returned error to wrap the docker failure, got %v", err)
	}

	if _, stillThere := workbenchesRepo.rows[workbenchKey{vaultID, failingUser}]; !stillThere {
		t.Fatal("failing row should not have been deleted since its volume removal errored")
	}

	if _, stillThere := workbenchesRepo.rows[workbenchKey{vaultID, okUser}]; stillThere {
		t.Fatal("the other row should still have been torn down despite the first row's failure")
	}
}

// notRunningStatuses is every domain.WorkbenchStatus other than 'running' — the four terminal-tab
// methods must reject all of them with user_errors.WorkbenchNotRunning before ever reaching the
// docker client.
var notRunningStatuses = []domain.WorkbenchStatus{
	domain.WorkbenchStatusCreated,
	domain.WorkbenchStatusConfiguring,
	domain.WorkbenchStatusStopped,
	domain.WorkbenchStatusRemoved,
}

// TestListTerminalTabs_NotRunning confirms ListTerminalTabs rejects every non-running workbench
// status with user_errors.WorkbenchNotRunning and never reaches the docker client.
func TestListTerminalTabs_NotRunning(t *testing.T) {
	for _, status := range notRunningStatuses {
		t.Run(string(status), func(t *testing.T) {
			client := &fakeDockerClient{}

			svc := newTestService(t, status, client)

			_, err := svc.ListTerminalTabs(withUser(uuid.New()), uuid.New())
			if !errors.Is(err, user_errors.WorkbenchNotRunning) {
				t.Fatalf("expected user_errors.WorkbenchNotRunning, got %v", err)
			}

			if client.listTmuxWindowsCalls != 0 {
				t.Fatalf("ListTmuxWindows called %d times for a non-running workbench, want 0", client.listTmuxWindowsCalls)
			}
		})
	}
}

// TestListTerminalTabs_Running confirms the happy path passes the docker client's tabs straight
// through unwrapped.
func TestListTerminalTabs_Running(t *testing.T) {
	wantTabs := []domain.TerminalTab{
		{ID: "@1", Name: "claude", Active: true},
		{ID: "@2", Name: "bash", Active: false},
	}

	client := &fakeDockerClient{
		listTmuxWindowsFunc: func(ctx context.Context, containerID string) ([]domain.TerminalTab, error) {
			if containerID != testContainerId {
				t.Fatalf("containerID = %q, want %q", containerID, testContainerId)
			}

			return wantTabs, nil
		},
	}

	svc := newTestService(t, domain.WorkbenchStatusRunning, client)

	tabs, err := svc.ListTerminalTabs(withUser(uuid.New()), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(tabs, wantTabs) {
		t.Fatalf("tabs = %#v, want %#v", tabs, wantTabs)
	}

	if client.listTmuxWindowsCalls != 1 {
		t.Fatalf("ListTmuxWindows called %d times, want 1", client.listTmuxWindowsCalls)
	}
}

// TestCreateTerminalTab_NotRunning confirms CreateTerminalTab rejects every non-running workbench
// status with user_errors.WorkbenchNotRunning and never reaches the docker client.
func TestCreateTerminalTab_NotRunning(t *testing.T) {
	for _, status := range notRunningStatuses {
		t.Run(string(status), func(t *testing.T) {
			client := &fakeDockerClient{}

			svc := newTestService(t, status, client)

			_, err := svc.CreateTerminalTab(withUser(uuid.New()), uuid.New())
			if !errors.Is(err, user_errors.WorkbenchNotRunning) {
				t.Fatalf("expected user_errors.WorkbenchNotRunning, got %v", err)
			}

			if client.newTmuxWindowCalls != 0 {
				t.Fatalf("NewTmuxWindow called %d times for a non-running workbench, want 0", client.newTmuxWindowCalls)
			}
		})
	}
}

// TestCreateTerminalTab_Running confirms the happy path passes the docker client's new tab
// straight through unwrapped.
func TestCreateTerminalTab_Running(t *testing.T) {
	wantTab := domain.TerminalTab{ID: "@3", Active: true}

	client := &fakeDockerClient{
		newTmuxWindowFunc: func(ctx context.Context, containerID string) (domain.TerminalTab, error) {
			if containerID != testContainerId {
				t.Fatalf("containerID = %q, want %q", containerID, testContainerId)
			}

			return wantTab, nil
		},
	}

	svc := newTestService(t, domain.WorkbenchStatusRunning, client)

	tab, err := svc.CreateTerminalTab(withUser(uuid.New()), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tab != wantTab {
		t.Fatalf("tab = %#v, want %#v", tab, wantTab)
	}

	if client.newTmuxWindowCalls != 1 {
		t.Fatalf("NewTmuxWindow called %d times, want 1", client.newTmuxWindowCalls)
	}
}

// TestSelectTerminalTab_NotRunning confirms SelectTerminalTab rejects every non-running workbench
// status with user_errors.WorkbenchNotRunning and never reaches the docker client.
func TestSelectTerminalTab_NotRunning(t *testing.T) {
	for _, status := range notRunningStatuses {
		t.Run(string(status), func(t *testing.T) {
			client := &fakeDockerClient{}

			svc := newTestService(t, status, client)

			err := svc.SelectTerminalTab(withUser(uuid.New()), uuid.New(), "@1")
			if !errors.Is(err, user_errors.WorkbenchNotRunning) {
				t.Fatalf("expected user_errors.WorkbenchNotRunning, got %v", err)
			}

			if client.selectTmuxWindowCalls != 0 {
				t.Fatalf("SelectTmuxWindow called %d times for a non-running workbench, want 0", client.selectTmuxWindowCalls)
			}
		})
	}
}

// TestSelectTerminalTab_Running confirms the happy path delegates to the docker client with the
// caller-supplied tab id, unwrapped.
func TestSelectTerminalTab_Running(t *testing.T) {
	const wantTabID = "@2"

	var gotContainerID, gotWindowID string

	client := &fakeDockerClient{
		selectTmuxWindowFunc: func(ctx context.Context, containerID, windowID string) error {
			gotContainerID = containerID
			gotWindowID = windowID

			return nil
		},
	}

	svc := newTestService(t, domain.WorkbenchStatusRunning, client)

	err := svc.SelectTerminalTab(withUser(uuid.New()), uuid.New(), wantTabID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotContainerID != testContainerId {
		t.Fatalf("containerID = %q, want %q", gotContainerID, testContainerId)
	}

	if gotWindowID != wantTabID {
		t.Fatalf("windowID = %q, want %q", gotWindowID, wantTabID)
	}

	if client.selectTmuxWindowCalls != 1 {
		t.Fatalf("SelectTmuxWindow called %d times, want 1", client.selectTmuxWindowCalls)
	}
}

// TestCloseTerminalTab_NotRunning confirms CloseTerminalTab rejects every non-running workbench
// status with user_errors.WorkbenchNotRunning and never reaches the docker client at all — not
// even ListTmuxWindows.
func TestCloseTerminalTab_NotRunning(t *testing.T) {
	for _, status := range notRunningStatuses {
		t.Run(string(status), func(t *testing.T) {
			client := &fakeDockerClient{}

			svc := newTestService(t, status, client)

			err := svc.CloseTerminalTab(withUser(uuid.New()), uuid.New(), "@1")
			if !errors.Is(err, user_errors.WorkbenchNotRunning) {
				t.Fatalf("expected user_errors.WorkbenchNotRunning, got %v", err)
			}

			if client.listTmuxWindowsCalls != 0 {
				t.Fatalf("ListTmuxWindows called %d times for a non-running workbench, want 0", client.listTmuxWindowsCalls)
			}

			if client.killTmuxWindowCalls != 0 {
				t.Fatalf("KillTmuxWindow called %d times for a non-running workbench, want 0", client.killTmuxWindowCalls)
			}
		})
	}
}

// TestCloseTerminalTab_LastWindowRefused confirms CloseTerminalTab refuses to close a workbench's
// only remaining tmux window with user_errors.WorkbenchCannotCloseLastTab, and never calls
// KillTmuxWindow.
func TestCloseTerminalTab_LastWindowRefused(t *testing.T) {
	client := &fakeDockerClient{
		listTmuxWindowsFunc: func(ctx context.Context, containerID string) ([]domain.TerminalTab, error) {
			return []domain.TerminalTab{{ID: "@1", Name: "claude", Active: true}}, nil
		},
	}

	svc := newTestService(t, domain.WorkbenchStatusRunning, client)

	err := svc.CloseTerminalTab(withUser(uuid.New()), uuid.New(), "@1")
	if !errors.Is(err, user_errors.WorkbenchCannotCloseLastTab) {
		t.Fatalf("expected user_errors.WorkbenchCannotCloseLastTab, got %v", err)
	}

	if client.killTmuxWindowCalls != 0 {
		t.Fatalf("KillTmuxWindow called %d times when closing the last tab, want 0", client.killTmuxWindowCalls)
	}
}

// TestCloseTerminalTab_MultipleWindowsCloses confirms CloseTerminalTab proceeds to KillTmuxWindow
// (happy path, no error) once two or more tmux windows remain.
func TestCloseTerminalTab_MultipleWindowsCloses(t *testing.T) {
	const wantTabID = "@2"

	var gotContainerID, gotWindowID string

	client := &fakeDockerClient{
		listTmuxWindowsFunc: func(ctx context.Context, containerID string) ([]domain.TerminalTab, error) {
			tabs := []domain.TerminalTab{
				{ID: "@1", Name: "claude", Active: true},
				{ID: "@2", Name: "bash", Active: false},
			}

			return tabs, nil
		},
		killTmuxWindowFunc: func(ctx context.Context, containerID, windowID string) error {
			gotContainerID = containerID
			gotWindowID = windowID

			return nil
		},
	}

	svc := newTestService(t, domain.WorkbenchStatusRunning, client)

	err := svc.CloseTerminalTab(withUser(uuid.New()), uuid.New(), wantTabID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotContainerID != testContainerId {
		t.Fatalf("containerID = %q, want %q", gotContainerID, testContainerId)
	}

	if gotWindowID != wantTabID {
		t.Fatalf("windowID = %q, want %q", gotWindowID, wantTabID)
	}

	if client.killTmuxWindowCalls != 1 {
		t.Fatalf("KillTmuxWindow called %d times, want 1", client.killTmuxWindowCalls)
	}
}
