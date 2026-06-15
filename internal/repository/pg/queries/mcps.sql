-- name: GetMcpDefinition :one
SELECT name, author, description, tools, created_at
FROM mcps
WHERE name = $1;

-- name: UpsertMcpDefinition :one
INSERT INTO mcps (name, author, description, tools)
VALUES ($1, $2, $3, $4)
ON CONFLICT (name) DO UPDATE
    SET author      = EXCLUDED.author,
        description = EXCLUDED.description,
        tools       = EXCLUDED.tools
RETURNING name, author, description, tools, created_at;

-- name: ListMcpDefinitions :many
SELECT name, author, description, tools, created_at
FROM mcps
ORDER BY name;

-- name: DeleteMcpDefinition :exec
DELETE FROM mcps WHERE name = $1;
