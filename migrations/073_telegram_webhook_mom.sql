-- +goose Up
-- Extends the telegram MoM (seeded in migration 065) for the inbound webhook relay
-- (internal/transport/telegram_webhook): send_message gains an optional reply_markup param so
-- permission_request events can render as an inline keyboard, and two new tools —
-- edit_message_text (coalescing assistant_text_delta edits, and updating a permission-request
-- message once a decision is made) and answer_callback_query (dismissing the button spinner
-- after a callback_query) — round out what the relay needs. set_webhook is seeded too, for the
-- day a public base URL config exists to call it with (see AddTelegramConnection's comment) —
-- unused by any code path today.
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

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'telegram',
    'edit_message_text',
    'Edit the text (and optionally the inline keyboard) of a previously sent message via the connected Telegram bot',
    '{
        "properties": {
            "chat_id":    { "type": "string", "description": "Target chat/user/channel id" },
            "message_id": { "type": "string", "description": "Id of the message to edit" },
            "text":       { "type": "string", "description": "New message text" },
            "reply_markup": {
                "type": "string",
                "description": "Optional JSON-serialized InlineKeyboardMarkup object; omit to leave the keyboard as-is, pass \"{}\" to clear it"
            }
        },
        "required": ["chat_id", "message_id", "text"]
    }'::jsonb,
    '{
        "properties": {
            "ok":     { "type": "boolean" },
            "result": { "type": "object", "properties": {
                "message_id": { "type": "integer" },
                "text":       { "type": "string" }
            } }
        }
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.telegram.org/bot${{secrets.bot_token}}/editMessageText",
            "headers": { "Content-Type": "application/json" },
            "body": {
                "chat_id": "${{params.chat_id}}",
                "message_id": "${{params.message_id}}",
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

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'telegram',
    'answer_callback_query',
    'Acknowledge an inline-keyboard button press via the connected Telegram bot, dismissing its loading spinner',
    '{
        "properties": {
            "callback_query_id": { "type": "string", "description": "Id of the callback query being answered" },
            "text": { "type": "string", "description": "Optional short notification text shown to the user" }
        },
        "required": ["callback_query_id"]
    }'::jsonb,
    '{
        "properties": {
            "ok": { "type": "boolean" }
        }
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.telegram.org/bot${{secrets.bot_token}}/answerCallbackQuery",
            "headers": { "Content-Type": "application/json" },
            "body": {
                "callback_query_id": "${{params.callback_query_id}}",
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

INSERT INTO mcp_tools (mcp_name, name, description, input_schema, output_schema, action)
VALUES (
    'telegram',
    'set_webhook',
    'Register this Artel connection''s inbound webhook URL with Telegram for the connected bot',
    '{
        "properties": {
            "url": { "type": "string", "description": "Public callback URL, e.g. https://<host>/webhooks/telegram/<connection_id>" },
            "secret_token": { "type": "string", "description": "Secret Telegram echoes back on X-Telegram-Bot-Api-Secret-Token for every delivery" }
        },
        "required": ["url", "secret_token"]
    }'::jsonb,
    '{
        "properties": {
            "ok":          { "type": "boolean" },
            "description": { "type": "string" }
        }
    }'::jsonb,
    '{
        "http": {
            "method": "POST",
            "url": "https://api.telegram.org/bot${{secrets.bot_token}}/setWebhook",
            "headers": { "Content-Type": "application/json" },
            "body": {
                "url": "${{params.url}}",
                "secret_token": "${{params.secret_token}}"
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
DELETE FROM mcp_tools WHERE mcp_name = 'telegram' AND name = 'set_webhook';
DELETE FROM mcp_tools WHERE mcp_name = 'telegram' AND name = 'answer_callback_query';
DELETE FROM mcp_tools WHERE mcp_name = 'telegram' AND name = 'edit_message_text';
-- send_message reverts to its pre-073 shape (no reply_markup) by re-running 065's insert.
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
