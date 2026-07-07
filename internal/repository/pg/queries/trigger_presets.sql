-- name: ListTriggerPresets :many
SELECT key, category, label, description, provider, payload_schema, default_matchers, created_at
FROM trigger_presets
ORDER BY category, label;

-- name: GetTriggerPresetByKey :one
SELECT key, category, label, description, provider, payload_schema, default_matchers, created_at
FROM trigger_presets
WHERE key = $1;
