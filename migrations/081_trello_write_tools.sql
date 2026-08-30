-- +goose NO TRANSACTION
-- +goose Up
-- Adds the Trello write tools that round out the board-editing surface on top of migration 046
-- (create_card/move_card/archive_card/...) and 049 (get_card/list_card_comments): editing a card's
-- title/description, creating and renaming/archiving lists (columns), creating a label, and
-- attaching a link to a card. All use Trello's query-param auth convention (key/token as query, no
-- JSON body) to match the existing trello tools, and plain hardcoded https://api.trello.com/1/...
-- literal URLs.
--
-- UPSERT only (ON CONFLICT (mcp_name, name) DO UPDATE) — never touches existing trello tool rows.
--
-- The file runs WITHOUT a transaction (see the NO TRANSACTION directive on line 1): the five
-- ALTER TYPE ... ADD VALUE statements below extend the mcp_tool_name enum from migration 080, and
-- ADD VALUE cannot run inside a transaction block (and its new label isn't usable in the same
-- transaction even where it can). IF NOT EXISTS keeps the migration re-runnable.

ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'trello.update_card';
ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'trello.create_list';
ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'trello.update_list';
ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'trello.create_label';
ALTER TYPE mcp_tool_name ADD VALUE IF NOT EXISTS 'trello.add_attachment';

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'update_card',
    'Update a Trello card''s title and/or description. Provide at least one of name, desc.',
    '{
        "properties": {
            "card_id": { "type": "string", "description": "Trello card id or shortLink" },
            "name":    { "type": "string", "description": "New card title (optional)" },
            "desc":    { "type": "string", "description": "New card description (optional)" }
        },
        "required": ["card_id"]
    }'::jsonb,
    '{
        "properties": {
            "id":     { "type": "string" },
            "name":   { "type": "string" },
            "desc":   { "type": "string" },
            "idList": { "type": "string" },
            "closed": { "type": "boolean" },
            "url":    { "type": "string" }
        },
        "required": ["id", "name"]
    }'::jsonb,
    '{
        "http": {
            "method": "PUT",
            "url": "https://api.trello.com/1/cards/${{params.card_id}}",
            "headers": { "Accept": "application/json" },
            "query": {
                "name":  "${{params.name}}",
                "desc":  "${{params.desc}}",
                "key":   "__secrets.api_key",
                "token": "__secrets.api_token"
            },
            "credentials": "trello"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'create_list',
    'Create a new list (column) on a Trello board',
    '{
        "properties": {
            "board_id": { "type": "string", "description": "Trello board id the list is created on" },
            "name":     { "type": "string", "description": "List (column) name" },
            "pos":      { "type": "string", "description": "Position: \"top\", \"bottom\", or a positive number (default bottom)" }
        },
        "required": ["board_id", "name"]
    }'::jsonb,
    '{
        "properties": {
            "id":      { "type": "string" },
            "name":    { "type": "string" },
            "closed":  { "type": "boolean" },
            "idBoard": { "type": "string" }
        },
        "required": ["id", "name", "idBoard"]
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.trello.com/1/lists",
            "headers": { "Accept": "application/json" },
            "query": {
                "name":    "${{params.name}}",
                "idBoard": "${{params.board_id}}",
                "pos":     "${{params.pos}}",
                "key":     "__secrets.api_key",
                "token":   "__secrets.api_token"
            },
            "credentials": "trello"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'update_list',
    'Rename a Trello list (column) and/or archive it (closed=true) or restore it (closed=false)',
    '{
        "properties": {
            "list_id": { "type": "string", "description": "Trello list (column) id" },
            "name":    { "type": "string", "description": "New list name (optional)" },
            "closed":  { "type": "boolean", "description": "true to archive the list, false to restore it (optional)" }
        },
        "required": ["list_id"]
    }'::jsonb,
    '{
        "properties": {
            "id":      { "type": "string" },
            "name":    { "type": "string" },
            "closed":  { "type": "boolean" },
            "idBoard": { "type": "string" }
        },
        "required": ["id", "name", "closed"]
    }'::jsonb,
    '{
        "http": {
            "method": "PUT",
            "url": "https://api.trello.com/1/lists/${{params.list_id}}",
            "headers": { "Accept": "application/json" },
            "query": {
                "name":   "${{params.name}}",
                "closed": "${{params.closed}}",
                "key":    "__secrets.api_key",
                "token":  "__secrets.api_token"
            },
            "credentials": "trello"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'create_label',
    'Create a new label on a Trello board',
    '{
        "properties": {
            "board_id": { "type": "string", "description": "Trello board id the label is created on" },
            "name":     { "type": "string", "description": "Label name" },
            "color":    { "type": "string", "description": "Label color", "enum": ["yellow", "purple", "blue", "red", "green", "orange", "black", "sky", "pink", "lime"] }
        },
        "required": ["board_id", "name", "color"]
    }'::jsonb,
    '{
        "properties": {
            "id":      { "type": "string" },
            "name":    { "type": "string" },
            "color":   { "type": "string" },
            "idBoard": { "type": "string" }
        },
        "required": ["id", "name", "color"]
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.trello.com/1/labels",
            "headers": { "Accept": "application/json" },
            "query": {
                "name":    "${{params.name}}",
                "color":   "${{params.color}}",
                "idBoard": "${{params.board_id}}",
                "key":     "__secrets.api_key",
                "token":   "__secrets.api_token"
            },
            "credentials": "trello"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'add_attachment',
    'Attach a link (URL) to a Trello card. File uploads are not supported — pass a publicly reachable http/https URL.',
    '{
        "properties": {
            "card_id": { "type": "string", "description": "Trello card id or shortLink" },
            "url":     { "type": "string", "description": "URL to attach (http/https)" },
            "name":    { "type": "string", "description": "Display name for the attachment (optional, defaults to the URL)" }
        },
        "required": ["card_id", "url"]
    }'::jsonb,
    '{
        "properties": {
            "id":       { "type": "string" },
            "name":     { "type": "string" },
            "url":      { "type": "string" },
            "date":     { "type": "string" },
            "mimeType": { "type": "string" },
            "bytes":    { "type": "number" }
        },
        "required": ["id", "url"]
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.trello.com/1/cards/${{params.card_id}}/attachments",
            "headers": { "Accept": "application/json" },
            "query": {
                "url":   "${{params.url}}",
                "name":  "${{params.name}}",
                "key":   "__secrets.api_key",
                "token": "__secrets.api_token"
            },
            "credentials": "trello"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

-- +goose Down
-- Only the tool rows are removed. The five 'trello.*' labels this migration adds to the mcp_tool_name
-- enum are intentionally left in place: Postgres cannot drop an enum value. This is the one spot
-- the profile's "every migration reverses exactly" rule is knowingly broken — see
-- docs/profile-drift.md.
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'add_attachment';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'create_label';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'update_list';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'create_list';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'update_card';
