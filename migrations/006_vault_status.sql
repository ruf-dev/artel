-- +goose Up
ALTER TABLE vaults ADD COLUMN status TEXT NOT NULL DEFAULT 'ready';

-- +goose Down
ALTER TABLE vaults DROP COLUMN status;
