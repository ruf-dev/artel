-- name: InsertTrigger :one
INSERT INTO triggers (user_id, name, kind, source, config, payload_schema, secret_hash, token_suffix, matchers, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, trigger_uuid, user_id, name, kind, source, config, payload_schema, secret_hash, enabled, created_at, matchers, token_suffix;

-- name: GetTrigger :one
SELECT id, trigger_uuid, user_id, name, kind, source, config, payload_schema, secret_hash, enabled, created_at, matchers, token_suffix
FROM triggers
WHERE id = $1;

-- name: GetTriggerByTriggerUuid :one
-- Webhook routing lookup: the fired webhook only knows the trigger's rotatable routing id.
SELECT id, trigger_uuid, user_id, name, kind, source, config, payload_schema, secret_hash, enabled, created_at, matchers, token_suffix
FROM triggers
WHERE trigger_uuid = $1;

-- name: ListTriggersByUser :many
SELECT id, trigger_uuid, user_id, name, kind, source, config, payload_schema, secret_hash, enabled, created_at, matchers, token_suffix
FROM triggers
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: SetTriggerEnabled :exec
UPDATE triggers
SET enabled = $2
WHERE id = $1;

-- name: DeleteTrigger :exec
DELETE FROM triggers
WHERE id = $1;

-- name: RotateTriggerSecret :one
-- Invalidates the trigger's current webhook URL/token by overwriting trigger_uuid,
-- secret_hash, and token_suffix in place, keyed by the trigger's stable primary key id.
UPDATE triggers
SET trigger_uuid = $2,
    secret_hash   = $3,
    token_suffix  = $4
WHERE id = $1
RETURNING id, trigger_uuid, user_id, name, kind, source, config, payload_schema, secret_hash, enabled, created_at, matchers, token_suffix;

-- name: LinkTriggerToTract :exec
INSERT INTO tract_trigger_links (tract_id, trigger_id, filters)
VALUES ($1, $2, $3)
ON CONFLICT (tract_id, trigger_id) DO UPDATE
    SET filters = EXCLUDED.filters;

-- name: UnlinkTriggerFromTract :exec
DELETE FROM tract_trigger_links
WHERE trigger_id = $1
  AND tract_id = $2;

-- name: ListTriggerLinksByTract :many
-- tractUuid's linked triggers, Trigger populated — the tract editor's "wired up triggers" view.
SELECT t.id, t.trigger_uuid, t.user_id, t.name, t.kind, t.source, t.config, t.payload_schema,
       t.secret_hash, t.enabled, t.created_at, t.matchers, t.token_suffix, l.filters
FROM tract_trigger_links l
         JOIN triggers t ON t.id = l.trigger_id
WHERE l.tract_id = $1;

-- name: ListTriggerLinksByTrigger :many
-- triggerUuid's linked tracts, Tract populated — the webhook handler's fan-out lookup.
SELECT tr.id, tr.user_id, tr.name, tr.description, tr.enabled, tr.definition, tr.created_at, tr.updated_at, l.filters
FROM tract_trigger_links l
         JOIN tracts tr ON tr.id = l.tract_id
WHERE l.trigger_id = $1;

-- name: InsertTriggerProviderLink :exec
INSERT INTO trigger_provider_links (trigger_id, external_connection_id)
VALUES ($1, $2)
ON CONFLICT (trigger_id, external_connection_id) DO NOTHING;

-- name: ListTriggersByExternalConnection :many
-- The gitlab_webhook handler's fan-out lookup: every trigger sharing one provider connection.
SELECT t.id, t.trigger_uuid, t.user_id, t.name, t.kind, t.source, t.config, t.payload_schema,
       t.secret_hash, t.enabled, t.created_at, t.matchers, t.token_suffix
FROM trigger_provider_links l
         JOIN triggers t ON t.id = l.trigger_id
WHERE l.external_connection_id = $1;
