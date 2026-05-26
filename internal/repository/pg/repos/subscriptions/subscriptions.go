package subscriptions

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type SubscriptionsRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *SubscriptionsRepo {
	return &SubscriptionsRepo{q: q}
}

func (r *SubscriptionsRepo) Upsert(ctx context.Context, userID uuid.UUID, active bool) (domain.Subscription, error) {
	params := artel_q.UpsertSubscriptionParams{
		UserID: userID,
		Active: active,
	}
	row, err := r.q.UpsertSubscription(ctx, params)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error upserting subscription")
	}

	s := domain.Subscription{
		Uuid:      row.ID,
		UserUuid:  row.UserID,
		Active:    row.Active,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return s, nil
}

func (r *SubscriptionsRepo) GetByUser(ctx context.Context, userID uuid.UUID) (sql.Null[domain.Subscription], error) {
	row, err := r.q.GetSubscriptionByUser(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.Null[domain.Subscription]{}, nil
	}
	if err != nil {
		return sql.Null[domain.Subscription]{}, rerrors.Wrap(err, "error getting subscription by user")
	}

	return sql.Null[domain.Subscription]{
		V: domain.Subscription{
			Uuid:      row.ID,
			UserUuid:  row.UserID,
			Active:    row.Active,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		},
		Valid: true,
	}, nil
}
