-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, password_hash, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE id = $1;

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
INSERT INTO users (username) VALUES ($1) RETURNING id, email, username, password_hash, created_at, updated_at;

-- name: UpsertTelegramIdentity :exec
INSERT INTO identities_telegram (user_id, telegram_id, photo_url)
VALUES ($1, $2, $3)
ON CONFLICT (telegram_id) DO UPDATE
    SET photo_url = EXCLUDED.photo_url, updated_at = NOW();

-- name: GetTelegramPhotoUrlByUserId :one
SELECT photo_url FROM identities_telegram WHERE user_id = $1;

-- name: GetUserDetails :one
SELECT u.id, u.username, u.email, u.created_at, u.updated_at, u.password_hash,
       up.is_administrator, up.has_emails, up.has_task_trackers, up.has_notes,
       s.active AS subscription_active
FROM users u
LEFT JOIN user_permissions up ON up.user_id = u.id
LEFT JOIN subscriptions s ON s.user_id = u.id
WHERE u.id = $1;
