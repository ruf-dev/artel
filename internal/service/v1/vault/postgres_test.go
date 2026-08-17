package vault

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// --- fakes -------------------------------------------------------------

// fakeVaultMembersRepo is a hand-rolled fake of repository.VaultMembers.
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

// fakePostgresInstancesRepo is a hand-rolled fake of repository.PostgresInstances.
type fakePostgresInstancesRepo struct {
	pickForUserFunc func(ctx context.Context, userID uuid.UUID) (domain.PostgresInstance, error)
	getFunc         func(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error)
}

func (f *fakePostgresInstancesRepo) Register(ctx context.Context, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Get(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
	return f.getFunc(ctx, id)
}

func (f *fakePostgresInstancesRepo) RandomPick(ctx context.Context) (domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) PickForUser(ctx context.Context, userID uuid.UUID) (domain.PostgresInstance, error) {
	return f.pickForUserFunc(ctx, userID)
}

func (f *fakePostgresInstancesRepo) GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.PostgresInstance], error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) List(ctx context.Context) ([]domain.PostgresInstance, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Update(ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) error {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) RegisterOwned(ctx context.Context, ownerUserID uuid.UUID, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string) (uuid.UUID, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Delete(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) Exists(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (f *fakePostgresInstancesRepo) WithTx(tx postgres.DB) repository.PostgresInstances {
	panic("not implemented")
}

// fakeVaultPostgresRepo is a hand-rolled fake of repository.VaultPostgresDatabases.
type fakeVaultPostgresRepo struct {
	createFunc     func(ctx context.Context, vaultID, postgresInstanceID uuid.UUID, dbName, roleUsername string, rolePasswordPlain []byte) (domain.VaultPostgresDatabase, error)
	getByVaultFunc func(ctx context.Context, vaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error)
	markReadyFunc  func(ctx context.Context, vaultID uuid.UUID) error
	markErrorFunc  func(ctx context.Context, vaultID uuid.UUID, message string) error
	deleteFunc     func(ctx context.Context, vaultID uuid.UUID) error
}

func (f *fakeVaultPostgresRepo) Create(ctx context.Context, vaultID, postgresInstanceID uuid.UUID, dbName, roleUsername string, rolePasswordPlain []byte) (domain.VaultPostgresDatabase, error) {
	return f.createFunc(ctx, vaultID, postgresInstanceID, dbName, roleUsername, rolePasswordPlain)
}

func (f *fakeVaultPostgresRepo) GetByVaultID(ctx context.Context, vaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
	return f.getByVaultFunc(ctx, vaultID)
}

func (f *fakeVaultPostgresRepo) MarkReady(ctx context.Context, vaultID uuid.UUID) error {
	return f.markReadyFunc(ctx, vaultID)
}

func (f *fakeVaultPostgresRepo) MarkError(ctx context.Context, vaultID uuid.UUID, message string) error {
	return f.markErrorFunc(ctx, vaultID, message)
}

func (f *fakeVaultPostgresRepo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	return f.deleteFunc(ctx, vaultID)
}

func (f *fakeVaultPostgresRepo) WithTx(tx postgres.DB) repository.VaultPostgresDatabases {
	panic("not implemented")
}

// fakeSubscriptionService is a hand-rolled fake of service.SubscriptionService.
type fakeSubscriptionService struct {
	checkFeatureFunc func(ctx context.Context, userUuid uuid.UUID, feature domain.SubscriptionFeature) error
}

func (f *fakeSubscriptionService) CheckActive(ctx context.Context, userUuid uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeSubscriptionService) HasFeature(ctx context.Context, userUuid uuid.UUID, feature domain.SubscriptionFeature) (bool, error) {
	panic("not implemented")
}

func (f *fakeSubscriptionService) CheckFeature(ctx context.Context, userUuid uuid.UUID, feature domain.SubscriptionFeature) error {
	return f.checkFeatureFunc(ctx, userUuid, feature)
}

func (f *fakeSubscriptionService) GetEffective(ctx context.Context, userUuid uuid.UUID) (domain.EffectiveSubscription, error) {
	panic("not implemented")
}

func (f *fakeSubscriptionService) GetUsage(ctx context.Context, userUuid uuid.UUID) (domain.StorageUsage, error) {
	panic("not implemented")
}

func (f *fakeSubscriptionService) CheckStorageQuota(ctx context.Context, userUuid uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeSubscriptionService) CheckSkillLimit(ctx context.Context, vaultUuid uuid.UUID, wantHotPlug bool) error {
	panic("not implemented")
}

// fakeAdminClient is a hand-rolled fake of the local adminClient interface.
type fakeAdminClient struct {
	ensureRoleFunc     func(ctx context.Context, roleName, password string) error
	ensureDatabaseFunc func(ctx context.Context, dbName, ownerRole string) error
	dropDatabaseFunc   func(ctx context.Context, dbName string) error
	dropRoleFunc       func(ctx context.Context, roleName string) error
	closeFunc          func() error
}

func (f *fakeAdminClient) EnsureRole(ctx context.Context, roleName, password string) error {
	return f.ensureRoleFunc(ctx, roleName, password)
}

func (f *fakeAdminClient) EnsureDatabase(ctx context.Context, dbName, ownerRole string) error {
	return f.ensureDatabaseFunc(ctx, dbName, ownerRole)
}

func (f *fakeAdminClient) DropDatabase(ctx context.Context, dbName string) error {
	if f.dropDatabaseFunc == nil {
		return nil
	}

	return f.dropDatabaseFunc(ctx, dbName)
}

func (f *fakeAdminClient) DropRole(ctx context.Context, roleName string) error {
	if f.dropRoleFunc == nil {
		return nil
	}

	return f.dropRoleFunc(ctx, roleName)
}

func (f *fakeAdminClient) Close() error {
	if f.closeFunc == nil {
		return nil
	}

	return f.closeFunc()
}

// --- tests ---------------------------------------------------------------

func withUser(userUuid uuid.UUID) context.Context {
	uc := user_context.UserContext{UserUuid: userUuid, UserName: "tester"}

	return user_context.WithUserContext(context.Background(), uc)
}

func TestEnablePostgresDatabase_AlreadyEnabled(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{VaultUuid: gotVaultID, UserUuid: gotUserID, Role: artel_q.VaultRoleOwner}, nil
		},
	}

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{
				Valid: true,
				V:     domain.VaultPostgresDatabase{VaultUuid: gotVaultID, Status: domain.VaultPostgresStatusReady},
			}, nil
		},
	}

	svc := &Service{
		vaultMembersRepo:  membersRepo,
		vaultPostgresRepo: vpgRepo,
	}

	_, err := svc.EnablePostgresDatabase(withUser(userUuid), vaultID)
	if !errors.Is(err, user_errors.PostgresDatabaseAlreadyEnabled) {
		t.Fatalf("expected PostgresDatabaseAlreadyEnabled, got %v", err)
	}
}

func TestEnablePostgresDatabase_FeatureNotEnabled(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()
	wantErr := user_errors.FeatureNotEnabled

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{VaultUuid: gotVaultID, UserUuid: gotUserID, Role: artel_q.VaultRoleOwner}, nil
		},
	}

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{Valid: false}, nil
		},
	}

	subs := &fakeSubscriptionService{
		checkFeatureFunc: func(ctx context.Context, gotUserUuid uuid.UUID, feature domain.SubscriptionFeature) error {
			if feature != domain.FeaturePostgres {
				t.Fatalf("expected FeaturePostgres, got %v", feature)
			}

			return wantErr
		},
	}

	svc := &Service{
		vaultMembersRepo:  membersRepo,
		vaultPostgresRepo: vpgRepo,
		subscriptions:     subs,
	}

	_, err := svc.EnablePostgresDatabase(withUser(userUuid), vaultID)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected FeatureNotEnabled, got %v", err)
	}
}

func TestEnablePostgresDatabase_Success(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()
	instance := domain.PostgresInstance{Uuid: uuid.New(), Host: "db.internal"}

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{VaultUuid: gotVaultID, UserUuid: gotUserID, Role: artel_q.VaultRoleOwner}, nil
		},
	}

	instancesRepo := &fakePostgresInstancesRepo{
		pickForUserFunc: func(ctx context.Context, gotUserID uuid.UUID) (domain.PostgresInstance, error) {
			return instance, nil
		},
	}

	var markedReady bool

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{Valid: false}, nil
		},
		createFunc: func(ctx context.Context, gotVaultID, gotInstanceID uuid.UUID, dbName, roleUsername string, rolePasswordPlain []byte) (domain.VaultPostgresDatabase, error) {
			if gotVaultID != vaultID {
				t.Fatalf("unexpected vault id: %s", gotVaultID)
			}

			if gotInstanceID != instance.Uuid {
				t.Fatalf("unexpected instance id: %s", gotInstanceID)
			}

			return domain.VaultPostgresDatabase{
				VaultUuid: gotVaultID, PostgresInstanceUuid: gotInstanceID,
				DatabaseName: dbName, RoleUsername: roleUsername, Status: domain.VaultPostgresStatusProvisioning,
			}, nil
		},
		markReadyFunc: func(ctx context.Context, gotVaultID uuid.UUID) error {
			markedReady = true

			return nil
		},
	}

	subs := &fakeSubscriptionService{
		checkFeatureFunc: func(ctx context.Context, gotUserUuid uuid.UUID, feature domain.SubscriptionFeature) error {
			return nil
		},
	}

	var ensuredRole, ensuredDB bool

	svc := &Service{
		vaultMembersRepo:      membersRepo,
		postgresInstancesRepo: instancesRepo,
		vaultPostgresRepo:     vpgRepo,
		subscriptions:         subs,
		newAdminClient: func(gotInstance domain.PostgresInstance) (adminClient, error) {
			if gotInstance.Uuid != instance.Uuid {
				t.Fatalf("unexpected instance passed to newAdminClient: %s", gotInstance.Uuid)
			}

			return &fakeAdminClient{
				ensureRoleFunc: func(ctx context.Context, roleName, password string) error {
					ensuredRole = true

					return nil
				},
				ensureDatabaseFunc: func(ctx context.Context, dbName, ownerRole string) error {
					ensuredDB = true

					return nil
				},
			}, nil
		},
	}

	got, err := svc.EnablePostgresDatabase(withUser(userUuid), vaultID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ensuredRole || !ensuredDB {
		t.Fatalf("expected EnsureRole and EnsureDatabase to be called, got role=%v db=%v", ensuredRole, ensuredDB)
	}

	if !markedReady {
		t.Fatal("expected vaultPostgresRepo.MarkReady to be called")
	}

	if got.Status != domain.VaultPostgresStatusReady {
		t.Fatalf("expected returned status ready, got %v", got.Status)
	}
}

func TestEnablePostgresDatabase_ProvisionFailureMarksError(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()
	instance := domain.PostgresInstance{Uuid: uuid.New()}
	provisionErr := errors.New("connection refused")

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleOwner}, nil
		},
	}

	instancesRepo := &fakePostgresInstancesRepo{
		pickForUserFunc: func(ctx context.Context, gotUserID uuid.UUID) (domain.PostgresInstance, error) {
			return instance, nil
		},
	}

	var markedErrorMsg string

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{Valid: false}, nil
		},
		createFunc: func(ctx context.Context, gotVaultID, gotInstanceID uuid.UUID, dbName, roleUsername string, rolePasswordPlain []byte) (domain.VaultPostgresDatabase, error) {
			return domain.VaultPostgresDatabase{Status: domain.VaultPostgresStatusProvisioning}, nil
		},
		markErrorFunc: func(ctx context.Context, gotVaultID uuid.UUID, message string) error {
			markedErrorMsg = message

			return nil
		},
	}

	subs := &fakeSubscriptionService{
		checkFeatureFunc: func(ctx context.Context, gotUserUuid uuid.UUID, feature domain.SubscriptionFeature) error {
			return nil
		},
	}

	svc := &Service{
		vaultMembersRepo:      membersRepo,
		postgresInstancesRepo: instancesRepo,
		vaultPostgresRepo:     vpgRepo,
		subscriptions:         subs,
		newAdminClient: func(gotInstance domain.PostgresInstance) (adminClient, error) {
			return &fakeAdminClient{
				ensureRoleFunc: func(ctx context.Context, roleName, password string) error {
					return provisionErr
				},
			}, nil
		},
	}

	_, err := svc.EnablePostgresDatabase(withUser(userUuid), vaultID)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if markedErrorMsg == "" {
		t.Fatal("expected vaultPostgresRepo.MarkError to be called with a message")
	}
}

func TestGetPostgresDatabase_NotMember(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()
	wantErr := errors.New("not found")

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, wantErr
		},
	}

	svc := &Service{vaultMembersRepo: membersRepo}

	_, err := svc.GetPostgresDatabase(withUser(userUuid), vaultID)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped membership error, got %v", err)
	}
}

func TestGetPostgresDatabase_NotEnabled(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, nil
		},
	}

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{Valid: false}, nil
		},
	}

	svc := &Service{vaultMembersRepo: membersRepo, vaultPostgresRepo: vpgRepo}

	got, err := svc.GetPostgresDatabase(withUser(userUuid), vaultID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Valid {
		t.Fatal("expected Valid: false when postgres was never enabled")
	}
}

func TestDisablePostgresDatabase_NotOwner(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleReader}, nil
		},
	}

	svc := &Service{vaultMembersRepo: membersRepo}

	err := svc.DisablePostgresDatabase(withUser(userUuid), vaultID)
	if !errors.Is(err, user_errors.NotVaultOwner) {
		t.Fatalf("expected NotVaultOwner, got %v", err)
	}
}

func TestDisablePostgresDatabase_NoopWhenNotEnabled(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleOwner}, nil
		},
	}

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{Valid: false}, nil
		},
	}

	svc := &Service{vaultMembersRepo: membersRepo, vaultPostgresRepo: vpgRepo}

	err := svc.DisablePostgresDatabase(withUser(userUuid), vaultID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDisablePostgresDatabase_Success(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()
	instance := domain.PostgresInstance{Uuid: uuid.New()}
	existing := domain.VaultPostgresDatabase{
		VaultUuid: vaultID, PostgresInstanceUuid: instance.Uuid,
		DatabaseName: "vault_x", RoleUsername: "vaultpg_x", Status: domain.VaultPostgresStatusReady,
	}

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleOwner}, nil
		},
	}

	instancesRepo := &fakePostgresInstancesRepo{
		getFunc: func(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
			return instance, nil
		},
	}

	var deletedVaultID uuid.UUID

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{Valid: true, V: existing}, nil
		},
		deleteFunc: func(ctx context.Context, gotVaultID uuid.UUID) error {
			deletedVaultID = gotVaultID

			return nil
		},
	}

	var droppedDB, droppedRole bool

	svc := &Service{
		vaultMembersRepo:      membersRepo,
		postgresInstancesRepo: instancesRepo,
		vaultPostgresRepo:     vpgRepo,
		newAdminClient: func(gotInstance domain.PostgresInstance) (adminClient, error) {
			return &fakeAdminClient{
				dropDatabaseFunc: func(ctx context.Context, dbName string) error {
					droppedDB = true

					if dbName != existing.DatabaseName {
						t.Fatalf("unexpected db name: %s", dbName)
					}

					return nil
				},
				dropRoleFunc: func(ctx context.Context, roleName string) error {
					droppedRole = true

					if roleName != existing.RoleUsername {
						t.Fatalf("unexpected role name: %s", roleName)
					}

					return nil
				},
			}, nil
		},
	}

	err := svc.DisablePostgresDatabase(withUser(userUuid), vaultID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !droppedDB || !droppedRole {
		t.Fatalf("expected DropDatabase and DropRole to be called, got db=%v role=%v", droppedDB, droppedRole)
	}

	if deletedVaultID != vaultID {
		t.Fatalf("expected vaultPostgresRepo.Delete to be called with %s, got %s", vaultID, deletedVaultID)
	}
}

// TestDisablePostgresDatabase_ToleratesDropFailure confirms that DropDatabase/DropRone errors
// don't block deleting the row — best-effort per DisablePostgresDatabase's doc comment, since
// vaultpg.AdminClient's DropDatabase/DropRole are already IF EXISTS/idempotent and a failure here
// most likely means the target server is unreachable, not that the drop itself was rejected.
func TestDisablePostgresDatabase_ToleratesDropFailure(t *testing.T) {
	userUuid := uuid.New()
	vaultID := uuid.New()
	instance := domain.PostgresInstance{Uuid: uuid.New()}
	existing := domain.VaultPostgresDatabase{VaultUuid: vaultID, PostgresInstanceUuid: instance.Uuid}

	membersRepo := &fakeVaultMembersRepo{
		getFunc: func(ctx context.Context, gotVaultID, gotUserID uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleOwner}, nil
		},
	}

	instancesRepo := &fakePostgresInstancesRepo{
		getFunc: func(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
			return instance, nil
		},
	}

	deleteCalled := false

	vpgRepo := &fakeVaultPostgresRepo{
		getByVaultFunc: func(ctx context.Context, gotVaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
			return sql.Null[domain.VaultPostgresDatabase]{Valid: true, V: existing}, nil
		},
		deleteFunc: func(ctx context.Context, gotVaultID uuid.UUID) error {
			deleteCalled = true

			return nil
		},
	}

	svc := &Service{
		vaultMembersRepo:      membersRepo,
		postgresInstancesRepo: instancesRepo,
		vaultPostgresRepo:     vpgRepo,
		newAdminClient: func(gotInstance domain.PostgresInstance) (adminClient, error) {
			return &fakeAdminClient{
				dropDatabaseFunc: func(ctx context.Context, dbName string) error {
					return errors.New("connection refused")
				},
				dropRoleFunc: func(ctx context.Context, roleName string) error {
					return errors.New("connection refused")
				},
			}, nil
		},
	}

	err := svc.DisablePostgresDatabase(withUser(userUuid), vaultID)
	if err != nil {
		t.Fatalf("expected drop failures to be tolerated, got error: %v", err)
	}

	if !deleteCalled {
		t.Fatal("expected vaultPostgresRepo.Delete to still be called despite drop failures")
	}
}

func TestVaultPostgresIdentifiers(t *testing.T) {
	id := uuid.MustParse("11111111-2222-3333-4444-555555555555")

	dbName := vaultPostgresDatabaseName(id)
	if dbName != "vault_11111111222233334444555555555555" {
		t.Fatalf("unexpected db name: %s", dbName)
	}

	roleName := vaultPostgresRoleUsername(id)
	if roleName != "vaultpg_11111111222233334444555555555555" {
		t.Fatalf("unexpected role name: %s", roleName)
	}

	if len(dbName) > 63 || len(roleName) > 63 {
		t.Fatalf("identifier exceeds postgres's 63-byte limit: db=%d role=%d", len(dbName), len(roleName))
	}
}
