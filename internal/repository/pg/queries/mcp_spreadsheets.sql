-- name: InsertMcpSpreadsheet :one
INSERT INTO mcp_spreadsheets (user_id, external_connection_id, spreadsheet_id, name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, spreadsheet_id) DO UPDATE
    SET name = EXCLUDED.name
RETURNING id, user_id, external_connection_id, spreadsheet_id, name, created_at;

-- name: ListMcpSpreadsheetsByUser :many
SELECT id, user_id, external_connection_id, spreadsheet_id, name, created_at
FROM mcp_spreadsheets
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: DeleteMcpSpreadsheet :exec
DELETE FROM mcp_spreadsheets
WHERE user_id = $1 AND spreadsheet_id = $2;
