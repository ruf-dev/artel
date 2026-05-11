-- +goose Up

ALTER TABLE users ADD COLUMN username TEXT NOT NULL DEFAULT '';
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

CREATE TABLE identities_telegram (
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    telegram_id TEXT        NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS identities_telegram;
ALTER TABLE users DROP COLUMN IF EXISTS username;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;