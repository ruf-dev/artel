-- +goose Up
ALTER TABLE user_permissions ADD COLUMN has_emails BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE user_permissions DROP COLUMN has_emails;
