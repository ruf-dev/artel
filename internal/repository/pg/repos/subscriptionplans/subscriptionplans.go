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

	features, err := unmarshalFeatureSet(row.Features)
	if err != nil {
		return domain.SubscriptionPlan{}, rerrors.Wrap(err, "error decoding plan features")
	}

	plan := domain.SubscriptionPlan{
		PlanKey:         row.PlanKey,
		CouchQuotaBytes: row.CouchQuotaBytes,
		S3QuotaBytes:    row.S3QuotaBytes,
		Features:        features,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
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
