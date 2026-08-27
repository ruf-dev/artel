-- +goose Up

ALTER TABLE system_settings ADD COLUMN system_prompt TEXT NOT NULL DEFAULT 'You are Artel''s in-app assistant, helping a user work with their vault. You may call the provided tools to read or change the user''s data. Every tool call is shown to the user and requires their approval, so prefer one focused call at a time and explain what you are doing.';

CREATE TABLE user_settings
(
    user_id    UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    user_prompt TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE vaults ADD COLUMN prompt TEXT NOT NULL DEFAULT '';
ALTER TABLE vaults ADD COLUMN use_system_prompt BOOLEAN NOT NULL DEFAULT TRUE;

-- +goose Down

ALTER TABLE vaults DROP COLUMN use_system_prompt;
ALTER TABLE vaults DROP COLUMN prompt;
DROP TABLE user_settings;
ALTER TABLE system_settings DROP COLUMN system_prompt;
