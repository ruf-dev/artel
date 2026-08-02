-- +goose Up

ALTER TABLE vaults ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE vaults ADD COLUMN slug      TEXT UNIQUE;

-- +goose Down

ALTER TABLE vaults DROP COLUMN is_public;
ALTER TABLE vaults DROP COLUMN slug;
