-- name: GetUserPermissions :one
SELECT user_id, is_administrator
FROM user_permissions
WHERE user_id = $1;

-- name: UpsertUserPermissions :one
INSERT INTO user_permissions (user_id, is_administrator)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
    SET is_administrator = EXCLUDED.is_administrator
RETURNING user_id, is_administrator;

-- name: CreateDefaultUserPermissions :exec
INSERT INTO user_permissions (user_id, is_administrator)
VALUES ($1, FALSE)
ON CONFLICT DO NOTHING;
