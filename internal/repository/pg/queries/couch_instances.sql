-- name: RegisterCouchInstance :one
INSERT INTO couch_instances (url, username, password_enc)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetCouchInstance :one
SELECT id, url, username, created_at FROM couch_instances WHERE id = $1;

-- name: GetCouchInstanceWithCreds :one
SELECT id, url, username, password_enc, created_at FROM couch_instances WHERE id = $1;

-- name: RandomPickCouchInstance :one
SELECT id, url, username, password_enc, created_at FROM couch_instances WHERE owner_user_id IS NULL ORDER BY RANDOM() LIMIT 1;

-- name: PickOwnedCouchInstance :one
SELECT id, url, username, password_enc, created_at FROM couch_instances WHERE owner_user_id = $1 LIMIT 1;

-- name: RegisterOwnedCouchInstance :one
INSERT INTO couch_instances (url, username, password_enc, owner_user_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: DeleteOwnedCouchInstanceIfUnreferenced :exec
DELETE FROM couch_instances
WHERE owner_user_id = $1
  AND NOT EXISTS (SELECT 1 FROM vaults WHERE vaults.couch_instance_id = couch_instances.id);

-- name: ListCouchInstances :many
SELECT id, url, username, created_at FROM couch_instances ORDER BY created_at DESC;

-- name: UpdateCouchInstance :exec
UPDATE couch_instances
SET url = $2, username = $3, password_enc = $4
WHERE id = $1;

-- name: DeleteCouchInstance :exec
DELETE FROM couch_instances WHERE id = $1;

-- name: CouchInstanceExists :one
SELECT EXISTS(SELECT 1 FROM couch_instances);
