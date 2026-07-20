-- name: InsertTractTemplate :one
INSERT INTO tract_templates (source_tract_id, owner_id, name, description, definition, category)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, source_tract_id, owner_id, name, description, definition, category, install_count, published_at, updated_at;

-- name: GetTractTemplate :one
SELECT id, source_tract_id, owner_id, name, description, definition, category, install_count, published_at, updated_at
FROM tract_templates
WHERE id = $1;

-- name: ListTractTemplates :many
SELECT id, source_tract_id, owner_id, name, description, definition, category, install_count, published_at, updated_at
FROM tract_templates
WHERE ($1::text = '' OR category = $1)
ORDER BY install_count DESC, published_at DESC;

-- name: ListTractTemplatesByOwner :many
SELECT id, source_tract_id, owner_id, name, description, definition, category, install_count, published_at, updated_at
FROM tract_templates
WHERE owner_id = $1
ORDER BY published_at DESC;

-- name: DeleteTractTemplate :exec
DELETE FROM tract_templates
WHERE id = $1;

-- name: IncrementTractTemplateInstallCount :exec
UPDATE tract_templates
SET install_count = install_count + 1,
    updated_at    = NOW()
WHERE id = $1;
