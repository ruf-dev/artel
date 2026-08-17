-- name: UpsertSubscription :one
INSERT INTO subscriptions (user_id, active)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET active = EXCLUDED.active
RETURNING user_id, active, plan_key, feature_overrides, couch_quota_override_bytes, s3_quota_override_bytes,
    max_hot_plug_skills_override, max_total_skills_override;

-- name: GetSubscriptionByUser :one
SELECT user_id, active, plan_key, feature_overrides, couch_quota_override_bytes, s3_quota_override_bytes,
    max_hot_plug_skills_override, max_total_skills_override
FROM subscriptions
WHERE user_id = $1;

-- name: CreateDefaultSubscription :exec
INSERT INTO subscriptions (user_id, active)
VALUES ($1, FALSE)
ON CONFLICT DO NOTHING;

-- name: GetSubscriptionWithPlan :one
SELECT
    s.user_id, s.active, s.plan_key, s.feature_overrides,
    s.couch_quota_override_bytes, s.s3_quota_override_bytes,
    s.max_hot_plug_skills_override, s.max_total_skills_override,
    p.couch_quota_bytes AS plan_couch_quota_bytes,
    p.s3_quota_bytes AS plan_s3_quota_bytes,
    p.max_hot_plug_skills AS plan_max_hot_plug_skills,
    p.max_total_skills AS plan_max_total_skills,
    p.features AS plan_features
FROM subscriptions s
JOIN subscription_plans p ON p.plan_key = s.plan_key
WHERE s.user_id = $1;

-- name: UpsertSubscriptionOverrides :one
UPDATE subscriptions
SET plan_key = $2,
    feature_overrides = $3,
    couch_quota_override_bytes = $4,
    s3_quota_override_bytes = $5,
    max_hot_plug_skills_override = $6,
    max_total_skills_override = $7
WHERE user_id = $1
RETURNING user_id, active, plan_key, feature_overrides, couch_quota_override_bytes, s3_quota_override_bytes,
    max_hot_plug_skills_override, max_total_skills_override;
