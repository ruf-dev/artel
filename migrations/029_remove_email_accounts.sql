-- +goose Up
ALTER TABLE mcp_keys DROP COLUMN IF EXISTS email_account_id;
DROP TABLE IF EXISTS email_accounts;

-- +goose Down
-- intentionally not reversible
