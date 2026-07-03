-- name: InsertTract :one
INSERT INTO tracts (user_id, name, description, enabled, definition)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, description, enabled, definition, created_at, updated_at;

-- name: GetTract :one
SELECT id, user_id, name, description, enabled, definition, created_at, updated_at
FROM tracts
WHERE id = $1;

-- name: ListTractsByUser :many
SELECT id, user_id, name, description, enabled, definition, created_at, updated_at
FROM tracts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateTract :one
UPDATE tracts
SET name        = $2,
    description = $3,
    definition  = $4,
    updated_at  = NOW()
WHERE id = $1
RETURNING id, user_id, name, description, enabled, definition, created_at, updated_at;

-- name: SetTractEnabled :exec
UPDATE tracts
SET enabled    = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteTract :exec
DELETE FROM tracts
WHERE id = $1;

-- name: InsertTractRun :one
INSERT INTO tract_runs (tract_id, trigger_id, started_by, trigger_output)
VALUES ($1, $2, $3, $4)
RETURNING id, tract_id, trigger_id, started_by, status, trigger_output, error, created_at, updated_at;

-- name: GetTractRun :one
SELECT id, tract_id, trigger_id, started_by, status, trigger_output, error, created_at, updated_at
FROM tract_runs
WHERE id = $1;

-- name: ListTractRunsByTract :many
SELECT id, tract_id, trigger_id, started_by, status, trigger_output, error, created_at, updated_at
FROM tract_runs
WHERE tract_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: UpdateTractRunStatus :exec
UPDATE tract_runs
SET status     = $2,
    error      = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: SweepStaleTractRuns :exec
UPDATE tract_runs
SET status     = 'failed',
    error      = 'run did not finish before the server restarted',
    updated_at = NOW()
WHERE status = 'running'
  AND created_at < $1;

-- name: InsertTractRunStep :one
INSERT INTO tract_run_steps (run_id, step_id, step_name, step_type, input)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, run_id, step_id, step_name, step_type, input, output, status, error, started_at, finished_at;

-- name: UpdateTractRunStepFinish :exec
UPDATE tract_run_steps
SET output      = $2,
    status      = $3,
    error       = $4,
    finished_at = NOW()
WHERE id = $1;

-- name: ListTractRunStepsByRun :many
SELECT id, run_id, step_id, step_name, step_type, input, output, status, error, started_at, finished_at
FROM tract_run_steps
WHERE run_id = $1
ORDER BY started_at ASC;

-- name: SweepStaleTractRunSteps :exec
UPDATE tract_run_steps
SET status      = 'failed',
    error       = 'step did not finish before the server restarted',
    finished_at = NOW()
WHERE status = 'running'
  AND started_at < $1;
