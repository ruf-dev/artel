-- name: CreateSession :one
INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING id, user_id, token, expires_at, created_at;

-- name: GetSessionByToken :one
SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = $1;
