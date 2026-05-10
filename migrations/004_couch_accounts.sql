-- +goose Up

DROP TABLE IF EXISTS couch_credentials;

CREATE TABLE couch_accounts (
    id                 UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id            UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    couch_instance_id  UUID        NOT NULL REFERENCES couch_instances (id) ON DELETE CASCADE,
    couch_username     TEXT        NOT NULL,
    couch_password_enc BYTEA       NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, couch_instance_id)
);

ALTER TABLE vaults ADD COLUMN couch_instance_id UUID REFERENCES couch_instances (id);

-- +goose Down

ALTER TABLE vaults DROP COLUMN IF EXISTS couch_instance_id;
DROP TABLE IF EXISTS couch_accounts;

CREATE TABLE couch_credentials (
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    vault_id     UUID        NOT NULL REFERENCES vaults (id) ON DELETE CASCADE,
    host         TEXT        NOT NULL,
    username     TEXT        NOT NULL,
    password_enc BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
