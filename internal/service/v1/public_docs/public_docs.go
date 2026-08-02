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
	"go.redsock.ru/rerrors"
)

type Service struct {
	vaults         repository.Vaults
	couchInstances repository.CouchInstances
}

func New(repo repository.Repo) *Service {
	return &Service{
		vaults:         repo.Vaults(),
		couchInstances: repo.CouchInstances(),
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

func (s *Service) GetVaultBySlug(ctx context.Context, slug string) (domain.Vault, error) {
	vault, err := s.resolvePublicVault(ctx, slug)
	if err != nil {
		return domain.Vault{}, err
	}

	return vault, nil
}

func (s *Service) ListFolders(ctx context.Context, slug string) ([]string, error) {
	vault, err := s.resolvePublicVault(ctx, slug)
	if err != nil {
		return nil, err
	}

	client, err := s.liveSyncClientForPublicVault(ctx, vault)
	if err != nil {
		return nil, err
	}

	return client.ListFolders(ctx)
}

func (s *Service) ListNotes(ctx context.Context, slug string) ([]couchdb.NoteEntry, error) {
	vault, err := s.resolvePublicVault(ctx, slug)
	if err != nil {
		return nil, err
	}

	client, err := s.liveSyncClientForPublicVault(ctx, vault)
	if err != nil {
		return nil, err
	}

	return client.ListNotes(ctx)
}

func (s *Service) GetNote(ctx context.Context, slug, path string) (couchdb.NoteDoc, error) {
	vault, err := s.resolvePublicVault(ctx, slug)
	if err != nil {
		return couchdb.NoteDoc{}, err
	}

	client, err := s.liveSyncClientForPublicVault(ctx, vault)
	if err != nil {
		return couchdb.NoteDoc{}, err
	}

	return client.ReadNote(ctx, path)
}

func (s *Service) ListTags(ctx context.Context, slug string) ([]string, error) {
	vault, err := s.resolvePublicVault(ctx, slug)
	if err != nil {
		return nil, err
	}

	client, err := s.liveSyncClientForPublicVault(ctx, vault)
	if err != nil {
		return nil, err
	}

	tags, err := client.ListTags(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list tags")
	}

	return tags, nil
}
