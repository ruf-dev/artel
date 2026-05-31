-- +goose Up

CREATE TYPE vault_role AS ENUM ('owner', 'reader', 'maintainer');

ALTER TABLE vault_members DROP CONSTRAINT vault_members_role_check;
UPDATE vault_members SET role = 'reader' WHERE role = 'member';
ALTER TABLE vault_members
    ALTER COLUMN role TYPE vault_role USING role::vault_role;

CREATE TABLE vault_invites (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id   UUID        NOT NULL REFERENCES vaults (id) ON DELETE CASCADE,
    created_by UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       vault_role  NOT NULL CHECK (role IN ('reader', 'maintainer')),
    token      TEXT        NOT NULL UNIQUE,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS vault_invites;

ALTER TABLE vault_members
    ALTER COLUMN role TYPE TEXT USING role::TEXT;

UPDATE vault_members SET role = 'member' WHERE role = 'reader';
ALTER TABLE vault_members
    ADD CONSTRAINT vault_members_role_check CHECK (role IN ('owner', 'member'));

DROP TYPE IF EXISTS vault_role;
