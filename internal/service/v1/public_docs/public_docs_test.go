package public_docs

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/public_docs/githubdocs"
)

// --- fakes -------------------------------------------------------------

// fakeSystemSettingsRepo is a hand-rolled fake of repository.SystemSettingsRepo.
type fakeSystemSettingsRepo struct {
	getFunc func(ctx context.Context) (domain.SystemSettings, error)
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
	panic("not implemented")
}

func (f *fakeSystemSettingsRepo) UpdateDefaultDocsSource(ctx context.Context, source domain.DocsSource) error {
	panic("not implemented")
}

func (f *fakeSystemSettingsRepo) WithTx(tx *sql.Tx) repository.SystemSettingsRepo {
	panic("not implemented")
}

// fakeVaultsRepo is a hand-rolled fake of repository.Vaults.
type fakeVaultsRepo struct {
	getByIDFunc   func(ctx context.Context, id uuid.UUID) (domain.Vault, error)
	getBySlugFunc func(ctx context.Context, slug string) (domain.Vault, error)
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
	if f.getBySlugFunc == nil {
		panic("not implemented")
	}

	return f.getBySlugFunc(ctx, slug)
}

func (f *fakeVaultsRepo) WithTx(tx postgres.DB) repository.Vaults {
	panic("not implemented")
}

// fakeCouchInstances is a hand-rolled fake of repository.CouchInstances.
type fakeCouchInstances struct {
	getFunc func(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error)
}

func (f *fakeCouchInstances) Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error) {
	panic("not implemented")
}

func (f *fakeCouchInstances) Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error) {
	return f.getFunc(ctx, id)
}

func (f *fakeCouchInstances) RandomPick(ctx context.Context) (domain.CouchInstanceWithAccount, error) {
	panic("not implemented")
}

func (f *fakeCouchInstances) PickForUser(ctx context.Context, userID uuid.UUID) (domain.CouchInstanceWithAccount, error) {
	panic("not implemented")
}

func (f *fakeCouchInstances) GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.CouchInstance], error) {
	panic("not implemented")
}

func (f *fakeCouchInstances) List(ctx context.Context) ([]domain.CouchInstance, error) {
	panic("not implemented")
}

func (f *fakeCouchInstances) Update(ctx context.Context, id uuid.UUID, url, username string, passwordPlain []byte) error {
	panic("not implemented")
}

func (f *fakeCouchInstances) RegisterOwned(
	ctx context.Context, ownerUserID uuid.UUID, url, username string, passwordPlain []byte,
) (uuid.UUID, error) {
	panic("not implemented")
}

func (f *fakeCouchInstances) Delete(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeCouchInstances) DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error {
	panic("not implemented")
}

func (f *fakeCouchInstances) Exists(ctx context.Context) (bool, error) {
	panic("not implemented")
}

func (f *fakeCouchInstances) WithTx(tx postgres.DB) repository.CouchInstances {
	panic("not implemented")
}

// fakeDocsResolver is a hand-rolled fake of DocsResolver, matching fakeSystemSettingsRepo's
// function-field style.
type fakeDocsResolver struct {
	listFoldersFunc func(ctx context.Context) ([]string, error)
	listNotesFunc   func(ctx context.Context) ([]couchdb.NoteEntry, error)
	getNoteFunc     func(ctx context.Context, path string) (couchdb.NoteDoc, error)
	listTagsFunc    func(ctx context.Context) ([]string, error)
}

func (f *fakeDocsResolver) ListFolders(ctx context.Context) ([]string, error) {
	return f.listFoldersFunc(ctx)
}

func (f *fakeDocsResolver) ListNotes(ctx context.Context) ([]couchdb.NoteEntry, error) {
	return f.listNotesFunc(ctx)
}

func (f *fakeDocsResolver) GetNote(ctx context.Context, path string) (couchdb.NoteDoc, error) {
	return f.getNoteFunc(ctx, path)
}

func (f *fakeDocsResolver) ListTags(ctx context.Context) ([]string, error) {
	return f.listTagsFunc(ctx)
}

// --- tests ---------------------------------------------------------------

func TestGetDefaultVault_UnsetReturnsNotFound(t *testing.T) {
	settingsRepo := &fakeSystemSettingsRepo{
		getFunc: func(ctx context.Context) (domain.SystemSettings, error) {
			return domain.SystemSettings{}, nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			t.Fatal("expected vaults.GetByID not to be called when no default vault is configured")

			return domain.Vault{}, nil
		},
	}

	svc := &Service{vaults: vaultsRepo, systemSettings: settingsRepo}

	_, err := svc.GetDefaultVault(context.Background())
	if !errors.Is(err, user_errors.NotFound) {
		t.Fatalf("expected user_errors.NotFound, got %v", err)
	}
}

func TestGetDefaultVault_UnpublishedReturnsNotFound(t *testing.T) {
	vaultUuid := uuid.New()

	settingsRepo := &fakeSystemSettingsRepo{
		getFunc: func(ctx context.Context) (domain.SystemSettings, error) {
			return domain.SystemSettings{DefaultDocsVaultID: &vaultUuid}, nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			return domain.Vault{Uuid: id, IsPublic: false}, nil
		},
	}

	svc := &Service{vaults: vaultsRepo, systemSettings: settingsRepo}

	_, err := svc.GetDefaultVault(context.Background())
	if !errors.Is(err, user_errors.NotFound) {
		t.Fatalf("expected user_errors.NotFound, got %v", err)
	}
}

func TestGetDefaultVault_DeletedVaultReturnsNotFound(t *testing.T) {
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

	svc := &Service{vaults: vaultsRepo, systemSettings: settingsRepo}

	_, err := svc.GetDefaultVault(context.Background())
	if !errors.Is(err, user_errors.NotFound) {
		t.Fatalf("expected user_errors.NotFound, got %v", err)
	}
}

func TestGetDefaultVault_Success(t *testing.T) {
	vaultUuid := uuid.New()
	want := domain.Vault{Uuid: vaultUuid, Name: "Docs", Slug: "docs", IsPublic: true}

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

	svc := &Service{vaults: vaultsRepo, systemSettings: settingsRepo}

	got, err := svc.GetDefaultVault(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestGetDefaultVault_GithubSource_ReturnsIdentityVault covers the DocsSourceGithub branch:
// GetDefaultVault must return the GitHub identity vault directly, without ever resolving
// DefaultDocsVaultID via vaults.GetByID.
func TestGetDefaultVault_GithubSource_ReturnsIdentityVault(t *testing.T) {
	settingsRepo := &fakeSystemSettingsRepo{
		getFunc: func(ctx context.Context) (domain.SystemSettings, error) {
			return domain.SystemSettings{DefaultDocsSource: domain.DocsSourceGithub}, nil
		},
	}

	vaultsRepo := &fakeVaultsRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
			t.Fatal("expected vaults.GetByID not to be called when DefaultDocsSource is github")

			return domain.Vault{}, nil
		},
	}

	svc := &Service{vaults: vaultsRepo, systemSettings: settingsRepo}

	got, err := svc.GetDefaultVault(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := githubdocs.IdentityVault()
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestResolveDocsSource_GithubSlug_ReturnsIdentityVault covers resolveDocsSource's reserved-slug
// branch: it must return the GitHub identity vault and the service's githubDocs resolver,
// without touching vaults or couchInstances at all.
func TestResolveDocsSource_GithubSlug_ReturnsIdentityVault(t *testing.T) {
	vaultsRepo := &fakeVaultsRepo{}
	fakeResolver := &fakeDocsResolver{}

	svc := &Service{vaults: vaultsRepo, githubDocs: fakeResolver}

	vault, resolver, err := svc.resolveDocsSource(context.Background(), domain.ReservedGithubDocsSlug)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := githubdocs.IdentityVault()
	if vault != want {
		t.Fatalf("expected identity vault %v, got %v", want, vault)
	}

	if resolver != DocsResolver(fakeResolver) {
		t.Fatalf("expected the service's githubDocs resolver to be returned")
	}
}

// TestResolveDocsSource_RegularSlug_UsesCouchResolver is a regression test confirming a normal
// published-vault slug still resolves through the CouchDB path (resolvePublicVault +
// liveSyncClientForPublicVault wrapped in couchResolver) unchanged by the reserved-slug branch.
func TestResolveDocsSource_RegularSlug_UsesCouchResolver(t *testing.T) {
	vaultUuid := uuid.New()
	couchInstanceUuid := uuid.New()
	want := domain.Vault{
		Uuid:              vaultUuid,
		Slug:              "docs",
		IsPublic:          true,
		CouchInstanceUuid: couchInstanceUuid,
		CouchDBName:       "vault_docs",
	}

	vaultsRepo := &fakeVaultsRepo{
		getBySlugFunc: func(ctx context.Context, slug string) (domain.Vault, error) {
			return want, nil
		},
	}

	couchInstances := &fakeCouchInstances{
		getFunc: func(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error) {
			return domain.CouchInstance{Uuid: id, Url: "http://couch.local"}, nil
		},
	}

	svc := &Service{vaults: vaultsRepo, couchInstances: couchInstances}

	vault, resolver, err := svc.resolveDocsSource(context.Background(), "docs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vault != want {
		t.Fatalf("expected %v, got %v", want, vault)
	}

	if _, ok := resolver.(couchResolver); !ok {
		t.Fatalf("expected a couchResolver, got %T", resolver)
	}
}
