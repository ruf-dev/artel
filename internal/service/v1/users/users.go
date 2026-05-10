package users

import (
	"context"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/service"
)

type Service struct {
	db *couchdb.Client
}

func New(db *couchdb.Client) *Service {
	return &Service{db: db}
}

func (s *Service) CreateUser(ctx context.Context, username, password string, roles []string) error {
	err := s.db.CreateUser(ctx, username, password, roles)
	if err != nil {
		return rerrors.Wrap(err, "create user")
	}
	return nil
}

func (s *Service) GetUser(ctx context.Context, username string) (service.User, error) {
	info, err := s.db.GetUser(ctx, username)
	if err != nil {
		return service.User{}, rerrors.Wrap(err, "get user")
	}

	u := service.User{
		Username: info.Name,
		Roles:    info.Roles,
	}
	return u, nil
}

func (s *Service) UpdateUser(ctx context.Context, username, password string, roles []string) error {
	err := s.db.UpdateUser(ctx, username, password, roles)
	if err != nil {
		return rerrors.Wrap(err, "update user")
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, username string) error {
	err := s.db.DeleteUser(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "delete user")
	}
	return nil
}
