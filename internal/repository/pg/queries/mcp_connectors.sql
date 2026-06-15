-- name: InsertMcpConnector :one
INSERT INTO mcp_connectors (mcp_key_id, mcp_name, external_connection_id)
VALUES ($1, $2, $3)
ON CONFLICT (mcp_key_id, mcp_name) DO UPDATE
    SET external_connection_id = EXCLUDED.external_connection_id
RETURNING id, mcp_key_id, mcp_name, external_connection_id, created_at;

-- name: GetMcpConnector :one
SELECT id, mcp_key_id, mcp_name, external_connection_id, created_at
FROM mcp_connectors
WHERE mcp_key_id = $1 AND mcp_name = $2;

-- name: ListMcpConnectorsByKey :many
SELECT id, mcp_key_id, mcp_name, external_connection_id, created_at
FROM mcp_connectors
WHERE mcp_key_id = $1
ORDER BY created_at ASC;

-- name: DeleteMcpConnector :exec
DELETE FROM mcp_connectors
WHERE mcp_key_id = $1 AND mcp_name = $2;
