-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, username, photo_url, password_hash, created_at, updated_at FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserByTelegramId :one
SELECT u.id, u.email, u.username, u.password_hash, u.created_at, u.updated_at
FROM users u
JOIN identities_telegram i ON u.id = i.user_id
WHERE i.telegram_id = $1
FOR UPDATE;

-- name: CreateByUsername :one
INSERT INTO users (username, photo_url) VALUES ($1, $2) RETURNING id, email, username, photo_url, password_hash, created_at, updated_at;

-- name: UpsertTelegramIdentity :exec
INSERT INTO identities_telegram (user_id, telegram_id)
VALUES ($1, $2)
ON CONFLICT (telegram_id) DO UPDATE SET updated_at = NOW();

-- name: GetTelegramPhotoUrlByUserId :one
SELECT photo_url FROM users WHERE id = $1;

-- name: GetTelegramIdByUserId :one
SELECT telegram_id FROM identities_telegram WHERE user_id = $1;

-- name: UpdateUserPhotoUrl :exec
UPDATE users SET photo_url = $2 WHERE id = $1;

-- name: UpdateUserPasswordHash :exec
UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1;

-- name: GetUserDetails :one
-- Feature flags and subscription state are no longer joined here — callers fetch those via
-- SubscriptionService.GetEffective, since they now come from a plan+override merge rather than
-- raw columns on this row.
SELECT u.id, u.username, u.email, u.photo_url, u.created_at, u.updated_at, u.password_hash,
       up.is_administrator
FROM users u
LEFT JOIN user_permissions up ON up.user_id = u.id
WHERE u.id = $1;
