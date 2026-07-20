-- +goose Up
-- A published, immutable snapshot of a tract — browsable/copyable by any user. source_tract_id
-- is nullable and ON DELETE SET NULL deliberately: a published template must survive the
-- original tract being deleted (it's a frozen copy, not a live view), see domain.TractTemplate.
CREATE TABLE tract_templates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_tract_id UUID REFERENCES tracts(id) ON DELETE SET NULL,
    owner_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    definition      JSONB NOT NULL, -- {steps: [...]} snapshot; connection_uuid stripped to nil at publish time, see domain.TractTemplate
    category        TEXT NOT NULL DEFAULT '',
    install_count   INT NOT NULL DEFAULT 0,
    published_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON tract_templates (category);
CREATE INDEX ON tract_templates (owner_id);

-- +goose Down
DROP TABLE tract_templates;
