-- name: CreateVaultPostgresDatabase :one
INSERT INTO vault_postgres_databases (vault_id, postgres_instance_id, database_name, role_username, role_password_enc)
VALUES ($1, $2, $3, $4, $5)
RETURNING vault_id, postgres_instance_id, database_name, role_username, role_password_enc, status, error_message, created_at, updated_at;

-- name: GetVaultPostgresDatabaseByVaultID :one
SELECT vault_id, postgres_instance_id, database_name, role_username, role_password_enc, status, error_message, created_at, updated_at
FROM vault_postgres_databases WHERE vault_id = $1;

-- name: MarkVaultPostgresDatabaseReady :exec
UPDATE vault_postgres_databases
SET status = 'ready', error_message = '', updated_at = NOW()
WHERE vault_id = $1;

-- name: MarkVaultPostgresDatabaseError :exec
UPDATE vault_postgres_databases
SET status = 'error', error_message = $2, updated_at = NOW()
WHERE vault_id = $1;

-- name: DeleteVaultPostgresDatabase :exec
DELETE FROM vault_postgres_databases WHERE vault_id = $1;
