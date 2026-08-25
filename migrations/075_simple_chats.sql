-- +goose Up

CREATE TABLE simple_chats (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id         UUID        NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id          UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title            TEXT,
    model            TEXT        NOT NULL,
    vault_access     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX simple_chats_vault_id_user_id_idx ON simple_chats(vault_id, user_id);

-- +goose Down

DROP TABLE IF EXISTS simple_chats;
