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
	couchAccounts  repository.CouchAccounts
	userPerms      repository.UserPermissionsRepo
}

func New(repo repository.Repo) *Service {
	return &Service{
		vaults:         repo.Vaults(),
		vaultMembers:   repo.VaultMembers(),
		couchInstances: repo.CouchInstances(),
		couchAccounts:  repo.CouchAccounts(),
		userPerms:      repo.UserPermissions(),
	}
}

func (s *Service) liveSyncClient(ctx context.Context, vaultID uuid.UUID) (*couchdb.LiveSyncClient, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	perms, err := s.userPerms.Get(ctx, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting permissions")
	}
	if !perms.HasNotes {
		return nil, rerrors.Wrap(user_errors.NotesNotEnabled)
	}

	_, err = s.vaultMembers.Get(ctx, vaultID, uc.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.Unauthenticated, "not a vault member")
	}

	vault, err := s.vaults.GetByID(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "get vault")
	}

	instance, err := s.couchInstances.Get(ctx, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "get couch instance")
	}

	account, err := s.couchAccounts.GetByUserAndInstance(ctx, uc.UserUuid, vault.CouchInstanceUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "get couch account")
	}

	client := couchdb.NewLiveSyncClient(
		instance.Url,
		vault.CouchDBName,
		account.CouchUsername,
		account.CouchPassword,
	)

	return client, nil
}

func (s *Service) ListFolders(ctx context.Context, vaultID uuid.UUID) ([]string, error) {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return client.ListFolders(ctx)
}

func (s *Service) ListNotes(ctx context.Context, vaultID uuid.UUID) ([]couchdb.NoteEntry, error) {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return client.ListNotes(ctx)
}

func (s *Service) GetNote(ctx context.Context, vaultID uuid.UUID, path string) (couchdb.NoteDoc, error) {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return couchdb.NoteDoc{}, err
	}

	return client.ReadNote(ctx, path)
}

func (s *Service) ListTags(ctx context.Context, vaultID uuid.UUID) ([]string, error) {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return nil, err
	}

	return client.ListTags(ctx)
}

func (s *Service) SaveNote(ctx context.Context, vaultID uuid.UUID, path, content string) error {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return err
	}

	err = client.WriteNote(ctx, path, content)
	if err != nil {
		return rerrors.Wrap(err, "error writing note")
	}
	return nil
}

func (s *Service) MoveNote(ctx context.Context, vaultID uuid.UUID, oldPath, newPath string) error {
	client, err := s.liveSyncClient(ctx, vaultID)
	if err != nil {
		return err
	}

	err = client.MoveNote(ctx, oldPath, newPath)
	if err != nil {
		return rerrors.Wrap(err, "error moving note")
	}
	return nil
}
