package subscriptionplans

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) WithTx(tx *sql.Tx) repository.SubscriptionPlansRepo {
	return &Repo{q: r.q.WithTx(tx)}
}

func (r *Repo) Get(ctx context.Context, planKey string) (domain.SubscriptionPlan, error) {
	row, err := r.q.GetSubscriptionPlan(ctx, planKey)
	if err != nil {
		return domain.SubscriptionPlan{}, rerrors.Wrap(err, "error getting subscription plan")
	}

	plan, err := planFromRow(row)
	if err != nil {
		return domain.SubscriptionPlan{}, rerrors.Wrap(err, "error decoding plan features")
	}

	return plan, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.SubscriptionPlan, error) {
	rows, err := r.q.ListSubscriptionPlans(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing subscription plans")
	}

	plans := make([]domain.SubscriptionPlan, len(rows))

	for i, row := range rows {
		plan, planErr := planFromListRow(row)
		if planErr != nil {
			return nil, rerrors.Wrap(planErr, "error decoding plan features")
		}

		plans[i] = plan
	}

	return plans, nil
}

func planFromRow(row artel_q.GetSubscriptionPlanRow) (domain.SubscriptionPlan, error) {
	featureSet, err := unmarshalFeatureSet(row.Features)
	if err != nil {
		return domain.SubscriptionPlan{}, err
	}

	plan := domain.SubscriptionPlan{
		PlanKey:          row.PlanKey,
		CouchQuotaBytes:  row.CouchQuotaBytes,
		S3QuotaBytes:     row.S3QuotaBytes,
		MaxHotPlugSkills: int(row.MaxHotPlugSkills),
		MaxTotalSkills:   int(row.MaxTotalSkills),
		Features:         featureSet,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}

	return plan, nil
}

func planFromListRow(row artel_q.ListSubscriptionPlansRow) (domain.SubscriptionPlan, error) {
	featureSet, err := unmarshalFeatureSet(row.Features)
	if err != nil {
		return domain.SubscriptionPlan{}, err
	}

	plan := domain.SubscriptionPlan{
		PlanKey:          row.PlanKey,
		CouchQuotaBytes:  row.CouchQuotaBytes,
		S3QuotaBytes:     row.S3QuotaBytes,
		MaxHotPlugSkills: int(row.MaxHotPlugSkills),
		MaxTotalSkills:   int(row.MaxTotalSkills),
		Features:         featureSet,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}

	return plan, nil
}

// unmarshalFeatureSet decodes subscription_plans.features — domain.FeatureSet's json tags
// match the JSONB column's key names directly.
func unmarshalFeatureSet(raw json.RawMessage) (domain.FeatureSet, error) {
	var featureSet domain.FeatureSet

	err := json.Unmarshal(raw, &featureSet)
	if err != nil {
		return domain.FeatureSet{}, err
	}

	return featureSet, nil
}
