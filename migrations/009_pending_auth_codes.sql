-- +goose Up

CREATE TABLE pending_auth_codes (
    code          TEXT        PRIMARY KEY,
    raw_token     TEXT        NOT NULL,
    code_challenge TEXT       NOT NULL,
    redirect_uri  TEXT        NOT NULL,
    client_id     TEXT        NOT NULL,
    expires_at    TIMESTAMPTZ NOT NULL
);

-- +goose Down

DROP TABLE IF EXISTS pending_auth_codes;