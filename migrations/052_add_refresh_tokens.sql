-- +goose Up

ALTER TABLE sessions
    ADD COLUMN refresh_token TEXT UNIQUE,
    ADD COLUMN refresh_expires_at TIMESTAMPTZ;

-- +goose Down

ALTER TABLE sessions
    DROP COLUMN IF EXISTS refresh_token,
    DROP COLUMN IF EXISTS refresh_expires_at;
