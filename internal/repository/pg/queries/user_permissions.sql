-- name: GetUserPermissions :one
SELECT user_id, is_administrator, has_emails
FROM user_permissions
WHERE user_id = $1;

-- name: UpsertUserPermissions :one
INSERT INTO user_permissions (user_id, is_administrator, has_emails)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE
    SET is_administrator = EXCLUDED.is_administrator,
        has_emails = EXCLUDED.has_emails
RETURNING user_id, is_administrator, has_emails;

-- name: CreateDefaultUserPermissions :exec
INSERT INTO user_permissions (user_id, is_administrator, has_emails)
VALUES ($1, FALSE, FALSE)
ON CONFLICT DO NOTHING;
