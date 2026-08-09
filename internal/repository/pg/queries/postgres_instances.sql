-- name: RegisterPostgresInstance :one
INSERT INTO postgres_instances (host, port, admin_database, username, password_enc, ssl_mode)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetPostgresInstance :one
SELECT id, host, port, admin_database, username, password_enc, ssl_mode, owner_user_id, created_at
FROM postgres_instances WHERE id = $1;

-- name: RandomPickPostgresInstance :one
SELECT id, host, port, admin_database, username, password_enc, ssl_mode, owner_user_id, created_at
FROM postgres_instances ORDER BY RANDOM() LIMIT 1;

-- name: PickOwnedPostgresInstance :one
SELECT id, host, port, admin_database, username, password_enc, ssl_mode, owner_user_id, created_at
FROM postgres_instances WHERE owner_user_id = $1 LIMIT 1;

-- name: RegisterOwnedPostgresInstance :one
INSERT INTO postgres_instances (host, port, admin_database, username, password_enc, ssl_mode, owner_user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: DeleteOwnedPostgresInstanceIfUnreferenced :exec
DELETE FROM postgres_instances
WHERE owner_user_id = $1
  AND NOT EXISTS (
      SELECT 1 FROM vault_postgres_databases
      WHERE vault_postgres_databases.postgres_instance_id = postgres_instances.id
  );

-- name: ListPostgresInstances :many
SELECT id, host, port, admin_database, username, password_enc, ssl_mode, owner_user_id, created_at
FROM postgres_instances ORDER BY created_at DESC;

-- name: UpdatePostgresInstance :exec
UPDATE postgres_instances
SET host = $2, port = $3, admin_database = $4, username = $5, password_enc = $6, ssl_mode = $7
WHERE id = $1;

-- name: DeletePostgresInstance :exec
DELETE FROM postgres_instances WHERE id = $1;

-- name: PostgresInstanceExists :one
SELECT EXISTS(SELECT 1 FROM postgres_instances);
