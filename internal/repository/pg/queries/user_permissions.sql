-- name: GetUserPermissions :one
SELECT user_id, is_administrator, has_emails, has_task_trackers
FROM user_permissions
WHERE user_id = $1;

-- name: UpsertUserPermissions :one
INSERT INTO user_permissions (user_id, is_administrator, has_emails, has_task_trackers)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
    SET is_administrator = EXCLUDED.is_administrator,
        has_emails = EXCLUDED.has_emails,
        has_task_trackers = EXCLUDED.has_task_trackers
RETURNING user_id, is_administrator, has_emails, has_task_trackers;

-- name: CreateDefaultUserPermissions :exec
INSERT INTO user_permissions (user_id, is_administrator, has_emails, has_task_trackers)
VALUES ($1, FALSE, FALSE, FALSE)
ON CONFLICT DO NOTHING;
