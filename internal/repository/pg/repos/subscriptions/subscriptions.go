package subscriptions

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type SubscriptionsRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *SubscriptionsRepo {
	return &SubscriptionsRepo{q: q}
}

func (r *SubscriptionsRepo) WithTx(tx *sql.Tx) repository.Subscriptions {
	return &SubscriptionsRepo{q: r.q.WithTx(tx)}
}

func (r *SubscriptionsRepo) Upsert(ctx context.Context, userID uuid.UUID, active bool) (domain.Subscription, error) {
	row, err := r.q.UpsertSubscription(ctx, artel_q.UpsertSubscriptionParams{
		UserID: userID,
		Active: active,
	})
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error upserting subscription")
	}

	return domain.Subscription{
		UserUuid: row.UserID,
		Active:   row.Active,
	}, nil
}

func (r *SubscriptionsRepo) GetByUser(ctx context.Context, userID uuid.UUID) (domain.Subscription, error) {
	row, err := r.q.GetSubscriptionByUser(ctx, userID)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error getting subscription by user")
	}

	return domain.Subscription{
		UserUuid: row.UserID,
		Active:   row.Active,
	}, nil
}

func (r *SubscriptionsRepo) CreateDefault(ctx context.Context, userID uuid.UUID) error {
	err := r.q.CreateDefaultSubscription(ctx, userID)
	if err != nil {
		return rerrors.Wrap(err, "create default subscription")
	}

	return nil
}
