-- name: RegisterS3Instance :one
INSERT INTO s3_instances (endpoint, region, access_key, secret_key_enc, use_ssl, path_style)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetS3Instance :one
SELECT id, endpoint, region, access_key, use_ssl, path_style, created_at FROM s3_instances WHERE id = $1;

-- name: GetS3InstanceWithCreds :one
SELECT id, endpoint, region, access_key, secret_key_enc, use_ssl, path_style, created_at FROM s3_instances WHERE id = $1;

-- name: ListS3Instances :many
SELECT id, endpoint, region, access_key, use_ssl, path_style, created_at FROM s3_instances ORDER BY created_at DESC;

-- name: UpdateS3Instance :exec
UPDATE s3_instances
SET endpoint = $2, region = $3, access_key = $4, secret_key_enc = $5, use_ssl = $6, path_style = $7
WHERE id = $1;

-- name: DeleteS3Instance :exec
DELETE FROM s3_instances WHERE id = $1;

-- name: S3InstanceExists :one
SELECT EXISTS(SELECT 1 FROM s3_instances);
