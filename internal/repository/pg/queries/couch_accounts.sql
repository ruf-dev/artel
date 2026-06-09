-- name: UpsertCouchAccount :exec
INSERT INTO couch_accounts (user_id, couch_instance_id, couch_username, couch_password_enc)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING;


-- name: GetCouchAccountByUserAndInstance :one
SELECT id, user_id, couch_instance_id, couch_username, couch_password_enc, created_at
FROM couch_accounts
WHERE user_id = $1
  AND couch_instance_id = $2;

-- name: ListCouchAccountsByUser :many
SELECT id, user_id, couch_instance_id, couch_username, couch_password_enc, created_at
FROM couch_accounts
WHERE user_id = $1;

-- name: DeleteCouchAccount :exec
DELETE
FROM couch_accounts
WHERE id = $1;
