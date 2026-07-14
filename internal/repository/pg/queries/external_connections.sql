-- name: UpsertExternalConnection :one
INSERT INTO external_connections (user_id, provider, provider_type, credentials_enc, metadata)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, provider) WHERE provider <> 'trello' DO UPDATE
    SET provider_type   = EXCLUDED.provider_type,
        credentials_enc = EXCLUDED.credentials_enc,
        metadata        = EXCLUDED.metadata,
        updated_at      = NOW()
RETURNING id, user_id, provider, provider_type, credentials_enc, metadata, created_at, updated_at;

-- name: InsertExternalConnection :one
INSERT INTO external_connections (user_id, provider, provider_type, credentials_enc, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, provider, provider_type, credentials_enc, metadata, created_at, updated_at;

-- name: GetExternalConnectionByUserAndProvider :one
SELECT id, user_id, provider, provider_type, credentials_enc, metadata, created_at, updated_at
FROM external_connections
WHERE user_id = $1 AND provider = $2;

-- name: ListExternalConnectionsByUser :many
SELECT id, user_id, provider, provider_type, credentials_enc, metadata, created_at, updated_at
FROM external_connections
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetExternalConnectionByID :one
SELECT id, user_id, provider, provider_type, credentials_enc, metadata, created_at, updated_at
FROM external_connections
WHERE id = $1;

-- name: DeleteExternalConnection :exec
DELETE FROM external_connections
WHERE user_id = $1 AND provider = $2;

-- name: DeleteExternalConnectionByID :exec
DELETE FROM external_connections
WHERE id = $1 AND user_id = $2;
