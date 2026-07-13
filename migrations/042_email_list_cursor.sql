-- +goose Up
-- Adds UID-cursor pagination (after_uid/before_uid) to email/list_emails' input_schema, so
-- callers can page forward through mail already seen or start from the oldest message instead of
-- always getting the newest `limit`. See internal/clients/imap/client.go (ListEmails) and
-- internal/service/v1/mcp/executors/email.go for the executor/client side.
UPDATE mcp_tools
SET input_schema = '{
    "properties": {
        "limit":      { "type": "integer", "description": "Maximum number of emails to return (default 20)" },
        "after_uid":  { "type": "string", "description": "Only return emails with UID greater than this value, oldest-first. Pass \"0\" to start from the oldest email in the inbox. To page forward, pass the id of the last (newest) email from the previous page. Mutually exclusive with before_uid." },
        "before_uid": { "type": "string", "description": "Only return emails with UID less than this value, also oldest-first. To page further back, pass the id of the first (oldest) email from the previous page. Mutually exclusive with after_uid." }
    },
    "required": []
}'::jsonb
WHERE mcp_name = 'email' AND name = 'list_emails';

-- +goose Down
UPDATE mcp_tools
SET input_schema = '{
    "properties": {
        "limit": { "type": "integer", "description": "Max number of emails to return (default 20)" }
    },
    "required": []
}'::jsonb
WHERE mcp_name = 'email' AND name = 'list_emails';
