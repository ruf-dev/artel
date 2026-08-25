-- +goose Up

-- role is a plain TEXT column ('user' | 'assistant' | 'tool'), not a Postgres ENUM: this
-- schema already mixes both conventions (workbench_status/workbench_auth_mode are ENUMs,
-- while e.g. vault_postgres_databases.status and workbenches columns elsewhere are plain
-- TEXT) and there's no cross-table reuse or admin-facing filter that would benefit from an
-- ENUM's type safety here — the value only ever round-trips through domain.SimpleChatRole*
-- constants in Go.
CREATE TABLE simple_chat_messages (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id       UUID        NOT NULL REFERENCES simple_chats(id) ON DELETE CASCADE,
    role          TEXT        NOT NULL,
    content       TEXT        NOT NULL DEFAULT '',
    tool_call_id  TEXT,
    tool_name     TEXT,
    tool_input    JSONB,
    is_error      BOOLEAN     NOT NULL DEFAULT FALSE,
    model         TEXT,
    seq           BIGINT      NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX simple_chat_messages_chat_id_seq_idx ON simple_chat_messages(chat_id, seq);

-- +goose Down

DROP TABLE IF EXISTS simple_chat_messages;
