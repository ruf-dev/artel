-- +goose Up
CREATE TABLE tracts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    enabled     BOOL NOT NULL DEFAULT TRUE,
    definition  JSONB NOT NULL DEFAULT '{}', -- {steps: [...]} step tree, see domain.TractDefinition
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON tracts (user_id);

-- Standalone, reusable trigger — not owned by a single tract. trigger_uuid is the rotatable
-- webhook routing id embedded in the public webhook URL (separate from id, the stable
-- owner-facing primary key), so rotating it invalidates old webhook URLs without disturbing
-- id-keyed references (tract links, etc).
CREATE TABLE triggers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trigger_uuid   UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name           TEXT NOT NULL,
    kind           TEXT NOT NULL,                    -- 'webhook' | 'manual'
    source         TEXT NOT NULL DEFAULT 'generic',  -- preset key: 'gitlab_push' | 'generic'
    config         JSONB NOT NULL DEFAULT '{}',
    payload_schema JSONB NOT NULL DEFAULT '{}',      -- schema node ({properties, required} shape, same as mcp_tools.input_schema)
    secret_hash    BYTEA,                            -- sha256(webhook token); NULL for manual
    enabled        BOOL NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON triggers (user_id);

-- Fans a trigger out to N tracts (and a tract can be wired to N triggers), each link carrying
-- its own filter conditions gating whether a given delivery starts a run for that tract.
CREATE TABLE tract_trigger_links (
    tract_id   UUID NOT NULL REFERENCES tracts(id) ON DELETE CASCADE,
    trigger_id UUID NOT NULL REFERENCES triggers(id) ON DELETE CASCADE,
    filters    JSONB NOT NULL DEFAULT '[]',          -- [{left, op, right}], AND semantics
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tract_id, trigger_id)
);
CREATE INDEX ON tract_trigger_links (trigger_id);

CREATE TYPE tract_run_status AS ENUM ('running', 'done', 'failed');

CREATE TABLE tract_runs (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tract_id       UUID NOT NULL REFERENCES tracts(id) ON DELETE CASCADE,
    trigger_id     UUID REFERENCES triggers(id) ON DELETE SET NULL,
    started_by     TEXT NOT NULL,                  -- 'webhook' | 'manual' | 'mcp'
    status         tract_run_status NOT NULL DEFAULT 'running',
    trigger_output JSONB NOT NULL DEFAULT '{}',    -- normalized payload / manual form values
    error          TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON tract_runs (tract_id, created_at DESC);

CREATE TYPE tract_run_step_status AS ENUM ('running', 'done', 'failed');

-- One row per EXECUTED step: the step "pattern" (identified by its JSONB tree node id,
-- step_id) instantiated with run data (audit + state store).
CREATE TABLE tract_run_steps (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id      UUID NOT NULL REFERENCES tract_runs(id) ON DELETE CASCADE,
    step_id     TEXT NOT NULL,                     -- TractStep.Id — steps are JSONB tree nodes, not rows
    step_name   TEXT NOT NULL,                    -- denormalized: survives definition edits, keys template state
    step_type   TEXT NOT NULL,
    input       JSONB,                            -- rendered params (action) / rendered conditions (condition)
    output      JSONB,                            -- parsed tool result; condition: {"result": bool}
    status      tract_run_step_status NOT NULL DEFAULT 'running',
    error       TEXT,
    started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ
);
CREATE INDEX ON tract_run_steps (run_id, started_at);

-- +goose Down
DROP TABLE tract_run_steps;
DROP TYPE tract_run_step_status;
DROP TABLE tract_runs;
DROP TYPE tract_run_status;
DROP TABLE tract_trigger_links;
DROP TABLE triggers;
DROP TABLE tracts;
