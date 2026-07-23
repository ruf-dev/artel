-- +goose Up
CREATE TYPE workbench_status AS ENUM (
    'created',
    'running',
    'stopped',
    'removed'
);

CREATE TYPE workbench_auth_mode AS ENUM (
    'api_key',
    'subscription_login'
);

CREATE TABLE workbenches
(
    id             UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id       UUID                 NOT NULL REFERENCES vaults (id) ON DELETE CASCADE,
    user_id        UUID                 NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status         workbench_status     NOT NULL DEFAULT 'created',
    auth_mode      workbench_auth_mode,
    container_id   TEXT,
    volume_name    TEXT                 NOT NULL,
    created_at     TIMESTAMPTZ          NOT NULL DEFAULT NOW(),
    started_at     TIMESTAMPTZ,
    stopped_at     TIMESTAMPTZ,
    UNIQUE(vault_id)
);

CREATE INDEX ON workbenches(user_id);

-- +goose Down
DROP TABLE workbenches;
DROP TYPE workbench_auth_mode;
DROP TYPE workbench_status;
