-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS photo_url TEXT NOT NULL DEFAULT '';

ALTER TABLE identities_telegram DROP COLUMN IF EXISTS photo_url;

-- +goose Down
ALTER TABLE identities_telegram ADD COLUMN IF NOT EXISTS photo_url TEXT NOT NULL DEFAULT '';

ALTER TABLE users DROP COLUMN IF EXISTS photo_url;