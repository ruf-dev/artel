-- +goose Up
-- Rounds out the trello MoM (seeded by migration 035, extended by 044) with the tools needed to
-- actually work a board end to end: list_boards/list_lists/list_cards for browsing, move_card for
-- dragging a card between columns, list_labels/add_label_to_card for labelling, create_card, and
-- archive_card. All use Trello's query-param auth convention (key/token as query, no JSON body)
-- to match the existing add_comment/get_current_user tools rather than gitlab's JSON-body style —
-- Trello's REST API accepts query params uniformly across GET/POST/PUT.
--
-- Terminology: Trello's "board" maps to what was asked for as "table", and "list" (a board's
-- columns) maps to "column".
--
-- UPSERT only (ON CONFLICT (mcp_name, name) DO UPDATE) — never touches existing trello tool rows.

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'list_boards',
    'List Trello boards (tables) the current user belongs to',
    '{
        "properties": {
            "filter": { "type": "string", "description": "Which boards to return (default open)", "enum": ["open", "closed", "all"] }
        },
        "required": []
    }'::jsonb,
    '{
        "properties": {
            "id":     { "type": "string" },
            "name":   { "type": "string" },
            "desc":   { "type": "string" },
            "closed": { "type": "boolean" },
            "url":    { "type": "string" }
        },
        "required": ["id", "name"]
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.trello.com/1/members/me/boards",
            "headers": { "Accept": "application/json" },
            "query": {
                "filter": "${{params.filter}}",
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
    'list_lists',
    'List the lists (columns) on a Trello board',
    '{
        "properties": {
            "board_id": { "type": "string", "description": "Trello board id" },
            "filter":   { "type": "string", "description": "Which lists to return (default open)", "enum": ["open", "closed", "all"] }
        },
        "required": ["board_id"]
    }'::jsonb,
    '{
        "properties": {
            "id":      { "type": "string" },
            "name":    { "type": "string" },
            "closed":  { "type": "boolean" },
            "idBoard": { "type": "string" }
        },
        "required": ["id", "name"]
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.trello.com/1/boards/${{params.board_id}}/lists",
            "headers": { "Accept": "application/json" },
            "query": {
                "filter": "${{params.filter}}",
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
    'list_cards',
    'List the cards on a Trello board, across all its lists (columns)',
    '{
        "properties": {
            "board_id": { "type": "string", "description": "Trello board id" }
        },
        "required": ["board_id"]
    }'::jsonb,
    '{
        "properties": {
            "id":      { "type": "string" },
            "name":    { "type": "string" },
            "desc":    { "type": "string" },
            "idList":  { "type": "string", "description": "Id of the list (column) this card is in" },
            "idBoard": { "type": "string" },
            "closed":  { "type": "boolean" },
            "url":     { "type": "string" },
            "labels":  { "type": "array", "items": { "type": "string" } }
        },
        "required": ["id", "name", "idList"]
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.trello.com/1/boards/${{params.board_id}}/cards",
            "headers": { "Accept": "application/json" },
            "query": {
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
    'move_card',
    'Move a card to a different list (column), optionally another board',
    '{
        "properties": {
            "card_id": { "type": "string", "description": "Trello card id to move" },
            "list_id": { "type": "string", "description": "Destination list (column) id" }
        },
        "required": ["card_id", "list_id"]
    }'::jsonb,
    '{
        "properties": {
            "id":     { "type": "string" },
            "name":   { "type": "string" },
            "idList": { "type": "string" }
        },
        "required": ["id", "idList"]
    }'::jsonb,
    '{
        "http": {
            "method": "PUT",
            "url": "https://api.trello.com/1/cards/${{params.card_id}}",
            "headers": { "Accept": "application/json" },
            "query": {
                "idList": "${{params.list_id}}",
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
    'list_labels',
    'List the labels defined on a Trello board, with their ids (needed to assign labels to a card)',
    '{
        "properties": {
            "board_id": { "type": "string", "description": "Trello board id" }
        },
        "required": ["board_id"]
    }'::jsonb,
    '{
        "properties": {
            "id":    { "type": "string" },
            "name":  { "type": "string" },
            "color": { "type": "string" }
        },
        "required": ["id", "color"]
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.trello.com/1/boards/${{params.board_id}}/labels",
            "headers": { "Accept": "application/json" },
            "query": {
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

-- add_label_to_card's response is Trello's raw JSON array of the card's label ids (a bare array
-- of strings), which domain.ToolSchema can't describe (same reasoning as email/list_email_folders
-- in migration 030) so output_schema is left undeclared ('{}').
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'add_label_to_card',
    'Assign an existing label to a Trello card (get the label id from list_labels)',
    '{
        "properties": {
            "card_id":  { "type": "string", "description": "Trello card id" },
            "label_id": { "type": "string", "description": "Label id from list_labels" }
        },
        "required": ["card_id", "label_id"]
    }'::jsonb,
    '{}'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.trello.com/1/cards/${{params.card_id}}/idLabels",
            "headers": { "Accept": "application/json" },
            "query": {
                "value": "${{params.label_id}}",
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
    'create_card',
    'Create a new card in a Trello list (column)',
    '{
        "properties": {
            "list_id": { "type": "string", "description": "Id of the list (column) to create the card in" },
            "name":    { "type": "string", "description": "Card title" },
            "desc":    { "type": "string", "description": "Card description (optional)" }
        },
        "required": ["list_id", "name"]
    }'::jsonb,
    '{
        "properties": {
            "id":     { "type": "string" },
            "name":   { "type": "string" },
            "desc":   { "type": "string" },
            "idList": { "type": "string" },
            "url":    { "type": "string" }
        },
        "required": ["id", "name", "idList"]
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.trello.com/1/cards",
            "headers": { "Accept": "application/json" },
            "query": {
                "idList": "${{params.list_id}}",
                "name":   "${{params.name}}",
                "desc":   "${{params.desc}}",
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
    'archive_card',
    'Archive (close) a Trello card',
    '{
        "properties": {
            "card_id": { "type": "string", "description": "Trello card id to archive" }
        },
        "required": ["card_id"]
    }'::jsonb,
    '{
        "properties": {
            "id":     { "type": "string" },
            "name":   { "type": "string" },
            "closed": { "type": "boolean" }
        },
        "required": ["id", "closed"]
    }'::jsonb,
    '{
        "http": {
            "method": "PUT",
            "url": "https://api.trello.com/1/cards/${{params.card_id}}",
            "headers": { "Accept": "application/json" },
            "query": {
                "closed": "true",
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

-- +goose Down
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'archive_card';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'create_card';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'add_label_to_card';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'list_labels';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'move_card';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'list_cards';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'list_lists';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'list_boards';
