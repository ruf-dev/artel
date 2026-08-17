-- name: GetSubscriptionPlan :one
SELECT plan_key, couch_quota_bytes, s3_quota_bytes, max_hot_plug_skills, max_total_skills, features, created_at, updated_at
FROM subscription_plans
WHERE plan_key = $1;

-- name: ListSubscriptionPlans :many
SELECT plan_key, couch_quota_bytes, s3_quota_bytes, max_hot_plug_skills, max_total_skills, features, created_at, updated_at
FROM subscription_plans
ORDER BY plan_key;
