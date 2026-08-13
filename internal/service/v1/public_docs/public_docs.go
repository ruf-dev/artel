// Package public_docs implements the anonymous, auth-exempt read surface backing
// PublicDocsAPI (see internal/api/server/artel_api/public_docs_grpc.pb.go's doc comment and
// docs/architecture.md). Every method here resolves its target vault by slug and only ever
// returns data for a published vault (domain.Vault.IsPublic) — it must never call
// user_context.GetUserContext, and must never depend on repository.VaultMembers or
// service.SubscriptionService, since callers here are never authenticated.
package public_docs

import (
	"context"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/service/v1/public_docs/githubdocs"
	"go.redsock.ru/rerrors"
)

type Service struct {
	vaults         repository.Vaults
	couchInstances repository.CouchInstances
	systemSettings repository.SystemSettingsRepo
	githubDocs     DocsResolver
}

func New(repo repository.Repo) *Service {
	return &Service{
		vaults:         repo.Vaults(),
		couchInstances: repo.CouchInstances(),
		systemSettings: repo.SystemSettings(),
		githubDocs:     githubdocs.NewResolver(),
	}
}

// resolvePublicVault resolves slug to a vault, returning user_errors.NotFound both when the
// slug is unknown and when it resolves to a vault that isn't published — a caller must not be
// able to distinguish "wrong slug" from "vault exists but isn't public".
func (s *Service) resolvePublicVault(ctx context.Context, slug string) (domain.Vault, error) {
	vault, err := s.vaults.GetBySlug(ctx, slug)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(user_errors.NotFound)
	}

	if !vault.IsPublic {
		return domain.Vault{}, rerrors.Wrap(user_errors.NotFound)
	}

	return vault, nil
}

// liveSyncClientForPublicVault builds a couchdb.LiveSyncClient for an already-resolved public
// vault.
//
// Judgment call: notes.Service.liveSyncClientForVault authenticates to CouchDB using the
// CALLING USER's own per-user couch account (repository.CouchAccounts.GetByUserAndInstance),
// which doesn't exist for an anonymous caller here. Instead this uses the vault's CouchInstance
// admin credentials (domain.CouchInstance.Username/Password) — the same admin credentials
// vault.New's newCouchClient uses to create databases/users. CouchDB admin credentials
// authenticate against every database on the instance regardless of that database's _security
// object (see couchClient.SetDatabaseSecurity in internal/service/v1/vault/vault.go, which is
// what scopes non-admin per-user access), so this is sufficient to read a published vault's
// content without any per-user couch account — exactly the vault-level (not user-level)
// credential this anonymous path needs.
func (s *Service) liveSyncClientForPublicVault(ctx context.Context, vault domain.Vault) (*couchdb.LiveSyncClient, error) {
	instance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting couch instance")
	}

	client := couchdb.NewLiveSyncClient(
		instance.Url,
		vault.CouchDBName,
		instance.Username,
		instance.Password,
	)

	return client, nil
}

// resolveDocsSource returns (vault, resolver, err) for slug — either a real published vault +
// couchResolver, or the synthetic GitHub identity vault + s.githubDocs when slug is the
// reserved GitHub docs slug. Centralizes the only CouchDB-vs-GitHub branch in this service.
func (s *Service) resolveDocsSource(ctx context.Context, slug string) (domain.Vault, DocsResolver, error) {
	if slug == domain.ReservedGithubDocsSlug {
		return githubdocs.IdentityVault(), s.githubDocs, nil
	}

	vault, err := s.resolvePublicVault(ctx, slug)
	if err != nil {
		return domain.Vault{}, nil, err
	}

	client, err := s.liveSyncClientForPublicVault(ctx, vault)
	if err != nil {
		return domain.Vault{}, nil, err
	}

	resolver := couchResolver{client: client}

	return vault, resolver, nil
}

func (s *Service) GetVaultBySlug(ctx context.Context, slug string) (domain.Vault, error) {
	vault, _, err := s.resolveDocsSource(ctx, slug)
	if err != nil {
		return domain.Vault{}, err
	}

	return vault, nil
}

func (s *Service) ListFolders(ctx context.Context, slug string) ([]string, error) {
	_, resolver, err := s.resolveDocsSource(ctx, slug)
	if err != nil {
		return nil, err
	}

	return resolver.ListFolders(ctx)
}

func (s *Service) ListNotes(ctx context.Context, slug string) ([]couchdb.NoteEntry, error) {
	_, resolver, err := s.resolveDocsSource(ctx, slug)
	if err != nil {
		return nil, err
	}

	return resolver.ListNotes(ctx)
}

func (s *Service) GetNote(ctx context.Context, slug, path string) (couchdb.NoteDoc, error) {
	_, resolver, err := s.resolveDocsSource(ctx, slug)
	if err != nil {
		return couchdb.NoteDoc{}, err
	}

	return resolver.GetNote(ctx, path)
}

// GetDefaultVault resolves the admin-configured default `/docs` vault
// (SystemSettings.DefaultDocsVaultID). It returns user_errors.NotFound both when unset and when
// the configured vault is no longer published — same "can't distinguish reason" contract
// resolvePublicVault uses for an unknown/unpublished slug.
func (s *Service) GetDefaultVault(ctx context.Context) (domain.Vault, error) {
	settings, err := s.systemSettings.Get(ctx)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting system settings")
	}

	if settings.DefaultDocsSource == domain.DocsSourceGithub {
		return githubdocs.IdentityVault(), nil
	}

	if settings.DefaultDocsVaultID == nil {
		return domain.Vault{}, rerrors.Wrap(user_errors.NotFound)
	}

	vault, err := s.vaults.GetByID(ctx, *settings.DefaultDocsVaultID)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(user_errors.NotFound)
	}

	if !vault.IsPublic {
		return domain.Vault{}, rerrors.Wrap(user_errors.NotFound)
	}

	return vault, nil
}

func (s *Service) ListTags(ctx context.Context, slug string) ([]string, error) {
	_, resolver, err := s.resolveDocsSource(ctx, slug)
	if err != nil {
		return nil, err
	}

	tags, err := resolver.ListTags(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list tags")
	}

	return tags, nil
}
