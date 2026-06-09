-- +goose Up
ALTER TABLE user_permissions ADD COLUMN has_notes BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose Down
SELECT 1;