-- +goose Up

ALTER TABLE mcp_keys
    ADD COLUMN email_account_id UUID REFERENCES email_accounts(id) ON DELETE SET NULL,
    ADD COLUMN last_accessed_at TIMESTAMPTZ;

-- +goose Down

ALTER TABLE mcp_keys
    DROP COLUMN IF EXISTS email_account_id,
    DROP COLUMN IF EXISTS last_accessed_at;
