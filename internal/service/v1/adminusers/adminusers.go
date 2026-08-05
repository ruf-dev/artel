package adminusers

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
	"go.redsock.ru/rerrors"
)

type Service struct {
	users         repository.Users
	sessions      repository.Sessions
	subscriptions service.SubscriptionService
	auth          service.AuthService
}

func New(
	users repository.Users, sessions repository.Sessions, subscriptions service.SubscriptionService,
	auth service.AuthService,
) *Service {
	adminUsers := &Service{users: users, sessions: sessions, subscriptions: subscriptions, auth: auth}

	return adminUsers
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

	effective, err := s.subscriptions.GetEffective(ctx, userUuid)
	if err != nil {
		return domain.UserDetails{}, rerrors.Wrap(err, "error getting effective subscription")
	}

	details.EffectiveSubscription = effective

	return details, nil
}

func (s *Service) GetUserSessions(ctx context.Context, userUuid uuid.UUID) ([]domain.Session, error) {
	sessions, err := s.sessions.GetByUserID(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting user sessions")
	}

	return sessions, nil
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (domain.User, error) {
	user, err := s.auth.CreateUserUnchecked(ctx, email, password)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error creating user")
	}

	return user, nil
}

func (s *Service) ChangePassword(ctx context.Context, userUuid uuid.UUID, newPassword string) error {
	err := s.auth.ChangePassword(ctx, userUuid, newPassword)
	if err != nil {
		return rerrors.Wrap(err, "error changing password")
	}

	return nil
}
