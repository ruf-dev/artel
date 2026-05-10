-- name: StoreCouchCred :exec
INSERT INTO couch_credentials (vault_id, host, username, password_enc) VALUES ($1, $2, $3, $4);

-- name: LoadCouchCred :one
SELECT id, vault_id, host, username, password_enc, created_at FROM couch_credentials WHERE vault_id = $1;

-- name: DeleteCouchCred :exec
DELETE FROM couch_credentials WHERE vault_id = $1;
