-- +goose Up
ALTER TABLE user_permissions
    ADD COLUMN IF NOT EXISTS has_spreadsheets BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
SELECT 1;
