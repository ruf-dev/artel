-- +goose Up
-- Last 4 hex chars of the raw webhook token, kept alongside secret_hash so the UI can show a
-- masked "hook info" reminder (••••1234) after the one-time token reveal is gone. NULL for
-- manual/provider-linked triggers, same as secret_hash.
ALTER TABLE triggers ADD COLUMN token_suffix TEXT;

-- +goose Down
ALTER TABLE triggers DROP COLUMN token_suffix;
