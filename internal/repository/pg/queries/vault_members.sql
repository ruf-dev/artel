-- name: AddVaultMember :exec
INSERT INTO vault_members (vault_id, user_id, role) VALUES ($1, $2, $3)
ON CONFLICT (vault_id, user_id) DO NOTHING;

-- name: RemoveVaultMember :exec
DELETE FROM vault_members WHERE vault_id = $1 AND user_id = $2;

-- name: GetVaultMembership :one
SELECT id, vault_id, user_id, role, created_at FROM vault_members WHERE vault_id = $1 AND user_id = $2;

-- name: ListVaultMembers :many
SELECT id, vault_id, user_id, role, created_at FROM vault_members WHERE vault_id = $1;
