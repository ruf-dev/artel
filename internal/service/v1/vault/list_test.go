package vault

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
)

// fakeCouchInstances is a hand-written repository.CouchInstances exposing only Get's behavior —
// the rest of the interface is unused by TestListVaults_Role.
type fakeCouchInstances struct{}

func (f *fakeCouchInstances) Register(context.Context, string, string, []byte) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (f *fakeCouchInstances) Get(context.Context, uuid.UUID) (domain.CouchInstance, error) {
	return domain.CouchInstance{Uuid: uuid.New(), Url: "http://couch.test"}, nil
}

func (f *fakeCouchInstances) RandomPick(context.Context) (domain.CouchInstanceWithAccount, error) {
	return domain.CouchInstanceWithAccount{}, nil
}

func (f *fakeCouchInstances) List(context.Context) ([]domain.CouchInstance, error) { return nil, nil }

func (f *fakeCouchInstances) Update(context.Context, uuid.UUID, string, string, []byte) error {
	return nil
}

func (f *fakeCouchInstances) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeCouchInstances) Exists(context.Context) (bool, error) { return false, nil }

func (f *fakeCouchInstances) WithTx(sqldb.DB) repository.CouchInstances { return f }

// fakeCouchAccounts is a hand-written repository.CouchAccounts exposing only
// GetByUserAndInstance's behavior — the rest of the interface is unused by TestListVaults_Role.
type fakeCouchAccounts struct{}

func (f *fakeCouchAccounts) Upsert(context.Context, uuid.UUID, uuid.UUID, string, string) (domain.CouchAccount, error) {
	return domain.CouchAccount{}, nil
}

func (f *fakeCouchAccounts) GetByUserAndInstance(context.Context, uuid.UUID, uuid.UUID) (domain.CouchAccount, error) {
	return domain.CouchAccount{CouchUsername: "u", CouchPassword: "p"}, nil
}

func (f *fakeCouchAccounts) ListByUser(context.Context, uuid.UUID) ([]domain.CouchAccount, error) {
	return nil, nil
}

func (f *fakeCouchAccounts) UpdatePassword(context.Context, string, uuid.UUID, string) error {
	return nil
}

func (f *fakeCouchAccounts) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeCouchAccounts) WithTx(sqldb.DB) repository.CouchAccounts { return f }

// TestListVaults_Role asserts that ListVaults surfaces the caller's membership role (as set by
// the repo's ListByMembership, per-vault) unchanged in each returned domain.Vault: "owner" for
// an owned vault, "reader" for a non-owner member, and "" for a public vault the caller is not
// a member of. This is the read model the frontend uses to distinguish "I can edit this public
// vault" from "I'm just a public viewer".
func TestListVaults_Role(t *testing.T) {
	userID := uuid.New()

	ownedVault := domain.Vault{
		Uuid:              uuid.New(),
		CouchInstanceUuid: uuid.New(),
		Name:              "owned",
		MyRole:            "owner",
	}
	memberVault := domain.Vault{
		Uuid:              uuid.New(),
		CouchInstanceUuid: uuid.New(),
		Name:              "member",
		MyRole:            "reader",
	}
	publicVault := domain.Vault{
		Uuid:              uuid.New(),
		CouchInstanceUuid: uuid.New(),
		Name:              "public",
		IsPublic:          true,
		MyRole:            "", // caller has no vault_members row on this vault
	}

	vaults := &fakeVaults{
		listByMembership: func(context.Context, uuid.UUID) ([]domain.Vault, error) {
			return []domain.Vault{ownedVault, memberVault, publicVault}, nil
		},
	}

	svc := &Service{
		vaultsRepo:         vaults,
		couchAccountsRepo:  &fakeCouchAccounts{},
		couchInstancesRepo: &fakeCouchInstances{},
	}

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	got, err := svc.ListVaults(ctx)
	require.NoError(t, err)
	require.Len(t, got, 3)

	require.Equal(t, "owner", got[0].MyRole)
	require.Equal(t, "reader", got[1].MyRole)
	require.Equal(t, "", got[2].MyRole)
	require.True(t, got[2].IsPublic)
}
