# Data Model & Lifecycle

## Table

New migration (next free number as of this writing is `054_*.sql` — confirm against
`migrations/` at implementation time, the folder moves fast):

```sql
-- +goose Up
CREATE TYPE workbench_status AS ENUM (
    'created',
    'running',
    'stopped',
    'removed'
);

CREATE TYPE workbench_auth_mode AS ENUM (
    'api_key',
    'subscription_login'
);

CREATE TABLE workbenches
(
    id             UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id       UUID                 NOT NULL REFERENCES vaults (id) ON DELETE CASCADE,
    user_id        UUID                 NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status         workbench_status     NOT NULL DEFAULT 'created',
    auth_mode      workbench_auth_mode,
    container_id   TEXT,
    volume_name    TEXT                 NOT NULL,
    created_at     TIMESTAMPTZ          NOT NULL DEFAULT NOW(),
    started_at     TIMESTAMPTZ,
    stopped_at     TIMESTAMPTZ,
    UNIQUE(vault_id)
);

CREATE INDEX ON workbenches(user_id);

-- +goose Down
DROP TABLE workbenches;
DROP TYPE workbench_auth_mode;
DROP TYPE workbench_status;
```

Notes on the shape, following the `vaults`/`external_connections` conventions already in the repo
(see `migrations/001_initial_schema.sql`, `migrations/024_external_connections.sql`):

- `UNIQUE(vault_id)` — one workbench per vault, matching "created alongside the vault." If a user
  ever needs more than one workbench per vault, that's a deliberate later decision, not an
  oversight — don't relax this preemptively.
- `auth_mode` is nullable and only set at `Start` time (a workbench sitting in `created` hasn't
  chosen how it'll authenticate yet).
- `container_id` is nullable for the same reason: the `docker create` call happens at workbench
  creation (see below), so it's actually populated immediately — nullable only to cover a
  creation-failed/retry edge case cleanly rather than using a sentinel empty string.
- No `image_tag`/resource-limit columns yet — hardcode the image and limits in the service for the
  prototype; promote to columns only once there's a reason to vary them per workbench (e.g. a
  paid-tier resource multiplier).

## State machine

```
created -> running -> stopped -> running (restart)
created -> removed
stopped -> removed
```

| Status | Docker state | Meaning |
|---|---|---|
| `created` | container exists via `docker create`, never started | Vault exists, workbench reserved, nothing running, nothing billed |
| `running` | container started | `claude` session live inside, reachable |
| `stopped` | container stopped, volume retained | Paused — user or idle-timeout stopped it, notes/state preserved |
| `removed` | container + volume deleted | Vault was deleted, or workbench explicitly torn down |

Deliberately mapping 1:1 onto real Docker container states (`docker create` really does leave a
container in a stopped, non-running state distinct from one that's been started and stopped) —
the DB row is a thin reflection of Docker's own lifecycle, not a parallel state machine that can
drift from it. Reconciliation (what happens if the DB says `running` but the container died) is a
later hardening concern — not designed here, flag it in the task breakdown instead of guessing.

## Hook points

**Vault creation** — `VaultService.CreateVault` is untouched. The transport handler that calls it
(wherever `CreateVault` RPC is currently handled) additionally calls
`WorkbenchService.CreateWorkbench(ctx, vaultID)` right after `CreateVault` succeeds. Keeping this
a sibling call at the handler rather than a dependency inside `VaultService` avoids giving the
vault package a compile-time dependency on Docker — mirrors how `LinkS3Bucket` is already a
separate, explicit call rather than something `CreateVault` triggers automatically.

`WorkbenchService.CreateWorkbench`:
1. Generate a volume name (`workbench-<vault_id>`).
2. `docker volume create`.
3. `docker create` with that volume mounted, the workbench image, no auth env vars yet (none
   exist until `Start` picks a mode), container **not started**.
4. Insert the `workbenches` row with `status='created'`, `container_id` from the create call.

If step 2/3 fails, the vault itself still exists — surface the workbench creation failure
separately (e.g. a toast) rather than rolling back vault creation; a vault without a workbench
should be a recoverable, retryable state (`CreateWorkbench` can be re-invoked idempotently keyed
on `vault_id`), not a reason to fail the whole vault flow.

**Vault deletion** — same pattern in reverse: `VaultService.DeleteVault`'s call site also invokes
`WorkbenchService.DeleteWorkbench(ctx, vaultID)`. The DB row cascades automatically via
`ON DELETE CASCADE`, but the actual Docker container and volume do not — `DeleteWorkbench` must
explicitly `docker stop` (if running) + `docker rm` + `docker volume rm` *before* the vault delete
completes, or on failure leave an orphaned container that nothing in Postgres points to anymore.
Ordering matters here more than in most cascades — get this right before shipping, an orphaned
container is a live `claude` process with a live API key nobody can reach anymore to stop.

## Interface sketch

Following `service.VaultService`'s shape (`internal/service/interfaces.go`):

```go
type WorkbenchService interface {
    CreateWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error)
    GetWorkbench(ctx context.Context, vaultID uuid.UUID) (domain.Workbench, error)
    StartWorkbench(ctx context.Context, vaultID uuid.UUID, authMode domain.WorkbenchAuthMode) (domain.Workbench, error)
    StopWorkbench(ctx context.Context, vaultID uuid.UUID) error
    DeleteWorkbench(ctx context.Context, vaultID uuid.UUID) error
}
```

Keyed by `vault_id` throughout rather than a separate workbench ID in the public API, since the
`UNIQUE(vault_id)` constraint makes the vault the natural handle — the frontend already has the
vault ID everywhere it would show workbench UI.
