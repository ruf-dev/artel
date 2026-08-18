-- +goose Up

ALTER TABLE workbenches DROP CONSTRAINT workbenches_vault_id_key;

ALTER TABLE workbenches ADD CONSTRAINT workbenches_vault_id_user_id_key
    UNIQUE (vault_id, user_id);

-- +goose Down

ALTER TABLE workbenches DROP CONSTRAINT workbenches_vault_id_user_id_key;

ALTER TABLE workbenches ADD CONSTRAINT workbenches_vault_id_key
    UNIQUE (vault_id);
