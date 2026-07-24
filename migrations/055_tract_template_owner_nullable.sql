-- +goose Up
-- Built-in templates (seeded via migration, e.g. 056) are owned by nobody — a fresh install has
-- no user rows to reference, so owner_id must allow NULL. NULL means "built-in", see
-- domain.TractTemplate doc comment.
ALTER TABLE tract_templates ALTER COLUMN owner_id DROP NOT NULL;

-- +goose Down
ALTER TABLE tract_templates ALTER COLUMN owner_id SET NOT NULL;
