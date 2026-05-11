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
