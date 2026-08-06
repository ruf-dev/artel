package vault

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// fakeVaults is a hand-written repository.Vaults exposing only PublishVault's behavior as
// configurable — the rest of the interface is unused by TestPublishVault.
type fakeVaults struct {
	publishVault     func(ctx context.Context, vaultID uuid.UUID, slug string) (domain.Vault, error)
	listByMembership func(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error)
}

func (f *fakeVaults) Upsert(context.Context, uuid.UUID, uuid.UUID, string, string, string, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) GetByID(context.Context, uuid.UUID) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) GetByNameAndUser(context.Context, uuid.UUID, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) UpdateStatus(context.Context, uuid.UUID, string) error { return nil }

func (f *fakeVaults) SetLiveSyncPassphrase(context.Context, uuid.UUID, string) error { return nil }

func (f *fakeVaults) ListByMembership(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error) {
	if f.listByMembership == nil {
		return nil, nil
	}

	return f.listByMembership(ctx, userID)
}

func (f *fakeVaults) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) LinkS3Bucket(context.Context, uuid.UUID, uuid.UUID, string) error { return nil }

func (f *fakeVaults) UnlinkS3Bucket(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) SetUseCouchDBForBinaries(context.Context, uuid.UUID, bool) error { return nil }

func (f *fakeVaults) PublishVault(ctx context.Context, vaultID uuid.UUID, slug string) (domain.Vault, error) {
	if f.publishVault == nil {
		return domain.Vault{}, nil
	}

	return f.publishVault(ctx, vaultID, slug)
}

func (f *fakeVaults) UnpublishVault(context.Context, uuid.UUID) error { return nil }

func (f *fakeVaults) GetBySlug(context.Context, string) (domain.Vault, error) {
	return domain.Vault{}, nil
}

func (f *fakeVaults) WithTx(postgres.DB) repository.Vaults { return f }

// fakeVaultMembers is a hand-written repository.VaultMembers exposing only Get's behavior as
// configurable — the rest of the interface is unused by TestPublishVault.
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

func (f *fakeVaultMembers) WithTx(postgres.DB) repository.VaultMembers { return f }

func newPublishTestService(vaults *fakeVaults, members *fakeVaultMembers) *Service {
	return &Service{
		vaultsRepo:       vaults,
		vaultMembersRepo: members,
	}
}

func TestPublishVault_OwnerSucceeds(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()

	published := domain.Vault{Uuid: vaultID, IsPublic: true, Slug: "my-vault"}

	vaults := &fakeVaults{
		publishVault: func(context.Context, uuid.UUID, string) (domain.Vault, error) {
			return published, nil
		},
	}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleOwner}, nil
		},
	}

	svc := newPublishTestService(vaults, members)
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	got, err := svc.PublishVault(ctx, vaultID, "my-vault")
	require.NoError(t, err)
	require.Equal(t, published, got)
}

func TestPublishVault_NonOwnerMemberRejected(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()

	vaults := &fakeVaults{}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleReader}, nil
		},
	}

	svc := newPublishTestService(vaults, members)
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	_, err := svc.PublishVault(ctx, vaultID, "my-vault")
	require.ErrorIs(t, err, user_errors.NotVaultOwner)
}

func TestPublishVault_NonMemberRejected(t *testing.T) {
	vaultID := uuid.New()
	userID := uuid.New()

	vaults := &fakeVaults{}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{}, rerrors.Wrap(user_errors.NotFound)
		},
	}

	svc := newPublishTestService(vaults, members)
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	_, err := svc.PublishVault(ctx, vaultID, "my-vault")
	require.Error(t, err)
}

func TestPublishVault_InvalidSlug(t *testing.T) {
	cases := []string{
		"",
		"AB",             // too short (and uppercase)
		"UpperCase",      // uppercase
		"has spaces",     // spaces
		"-leading-hyph",  // leading hyphen
		"trailing-hyph-", // trailing hyphen
		"ab",             // too short
	}

	vaultID := uuid.New()
	userID := uuid.New()

	vaults := &fakeVaults{
		publishVault: func(context.Context, uuid.UUID, string) (domain.Vault, error) {
			t.Fatal("repo PublishVault must not be called for an invalid slug")
			return domain.Vault{}, nil
		},
	}
	members := &fakeVaultMembers{
		get: func(context.Context, uuid.UUID, uuid.UUID) (domain.VaultMember, error) {
			return domain.VaultMember{Role: artel_q.VaultRoleOwner}, nil
		},
	}

	svc := newPublishTestService(vaults, members)
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userID})

	for _, slug := range cases {
		_, err := svc.PublishVault(ctx, vaultID, slug)
		require.ErrorIsf(t, err, user_errors.VaultSlugInvalid, "slug %q should be rejected", slug)
	}

	longSlug := ""
	for i := 0; i < 65; i++ {
		longSlug += "a"
	}

	_, err := svc.PublishVault(ctx, vaultID, longSlug)
	require.ErrorIs(t, err, user_errors.VaultSlugInvalid)
}
