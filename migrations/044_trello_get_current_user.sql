-- +goose Up
-- Adds a dead-simple "who am I" tool to the existing trello MoM (seeded by migration 035),
-- mirroring gitlab's get_current_user tool (migration 031) — good for a first sanity check
-- of a freshly connected trello credential.
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'trello',
    'get_current_user',
    'Get the current authenticated Trello member (you)',
    '{
        "properties": {},
        "required": []
    }'::jsonb,
    '{
        "properties": {
            "id":       { "type": "string", "description": "Trello member id" },
            "username": { "type": "string", "description": "Trello username" },
            "fullName": { "type": "string", "description": "Display name" }
        }
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.trello.com/1/members/me",
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

-- +goose Down
DELETE FROM mcp_tools WHERE mcp_name = 'trello' AND name = 'get_current_user';
