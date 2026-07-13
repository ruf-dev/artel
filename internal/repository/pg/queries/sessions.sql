-- name: CreateSession :one
INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING id, user_id, token, expires_at, created_at;

-- name: GetSessionByToken :one
SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = $1;

-- name: GetSessionWithUser :one
SELECT s.id AS session_id, s.token, s.expires_at AS session_expires_at, s.created_at AS session_created_at,
       u.id AS user_id, u.email, u.username, u.photo_url, u.password_hash,
       u.created_at AS user_created_at, u.updated_at AS user_updated_at
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.token = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = $1;

-- name: GetSessionsByUserID :many
SELECT * FROM sessions WHERE user_id = $1 ORDER BY created_at DESC;
