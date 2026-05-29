-- +goose Up
ALTER TABLE identities_telegram
    ADD COLUMN IF NOT EXISTS photo_url TEXT NOT NULL DEFAULT '';
-- +goose Down
SELECT 1;