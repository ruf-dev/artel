-- name: RegisterCouchInstance :one
INSERT INTO couch_instances (url, username, password_enc)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetCouchInstance :one
SELECT id, url, username, created_at FROM couch_instances WHERE id = $1;

-- name: GetCouchInstanceWithCreds :one
SELECT id, url, username, password_enc, created_at FROM couch_instances WHERE id = $1;

-- name: RandomPickCouchInstance :one
SELECT id, url, username, password_enc, created_at FROM couch_instances ORDER BY RANDOM() LIMIT 1;

-- name: ListCouchInstances :many
SELECT id, url, username, created_at FROM couch_instances ORDER BY created_at DESC;

-- name: UpdateCouchInstance :exec
UPDATE couch_instances
SET url = $2, username = $3, password_enc = $4
WHERE id = $1;

-- name: DeleteCouchInstance :exec
DELETE FROM couch_instances WHERE id = $1;
