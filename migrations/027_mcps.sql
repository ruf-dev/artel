-- +goose Up
CREATE TABLE mcps (
    name        TEXT        PRIMARY KEY,
    author      TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    tools       JSONB       NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO mcps (name, author, description, tools) VALUES (
    'email',
    'Artel',
    'Email handler. Can send, read, list emails in connected account.',
    '[
        {
            "api_description": {
                "name": "list_email_folders",
                "description": "List IMAP folders in the connected email account",
                "properties": {},
                "required": []
            },
            "action": { "imap": { "operation": "list_folders" } }
        },
        {
            "api_description": {
                "name": "list_emails",
                "description": "List recent emails from the inbox",
                "properties": {
                    "limit": { "type": "integer", "description": "Max number of emails to return (default 20)" }
                },
                "required": []
            },
            "action": { "imap": { "operation": "list_messages" } }
        },
        {
            "api_description": {
                "name": "read_email",
                "description": "Read the full content of an email by its UID",
                "properties": {
                    "id": { "type": "string", "description": "Email UID from list_emails" }
                },
                "required": ["id"]
            },
            "action": { "imap": { "operation": "fetch_message" } }
        },
        {
            "api_description": {
                "name": "send_email",
                "description": "Send an email from the connected account",
                "properties": {
                    "to":      { "type": "string", "description": "Recipient email address" },
                    "subject": { "type": "string", "description": "Email subject" },
                    "body":    { "type": "string", "description": "Email body (plain text)" }
                },
                "required": ["to", "subject", "body"]
            },
            "action": { "smtp": { "operation": "send" } }
        }
    ]'
) ON CONFLICT (name) DO NOTHING;

-- +goose Down
DROP TABLE mcps CASCADE;
