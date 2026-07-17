-- +goose Up
-- Adds the two Trello tools needed for a card-detail view: fetching a single card and listing
-- its comments. Rounds out list_boards/list_lists/list_cards (migration 046) so a caller can go
-- from "board" all the way down to "one card's full detail + comment thread" via MoM alone.
--
-- UPSERT only (ON CONFLICT (mcp_name, name) DO UPDATE) — never touches other existing trello
-- tool rows.

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'get_card',
    'Get the full details of a single Trello card by id or shortLink',
    '{
        "properties": {
            "id": { "type": "string", "description": "Trello card id or shortLink" }
        },
        "required": ["id"]
    }'::jsonb,
    '{
        "properties": {
            "id":        { "type": "string" },
            "name":      { "type": "string" },
            "desc":      { "type": "string" },
            "idList":    { "type": "string", "description": "Id of the list (column) this card is in" },
            "idBoard":   { "type": "string" },
            "closed":    { "type": "boolean" },
            "url":       { "type": "string" },
            "shortLink": { "type": "string" }
        },
        "required": ["id", "name"]
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.trello.com/1/cards/${{params.id}}",
            "headers": { "Accept": "application/json" },
            "query": {
                "fields": "id,name,desc,idList,idBoard,closed,url,shortLink",
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

-- list_card_comments' response is Trello's raw JSON array of action objects
-- ({id, data: {text}, date, memberCreator: {fullName}}), a bare array which domain.ToolSchema
-- can't describe (same reasoning as add_label_to_card in migration 046 and
-- email/list_email_folders in migration 030) so output_schema is left undeclared ('{}').
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'list_card_comments',
    'List the comments on a Trello card',
    '{
        "properties": {
            "id": { "type": "string", "description": "Trello card id or shortLink" }
        },
        "required": ["id"]
    }'::jsonb,
    '{}'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.trello.com/1/cards/${{params.id}}/actions",
            "headers": { "Accept": "application/json" },
            "query": {
                "filter": "commentCard",
                "fields": "data,date,memberCreator",
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
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'list_card_comments';
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'get_card';
