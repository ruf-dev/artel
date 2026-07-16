package adminsubscriptions

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"go.redsock.ru/rerrors"
)

type Service struct {
	plans         repository.SubscriptionPlansRepo
	subscriptions repository.Subscriptions
}

func New(plans repository.SubscriptionPlansRepo, subscriptions repository.Subscriptions) *Service {
	adminSubscriptions := &Service{plans: plans, subscriptions: subscriptions}

	return adminSubscriptions
}

func (s *Service) ListPlans(ctx context.Context) ([]domain.SubscriptionPlan, error) {
	plans, err := s.plans.List(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing subscription plans")
	}

	return plans, nil
}

func (s *Service) GetUserSubscription(
	ctx context.Context, userUuid uuid.UUID,
) (domain.Subscription, domain.EffectiveSubscription, error) {
	sub, err := s.subscriptions.GetByUser(ctx, userUuid)
	if err != nil {
		return domain.Subscription{}, domain.EffectiveSubscription{}, rerrors.Wrap(err, "error getting subscription")
	}

	effective, err := s.subscriptions.GetWithPlan(ctx, userUuid)
	if err != nil {
		return domain.Subscription{}, domain.EffectiveSubscription{}, rerrors.Wrap(err, "error getting effective subscription")
	}

	return sub, effective, nil
}

// UpdateUserSubscription writes the admin-supplied active flag, plan assignment, and overrides
// in place — the two writes aren't transactional, but they're both idempotent full-replace
// upserts on the same row, so a failure between them just leaves the next admin save to retry.
func (s *Service) UpdateUserSubscription(
	ctx context.Context, sub domain.Subscription,
) (domain.EffectiveSubscription, error) {
	_, err := s.subscriptions.Upsert(ctx, sub.UserUuid, sub.Active)
	if err != nil {
		return domain.EffectiveSubscription{}, rerrors.Wrap(err, "error upserting subscription active state")
	}

	_, err = s.subscriptions.UpsertOverrides(ctx, sub)
	if err != nil {
		return domain.EffectiveSubscription{}, rerrors.Wrap(err, "error upserting subscription overrides")
	}

	effective, err := s.subscriptions.GetWithPlan(ctx, sub.UserUuid)
	if err != nil {
		return domain.EffectiveSubscription{}, rerrors.Wrap(err, "error getting effective subscription")
	}

	return effective, nil
}
