-- name: AddVaultMember :exec
INSERT INTO vault_members (vault_id, user_id, role) VALUES ($1, $2, $3)
ON CONFLICT (vault_id, user_id) DO NOTHING;

-- name: RemoveVaultMember :exec
DELETE FROM vault_members WHERE vault_id = $1 AND user_id = $2;

-- name: GetVaultMembership :one
SELECT id, vault_id, user_id, role, created_at FROM vault_members WHERE vault_id = $1 AND user_id = $2;

-- name: ListVaultMembers :many
SELECT id, vault_id, user_id, role, created_at FROM vault_members WHERE vault_id = $1;

-- name: ListVaultMembersWithUsers :many
SELECT vm.id, vm.vault_id, vm.user_id, vm.role, vm.created_at,
       u.email, u.username
FROM vault_members vm
JOIN users u ON u.id = vm.user_id
WHERE vm.vault_id = $1;
