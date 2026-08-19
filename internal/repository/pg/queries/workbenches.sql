-- name: CreateWorkbench :one
INSERT INTO workbenches (vault_id, user_id, volume_name, status, container_id, docker_host_id)
VALUES ($1, $2, $3, 'configuring', NULL, $4)
RETURNING id, vault_id, user_id, status, auth_mode, container_id, volume_name, created_at, started_at, stopped_at, docker_host_id;

-- name: GetWorkbenchByVaultAndUser :one
SELECT id, vault_id, user_id, status, auth_mode, container_id, volume_name, created_at, started_at, stopped_at, docker_host_id
FROM workbenches
WHERE vault_id = $1
  AND user_id = $2;

-- name: MarkWorkbenchContainerCreated :exec
UPDATE workbenches
SET status = 'created', container_id = $2
WHERE vault_id = $1
  AND user_id = $3;

-- name: MarkWorkbenchConfiguring :exec
UPDATE workbenches
SET status = 'configuring'
WHERE vault_id = $1
  AND user_id = $2;

-- name: MarkWorkbenchRunning :exec
UPDATE workbenches
SET status = 'running', auth_mode = $2, started_at = NOW()
WHERE vault_id = $1
  AND user_id = $3;

-- name: MarkWorkbenchStopped :exec
UPDATE workbenches
SET status = 'stopped', stopped_at = NOW()
WHERE vault_id = $1
  AND user_id = $2;

-- name: MarkWorkbenchRemoved :exec
UPDATE workbenches
SET status = 'removed'
WHERE vault_id = $1
  AND user_id = $2;

-- name: DeleteWorkbench :exec
DELETE
FROM workbenches
WHERE vault_id = $1
  AND user_id = $2;

-- name: ListWorkbenchesByVaultID :many
SELECT id, vault_id, user_id, status, auth_mode, container_id, volume_name, created_at, started_at, stopped_at, docker_host_id
FROM workbenches
WHERE vault_id = $1;

-- name: GetMostRecentWorkbenchByUser :one
-- Used by the telegram webhook relay (internal/transport/telegram_webhook) to pick which of a
-- user's (possibly several) vault workbenches to relay chat into, v1 scope: just the one most
-- recently started (falling back to most recently created for a workbench never started).
SELECT id, vault_id, user_id, status, auth_mode, container_id, volume_name, created_at, started_at, stopped_at, docker_host_id
FROM workbenches
WHERE user_id = $1
ORDER BY started_at DESC NULLS LAST, created_at DESC
LIMIT 1;
