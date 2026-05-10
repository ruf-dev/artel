-- name: UpsertSubscription :one
INSERT INTO subscriptions (user_id, active)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET active = EXCLUDED.active, updated_at = NOW()
RETURNING id, user_id, active, created_at, updated_at;

-- name: GetSubscriptionByUser :one
SELECT id, user_id, active, created_at, updated_at FROM subscriptions WHERE user_id = $1;
