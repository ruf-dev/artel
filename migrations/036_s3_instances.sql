-- +goose Up

CREATE TABLE s3_instances (
    id             UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    endpoint       TEXT        NOT NULL,
    region         TEXT        NOT NULL DEFAULT '',
    access_key     TEXT        NOT NULL,
    secret_key_enc BYTEA       NOT NULL,
    use_ssl        BOOLEAN     NOT NULL DEFAULT true,
    path_style     BOOLEAN     NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS s3_instances;
