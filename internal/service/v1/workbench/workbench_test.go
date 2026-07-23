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
	"go.redsock.ru/rerrors"
)

// fakeWorkbenches is a hand-written repository.Workbenches for exercising Service's branching
// logic without a live Postgres.
type fakeWorkbenches struct {
	getByVaultID func(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error)
	create       func(ctx context.Context, vaultID, userID uuid.UUID, volumeName, containerID string) (domain.Workbench, error)
	delete       func(ctx context.Context, vaultID uuid.UUID) error
}

func (f *fakeWorkbenches) Create(
	ctx context.Context,
	vaultID, userID uuid.UUID,
	volumeName, containerID string,
) (domain.Workbench, error) {
	return f.create(ctx, vaultID, userID, volumeName, containerID)
}

func (f *fakeWorkbenches) GetByVaultID(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error) {
	return f.getByVaultID(ctx, vaultID)
}

func (f *fakeWorkbenches) MarkRunning(context.Context, uuid.UUID, domain.WorkbenchAuthMode) error {
	return nil
}

func (f *fakeWorkbenches) MarkStopped(context.Context, uuid.UUID) error {
	return nil
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

// fakeDocker is a hand-written dockerClient. Every method defaults to a no-op success; tests
// override only the fields they need to exercise.
type fakeDocker struct {
	createVolume    func(ctx context.Context, name string) error
	removeVolume    func(ctx context.Context, name string) error
	createContainer func(ctx context.Context, opts workbenchdocker.CreateOpts) (string, error)
	stopContainer   func(ctx context.Context, containerID string) error
	removeContainer func(ctx context.Context, containerID string) error
}

func newFakeDocker() *fakeDocker {
	return &fakeDocker{
		createVolume: func(context.Context, string) error { return nil },
		removeVolume: func(context.Context, string) error { return nil },
		createContainer: func(context.Context, workbenchdocker.CreateOpts) (string, error) {
			return "container-1", nil
		},
		stopContainer:   func(context.Context, string) error { return nil },
		removeContainer: func(context.Context, string) error { return nil },
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

func (f *fakeDocker) StopContainer(ctx context.Context, containerID string) error {
	return f.stopContainer(ctx, containerID)
}

func (f *fakeDocker) RemoveContainer(ctx context.Context, containerID string) error {
	return f.removeContainer(ctx, containerID)
}

var errBoom = rerrors.New("boom")

func TestCreateWorkbench_Idempotent_ReturnsExisting(t *testing.T) {
	vaultID := uuid.New()
	existing := domain.Workbench{Uuid: uuid.New(), VaultUuid: vaultID, VolumeName: "workbench-existing"}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return existing, nil
		},
		create: func(context.Context, uuid.UUID, uuid.UUID, string, string) (domain.Workbench, error) {
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

	svc := New(workbenches, vaults, docker)

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
	var gotVaultID, gotUserID uuid.UUID
	var gotVolumeName, gotContainerID string

	created := domain.Workbench{Uuid: uuid.New(), VaultUuid: vaultID, UserUuid: userID}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "error getting workbench by vault id")
		},
		create: func(_ context.Context, vaultID, userID uuid.UUID, volumeName, containerID string) (domain.Workbench, error) {
			gotVaultID = vaultID
			gotUserID = userID
			gotVolumeName = volumeName
			gotContainerID = containerID

			return created, nil
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

	svc := New(workbenches, vaults, docker)

	got, err := svc.CreateWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, created, got)

	expectedVolumeName := "workbench-" + vaultID.String()
	require.Equal(t, expectedVolumeName, gotCreateVolumeName)
	require.Equal(t, expectedVolumeName, gotCreateOpts.Name)
	require.Equal(t, expectedVolumeName, gotCreateOpts.VolumeName)
	require.Equal(t, vaultID, gotVaultID)
	require.Equal(t, userID, gotUserID)
	require.Equal(t, expectedVolumeName, gotVolumeName)
	require.Equal(t, "container-1", gotContainerID)
}

func TestCreateWorkbench_GetByVaultIDUnexpectedError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}
	vaults := &fakeVaults{}
	docker := newFakeDocker()

	svc := New(workbenches, vaults, docker)

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

	svc := New(workbenches, vaults, docker)

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_CreateVolumeError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
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

	svc := New(workbenches, vaults, docker)

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_CreateContainerError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
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

	svc := New(workbenches, vaults, docker)

	_, err := svc.CreateWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
}

func TestCreateWorkbench_RepoCreateError_Propagates(t *testing.T) {
	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return domain.Workbench{}, rerrors.Wrap(sql.ErrNoRows, "not found")
		},
		create: func(context.Context, uuid.UUID, uuid.UUID, string, string) (domain.Workbench, error) {
			return domain.Workbench{}, errBoom
		},
	}
	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, nil
		},
	}
	docker := newFakeDocker()

	svc := New(workbenches, vaults, docker)

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

	svc := New(workbenches, &fakeVaults{}, newFakeDocker())

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

	svc := New(workbenches, &fakeVaults{}, newFakeDocker())

	_, err := svc.GetWorkbench(context.Background(), uuid.New())
	require.Error(t, err)
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

	svc := New(workbenches, &fakeVaults{}, docker)

	err := svc.DeleteWorkbench(context.Background(), uuid.New())
	require.NoError(t, err)
}

func TestDeleteWorkbench_HappyPath_StopsRemovesAndDeletes(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", VolumeName: "workbench-vol"}

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

	svc := New(workbenches, &fakeVaults{}, docker)

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, wb.ContainerId, stoppedID)
	require.Equal(t, wb.ContainerId, removedContainerID)
	require.Equal(t, wb.VolumeName, removedVolumeName)
	require.Equal(t, vaultID, deletedVaultID)
}

func TestDeleteWorkbench_NoContainerYet_SkipsContainerCalls(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "", VolumeName: "workbench-vol"}

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

	svc := New(workbenches, &fakeVaults{}, docker)

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
}

func TestDeleteWorkbench_StopContainerNotFound_ContinuesTeardown(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", VolumeName: "workbench-vol"}

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

	svc := New(workbenches, &fakeVaults{}, docker)

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.Equal(t, wb.ContainerId, removedContainerID)
	require.True(t, deleted)
}

func TestDeleteWorkbench_StopContainerGenericError_Propagates(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, ContainerId: "container-1", VolumeName: "workbench-vol"}

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

	svc := New(workbenches, &fakeVaults{}, docker)

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.Error(t, err)
}

func TestDeleteWorkbench_RemoveVolumeNotFound_ContinuesToDelete(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, VolumeName: "workbench-vol"}

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

	svc := New(workbenches, &fakeVaults{}, docker)

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestDeleteWorkbench_RepoDeleteError_Propagates(t *testing.T) {
	vaultID := uuid.New()
	wb := domain.Workbench{VaultUuid: vaultID, VolumeName: "workbench-vol"}

	workbenches := &fakeWorkbenches{
		getByVaultID: func(context.Context, uuid.UUID) (domain.Workbench, error) {
			return wb, nil
		},
		delete: func(context.Context, uuid.UUID) error {
			return errBoom
		},
	}

	svc := New(workbenches, &fakeVaults{}, newFakeDocker())

	err := svc.DeleteWorkbench(context.Background(), vaultID)
	require.Error(t, err)
}
