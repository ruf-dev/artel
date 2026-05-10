package subscriptions

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
)

type SubscriptionsRepo struct {
	db sqldb.DB
}

func New(db sqldb.DB) *SubscriptionsRepo {
	return &SubscriptionsRepo{db: db}
}

func (r *SubscriptionsRepo) Upsert(ctx context.Context, userID uuid.UUID, active bool) (domain.Subscription, error) {
	query := `
		INSERT INTO subscriptions (user_id, active)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET active = EXCLUDED.active, updated_at = NOW()
		RETURNING id, user_id, active, created_at, updated_at`

	var s domain.Subscription
	row := r.db.QueryRowContext(ctx, query, userID, active)
	err := row.Scan(&s.Uuid, &s.UserUuid, &s.Active, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error upserting subscription")
	}

	return s, nil
}

func (r *SubscriptionsRepo) GetByUser(ctx context.Context, userID uuid.UUID) (domain.Subscription, error) {
	query := `SELECT id, user_id, active, created_at, updated_at FROM subscriptions WHERE user_id = $1`

	var s domain.Subscription
	row := r.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&s.Uuid, &s.UserUuid, &s.Active, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error getting subscription by user")
	}

	return s, nil
}
