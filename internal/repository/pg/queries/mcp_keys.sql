-- name: CreateMcpKey :one
INSERT INTO mcp_keys (id, vault_id, user_id, name, key_hash, key_preview)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at, email_account_id, last_accessed_at;

-- name: ListMcpKeysByVault :many
SELECT id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at, email_account_id, last_accessed_at
FROM mcp_keys
WHERE vault_id = $1
  AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: GetMcpKeyByID :one
SELECT id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at, email_account_id, last_accessed_at
FROM mcp_keys
WHERE id = $1;

-- name: ListActiveMcpKeys :many
SELECT id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at, email_account_id, last_accessed_at
FROM mcp_keys
WHERE vault_id = $1
  AND revoked_at IS NULL;

-- name: ListMcpKeysByUser :many
SELECT id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at, email_account_id, last_accessed_at
FROM mcp_keys
WHERE user_id = $1
  AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: RevokeMcpKey :exec
UPDATE mcp_keys
SET revoked_at = NOW()
WHERE id = $1
  AND revoked_at IS NULL;

-- name: SetMcpKeyAccess :exec
UPDATE mcp_keys
SET vault_id = $2, email_account_id = $3
WHERE id = $1
  AND user_id = $4;

-- name: TouchMcpKeyLastAccessed :exec
UPDATE mcp_keys
SET last_accessed_at = NOW()
WHERE id = $1;
