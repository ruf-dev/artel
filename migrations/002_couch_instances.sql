-- +goose Up

CREATE TABLE couch_instances (
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    url          TEXT        NOT NULL UNIQUE,
    username     TEXT        NOT NULL,
    password_enc BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS couch_instances;
