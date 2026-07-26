-- name: RegisterDockerHost :one
INSERT INTO docker_hosts (url)
VALUES ($1)
RETURNING id;

-- name: GetDockerHost :one
SELECT id, url, created_at FROM docker_hosts WHERE id = $1;

-- name: ListDockerHosts :many
SELECT id, url, created_at FROM docker_hosts ORDER BY created_at DESC;

-- name: UpdateDockerHost :exec
UPDATE docker_hosts
SET url = $2
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
