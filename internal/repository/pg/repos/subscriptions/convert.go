package subscriptions

import (
	"database/sql"
	"encoding/json"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

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

func unmarshalFeatureOverrides(raw json.RawMessage) (map[domain.SubscriptionFeature]bool, error) {
	if len(raw) == 0 {
		return map[domain.SubscriptionFeature]bool{}, nil
	}

	rawMap := map[string]bool{}

	err := json.Unmarshal(raw, &rawMap)
	if err != nil {
		return nil, err
	}

	overrides := make(map[domain.SubscriptionFeature]bool, len(rawMap))
	for key, value := range rawMap {
		overrides[domain.SubscriptionFeature(key)] = value
	}

	return overrides, nil
}

func marshalFeatureOverrides(overrides map[domain.SubscriptionFeature]bool) (json.RawMessage, error) {
	rawMap := make(map[string]bool, len(overrides))
	for key, value := range overrides {
		rawMap[string(key)] = value
	}

	encoded, err := json.Marshal(rawMap)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func int64PtrToNull(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}

	nullable := sql.NullInt64{Int64: *v, Valid: true}

	return nullable
}

func nullInt64ToPtr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}

	value := v.Int64

	return &value
}

// subscriptionFromRow converts a raw generated Subscription row into the domain type, decoding
// its JSONB feature_overrides column.
func subscriptionFromRow(row artel_q.Subscription) (domain.Subscription, error) {
	overrides, err := unmarshalFeatureOverrides(row.FeatureOverrides)
	if err != nil {
		return domain.Subscription{}, err
	}

	sub := domain.Subscription{
		UserUuid:                row.UserID,
		Active:                  row.Active,
		PlanKey:                 row.PlanKey,
		FeatureOverrides:        overrides,
		CouchQuotaOverrideBytes: nullInt64ToPtr(row.CouchQuotaOverrideBytes),
		S3QuotaOverrideBytes:    nullInt64ToPtr(row.S3QuotaOverrideBytes),
	}

	return sub, nil
}
