# Task Breakdown

Ordered so each task is independently reviewable/mergeable and later tasks build on earlier ones.
File paths reference the current codebase layout. None of this is implemented yet.

## Stage 0 — Spike (no merged code, findings only)

### Task 0 — Confirm the `claude` headless login shape — DONE (2026-07-22)
_Ref: [03_auth_and_login_flow.md](03_auth_and_login_flow.md)_

- [x] Ran the built `deploy/workbench/Dockerfile` image, no `ANTHROPIC_API_KEY` set, drove it via
      `docker exec` + `tmux capture-pane`/`send-keys`, captured the literal prompt/output.
- [x] Determined: neither guessed shape exactly — it's a hosted-callback OAuth flow
      (`redirect_uri=https://platform.claude.com/oauth/code/callback`, not localhost) where the CLI
      prints a URL and waits for the user to paste a code back; no inbound port, no polling needed
      on Artel's side. See "Task 0 findings" in `03_auth_and_login_flow.md` for the full writeup.
- [x] Findings written back into `03_auth_and_login_flow.md`. Stage 3 (Task 7) is unblocked.

## Stage 1 — Provisioning (container + volume lifecycle, no auth yet)

### Task 1 — Migration + domain — DONE
_Ref: [01_data_model_and_lifecycle.md](01_data_model_and_lifecycle.md)_

- [x] `migrations/054_workbenches.sql`: `workbench_status`, `workbench_auth_mode` enums,
      `workbenches` table.
- [x] `internal/repository/pg/queries/workbenches.sql` + `sqlc generate`: insert, get-by-vault,
      mark-running/stopped/removed, delete.
- [x] `internal/domain/workbench.go`: `Workbench` struct, `WorkbenchStatus`/`WorkbenchAuthMode`
      typed constants mirroring the enum values.
- [x] `internal/repository/pg/repos/workbenches/` repo — pure DB ops only, no Docker calls here.
      Wired into the `Repos` aggregator (`internal/repository/interfaces.go`,
      `internal/repository/pg/impl.go`) so the service layer (Task 4) can consume it.

### Task 2 — `internal/clients/workbenchdocker` package — DONE
_Ref: [02_docker_topology.md](02_docker_topology.md)_

- [x] Added `github.com/docker/docker` (`+incompatible`, standard for this SDK) dependency.
- [x] `client.go`/`container.go`/`volume.go`: `New(host string) (*Client, error)`,
      `CreateContainer`, `StartContainer`, `StopContainer`, `RemoveContainer`, `CreateVolume`,
      `RemoveVolume` per the sketch in `02_docker_topology.md`. Containers attach to `workbench-net`,
      mount the volume at `/workspace`, no exposed ports, labeled `artel.workbench=true`.
- [x] Unit tests against `httptest.Server` (no live Docker daemon required for `go test ./...`).
- [x] `deploy/workbench/Dockerfile` + `entrypoint.sh`: `node:20-slim`, `@anthropic-ai/claude-code` +
      `tmux` installed, entrypoint starts `claude` inside a named `workbench` tmux session. Built
      and smoke-tested against a live daemon during the Task 0 spike above.

### Task 3 — Config gating — DONE
_Ref: [02_docker_topology.md](02_docker_topology.md)_

- [x] `config/config.yaml`: add `WorkbenchDockerHost` under `environment:`.
- [x] `rscli-dev project tidy` to regenerate `internal/config/environment.go` — never hand-edit it.
- [x] `internal/app/custom.go`: construct `workbenchdocker.Client` + `WorkbenchService` only when
      `cfg.WorkbenchDockerHost != ""`; skip transport route registration otherwise.

### Task 4 — `WorkbenchService`: create/get/delete (no start yet) — DONE
_Ref: [01_data_model_and_lifecycle.md](01_data_model_and_lifecycle.md)_

- [x] `internal/service/interfaces.go`: add `WorkbenchService` interface (create/get/start/stop/
      delete — start can be stubbed to return "not implemented" until Task 6).
- [x] `internal/service/v1/workbench/workbench.go`: `CreateWorkbench` (volume create → container
      create → DB insert, idempotent on `vault_id`), `GetWorkbench`, `DeleteWorkbench` (stop if
      running → remove container → remove volume → DB delete relies on cascade but Docker
      teardown must happen explicitly first).
- [x] `internal/service/v1/impl.go`: wire into `Services` struct (nil when Docker isn't
      configured — callers must handle a nil `Services.Workbench`, check existing patterns for
      how optional services are already handled, e.g. anything gated by `SubscriptionsEnabled`).

### Task 5 — Hook into vault create/delete — DONE
_Ref: [01_data_model_and_lifecycle.md](01_data_model_and_lifecycle.md)_

- [x] Found the `CreateVault` RPC's transport handler; after a successful `VaultService.CreateVault`
      call, also call `WorkbenchService.CreateWorkbench` (no-op / clearly logged if
      `Services.Workbench` is nil, i.e. Docker not configured for this deployment) — failure here
      does not roll back the vault.
- [x] Same handler's `DeleteVault` path: calls `WorkbenchService.DeleteWorkbench` before the
      vault delete completes, and fails loudly (does not proceed to delete the vault) if it errors.
- [x] `api/grpc/vaults.proto`: `GetVault.Response` gained `workbench_exists`/`workbench_status`
      fields (extended the existing vault-detail response rather than a new RPC, mirroring how
      `s3_instance_id`/`s3_bucket_name` are already exposed). `moti g` + `bun gen` regenerated.

## Stage 2 — Start/stop with `api_key` auth mode

### Task 6 — `StartWorkbench`/`StopWorkbench`, API-key path only — DONE
_Ref: [03_auth_and_login_flow.md](03_auth_and_login_flow.md)_

- [x] `WorkbenchService.StartWorkbench(ctx, vaultID, authMode)`: for `api_key`, looks up + decrypts
      the user's Anthropic `external_connections` row via new `ExternalConnectionService.
      GetAnthropicApiKey` (thin public wrapper over the existing private decrypt logic), fails with
      new `user_errors.WorkbenchMissingAnthropicConnection` if absent, `docker start` with
      `ANTHROPIC_API_KEY` injected, updates `workbenches.status='running'`, `started_at`. A
      non-`api_key` mode fails fast with `user_errors.WorkbenchAuthModeNotImplemented` (Task 7).
- [x] `StopWorkbench`: `docker stop`, updates `status='stopped'`, `stopped_at`.
- [x] Transport RPCs `StartWorkbench`/`StopWorkbench` added to `vaults.proto`, wired into
      `VaultsImpl` (nil-checked `workbenchSvc`, same pattern as create/delete hooks). No frontend
      UI wiring yet — deliberately backend/proto only, per "minimal — a button, not a full page".
- [x] Validates the prototype end-to-end at the service/API layer (unit-tested with fakes); a live
      Docker-daemon smoke test (`docker exec`/logs, not yet a user-facing terminal) is still a
      manual follow-up, not automated here.

## Stage 3 — `subscription_login` auth mode

### Task 7 — Login-URL capture + relay
_Ref: [03_auth_and_login_flow.md](03_auth_and_login_flow.md)_

- [ ] Blocked on Task 0's findings. Implement whichever shape the spike confirms — do not start
      this task until Task 0 is written up.
- [ ] `StartWorkbench` for `authMode='subscription_login'`: `docker start` with no key, stream
      stdout, extract the login URL.
- [ ] New RPC/polling mechanism to surface the URL (and later, "authorized"/"failed" status) to
      the frontend.
- [ ] Frontend: minimal UI to display the URL and poll status — not a full workbench page, just
      enough to prove the flow.

## Stage 4 — deferred, not scoped in detail here

- **Remote terminal access** — the actual "reach the running `claude` session from any device"
  piece (PTY-over-gRPC/WebSocket bridge, `tmux` attach, auth through Artel's existing session
  auth). Prototype above only proves lifecycle + login; a real user-facing terminal is its own
  design pass once those are solid.
- **Notes/file sync onto the workbench volume** — explicitly out of scope per the README; a
  separate design document once the container/login mechanics are proven.
- **Resource quotas tied to billing/subscription tier.**
- **Reconciliation** — DB says `running` but the container died (crash, host restart, OOM kill);
  needs a health-check/reconcile loop, not designed here.
- **Docker topology Option C** (dedicated remote host with mTLS) — revisit only if capacity or
  isolation requirements outgrow the same-VM second-daemon setup from Stage 1.
