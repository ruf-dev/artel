package users

import (
	"context"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

type Service struct {
	userRepo repository.Users
}

func New(userRepo repository.Users) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) CreateUser(ctx context.Context, username, password string, roles []string) error {
	_, err := s.userRepo.Create(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "create user in postgres")
	}

	//TODO
	//err = s.couchClient.CreateUser(ctx, username, password, roles)
	//if err != nil {
	//	return rerrors.Wrap(err, "create user in couchdb")
	//}

	return nil
}

func (s *Service) GetUser(ctx context.Context, username string) (domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, username)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "get user from postgres")
	}

	//TODO
	//info, err := s.couchClient.GetUser(ctx, username)
	//if err != nil {
	//	return domain.User{}, rerrors.Wrap(err, "get user from couchdb")
	//}

	//user.Roles = info.Roles
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, username, password string, roles []string) error {
	//TODO
	//err := s.couchClient.UpdateUser(ctx, username, password, roles)
	//if err != nil {
	//	return rerrors.Wrap(err, "update user in couchdb")
	//}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, username string) error {
	user, err := s.userRepo.GetByEmail(ctx, username)
	if err != nil {
		return rerrors.Wrap(err, "get user from postgres")
	}

	//TODO
	//err = s.couchClient.DeleteUser(ctx, username)
	//if err != nil {
	//	return rerrors.Wrap(err, "delete user from couchdb")
	//}

	err = s.userRepo.Delete(ctx, user.Uuid)
	if err != nil {
		return rerrors.Wrap(err, "delete user from postgres")
	}

	return nil
}
