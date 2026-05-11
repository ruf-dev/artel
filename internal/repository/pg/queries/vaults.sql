-- name: CreateVault :one
INSERT INTO vaults (user_id, name, couch_db_name, couch_instance_id, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, name, couch_db_name, couch_instance_id, status, created_at;

-- name: GetVaultByID :one
SELECT id, user_id, name, couch_db_name, couch_instance_id, status, created_at
FROM vaults
WHERE id = $1;

-- name: GetVaultByNameAndUser :one
SELECT id, user_id, name, couch_db_name, couch_instance_id, status, created_at
FROM vaults
WHERE user_id = $1
  AND name = $2;

-- name: UpdateVaultStatus :exec
UPDATE vaults
SET status = $2
WHERE id = $1;

-- name: ListVaultsByMembership :many
SELECT v.id, v.user_id, v.name, v.couch_db_name, v.couch_instance_id, v.status, v.created_at
FROM vaults v
         JOIN vault_members vm ON vm.vault_id = v.id
WHERE vm.user_id = $1;

-- name: DeleteVault :exec
DELETE
FROM vaults
WHERE id = $1;
