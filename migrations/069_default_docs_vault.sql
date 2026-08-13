-- +goose Up

ALTER TABLE system_settings ADD COLUMN default_docs_vault_id UUID REFERENCES vaults(id) ON DELETE SET NULL;

-- +goose Down

ALTER TABLE system_settings DROP COLUMN default_docs_vault_id;
