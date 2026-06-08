package notes

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	vaults         repository.Vaults
	vaultMembers   repository.VaultMembers
	couchInstances repository.CouchInstances
}

func New(repo repository.Repo) *Service {
	return &Service{
		vaults:         repo.Vaults(),
		vaultMembers:   repo.VaultMembers(),
		couchInstances: repo.CouchInstances(),
	}
}

func (s *Service) ListFolders(ctx context.Context, vaultID uuid.UUID) ([]string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	vault, err := s.vaults.GetByID(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
	}

	couchWithAccount, err := s.couchInstances.Pick(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "pick couch instance")
	}

	client := couchdb.NewLiveSyncClient(
		couchWithAccount.Instance.Url,
		vault.CouchDBName,
		couchWithAccount.Instance.Username,
		couchWithAccount.Instance.Password,
	)

	return client.ListFolders(ctx)
}

func (s *Service) ListNotes(ctx context.Context, vaultID uuid.UUID) ([]couchdb.NoteEntry, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	vault, err := s.vaults.GetByID(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
	}

	couchWithAccount, err := s.couchInstances.Pick(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "pick couch instance")
	}

	client := couchdb.NewLiveSyncClient(
		couchWithAccount.Instance.Url,
		vault.CouchDBName,
		couchWithAccount.Instance.Username,
		couchWithAccount.Instance.Password,
	)

	return client.ListNotes(ctx)
}

func (s *Service) GetNote(ctx context.Context, vaultID uuid.UUID, path string) (couchdb.NoteDoc, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return couchdb.NoteDoc{}, rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return couchdb.NoteDoc{}, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	vault, err := s.vaults.GetByID(ctx, vaultID)
	if err != nil {
		return couchdb.NoteDoc{}, rerrors.Wrap(err, "get vault")
	}

	couchWithAccount, err := s.couchInstances.Pick(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return couchdb.NoteDoc{}, rerrors.Wrap(err, "pick couch instance")
	}

	client := couchdb.NewLiveSyncClient(
		couchWithAccount.Instance.Url,
		vault.CouchDBName,
		couchWithAccount.Instance.Username,
		couchWithAccount.Instance.Password,
	)

	return client.ReadNote(ctx, path)
}

func (s *Service) ListTags(ctx context.Context, vaultID uuid.UUID) ([]string, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	_, err := s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	vault, err := s.vaults.GetByID(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
	}

	couchWithAccount, err := s.couchInstances.Pick(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "pick couch instance")
	}

	client := couchdb.NewLiveSyncClient(
		couchWithAccount.Instance.Url,
		vault.CouchDBName,
		couchWithAccount.Instance.Username,
		couchWithAccount.Instance.Password,
	)

	return client.ListTags(ctx)
}
