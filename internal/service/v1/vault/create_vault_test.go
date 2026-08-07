package vault

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/repository/pg/tx_manager"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// spyCouchInstances is a hand-written repository.CouchInstances that records whether PickForUser
// or RandomPick was called and with which userID — used to verify CreateVault resolves its couch
// instance via PickForUser (preferring a BYOK instance, falling back to the shared pool) instead
// of RandomPick (which picked from the shared pool unconditionally, with no per-user preference).
type spyCouchInstances struct {
	pickForUserCalled bool
	pickForUserUserID uuid.UUID
	pickForUserErr    error

	randomPickCalled bool
}

func (f *spyCouchInstances) Register(context.Context, string, string, []byte) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (f *spyCouchInstances) Get(context.Context, uuid.UUID) (domain.CouchInstance, error) {
	return domain.CouchInstance{}, nil
}

func (f *spyCouchInstances) RandomPick(context.Context) (domain.CouchInstanceWithAccount, error) {
	f.randomPickCalled = true

	return domain.CouchInstanceWithAccount{}, nil
}

func (f *spyCouchInstances) PickForUser(_ context.Context, userID uuid.UUID) (domain.CouchInstanceWithAccount, error) {
	f.pickForUserCalled = true
	f.pickForUserUserID = userID

	if f.pickForUserErr != nil {
		return domain.CouchInstanceWithAccount{}, f.pickForUserErr
	}

	return domain.CouchInstanceWithAccount{}, nil
}

func (f *spyCouchInstances) GetOwned(context.Context, uuid.UUID) (sql.Null[domain.CouchInstance], error) {
	return sql.Null[domain.CouchInstance]{}, nil
}

func (f *spyCouchInstances) List(context.Context) ([]domain.CouchInstance, error) { return nil, nil }

func (f *spyCouchInstances) Update(context.Context, uuid.UUID, string, string, []byte) error {
	return nil
}

func (f *spyCouchInstances) RegisterOwned(context.Context, uuid.UUID, string, string, []byte) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (f *spyCouchInstances) Delete(context.Context, uuid.UUID) error { return nil }

func (f *spyCouchInstances) DeleteOwnedIfUnreferenced(context.Context, uuid.UUID) error { return nil }

func (f *spyCouchInstances) Exists(context.Context) (bool, error) { return false, nil }

func (f *spyCouchInstances) WithTx(postgres.DB) repository.CouchInstances { return f }

var _ repository.CouchInstances = (*spyCouchInstances)(nil)

// fakeCouchAccountsForCreateVault is a hand-written repository.CouchAccounts stub — every method
// is unused by TestCreateVault_UsesPickForUser since CreateVault returns before reaching
// ensureCouchUserExists in the case under test.
type fakeCouchAccountsForCreateVault struct{}

func (f *fakeCouchAccountsForCreateVault) Upsert(
	context.Context, uuid.UUID, uuid.UUID, string, string,
) (domain.CouchAccount, error) {
	return domain.CouchAccount{}, nil
}

func (f *fakeCouchAccountsForCreateVault) GetByUserAndInstance(
	context.Context, uuid.UUID, uuid.UUID,
) (domain.CouchAccount, error) {
	return domain.CouchAccount{}, nil
}

func (f *fakeCouchAccountsForCreateVault) ListByUser(context.Context, uuid.UUID) ([]domain.CouchAccount, error) {
	return nil, nil
}

func (f *fakeCouchAccountsForCreateVault) UpdatePassword(context.Context, string, uuid.UUID, string) error {
	return nil
}

func (f *fakeCouchAccountsForCreateVault) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeCouchAccountsForCreateVault) WithTx(postgres.DB) repository.CouchAccounts { return f }

var _ repository.CouchAccounts = (*fakeCouchAccountsForCreateVault)(nil)

// noopTxManager runs do against a nil *sql.Tx — safe here because every fake repo's WithTx
// ignores the tx argument it's handed, so no real database connection is needed to exercise
// CreateVault's control flow.
type noopTxManager struct{}

func (noopTxManager) Execute(do func(tx *sql.Tx) error) error {
	return do(nil)
}

var _ tx_manager.TxManager = noopTxManager{}

// TestCreateVault_UsesPickForUser verifies CreateVault resolves its couch instance through
// couchInstancesRepo.PickForUser(ctx, callerUserUuid) — which prefers a BYOK instance the caller
// owns before falling back to the shared admin pool — rather than the old RandomPick, which
// picked from the shared pool unconditionally with no per-user preference.
func TestCreateVault_UsesPickForUser(t *testing.T) {
	userID := uuid.New()

	spyInstances := &spyCouchInstances{pickForUserErr: user_errors.NotFound}

	svc := &Service{
		vaultsRepo:         &fakeVaults{},
		vaultMembersRepo:   &fakeVaultMembers{},
		couchAccountsRepo:  &fakeCouchAccountsForCreateVault{},
		couchInstancesRepo: spyInstances,
		txManager:          noopTxManager{},
	}

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID, UserName: "tester"})

	_, err := svc.CreateVault(ctx, "my-vault")

	require.ErrorIs(t, err, user_errors.NoCouchDbInstance)
	require.True(t, spyInstances.pickForUserCalled, "expected CreateVault to call PickForUser")
	require.Equal(t, userID, spyInstances.pickForUserUserID)
	require.False(t, spyInstances.randomPickCalled, "expected CreateVault not to call RandomPick directly")
}
