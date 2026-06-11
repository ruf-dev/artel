package adminusers

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

type Service struct {
	users    repository.Users
	sessions repository.Sessions
}

func New(users repository.Users, sessions repository.Sessions) *Service {
	return &Service{users: users, sessions: sessions}
}

func (s *Service) ListUsers(ctx context.Context, req domain.ListUsersReq) ([]domain.User, int64, error) {
	users, total, err := s.users.ListAll(ctx, req)
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "error listing users")
	}
	return users, total, nil
}

func (s *Service) GetUser(ctx context.Context, userUuid uuid.UUID) (domain.UserDetails, error) {
	details, err := s.users.GetDetailsById(ctx, userUuid)
	if err != nil {
		return domain.UserDetails{}, rerrors.Wrap(err, "error getting user details")
	}
	return details, nil
}

func (s *Service) GetUserSessions(ctx context.Context, userUuid uuid.UUID) ([]domain.Session, error) {
	sessions, err := s.sessions.GetByUserID(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting user sessions")
	}
	return sessions, nil
}
