-- name: GetMcpDefinition :one
SELECT name, author, description, created_at, owner_user_id, is_community
FROM mcps
WHERE name = $1;

-- name: UpsertMcpDefinition :one
INSERT INTO mcps (name, author, description, owner_user_id, is_community)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (name) DO UPDATE
    SET author        = EXCLUDED.author,
        description   = EXCLUDED.description,
        owner_user_id = EXCLUDED.owner_user_id,
        is_community  = EXCLUDED.is_community
RETURNING name, author, description, created_at, owner_user_id, is_community;

-- name: ListMcpDefinitions :many
SELECT name, author, description, created_at, owner_user_id, is_community
FROM mcps
ORDER BY name;

-- name: DeleteMcpDefinition :exec
DELETE FROM mcps WHERE name = $1;
