-- +goose Up

CREATE TABLE vault_members (
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    vault_id   UUID        NOT NULL REFERENCES vaults (id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       TEXT        NOT NULL CHECK (role IN ('owner', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (vault_id, user_id)
);

INSERT INTO vault_members (vault_id, user_id, role)
SELECT id, user_id, 'owner' FROM vaults;

-- +goose Down

DROP TABLE IF EXISTS vault_members;
