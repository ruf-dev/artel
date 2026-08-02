-- name: CreateVault :one
INSERT INTO vaults (user_id, name, couch_db_name, couch_instance_id, status, livesync_passphrase_enc, s3_instance_id, s3_bucket_name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, name, couch_db_name, couch_instance_id, status, livesync_passphrase_enc, s3_instance_id, s3_bucket_name, use_couchdb_for_binaries, created_at, is_public, slug;

-- name: GetVaultByID :one
SELECT id, user_id, name, couch_db_name, couch_instance_id, status, livesync_passphrase_enc, s3_instance_id, s3_bucket_name, use_couchdb_for_binaries, created_at, is_public, slug
FROM vaults
WHERE id = $1;

-- name: GetVaultByNameAndUser :one
SELECT id, user_id, name, couch_db_name, couch_instance_id, status, livesync_passphrase_enc, s3_instance_id, s3_bucket_name, use_couchdb_for_binaries, created_at, is_public, slug
FROM vaults
WHERE user_id = $1
  AND name = $2;

-- name: UpdateVaultStatus :exec
UPDATE vaults
SET status = $2
WHERE id = $1;

-- name: SetVaultLiveSyncPassphrase :exec
UPDATE vaults
SET livesync_passphrase_enc = $2
WHERE id = $1;

-- name: ListVaultsByMembership :many
SELECT v.id, v.user_id, v.name, v.couch_db_name, v.couch_instance_id, v.status, v.livesync_passphrase_enc, v.s3_instance_id, v.s3_bucket_name, v.use_couchdb_for_binaries, v.created_at, v.is_public, v.slug, vm.role AS member_role
FROM vaults v
         LEFT JOIN vault_members vm ON vm.vault_id = v.id AND vm.user_id = $1
WHERE vm.user_id IS NOT NULL
   OR v.is_public;

-- name: DeleteVault :exec
DELETE
FROM vaults
WHERE id = $1;

-- name: LinkVaultS3Bucket :exec
UPDATE vaults SET s3_instance_id = $2, s3_bucket_name = $3 WHERE id = $1;

-- name: UnlinkVaultS3Bucket :exec
UPDATE vaults SET s3_instance_id = NULL, s3_bucket_name = NULL WHERE id = $1;

-- name: SetVaultUseCouchDBForBinaries :exec
UPDATE vaults SET use_couchdb_for_binaries = $2 WHERE id = $1;

-- name: PublishVault :one
UPDATE vaults
SET is_public = true,
    slug      = $2
WHERE id = $1
RETURNING id, user_id, name, couch_db_name, couch_instance_id, status, livesync_passphrase_enc, s3_instance_id, s3_bucket_name, use_couchdb_for_binaries, created_at, is_public, slug;

-- name: UnpublishVault :exec
UPDATE vaults
SET is_public = false
WHERE id = $1;

-- name: GetVaultBySlug :one
SELECT id, user_id, name, couch_db_name, couch_instance_id, status, livesync_passphrase_enc, s3_instance_id, s3_bucket_name, use_couchdb_for_binaries, created_at, is_public, slug
FROM vaults
WHERE slug = $1
  AND is_public;
