-- +goose Up

ALTER TABLE users ADD COLUMN password_hash TEXT NOT NULL DEFAULT '';

CREATE TABLE sessions (
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS sessions;
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;
