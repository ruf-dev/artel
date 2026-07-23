-- name: CreateWorkbench :one
INSERT INTO workbenches (vault_id, user_id, volume_name, container_id)
VALUES ($1, $2, $3, $4)
RETURNING id, vault_id, user_id, status, auth_mode, container_id, volume_name, created_at, started_at, stopped_at;

-- name: GetWorkbenchByVaultID :one
SELECT id, vault_id, user_id, status, auth_mode, container_id, volume_name, created_at, started_at, stopped_at
FROM workbenches
WHERE vault_id = $1;

-- name: MarkWorkbenchRunning :exec
UPDATE workbenches
SET status = 'running', auth_mode = $2, started_at = NOW()
WHERE vault_id = $1;

-- name: MarkWorkbenchStopped :exec
UPDATE workbenches
SET status = 'stopped', stopped_at = NOW()
WHERE vault_id = $1;

-- name: MarkWorkbenchRemoved :exec
UPDATE workbenches
SET status = 'removed'
WHERE vault_id = $1;

-- name: DeleteWorkbench :exec
DELETE
FROM workbenches
WHERE vault_id = $1;
