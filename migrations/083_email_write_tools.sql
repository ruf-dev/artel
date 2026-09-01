-- +goose NO TRANSACTION
-- +goose Up
-- Adds the agentic email write tools on top of the read-only IMAP/SMTP surface from migration 027:
--   - move_email          — relocate a message between IMAP folders (imap operation move_message)
--   - create_email_folder — create a new IMAP mailbox            (imap operation create_folder)
--   - set_email_read      — add/remove the \Seen flag on a message (imap operation set_flags)
-- and threads a new optional `folder` param through list_emails and read_email so a caller can
-- page/read a mailbox other than INBOX. See internal/clients/imap/client.go (MoveMessage,
-- CreateFolder, SetSeen, ListEmailsOptions.Folder, ReadEmail's folder arg) and
-- internal/service/v1/mcp/executors/email.go for the client/executor side.
--
-- The file runs WITHOUT a transaction (see the NO TRANSACTION directive on line 1): the three
-- ALTER TYPE ... ADD VALUE statements extend the mcp_tool_name enum from migration 080, and
-- ADD VALUE cannot run inside a transaction block (its new label also isn't usable in the same
-- transaction). IF NOT EXISTS keeps the migration re-runnable. UPSERT / targeted UPDATE only —
-- existing email tool rows other than list_emails/read_email are untouched.

ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'email.move_email';
ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'email.create_email_folder';
ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'email.set_email_read';

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'email',
    'move_email',
    'Move an email to another IMAP folder',
    '{
        "properties": {
            "id":            { "type": "string", "description": "Email UID from list_emails" },
            "dest_folder":   { "type": "string", "description": "Target folder name, e.g. Archive" },
            "source_folder": { "type": "string", "description": "Folder the message is currently in (default INBOX)" }
        },
        "required": ["id", "dest_folder"]
    }'::jsonb,
    '{}'::jsonb,
    '{ "imap": { "operation": "move_message" } }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'email',
    'create_email_folder',
    'Create a new IMAP folder in the connected email account',
    '{
        "properties": {
            "name": { "type": "string", "description": "Folder name to create" }
        },
        "required": ["name"]
    }'::jsonb,
    '{}'::jsonb,
    '{ "imap": { "operation": "create_folder" } }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'email',
    'set_email_read',
    'Mark an email as read or unread',
    '{
        "properties": {
            "id":            { "type": "string", "description": "Email UID from list_emails" },
            "seen":          { "type": "boolean", "description": "true = mark read, false = mark unread" },
            "source_folder": { "type": "string", "description": "Folder the message is in (default INBOX)" }
        },
        "required": ["id", "seen"]
    }'::jsonb,
    '{}'::jsonb,
    '{ "imap": { "operation": "set_flags" } }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

-- Add the optional `folder` param to list_emails (preserves limit/after_uid/before_uid from
-- migration 042).
UPDATE mcp_tools
SET input_schema = '{
    "properties": {
        "limit":      { "type": "integer", "description": "Maximum number of emails to return (default 20)" },
        "after_uid":  { "type": "string", "description": "Only return emails with UID greater than this value, oldest-first. Pass \"0\" to start from the oldest email in the inbox. To page forward, pass the id of the last (newest) email from the previous page. Mutually exclusive with before_uid." },
        "before_uid": { "type": "string", "description": "Only return emails with UID less than this value, also oldest-first. To page further back, pass the id of the first (oldest) email from the previous page. Mutually exclusive with after_uid." },
        "folder":     { "type": "string", "description": "IMAP folder to read from (default INBOX)" }
    },
    "required": []
}'::jsonb
WHERE mcp_name = 'email' AND name = 'list_emails';

-- Add the optional `folder` param to read_email (preserves id from migration 027).
UPDATE mcp_tools
SET input_schema = '{
    "properties": {
        "id":     { "type": "string", "description": "Email UID from list_emails" },
        "folder": { "type": "string", "description": "IMAP folder to read from (default INBOX)" }
    },
    "required": ["id"]
}'::jsonb
WHERE mcp_name = 'email' AND name = 'read_email';

-- +goose Down
-- Only the tool rows are removed and the two input_schemas reverted. The three 'email.*' labels
-- this migration adds to the mcp_tool_name enum are intentionally left in place: Postgres cannot
-- drop an enum value. This is the same knowing break of the profile's "every migration reverses
-- exactly" rule as migration 081 — see docs/profile-drift.md.
DELETE FROM mcp_tools WHERE mcp_name = 'email' AND name IN ('move_email', 'create_email_folder', 'set_email_read');

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

UPDATE mcp_tools
SET input_schema = '{
    "properties": {
        "id": { "type": "string", "description": "Email UID from list_emails" }
    },
    "required": ["id"]
}'::jsonb
WHERE mcp_name = 'email' AND name = 'read_email';
