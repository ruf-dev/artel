-- +goose Up

ALTER TABLE vaults ADD COLUMN use_couchdb_for_binaries BOOLEAN NOT NULL DEFAULT true;

-- +goose Down

ALTER TABLE vaults DROP COLUMN IF EXISTS use_couchdb_for_binaries;
