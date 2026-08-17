-- +goose Up

-- Skill caps join the existing storage-quota columns on subscription_plans/subscriptions:
-- a plan-level default (max_hot_plug_skills/max_total_skills) plus a nullable per-user
-- override (NULL means "use the plan's default"), same convention as
-- couch_quota_override_bytes/s3_quota_override_bytes added in 047_subscription_plans.sql.
ALTER TABLE subscription_plans
    ADD COLUMN max_hot_plug_skills INT NOT NULL DEFAULT 5,
    ADD COLUMN max_total_skills    INT NOT NULL DEFAULT 20;

ALTER TABLE subscriptions
    ADD COLUMN max_hot_plug_skills_override INT,
    ADD COLUMN max_total_skills_override    INT;

-- +goose Down

ALTER TABLE subscriptions
    DROP COLUMN IF EXISTS max_total_skills_override,
    DROP COLUMN IF EXISTS max_hot_plug_skills_override;

ALTER TABLE subscription_plans
    DROP COLUMN IF EXISTS max_total_skills,
    DROP COLUMN IF EXISTS max_hot_plug_skills;
