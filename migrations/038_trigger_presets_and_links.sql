-- +goose Up

-- Declarative catalog of webhook trigger presets — replaces the hardcoded Go slice previously
-- returned by tract.ListTriggerSources(). `provider` links a preset to an external_connections
-- row (same user, matching provider); NULL means "standalone" (generic, user-declared schema,
-- own trigger_uuid/secret_hash webhook URL).
CREATE TABLE trigger_presets (
    key              TEXT PRIMARY KEY,
    category         TEXT NOT NULL,               -- 'gitlab' | 'generic' (UI grouping)
    label            TEXT NOT NULL,                -- 'Branch push' (UI display name)
    description      TEXT NOT NULL,
    provider         TEXT,                          -- external_connections.provider, NULL = standalone
    payload_schema   JSONB NOT NULL DEFAULT '{}',  -- same {properties,required} shape as triggers.payload_schema
    default_matchers JSONB NOT NULL DEFAULT '{}',  -- domain.TriggerMatchers, copied onto the trigger at creation
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO trigger_presets (key, category, label, description, provider, payload_schema, default_matchers) VALUES
('generic', 'generic', 'Generic webhook',
 'Payload shape is declared per-trigger (payload_schema), passed through unchanged',
 NULL, '{}', '{}'),
('gitlab_push', 'gitlab', 'Branch push',
 'Fires when commits are pushed to a branch',
 'gitlab',
 '{"properties": {
      "ref": {"type": "string", "description": "full git ref, e.g. refs/heads/feature-x"},
      "branch": {"type": "string", "description": "ref with refs/heads/ trimmed (added by the normalizer)"},
      "user_name": {"type": "string"},
      "project": {"type": "object", "properties": {
          "id": {"type": "integer"}, "name": {"type": "string"}, "path": {"type": "string"}
      }},
      "commits": {"type": "array", "items": {"type": "object", "properties": {
          "id": {"type": "string", "description": "commit sha"},
          "message": {"type": "string"},
          "url": {"type": "string"},
          "author": {"type": "object", "properties": {
              "name": {"type": "string"}, "email": {"type": "string"}
          }}
      }}},
      "total_commits_count": {"type": "integer"}
  }, "required": ["ref", "branch"]}'::jsonb,
 '{"check_headers": [{"header": "X-Gitlab-Event", "equals": "Push Hook"}]}'::jsonb
);

-- Per-trigger inbound matcher config (domain.TriggerMatchers), seeded from the chosen preset's
-- default_matchers at creation time. Evaluated against an inbound delivery to decide whether THIS
-- trigger (among possibly several sharing one provider link) should fire.
ALTER TABLE triggers ADD COLUMN matchers JSONB NOT NULL DEFAULT '{}';

-- Fans a trigger out to a shared provider connection (the "link") instead of it minting its own
-- trigger_uuid/secret_hash webhook URL. Composite PK keeps this many-to-many capable even though
-- the v1 product constraint ("one GitLab connection per user") comes from external_connections'
-- own UNIQUE(user_id, provider), not from this table.
CREATE TABLE trigger_provider_links (
    trigger_id             UUID NOT NULL REFERENCES triggers(id) ON DELETE CASCADE,
    external_connection_id UUID NOT NULL REFERENCES external_connections(id) ON DELETE CASCADE,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (trigger_id, external_connection_id)
);
CREATE INDEX ON trigger_provider_links (external_connection_id);

-- +goose Down
DROP TABLE trigger_provider_links;
ALTER TABLE triggers DROP COLUMN matchers;
DROP TABLE trigger_presets;
