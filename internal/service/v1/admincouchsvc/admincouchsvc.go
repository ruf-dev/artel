package admincouchsvc

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/repository"
)

type Service struct {
	couchInstancesRepo repository.CouchInstances
}

func New(repo repository.Repo) *Service {
	return &Service{
		couchInstancesRepo: repo.CouchInstances(),
	}
}

func (s *Service) ListUsers(ctx context.Context, instanceId string) ([]couchdb.UserListEntry, error) {
	client, err := s.resolveClient(ctx, instanceId)
	if err != nil {
		return nil, rerrors.Wrap(err, "resolving couch client")
	}

	users, err := client.ListUsers(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing couch users")
	}

	return users, nil
}

func (s *Service) DeleteUser(ctx context.Context, instanceId, username string) error {
	client, err := s.resolveClient(ctx, instanceId)
	if err != nil {
		return rerrors.Wrap(err, "resolving couch client")
	}

	err = client.DeleteUser(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "deleting couch user")
	}

	return nil
}

func (s *Service) ChangeUserPassword(ctx context.Context, instanceId, username, newPassword string) error {
	client, err := s.resolveClient(ctx, instanceId)
	if err != nil {
		return rerrors.Wrap(err, "resolving couch client")
	}

	existing, err := client.GetUser(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "getting couch user")
	}

	err = client.UpdateUser(ctx, username, newPassword, existing.Roles)
	if err != nil {
		return rerrors.Wrap(err, "updating couch user")
	}

	return nil
}

func (s *Service) GrantDatabaseAccess(ctx context.Context, instanceId, dbName, username string) error {
	client, err := s.resolveClient(ctx, instanceId)
	if err != nil {
		return rerrors.Wrap(err, "resolving couch client")
	}

	err = client.GrantDatabaseAccess(ctx, dbName, username)
	if err != nil {
		return rerrors.Wrap(err, "granting database access")
	}

	return nil
}

func (s *Service) RevokeDatabaseAccess(ctx context.Context, instanceId, dbName, username string) error {
	client, err := s.resolveClient(ctx, instanceId)
	if err != nil {
		return rerrors.Wrap(err, "resolving couch client")
	}

	err = client.RevokeDatabaseAccess(ctx, dbName, username)
	if err != nil {
		return rerrors.Wrap(err, "revoking database access")
	}

	return nil
}

func (s *Service) ListDatabases(ctx context.Context, instanceId string) ([]string, error) {
	client, err := s.resolveClient(ctx, instanceId)
	if err != nil {
		return nil, rerrors.Wrap(err, "resolving couch client")
	}

	databases, err := client.ListDatabases(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "listing databases")
	}

	return databases, nil
}

func (s *Service) GetUserDatabaseAccess(ctx context.Context, instanceId, username string) ([]string, error) {
	client, err := s.resolveClient(ctx, instanceId)
	if err != nil {
		return nil, rerrors.Wrap(err, "resolving couch client")
	}

	databases, err := client.GetUserDatabaseAccess(ctx, username)
	if err != nil {
		return nil, rerrors.Wrap(err, "getting user database access")
	}

	return databases, nil
}

func (s *Service) resolveClient(ctx context.Context, instanceId string) (*couchdb.Client, error) {
	uid, err := uuid.Parse(instanceId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parsing instance uuid")
	}

	instance, err := s.couchInstancesRepo.Get(ctx, uid)
	if err != nil {
		return nil, rerrors.Wrap(err, "getting couch instance")
	}

	cfg := couchdb.Config{
		BaseURL:  instance.Url,
		User:     instance.Username,
		Password: instance.Password,
	}

	client, err := couchdb.New(cfg)
	if err != nil {
		return nil, rerrors.Wrap(err, "creating couch client")
	}
	return client, nil
}
