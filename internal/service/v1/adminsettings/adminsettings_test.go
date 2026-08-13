package adminsettings

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// --- fakes -------------------------------------------------------------

// fakeSystemSettingsRepo is a hand-rolled fake of repository.SystemSettingsRepo.
type fakeSystemSettingsRepo struct {
	getFunc                     func(ctx context.Context) (domain.SystemSettings, error)
	updateDefaultDocsVaultFunc  func(ctx context.Context, vaultUuid *uuid.UUID) error
	updateDefaultDocsSourceFunc func(ctx context.Context, source domain.DocsSource) error
}

func (f *fakeSystemSettingsRepo) Get(ctx context.Context) (domain.SystemSettings, error) {
	return f.getFunc(ctx)
}

func (f *fakeSystemSettingsRepo) GetForUpdate(ctx context.Context) (domain.SystemSettings, error) {
	panic("not implemented")
}

func (f *fakeSystemSettingsRepo) UpdateAuthMethods(ctx context.Context, passwordEnabled, telegramEnabled bool) error {
	panic("not implemented")
}

func (f *fakeSystemSettingsRepo) UpdateRegistrationMode(ctx context.Context, mode domain.RegistrationMode) error {
	panic("not implemented")
}

func (f *fakeSystemSettingsRepo) SetSetupToken(ctx context.Context, tokenHash string, issuedAt time.Time) error {
	panic("not implemented")
}

func (f *fakeSystemSettingsRepo) CompleteSetup(ctx context.Context) error {
	panic("not implemented")
}

func (f *fakeSystemSettingsRepo) UpdateDefaultDocsVault(ctx context.Context, vaultUuid *uuid.UUID) error {
	return f.updateDefaultDocsVaultFunc(ctx, vaultUuid)
}

func (f *fakeSystemSettingsRepo) UpdateDefaultDocsSource(ctx context.Context, source domain.DocsSource) error {
	return f.updateDefaultDocsSourceFunc(ctx, source)
}

func (f *fakeSystemSettingsRepo) WithTx(tx *sql.Tx) repository.SystemSettingsRepo {
	panic("not implemented")
}

// fakeVaultsRepo is a hand-rolled fake of repository.Vaults.
type fakeVaultsRepo struct {
	getByIDFunc func(ctx context.Context, id uuid.UUID) (domain.Vault, error)
}

func (f *fakeVaultsRepo) Upsert(
	ctx context.Context, userID, couchInstanceID uuid.UUID, name, couchDBName, status, passphrase string,
) (domain.Vault, error) {
	panic("not implemented")
}

func (f *fakeVaultsRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
	return f.getByIDFunc(ctx, id)
}

func (f *fakeVaultsRepo) GetByNameAndUser(ctx context.Context, userID uuid.UUID, name string) (domain.Vault, error) {
	panic("not implemented")
}

func (f *fakeVaultsRepo) UpdateStatus(ctx context.Context, vaultID uuid.UUID, status string) error {
	panic("not implemented")
}

func (f *fakeVaultsRepo) SetLiveSyncPassphrase(ctx context.Context, vaultID uuid.UUID, passphrase string) error {
	panic("not implemented")
}

func (f *fakeVaultsRepo) ListByMembership(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error) {
	panic("not implemented")
}

func (f *fakeVaultsRepo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeVaultsRepo) LinkS3Bucket(ctx context.Context, vaultID, s3InstanceID uuid.UUID, bucketName string) error {
	panic("not implemented")
}

func (f *fakeVaultsRepo) UnlinkS3Bucket(ctx context.Context, vaultID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeVaultsRepo) SetUseCouchDBForBinaries(ctx context.Context, vaultID uuid.UUID, value bool) error {
	panic("not implemented")
}

func (f *fakeVaultsRepo) PublishVault(ctx context.Context, vaultID uuid.UUID, slug string) (domain.Vault, error) {
	panic("not implemented")
}

func (f *fakeVaultsRepo) UnpublishVault(ctx context.Context, vaultID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeVaultsRepo) GetBySlug(ctx context.Context, slug string) (domain.Vault, error) {
	panic("not implemented")
}

func (f *fakeVaultsRepo) WithTx(tx postgres.DB) repository.Vaults {
	panic("not implemented")
}

// --- tests ---------------------------------------------------------------

func TestUpdateDefaultDocsVault_ClearsWhenNil(t *testing.T) {
	var gotVaultUuid *uuid.UUID
	var gotSource domain.DocsSource
	calledGetByID := false

	settingsRepo := &fakeSystemSettingsRepo{
		updateDefaultDocsVaultFunc: func(ctx context.Context, vaultUuid *uuid.UUID) error {
			gotVaultUuid = vaultUuid

			return nil
		},
		updateDefaultDocsSourceFunc: func(ctx context.Context, source domain.DocsSource) error {
			gotSource = source

			return nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			calledGetByID = true

			return domain.Vault{}, nil
		},
	}

	svc := New(settingsRepo, vaultsRepo)

	err := svc.UpdateDefaultDocsVault(context.Background(), nil, domain.DocsSourceVault)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calledGetByID {
		t.Fatal("expected GetByID not to be called when clearing the default vault")
	}

	if gotVaultUuid != nil {
		t.Fatalf("expected nil vault uuid to be passed through, got %v", gotVaultUuid)
	}

	if gotSource != domain.DocsSourceVault {
		t.Fatalf("expected settingsRepo.UpdateDefaultDocsSource to be called with vault, got %v", gotSource)
	}
}

func TestUpdateDefaultDocsVault_NotFoundWhenVaultMissing(t *testing.T) {
	vaultUuid := uuid.New()

	settingsRepo := &fakeSystemSettingsRepo{
		updateDefaultDocsVaultFunc: func(ctx context.Context, vaultUuid *uuid.UUID) error {
			t.Fatal("expected settingsRepo.UpdateDefaultDocsVault not to be called when the vault doesn't exist")

			return nil
		},
		updateDefaultDocsSourceFunc: func(ctx context.Context, source domain.DocsSource) error {
			t.Fatal("expected settingsRepo.UpdateDefaultDocsSource not to be called when the vault doesn't exist")

			return nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, sql.ErrNoRows
		},
	}

	svc := New(settingsRepo, vaultsRepo)

	err := svc.UpdateDefaultDocsVault(context.Background(), &vaultUuid, domain.DocsSourceVault)
	if !errors.Is(err, user_errors.NotFound) {
		t.Fatalf("expected user_errors.NotFound, got %v", err)
	}
}

func TestUpdateDefaultDocsVault_Success(t *testing.T) {
	vaultUuid := uuid.New()
	var gotVaultUuid *uuid.UUID
	var gotSource domain.DocsSource

	settingsRepo := &fakeSystemSettingsRepo{
		updateDefaultDocsVaultFunc: func(ctx context.Context, vaultUuid *uuid.UUID) error {
			gotVaultUuid = vaultUuid

			return nil
		},
		updateDefaultDocsSourceFunc: func(ctx context.Context, source domain.DocsSource) error {
			gotSource = source

			return nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			return domain.Vault{Uuid: id}, nil
		},
	}

	svc := New(settingsRepo, vaultsRepo)

	err := svc.UpdateDefaultDocsVault(context.Background(), &vaultUuid, domain.DocsSourceVault)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotVaultUuid == nil || *gotVaultUuid != vaultUuid {
		t.Fatalf("expected settingsRepo.UpdateDefaultDocsVault to be called with %s, got %v", vaultUuid, gotVaultUuid)
	}

	if gotSource != domain.DocsSourceVault {
		t.Fatalf("expected settingsRepo.UpdateDefaultDocsSource to be called with vault, got %v", gotSource)
	}
}

// TestUpdateDefaultDocsVault_GithubSource_SkipsVaultCheck covers the DocsSourceGithub branch:
// the vault-existence check must be skipped entirely (no vaultsRepo.GetByID call), the default
// vault cleared, and the source persisted as github.
func TestUpdateDefaultDocsVault_GithubSource_SkipsVaultCheck(t *testing.T) {
	vaultUuid := uuid.New()
	var gotVaultUuid *uuid.UUID
	var gotVaultUuidSet bool
	var gotSource domain.DocsSource

	settingsRepo := &fakeSystemSettingsRepo{
		updateDefaultDocsVaultFunc: func(ctx context.Context, vaultUuid *uuid.UUID) error {
			gotVaultUuid = vaultUuid
			gotVaultUuidSet = true

			return nil
		},
		updateDefaultDocsSourceFunc: func(ctx context.Context, source domain.DocsSource) error {
			gotSource = source

			return nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			t.Fatal("expected vaultsRepo.GetByID not to be called when source is github")

			return domain.Vault{}, nil
		},
	}

	svc := New(settingsRepo, vaultsRepo)

	err := svc.UpdateDefaultDocsVault(context.Background(), &vaultUuid, domain.DocsSourceGithub)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !gotVaultUuidSet || gotVaultUuid != nil {
		t.Fatalf("expected settingsRepo.UpdateDefaultDocsVault to be called with nil, got %v", gotVaultUuid)
	}

	if gotSource != domain.DocsSourceGithub {
		t.Fatalf("expected settingsRepo.UpdateDefaultDocsSource to be called with github, got %v", gotSource)
	}
}

func TestGetSettings_NoDefaultVault(t *testing.T) {
	settingsRepo := &fakeSystemSettingsRepo{
		getFunc: func(ctx context.Context) (domain.SystemSettings, error) {
			return domain.SystemSettings{}, nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			t.Fatal("expected vaultsRepo.GetByID not to be called when no default vault is configured")

			return domain.Vault{}, nil
		},
	}

	svc := New(settingsRepo, vaultsRepo)

	_, vault, err := svc.GetSettings(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vault.Uuid != uuid.Nil {
		t.Fatalf("expected zero-value vault, got %v", vault)
	}
}

func TestGetSettings_ResolvesDefaultVault(t *testing.T) {
	vaultUuid := uuid.New()
	want := domain.Vault{Uuid: vaultUuid, Name: "Docs", Slug: "docs"}

	settingsRepo := &fakeSystemSettingsRepo{
		getFunc: func(ctx context.Context) (domain.SystemSettings, error) {
			return domain.SystemSettings{DefaultDocsVaultID: &vaultUuid}, nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			return want, nil
		},
	}

	svc := New(settingsRepo, vaultsRepo)

	_, vault, err := svc.GetSettings(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vault != want {
		t.Fatalf("expected %v, got %v", want, vault)
	}
}

// TestGetSettings_ToleratesDeletedDefaultVault covers the case where the FK's ON DELETE SET
// NULL hasn't caught up yet (or the row was removed out from under it) — GetSettings should
// still succeed with a zero-value vault rather than failing the whole settings page.
func TestGetSettings_ToleratesDeletedDefaultVault(t *testing.T) {
	vaultUuid := uuid.New()

	settingsRepo := &fakeSystemSettingsRepo{
		getFunc: func(ctx context.Context) (domain.SystemSettings, error) {
			return domain.SystemSettings{DefaultDocsVaultID: &vaultUuid}, nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			return domain.Vault{}, sql.ErrNoRows
		},
	}

	svc := New(settingsRepo, vaultsRepo)

	_, vault, err := svc.GetSettings(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vault.Uuid != uuid.Nil {
		t.Fatalf("expected zero-value vault, got %v", vault)
	}
}
