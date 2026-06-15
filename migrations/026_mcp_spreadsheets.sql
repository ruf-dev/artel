-- +goose Up
CREATE TABLE mcp_spreadsheets (
    id                     UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    external_connection_id UUID        NOT NULL REFERENCES external_connections(id) ON DELETE CASCADE,
    spreadsheet_id         TEXT        NOT NULL,
    name                   TEXT        NOT NULL,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, spreadsheet_id)
);
CREATE INDEX ON mcp_spreadsheets(user_id);

-- +goose Down
DROP TABLE mcp_spreadsheets;
