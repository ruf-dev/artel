package notes

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/subscription"
	"go.redsock.ru/rerrors"
)

// fakeVaults is a hand-written repository.Vaults exposing only GetByID's behavior as
// configurable — the rest of the interface is unused by resolveVaultForRead.
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

func (f *fakeVaults) UpdateStatus(context.Context, uuid.UUID, string) error { return nil }

func (f *fakeVaults) SetLiveSyncPassphrase(context.Context, uuid.UUID, string) error { return nil }

func (f *fakeVaults) ListByMembership(context.Context, uuid.UUID) ([]domain.Vault, error) {
	return nil, nil
}

func (f *fakeVaults) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) LinkS3Bucket(context.Context, uuid.UUID, uuid.UUID, string) error { return nil }

func (f *fakeVaults) UnlinkS3Bucket(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) SetUseCouchDBForBinaries(context.Context, uuid.UUID, bool) error { return nil }

func (f *fakeVaults) PublishVault(context.Context, uuid.UUID, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) UnpublishVault(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) GetBySlug(context.Context, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) WithTx(sqldb.DB) repository.Vaults { return f }

// fakeVaultMembers is a hand-written repository.VaultMembers exposing only Get's behavior as
// configurable — the rest of the interface is unused by resolveVaultForRead.
type fakeVaultMembers struct {
	get func(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error)
}

func (f *fakeVaultMembers) Add(context.Context, uuid.UUID, uuid.UUID, artel_q.VaultRole) error {
	return nil
}

func (f *fakeVaultMembers) Remove(context.Context, uuid.UUID, uuid.UUID) error { return nil }

func (f *fakeVaultMembers) Get(ctx context.Context, vaultID, userID uuid.UUID) (domain.VaultMember, error) {
	return f.get(ctx, vaultID, userID)
}

func (f *fakeVaultMembers) ListByVault(context.Context, uuid.UUID) ([]domain.VaultMember, error) {
	return nil, nil
}

func (f *fakeVaultMembers) ListByVaultWithUsers(context.Context, uuid.UUID) ([]domain.VaultMemberInfo, error) {
	return nil, nil
}

func (f *fakeVaultMembers) WithTx(sqldb.DB) repository.VaultMembers { return f }

func newResolveTestService(vaults *fakeVaults, members *fakeVaultMembers) *Service {
	return &Service{
		vaults:        vaults,
		vaultMembers:  members,
		subscriptions: subscription.NewFree(),
	}
}

func TestResolveVaultForRead_Member(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()

	vault := domain.Vault{Uuid: vaultID, IsPublic: false}

	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return vault, nil
		},
	}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{VaultUuid: vaultID, UserUuid: userID}, nil
		},
	}

	svc := newResolveTestService(vaults, members)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	got, err := svc.resolveVaultForRead(ctx, vaultID)
	require.NoError(t, err)
	require.Equal(t, vaultID, got.Uuid)
}

func TestResolveVaultForRead_NonMemberPublicVault(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()

	vault := domain.Vault{Uuid: vaultID, IsPublic: true}

	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return vault, nil
		},
	}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, rerrors.Wrap(user_errors.NotFound)
		},
	}

	svc := newResolveTestService(vaults, members)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	got, err := svc.resolveVaultForRead(ctx, vaultID)
	require.NoError(t, err)
	require.Equal(t, vaultID, got.Uuid)
	require.True(t, got.IsPublic)
}

func TestResolveVaultForRead_NonMemberPrivateVault(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()

	vault := domain.Vault{Uuid: vaultID, IsPublic: false}

	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return vault, nil
		},
	}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, rerrors.Wrap(user_errors.NotFound)
		},
	}

	svc := newResolveTestService(vaults, members)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	_, err := svc.resolveVaultForRead(ctx, vaultID)
	require.ErrorIs(t, err, user_errors.Unauthenticated)
}

func TestResolveVaultForRead_NonMemberUnknownVault(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()

	vaults := &fakeVaults{
		getByID: func(context.Context, uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, rerrors.Wrap(user_errors.NotFound)
		},
	}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, rerrors.Wrap(user_errors.NotFound)
		},
	}

	svc := newResolveTestService(vaults, members)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	_, err := svc.resolveVaultForRead(ctx, vaultID)
	require.ErrorIs(t, err, user_errors.Unauthenticated)
}

func TestResolveVaultForRead_Unauthenticated(t *testing.T) {
	svc := newResolveTestService(&fakeVaults{}, &fakeVaultMembers{})

	_, err := svc.resolveVaultForRead(context.Background(), uuid.New())
	require.ErrorIs(t, err, user_errors.Unauthenticated)
}
