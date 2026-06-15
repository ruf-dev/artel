-- +goose Up
CREATE TABLE mcp_connectors (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    mcp_key_id             UUID        NOT NULL REFERENCES mcp_keys(id) ON DELETE CASCADE,
    mcp_name               TEXT        NOT NULL REFERENCES mcps(name),
    external_connection_id UUID        NOT NULL REFERENCES external_connections(id) ON DELETE CASCADE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(mcp_key_id, mcp_name)
);
CREATE INDEX ON mcp_connectors(mcp_key_id);

-- +goose Down
DROP TABLE mcp_connectors;
