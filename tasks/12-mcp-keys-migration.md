# Task 12 — MCP Keys: DB Migration + sqlc Queries

## Goal

Add the `mcp_keys` table to Postgres and write sqlc queries for it.

## Migration

Create file `migrations/008_mcp_keys.sql` using goose format (like existing migrations).

```sql
-- +goose Up

CREATE TABLE mcp_keys (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id    UUID        NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    key_hash    BYTEA       NOT NULL,
    key_preview TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at  TIMESTAMPTZ
);

CREATE INDEX mcp_keys_vault_id_idx ON mcp_keys(vault_id);
CREATE INDEX mcp_keys_user_id_idx  ON mcp_keys(user_id);

-- +goose Down

DROP TABLE IF EXISTS mcp_keys;
```

## sqlc Queries

Create file `internal/repository/pg/queries/mcp_keys.sql`:

```sql
-- name: CreateMcpKey :one
INSERT INTO mcp_keys (vault_id, user_id, name, key_hash, key_preview)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at;

-- name: ListMcpKeysByVault :many
SELECT id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at
FROM mcp_keys
WHERE vault_id = $1
  AND revoked_at IS NULL
ORDER BY created_at DESC;

-- name: GetMcpKeyByID :one
SELECT id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at
FROM mcp_keys
WHERE id = $1;

-- name: ListActiveMcpKeys :many
SELECT id, vault_id, user_id, name, key_hash, key_preview, created_at, revoked_at
FROM mcp_keys
WHERE vault_id = $1
  AND revoked_at IS NULL;

-- name: RevokeMcpKey :exec
UPDATE mcp_keys
SET revoked_at = NOW()
WHERE id = $1
  AND revoked_at IS NULL;
```

## After Writing Files

Run `sqlc generate` from the repo root to regenerate `internal/repository/pg/generated/`.

## Verification

- `go build ./...` must pass with no errors.
- Generated file `internal/repository/pg/generated/mcp_keys.sql.go` must exist.
