package subscriptions

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
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
	params := artel_q.UpsertSubscriptionParams{
		UserID: userID,
		Active: active,
	}

	row, err := r.q.UpsertSubscription(ctx, params)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error upserting subscription")
	}

	sub, err := subscriptionFromRow(row)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error decoding subscription row")
	}

	return sub, nil
}

func (r *SubscriptionsRepo) GetByUser(ctx context.Context, userID uuid.UUID) (domain.Subscription, error) {
	row, err := r.q.GetSubscriptionByUser(ctx, userID)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error getting subscription by user")
	}

	sub, err := subscriptionFromRow(row)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error decoding subscription row")
	}

	return sub, nil
}

func (r *SubscriptionsRepo) CreateDefault(ctx context.Context, userID uuid.UUID) error {
	err := r.q.CreateDefaultSubscription(ctx, userID)
	if err != nil {
		return rerrors.Wrap(err, "create default subscription")
	}

	return nil
}

// GetWithPlan joins the user's subscription row to its plan and returns the merged view —
// the plan's defaults with the user's feature/quota overrides applied on top.
func (r *SubscriptionsRepo) GetWithPlan(ctx context.Context, userID uuid.UUID) (domain.EffectiveSubscription, error) {
	row, err := r.q.GetSubscriptionWithPlan(ctx, userID)
	if err != nil {
		return domain.EffectiveSubscription{}, rerrors.Wrap(err, "error getting subscription with plan")
	}

	overrides, err := unmarshalFeatureOverrides(row.FeatureOverrides)
	if err != nil {
		return domain.EffectiveSubscription{}, rerrors.Wrap(err, "error decoding feature overrides")
	}

	planFeatures, err := unmarshalFeatureSet(row.PlanFeatures)
	if err != nil {
		return domain.EffectiveSubscription{}, rerrors.Wrap(err, "error decoding plan features")
	}

	features := planFeatures
	for feature, value := range overrides {
		features = features.With(feature, value)
	}

	couchQuota := row.PlanCouchQuotaBytes
	if row.CouchQuotaOverrideBytes.Valid {
		couchQuota = row.CouchQuotaOverrideBytes.Int64
	}

	s3Quota := row.PlanS3QuotaBytes
	if row.S3QuotaOverrideBytes.Valid {
		s3Quota = row.S3QuotaOverrideBytes.Int64
	}

	effective := domain.EffectiveSubscription{
		UserUuid:        row.UserID,
		Active:          row.Active,
		PlanKey:         row.PlanKey,
		Features:        features,
		CouchQuotaBytes: couchQuota,
		S3QuotaBytes:    s3Quota,
	}

	return effective, nil
}

// UpsertOverrides overwrites a user's plan assignment and overrides in place — used by the
// migration backfill and any future admin override tooling, not a live per-request path.
func (r *SubscriptionsRepo) UpsertOverrides(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	overridesJSON, err := marshalFeatureOverrides(sub.FeatureOverrides)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error encoding feature overrides")
	}

	couchOverride := int64PtrToNull(sub.CouchQuotaOverrideBytes)
	s3Override := int64PtrToNull(sub.S3QuotaOverrideBytes)

	params := artel_q.UpsertSubscriptionOverridesParams{
		UserID:                  sub.UserUuid,
		PlanKey:                 sub.PlanKey,
		FeatureOverrides:        overridesJSON,
		CouchQuotaOverrideBytes: couchOverride,
		S3QuotaOverrideBytes:    s3Override,
	}

	row, err := r.q.UpsertSubscriptionOverrides(ctx, params)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error upserting subscription overrides")
	}

	updated, err := subscriptionFromRow(row)
	if err != nil {
		return domain.Subscription{}, rerrors.Wrap(err, "error decoding subscription row")
	}

	return updated, nil
}
