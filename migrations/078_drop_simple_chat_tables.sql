-- +goose Up

-- Simple Chat storage moved to CouchDB (one JSONL doc per chat under each vault's
-- .chat_history/<user_id>/<chat_id>.jsonl — see internal/service/v1/simplechat). The feature was
-- never in production, so these tables are dropped outright rather than migrated in place.
DROP TABLE IF EXISTS simple_chat_tool_allowances;
DROP TABLE IF EXISTS simple_chat_messages;
DROP TABLE IF EXISTS simple_chats;

-- +goose Down

CREATE TABLE simple_chats (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id         UUID        NOT NULL REFERENCES vaults(id) ON DELETE CASCADE,
    user_id          UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title            TEXT,
    model            TEXT        NOT NULL,
    vault_access     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX simple_chats_vault_id_user_id_idx ON simple_chats(vault_id, user_id);

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

CREATE TABLE simple_chat_tool_allowances (
    chat_id    UUID        NOT NULL REFERENCES simple_chats(id) ON DELETE CASCADE,
    tool_name  TEXT        NOT NULL,
    decision   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chat_id, tool_name)
);
