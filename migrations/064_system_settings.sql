-- +goose Up

CREATE TABLE system_settings
(
    id                     SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    setup_completed        BOOLEAN     NOT NULL DEFAULT FALSE,
    password_auth_enabled  BOOLEAN     NOT NULL DEFAULT TRUE,
    telegram_auth_enabled  BOOLEAN     NOT NULL DEFAULT TRUE,
    registration_mode      TEXT        NOT NULL DEFAULT 'admin_only'
        CHECK (registration_mode IN ('admin_only', 'self_register')),
    setup_token_hash       TEXT,
    setup_token_issued_at  TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO system_settings (id) VALUES (1);

-- +goose Down

DROP TABLE system_settings;
