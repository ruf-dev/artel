-- +goose Up
-- Seeds a new trello MoM: no migration seeded a `mcps` row for trello before this one (only
-- pkg/mom_examples/README.md had a worked example, get_board_members, which is NOT seeded here —
-- only add_comment, per the tract reference flow's task-tracker leg). Credential field naming
-- (api_key/api_token) matches the README's trello example and its "trello" provider convention.
--
-- mcps no longer has a `tools` column (dropped by migration 032) so this INSERT only sets
-- name/author/description; ON CONFLICT DO NOTHING so an existing trello row (seeded by another
-- session) is left alone.
INSERT INTO mcps (name, author, description)
VALUES ('trello', 'community', 'Trello project management integration')
ON CONFLICT (name) DO NOTHING;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'add_comment',
    'Add a comment to a Trello card',
    '{
        "properties": {
            "card_id": { "type": "string", "description": "Trello card id" },
            "text":    { "type": "string", "description": "Comment text" }
        },
        "required": ["card_id", "text"]
    }'::jsonb,
    '{
        "properties": {
            "id":   { "type": "string", "description": "Comment action id" },
            "data": { "type": "object", "properties": {
                "text": { "type": "string" },
                "card": { "type": "object", "properties": {
                    "id":   { "type": "string" },
                    "name": { "type": "string" }
                } }
            } }
        }
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.trello.com/1/cards/${{params.card_id}}/actions/comments",
            "headers": { "Accept": "application/json" },
            "query": {
                "text":  "${{params.text}}",
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
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'add_comment';
DELETE FROM mcps WHERE name = 'trello';
