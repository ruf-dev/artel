# E2E Business Coverage & Code Wellness Report

Date: 2026-08-08
Scope: business e2e test coverage vs. actual service capabilities, plus a code
wellness audit (duplication, goroutine safety, layering) across
storage/clients/service/transport. Research only — no code changes made.

## 1. Business e2e test coverage

### Coverage matrix

| Business flow | Test file(s) | Real infra? | Paid/quota path? | Verdict |
|---|---|---|---|---|
| Register + login (email/password) | `tests/e2e/e2e_test.go`, `byok_storage_test.go`, `quota_test.go`, `tract_e2e_test.go`, `gitlab_trigger_e2e_test.go` | Postgres, real gRPC+auth interceptor | quota_test only | Yes |
| Vault creation (pool CouchDB) | `e2e_test.go`, `quota_test.go` | Postgres+CouchDB | quota_test | Yes |
| Vault creation (BYOK CouchDB/S3) | `byok_storage_test.go` | Postgres+CouchDB+MinIO | no | Yes |
| Vault deletion (full service flow) | none (only `tests/vault/vault_test.go`, which tests raw `couchdb.Client.DeleteDatabase`, not `Vault.DeleteVault`) | partial | no | **Partial** |
| MCP key create/revoke | `e2e_test.go`, `quota_test.go` | Postgres | quota_test | Yes |
| Notes write/read via MCP JSON-RPC | `e2e_test.go`, `quota_test.go` | CouchDB | quota_test | Yes |
| CouchDB storage quota enforcement | `quota_test.go` | CouchDB | Yes | Yes |
| S3 storage quota enforcement | `quota_test.go` | MinIO | Yes | Yes |
| External connections: CouchDB/S3 BYOK (add, reject bad creds) | `byok_storage_test.go` | Real | no | Yes |
| External connections: OpenAI/Anthropic/Trello/Google | none | — | — | **No coverage** |
| MoM `http` action (GitLab tools) | `gitlab_trigger_e2e_test.go` | mocked GitLab HTTP server | no | Yes (mocked, not real GitLab) |
| MoM `imap`/`smtp` action (email) | none | — | — | **No coverage** |
| MoM ownership/community MoMs | `mcp_ownership_test.go` | Postgres only | no | Partial (repo-level, no gRPC/API surface) |
| Tract CRUD, manual trigger, run, step result | `tract_e2e_test.go` | Postgres+CouchDB | no | Yes |
| GitLab webhook trigger → tract run (push, MR-merged, non-matching action) | `gitlab_trigger_e2e_test.go` | Postgres, mocked GitLab | no | Yes |
| LiveSync CouchDB doc compatibility (write/read/delete/move/tags/folders) | `tests/livesync/livesync_test.go` | CouchDB | n/a | Yes, but **missing `//go:build e2e` tag** |
| Admin: users/settings/subscriptions/couch instances/S3 instances/docker hosts | none | — | — | **No coverage** |
| Setup wizard (initial system provisioning) | none (only passed as constructor arg) | — | — | **No coverage** |
| Prompt, tasktracker, public_docs, workbench services | none (only passed as constructor args) | — | — | **No coverage** |
| Auth failure paths (bad password, invalid/expired token, duplicate email) | none | — | — | **No coverage** |

### Top risk gaps (highest priority first)

1. **`tests/livesync/livesync_test.go` and `tests/vault/vault_test.go` have no
   `//go:build e2e` tag** — confirmed via `go list ./tests/...`, both packages
   are picked up by a plain `go test ./...`. This directly contradicts
   CLAUDE.md's "Unit tests: `go test ./...` — no external services required":
   these suites dial `localhost:15985` CouchDB unconditionally and will
   fail/hang in any environment without the docker-compose stack running.
   This is a build-hygiene bug, not just a coverage gap, and should be the
   first thing fixed.

2. **No e2e coverage of admin-side services** (`adminusers`,
   `adminsettings`, `adminsubscriptions`, `admincouchsvc`, `couchinstances`,
   `s3instances`, `dockerhosts`) — these are the operator-facing surface that
   provisions and manages the shared CouchDB/S3 pool the whole product
   depends on; a regression here silently breaks provisioning for every
   user, but nothing exercises the admin transport layer end-to-end
   (`internal/transport/admin_*_api`).

3. **MoM `imap`/`smtp` actions are completely untested e2e** — only the
   `http` action (GitLab) has a mocked-server e2e test; email is a named
   example flow in the README/spec (`pkg/mom_examples/README.md`) and a
   first-class `ToolAction` variant, but nothing exercises `EmailExecutor`
   against even a fake SMTP/IMAP server.

4. **Vault deletion has no real e2e path** — `Vault.DeleteVault` is only
   ever called from `t.Cleanup` blocks in other suites (never asserted on),
   and `tests/vault/vault_test.go` tests raw CouchDB database deletion, not
   the service-layer flow (ownership checks, S3 bucket cleanup, DB row
   removal). A regression that leaves orphaned CouchDB databases or S3
   buckets would go undetected.

5. **No negative/error-path coverage for auth** — wrong password,
   expired/invalid session token, duplicate email registration. Auth is the
   single gate in front of every paid flow; only the happy path is tested
   across all suites.

6. **BYOK/external-connections coverage stops at CouchDB and S3** —
   `AddOpenAIConnection`, `AddAnthropicConnection`, `AddTrelloConnection`,
   and the Google OAuth flow
   (`internal/service/v1/externalconnections/external_connections.go`) have
   zero e2e references; these are exactly the kind of "reachable from real
   API, not just Go code" surfaces `byok_storage_test.go` was written to
   validate for CouchDB/S3, but the pattern wasn't extended.

7. **Setup wizard flow is untested** — `internal/service/v1/setupwizard`
   only appears as a constructor argument wired into
   `auth_api.NewAuthImpl` in every suite; the actual first-run provisioning
   flow it implements has no test driving it.

8. **`mcp_ownership_test.go` only exercises the repository layer**, not the
   actual `McpKeysAPI`/`McpHandler` transport — so "community MoM"
   visibility/authorization at the API boundary (can another user call or
   list someone else's community MoM tool?) is unverified.

### Notes on test quality

- All suites are careful about real infra (no mocking of CouchDB/S3/Postgres)
  and generally strong on cleanup (`t.Cleanup`) and negative-path testing
  where it does exist (`byok_storage_test.go`'s bad-creds case,
  `quota_test.go`'s quota-exceeded case) — this is a genuine strength, not a
  smell.
- The bootstrap gating (`e2e_bootstrap` tag) correctly isolates one-time
  setup/teardown from being accidentally picked up by
  `go test -tags e2e ./tests/...`, but the missing tag on `livesync`/`vault`
  shows the same discipline wasn't applied uniformly.

## 2. Code wellness

**Method:** `docs/architecture.md`/`docs/go-style.md` reviewed for intended
layering, cross-checked against `graphify-out/GRAPH_REPORT.md` (no Go import
cycles detected; only frontend has them), then `internal/` and `cmd/` grepped
directly for `go func`, `context.With*`, and cross-layer imports, with every
hit verified by reading the actual code. `go vet ./...` is clean.

**Overall health by layer:** Storage (`internal/repository`) is mostly
disciplined but has 2 concrete violations of its own documented "pure DB ops"
rule. Clients (`internal/clients`) are clean and well-commented (the IMAP
concurrent-drain pattern is a genuinely good example), main duplication is
cosmetic logging boilerplate. Service (`internal/service/v1`) is the biggest
offender for duplication — several packages are copy-pasted per entity.
Transport is where the most consequential layering smell lives: two webhook
handlers and the MCP OAuth handler hold repository interfaces directly,
self-documented as intentional shortcuts rather than gaps in understanding.

### Top findings

1. **[HIGH, concurrency]** `internal/service/v1/tract/dispatch.go:62` —
   `go startRun(baseCtx, ...)` fans out one unbounded, un-recovered goroutine
   per matched trigger-link on every webhook delivery (gitlab_webhook and
   tract_webhook both funnel through this). No `recover()` anywhere in the
   call chain, so a panic inside `StartRun`
   (`internal/service/v1/tract/engine.go:64`) crashes the whole process; no
   worker-pool/semaphore caps fan-out under a burst of webhook deliveries; no
   WaitGroup to drain in-flight runs on shutdown. Fix direction: wrap
   `startRun` in a `defer recover()`, bound concurrency.

2. **[HIGH, layering]** `internal/transport/gitlab_webhook/handler.go:36-38`,
   `internal/transport/tract_webhook/handler.go:32`,
   `internal/transport/mcp_api/oauth.go:35,42` — these transport handlers
   hold `repository.ExternalConnectionRepo`/`TriggersRepo`/`PendingAuthCodes`
   and call repo methods (`GetByID`, `GetByTriggerUuid`) directly, bypassing
   the service layer entirely — contradicts `docs/architecture.md`'s
   transport→service→clients flow. It's explicitly self-documented in a
   comment ("so it depends on repository.TriggersRepo rather than going
   through service.TractService"), i.e. a known, accepted shortcut rather
   than an oversight — worth a deliberate decision on whether to formalize
   or fix.

3. **[HIGH, duplication]**
   `internal/service/v1/externalconnections/external_connections.go:795-861`
   (`AddAnthropicConnection`) and `:1300-1366` (`AddOpenAIConnection`) are
   ~95% line-for-line identical, as are their paired
   `Check*Connection`/`stored*ApiKey`/`validate*Key` helpers — roughly 250+
   duplicated lines across the file's 1799 total. Fix direction: extract a
   provider-parametrized helper (provider name, credential struct,
   model-list validator as params).

4. **[MEDIUM, layering]** `internal/repository/pg/repos/sessions/sessions.go:54,69`
   and `internal/repository/pg/pg_err/err.go:14` — repository code directly
   returns `user_errors.Unauthenticated`/`user_errors.NotFound`
   (gRPC-coded, service-defined sentinels) instead of the
   `sql.Null[T]{Valid:false}` pattern `docs/go-style.md` mandates ("Absorb
   sql.ErrNoRows inside the repo... never errors.Is(sql.ErrNoRows) at the
   service layer"). `pg_err.UnwrapPgErr` is used this way in 16 repository
   files, so the inversion is load-bearing, not a one-off.

5. **[MEDIUM, duplication]**
   `internal/service/v1/couchinstances/couchinstances.go:48-97`,
   `internal/service/v1/dockerhosts/dockerhosts.go:35-84`,
   `internal/service/v1/s3instances/s3instances.go` (Get/List/Update/Delete/Has,
   ~12 functions total) — near-identical parse-uuid → call repo →
   `rerrors.Wrap` skeleton repeated per "instance registry" entity. Fix
   direction: a small generic CRUD helper (Go generics) or shared embedding.

6. **[MEDIUM, duplication]** `internal/clients/imap/client.go:48,240,334,367,405`
   — the same `opStart := time.Now(); defer func(){ if err != nil {...}
   else {...} }()` timing/logging boilerplate is repeated 5 times
   (`connect`, `ReadEmail`, `TestConnection`, `ListFolders`, and one more).
   Fix direction: extract `defer logOp(ctx, "op", addr, &err, start)`.

7. **[MEDIUM, concurrency]** `internal/service/v1/mcp/resolve_key.go:113` —
   `go func(){ ... TouchLastAccessed(context.Background(), ...) }()` fires
   on every `ResolveKey` call (i.e. essentially every MCP tool invocation),
   deliberately on `context.Background()` with a documented rationale. Low
   consequence if lost (just a last-accessed timestamp) but still unbounded
   and un-recovered; worth a shared "fire-and-forget with recover" helper if
   this pattern spreads further.

8. **[LOW, layering]** `internal/service/user_errors/user_errors.go:6,83,239,245`
   — the service-layer error-definitions package imports the generated gRPC
   package `internal/api/server/artel_api` (`pb.UserErrors_*`) to build
   precondition-failure details. A structural inversion (service depending
   on transport-layer generated code) but narrow and arguably necessary for
   client-matchable error codes — lowest priority of the set.

**Not flagged (checked, found clean):** the ~18 `go func` sites in
`internal/api/server/artel_api/*.pb.gw.go` are generated grpc-gateway
streaming code with matched `defer cancel()`, not actionable.
`internal/app/app.go`'s two lifecycle goroutines (interrupt-wait,
errgroup-wait) are correctly bounded, generated, and marked DO NOT EDIT. The
IMAP client's concurrent-drain goroutines (`client.go:196,387`) are a good
pattern, explicitly documented to avoid an unbuffered-channel deadlock —
cite as a positive example, not a smell.

## Suggested next steps

1. Fix the missing `//go:build e2e` tags on `tests/livesync/livesync_test.go`
   and `tests/vault/vault_test.go` — active bug affecting `go test ./...`.
2. Decide on the tract dispatch goroutine fan-out (recover + bound
   concurrency) before it's hit by a real webhook burst.
3. Deduplicate `AddAnthropicConnection`/`AddOpenAIConnection` and the
   per-entity instance-registry CRUD (couchinstances/dockerhosts/s3instances).
4. Add e2e coverage for vault deletion, admin services, and MoM
   imap/smtp — in that rough priority order.
