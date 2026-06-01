-- +goose Up
CREATE TABLE task_trackers (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type          TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    api_key_enc   BYTEA       NOT NULL,
    api_token_enc BYTEA       NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON task_trackers(user_id);

-- +goose Down
DROP TABLE task_trackers;
