-- name: CreatePendingAuthCode :exec
INSERT INTO pending_auth_codes (code, raw_token, code_challenge, redirect_uri, client_id, expires_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetPendingAuthCode :one
SELECT code, raw_token, code_challenge, redirect_uri, client_id, expires_at
FROM pending_auth_codes
WHERE code = $1;

-- name: DeletePendingAuthCode :exec
DELETE FROM pending_auth_codes WHERE code = $1;

-- name: DeleteExpiredPendingAuthCodes :exec
DELETE FROM pending_auth_codes WHERE expires_at < NOW();
