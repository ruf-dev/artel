-- +goose Up
ALTER TABLE vaults
    ADD COLUMN IF NOT EXISTS livesync_passphrase_enc BYTEA;

-- +goose Down
ALTER TABLE vaults
    DROP COLUMN IF EXISTS livesync_passphrase_enc;