-- +goose Up
-- Seeds a new telegram MoM: get_me validates a caller-supplied bot token before it is persisted
-- (see externalconnections.Service.validateTelegramToken), send_message is the first real tool a
-- tract/toolbox call can invoke against a saved telegram connection. Provider key ("telegram")
-- matches domain.ProviderTelegram, which dispatchHttp checks the saved external_connections row
-- against.
INSERT INTO mcps (name, author, description)
VALUES ('telegram', 'Artel', 'Telegram bot integration')
ON CONFLICT (name) DO NOTHING;

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'telegram',
    'get_me',
    'Get info about the connected Telegram bot (used to validate the bot token)',
    '{
        "properties": {},
        "required": []
    }'::jsonb,
    '{
        "properties": {
            "ok":     { "type": "boolean" },
            "result": { "type": "object", "properties": {
                "id":         { "type": "integer" },
                "username":   { "type": "string" },
                "first_name": { "type": "string" }
            } }
        }
    }'::jsonb,
    '{
        "http": {
            "method": "GET",
            "url": "https://api.telegram.org/bot${{secrets.bot_token}}/getMe",
            "credentials": "telegram"
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
    'telegram',
    'send_message',
    'Send a text message via the connected Telegram bot',
    '{
        "properties": {
            "chat_id": { "type": "string", "description": "Target chat/user/channel id" },
            "text":    { "type": "string", "description": "Message text" }
        },
        "required": ["chat_id", "text"]
    }'::jsonb,
    '{
        "properties": {
            "ok":     { "type": "boolean" },
            "result": { "type": "object", "properties": {
                "message_id": { "type": "integer" },
                "chat":       { "type": "object", "properties": {
                    "id": { "type": "integer" }
                } },
                "text":       { "type": "string" }
            } }
        }
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.telegram.org/bot${{secrets.bot_token}}/sendMessage",
            "headers": { "Content-Type": "application/json" },
            "body": {
                "chat_id": "${{params.chat_id}}",
                "text": "${{params.text}}"
            },
            "credentials": "telegram"
        }
    }'::jsonb
)
ON CONFLICT (mcp_name, name) DO UPDATE
    SET description   = EXCLUDED.description,
        input_schema  = EXCLUDED.input_schema,
        output_schema = EXCLUDED.output_schema,
        action        = EXCLUDED.action;

-- +goose Down
DELETE FROM mcp_tools WHERE mcp_name = 'telegram' AND name = 'send_message';
DELETE FROM mcp_tools WHERE mcp_name = 'telegram' AND name = 'get_me';
DELETE FROM mcps WHERE name = 'telegram';
