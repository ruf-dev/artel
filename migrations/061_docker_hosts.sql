-- +goose Up

CREATE TABLE docker_hosts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url        TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE workbenches ADD COLUMN docker_host_id UUID REFERENCES docker_hosts(id) ON DELETE RESTRICT;

-- +goose Down

ALTER TABLE workbenches DROP COLUMN docker_host_id;

DROP TABLE IF EXISTS docker_hosts;
