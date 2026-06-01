-- name: InsertTaskTracker :one
INSERT INTO task_trackers (user_id, type, name, api_key_enc, api_token_enc)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, type, name, api_key_enc, api_token_enc, created_at;

-- name: GetTaskTrackerByUuid :one
SELECT id, user_id, type, name, api_key_enc, api_token_enc, created_at
FROM task_trackers
WHERE id = $1;

-- name: ListTaskTrackersByUser :many
SELECT id, user_id, type, name, api_key_enc, api_token_enc, created_at
FROM task_trackers
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteTaskTracker :exec
DELETE FROM task_trackers WHERE id = $1;
