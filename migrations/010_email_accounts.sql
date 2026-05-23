-- +goose Up

CREATE TABLE email_accounts (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email        TEXT        NOT NULL,
    imap_host    TEXT        NOT NULL,
    imap_port    INT         NOT NULL,
    smtp_host    TEXT        NOT NULL,
    smtp_port    INT         NOT NULL,
    password_enc BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX email_accounts_user_id_idx ON email_accounts(user_id);

-- +goose Down

DROP TABLE IF EXISTS email_accounts;
