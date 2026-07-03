-- name: ListMcpToolsByMcpName :many
SELECT id, mcp_name, name, description, input_schema, output_schema, action, created_at
FROM mcp_tools
WHERE mcp_name = $1
ORDER BY name;

-- name: GetMcpTool :one
SELECT id, mcp_name, name, description, input_schema, output_schema, action, created_at
FROM mcp_tools
WHERE mcp_name = $1 AND name = $2;

-- name: ListAllMcpTools :many
SELECT id, mcp_name, name, description, input_schema, output_schema, action, created_at
FROM mcp_tools
ORDER BY mcp_name, name;

-- name: UpsertMcpTool :one
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action
RETURNING id, mcp_name, name, description, input_schema, output_schema, action, created_at;
