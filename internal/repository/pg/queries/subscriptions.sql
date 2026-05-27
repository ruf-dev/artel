-- name: UpsertSubscription :one
INSERT INTO subscriptions (user_id, active)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET active = EXCLUDED.active
RETURNING user_id, active;

-- name: GetSubscriptionByUser :one
SELECT user_id, active FROM subscriptions WHERE user_id = $1;

-- name: CreateDefaultSubscription :exec
INSERT INTO subscriptions (user_id, active)
VALUES ($1, FALSE)
ON CONFLICT DO NOTHING;
