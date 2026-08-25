-- name: CreateSimpleChat :one
INSERT INTO simple_chats (vault_id, user_id, model, vault_access)
VALUES ($1, $2, $3, $4)
RETURNING id, vault_id, user_id, title, model, vault_access, created_at, updated_at, last_activity_at;

-- name: GetSimpleChatByID :one
SELECT id, vault_id, user_id, title, model, vault_access, created_at, updated_at, last_activity_at
FROM simple_chats
WHERE id = $1;

-- name: ListSimpleChatsByVaultAndUser :many
SELECT id, vault_id, user_id, title, model, vault_access, created_at, updated_at, last_activity_at
FROM simple_chats
WHERE vault_id = $1
  AND user_id = $2
ORDER BY last_activity_at DESC;

-- name: UpdateSimpleChatLastActivity :exec
UPDATE simple_chats
SET last_activity_at = NOW(),
    updated_at        = NOW()
WHERE id = $1;

-- name: DeleteSimpleChat :exec
DELETE
FROM simple_chats
WHERE id = $1;
