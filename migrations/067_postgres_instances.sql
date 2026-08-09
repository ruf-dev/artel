-- +goose Up

CREATE TABLE postgres_instances (
    id             UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    host           TEXT        NOT NULL,
    port           INTEGER     NOT NULL DEFAULT 5432,
    admin_database TEXT        NOT NULL DEFAULT 'postgres',
    username       TEXT        NOT NULL,
    password_enc   BYTEA       NOT NULL,
    ssl_mode       TEXT        NOT NULL DEFAULT 'disable',
    owner_user_id  UUID        REFERENCES users(id) ON DELETE CASCADE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS postgres_instances;
