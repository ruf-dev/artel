-- name: InsertEmailAccount :one
INSERT INTO email_accounts (user_id, email, imap_host, imap_port, smtp_host, smtp_port, password_enc)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, email, imap_host, imap_port, smtp_host, smtp_port, password_enc, created_at;

-- name: GetEmailAccountByUuid :one
SELECT id, user_id, email, imap_host, imap_port, smtp_host, smtp_port, password_enc, created_at
FROM email_accounts
WHERE id = $1;

-- name: ListEmailAccountsByUser :many
SELECT id, user_id, email, imap_host, imap_port, smtp_host, smtp_port, password_enc, created_at
FROM email_accounts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteEmailAccount :exec
DELETE FROM email_accounts WHERE id = $1;
