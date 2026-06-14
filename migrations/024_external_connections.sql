-- +goose Up
CREATE TYPE external_provider_type AS ENUM (
    'google_oauth',
    'api_key',
    'password'
);

CREATE TABLE external_connections (
    id               UUID                   PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID                   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider         TEXT                   NOT NULL,
    provider_type    external_provider_type NOT NULL,
    credentials_enc  BYTEA                  NOT NULL,
    metadata         JSONB,
    created_at       TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, provider)
);

CREATE INDEX ON external_connections(user_id);

-- +goose Down
DROP TABLE external_connections;
DROP TYPE external_provider_type;
