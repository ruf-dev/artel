-- name: RegisterDockerHost :one
INSERT INTO docker_hosts (url, ca_cert_enc, client_cert_enc, client_key_enc)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetDockerHost :one
SELECT id, url, created_at FROM docker_hosts WHERE id = $1;

-- name: GetDockerHostWithCreds :one
SELECT id, url, created_at, ca_cert_enc, client_cert_enc, client_key_enc
FROM docker_hosts
WHERE id = $1;

-- name: ListDockerHosts :many
SELECT id, url, created_at FROM docker_hosts ORDER BY created_at DESC;

-- name: UpdateDockerHost :exec
-- ca_cert_enc/client_cert_enc/client_key_enc are three-way patch fields: the matching
-- update_* boolean is false when the caller didn't touch that cert (leave the stored value
-- alone), true with a NULL value when the caller cleared it, true with a value when the caller
-- set/replaced it. See internal/repository/pg/repos/dockerhosts/dockerhosts.go's Update.
UPDATE docker_hosts
SET url             = $2,
    ca_cert_enc     = CASE WHEN sqlc.arg(update_ca_cert)::boolean THEN sqlc.arg(ca_cert_enc)::bytea ELSE ca_cert_enc END,
    client_cert_enc = CASE WHEN sqlc.arg(update_client_cert)::boolean THEN sqlc.arg(client_cert_enc)::bytea ELSE client_cert_enc END,
    client_key_enc  = CASE WHEN sqlc.arg(update_client_key)::boolean THEN sqlc.arg(client_key_enc)::bytea ELSE client_key_enc END
WHERE id = $1;

-- name: DeleteDockerHost :exec
DELETE FROM docker_hosts WHERE id = $1;

-- name: DockerHostExists :one
SELECT EXISTS(SELECT 1 FROM docker_hosts);

-- name: PickLeastLoadedDockerHost :one
SELECT dh.id, dh.url, dh.created_at
FROM docker_hosts dh
LEFT JOIN workbenches w ON w.docker_host_id = dh.id AND w.status IN ('created','running')
GROUP BY dh.id
ORDER BY COUNT(w.id) ASC, dh.created_at ASC
LIMIT 1;
