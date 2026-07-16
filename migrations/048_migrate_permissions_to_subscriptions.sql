-- +goose Up

-- Backfill every user's current user_permissions feature booleans into subscriptions'
-- feature_overrides, as explicit true/false keys. This must write both true AND false values
-- (not just deviations from the 'basic' plan default) — 'basic' defaults these four features to
-- true, so a user who currently has e.g. has_notes = FALSE needs an explicit "notes": false
-- override, or they'd silently gain access to notes the moment subscriptions_enabled flips on.
UPDATE subscriptions s
SET feature_overrides = s.feature_overrides || jsonb_build_object(
    'emails', up.has_emails,
    'task_trackers', up.has_task_trackers,
    'notes', up.has_notes,
    'spreadsheets', up.has_spreadsheets
)
FROM user_permissions up
WHERE up.user_id = s.user_id;

-- is_administrator stays on user_permissions (unrelated to subscriptions — an internal-access
-- flag, not a purchasable feature); only the four feature columns are retired.
ALTER TABLE user_permissions
    DROP COLUMN IF EXISTS has_emails,
    DROP COLUMN IF EXISTS has_task_trackers,
    DROP COLUMN IF EXISTS has_notes,
    DROP COLUMN IF EXISTS has_spreadsheets;

-- +goose Down

-- Data loss on rollback is expected here: the dropped columns' values now live in
-- subscriptions.feature_overrides, not reconstructed back onto user_permissions.
ALTER TABLE user_permissions
    ADD COLUMN has_emails BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN has_task_trackers BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN has_notes BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN has_spreadsheets BOOLEAN NOT NULL DEFAULT FALSE;
