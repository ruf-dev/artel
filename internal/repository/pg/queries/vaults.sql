-- name: CreateVault :one
INSERT INTO vaults (user_id, name, couch_db_name) VALUES ($1, $2, $3) RETURNING id, user_id, name, couch_db_name, created_at;

-- name: GetVaultByName :one
SELECT id, user_id, name, couch_db_name, created_at FROM vaults WHERE name = $1;

-- name: ListAllVaults :many
SELECT id, user_id, name, couch_db_name, created_at FROM vaults;

-- name: ListVaultsByUser :many
SELECT id, user_id, name, couch_db_name, created_at FROM vaults WHERE user_id = $1;

-- name: DeleteVault :exec
DELETE FROM vaults WHERE id = $1;
