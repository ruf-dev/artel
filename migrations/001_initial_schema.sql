-- +goose Up

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    email      TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE subscriptions
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    active     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE vaults
(
    id            UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name          TEXT        NOT NULL,
    couch_db_name TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE couch_credentials
(
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    vault_id     UUID        NOT NULL REFERENCES vaults (id) ON DELETE CASCADE,
    host         TEXT        NOT NULL,
    username     TEXT        NOT NULL,
    password_enc BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS couch_credentials;
DROP TABLE IF EXISTS vaults;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS users;
