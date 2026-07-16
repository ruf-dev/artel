-- +goose Up

-- subscription_plans is a small, fixed catalog of named tiers. Keyed by a human-readable
-- plan_key (not a UUID) since rows are seeded/edited by migration, not created at runtime.
CREATE TABLE subscription_plans
(
    plan_key          TEXT PRIMARY KEY,
    couch_quota_bytes BIGINT      NOT NULL,
    s3_quota_bytes    BIGINT      NOT NULL,
    features          JSONB       NOT NULL DEFAULT '{}'::jsonb,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO subscription_plans (plan_key, couch_quota_bytes, s3_quota_bytes, features)
VALUES (
    'basic',
    5 * 1024 * 1024,
    20 * 1024 * 1024,
    '{
        "emails": true,
        "task_trackers": true,
        "notes": true,
        "spreadsheets": true,
        "connectors": true,
        "tract": false
    }'::jsonb
);

-- subscriptions gains a plan reference plus per-user overrides on the same row: feature
-- overrides are a sparse JSONB map (only keys present override the plan's features), and the
-- two quota overrides are nullable (NULL means "use the plan's default").
ALTER TABLE subscriptions
    ADD COLUMN plan_key TEXT NOT NULL DEFAULT 'basic' REFERENCES subscription_plans (plan_key),
    ADD COLUMN feature_overrides JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN couch_quota_override_bytes BIGINT,
    ADD COLUMN s3_quota_override_bytes BIGINT;

-- +goose Down

ALTER TABLE subscriptions
    DROP COLUMN IF EXISTS s3_quota_override_bytes,
    DROP COLUMN IF EXISTS couch_quota_override_bytes,
    DROP COLUMN IF EXISTS feature_overrides,
    DROP COLUMN IF EXISTS plan_key;

DROP TABLE IF EXISTS subscription_plans;
