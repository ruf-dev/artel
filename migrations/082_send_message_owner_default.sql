-- +goose Up
-- Makes telegram/send_message's chat_id optional so an agent can notify the key owner without
-- naming a recipient: chat_id drops out of "required" (only "text" stays), and the request body
-- falls back to the connected bot's captured owner chat id when the param is omitted
-- (${{params.chat_id | secrets.chat_id}}). Everything else is re-upserted verbatim from
-- migration 073. Callers that still pass chat_id (the telegram_webhook relay) are unaffected —
-- an explicit param always wins the fallback.
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'telegram',
    'send_message',
    'Send a text message via the connected Telegram bot',
    '{
        "properties": {
            "chat_id": { "type": "string", "description": "Target chat/user/channel id; defaults to the connected bot''s owner when omitted" },
            "text":    { "type": "string", "description": "Message text" },
            "reply_markup": {
                "type": "string",
                "description": "Optional JSON-serialized InlineKeyboardMarkup object"
            }
        },
        "required": ["text"]
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
                "chat_id": "${{params.chat_id | secrets.chat_id}}",
                "text": "${{params.text}}",
                "reply_markup": "${{params.reply_markup}}"
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
-- Restore send_message to its migration 073 shape: chat_id required, no secrets fallback.
INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'telegram',
    'send_message',
    'Send a text message via the connected Telegram bot',
    '{
        "properties": {
            "chat_id": { "type": "string", "description": "Target chat/user/channel id" },
            "text":    { "type": "string", "description": "Message text" },
            "reply_markup": {
                "type": "string",
                "description": "Optional JSON-serialized InlineKeyboardMarkup object"
            }
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
                "text": "${{params.text}}",
                "reply_markup": "${{params.reply_markup}}"
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
