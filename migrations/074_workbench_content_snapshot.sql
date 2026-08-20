-- +goose Up

-- content_snapshot is the path→mtime (unix ms) baseline captured the last time the vault was
-- materialized into the workbench container's workspace — see
-- internal/service/v1/workbench/workbench.go's materializeVault/syncWorkbenchToVault. Used on
-- stop to tell an untouched vault-side note apart from one edited elsewhere (e.g. the Obsidian
-- app) while the workbench was running.
--
-- IF NOT EXISTS: this was originally authored as 071_workbench_content_snapshot.sql and, before
-- being renumbered here, got applied by hand directly against the shared dev Postgres (outside
-- goose's own bookkeeping) after a numbering collision with a concurrently-developed branch that
-- had already taken 071-073 for an unrelated per-user workbenches refactor. A fresh
-- database/goose run under the new number 074 must not fail against that already-patched shared
-- instance.
ALTER TABLE workbenches ADD COLUMN IF NOT EXISTS content_snapshot JSONB;

-- +goose Down

ALTER TABLE workbenches DROP COLUMN IF EXISTS content_snapshot;
