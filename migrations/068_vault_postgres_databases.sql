-- +goose Up

-- status: 'provisioning' | 'ready' | 'error' (see domain.VaultPostgresStatus).
CREATE TABLE vault_postgres_databases (
    vault_id             UUID        PRIMARY KEY REFERENCES vaults (id) ON DELETE CASCADE,
    postgres_instance_id UUID        NOT NULL REFERENCES postgres_instances (id),
    database_name        TEXT        NOT NULL,
    role_username         TEXT       NOT NULL,
    role_password_enc     BYTEA      NOT NULL,
    status                TEXT       NOT NULL DEFAULT 'provisioning'
                               CHECK (status IN ('provisioning', 'ready', 'error')),
    error_message          TEXT      NOT NULL DEFAULT '',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS vault_postgres_databases;
