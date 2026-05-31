-- name: CreateVaultInvite :one
INSERT INTO vault_invites (vault_id, created_by, role, token)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetVaultInviteByToken :one
SELECT * FROM vault_invites WHERE token = $1;

-- name: ListVaultInvites :many
SELECT * FROM vault_invites WHERE vault_id = $1 ORDER BY created_at DESC;

-- name: RevokeVaultInvite :exec
UPDATE vault_invites SET revoked_at = NOW() WHERE id = $1 AND revoked_at IS NULL;
