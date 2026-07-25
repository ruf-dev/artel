-- +goose Up
ALTER TYPE workbench_status ADD VALUE 'configuring';

-- +goose Down
-- enum values cannot be dropped in Postgres; this migration is irreversible.
