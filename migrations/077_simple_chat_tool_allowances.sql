-- +goose Up

CREATE TABLE simple_chat_tool_allowances (
    chat_id    UUID        NOT NULL REFERENCES simple_chats(id) ON DELETE CASCADE,
    tool_name  TEXT        NOT NULL,
    decision   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, tool_name)
);

-- +goose Down

DROP TABLE IF EXISTS simple_chat_tool_allowances;
