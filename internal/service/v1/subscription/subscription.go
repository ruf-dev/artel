package subscription

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

type Service struct {
	subscriptions repository.Subscriptions
}

func New(subscriptions repository.Subscriptions) *Service {
	return &Service{subscriptions: subscriptions}
}

func (s *Service) CheckActive(ctx context.Context, userUuid uuid.UUID) error {
	sub, err := s.subscriptions.GetByUser(ctx, userUuid)
	if err != nil {
		return rerrors.Wrap(err, "get subscription")
	}
	if !sub.Valid || !sub.V.Active {
		return user_errors.NoActiveSubscription
	}

	return nil
}
