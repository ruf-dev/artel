# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**artel** is a Go service that automates Obsidian in-cloud vault creation. It provisions CouchDB instances, creates
users and resources within them, and exposes a React UI at `pkg/client/ArtelUI` (Bun + React). Generated
by [RedSock CLI (verv)](https://github.com/Red-Sock/verv).

### Frontend (pkg/client/ArtelUI)

### Code Generation

```bash
moti i          # Install protoc dependencies if new added to moti.yaml
moti g          # Generate protos for Go backend (internal/api/server/artel_api)
bun gen         # Generate TypeScript proto clients for frontend (run from pkg/client/ArtelUI)
sqlc generate   # Generate Go from SQL queries (internal/repository/pg/queries → internal/repository/pg/generated)
make lint       # golangci-lint run ./...
```

## Commit Message Convention

Commit messages must follow `[Area] Category: description`:

- `Area` — the part of the system changed, e.g. `Tract`, `Connections`, `Auth`. Free-form, pick per commit.
- `Category` — either a fixed keyword (`bug fixed`, `refactor`, `docs`, `chore`, `perf`), or for feature work, a
  Capitalized sub-area/feature name instead of the literal word "feature" (e.g. `Canvas`, `Connections`).

Examples:

```
[Tract] bug fixed: connections created on the caller side, not the receiver
[Tract] Canvas: added drag-select for nodes
```

This is enforced automatically by a `PreToolUse` hook in `.claude/settings.json` (`.claude/hooks/validate-commit-msg.py`),
which blocks non-conforming `git commit` calls. Since `.claude/settings.json` is committed to the repo, no per-clone
activation step is needed.

## Architecture

If touching Go backend files, read [docs/architecture.md](docs/architecture.md) for verv structure, configuration,
planned layers, service layout, and key dependencies.

## Testing

- Unit tests: `go test ./...` from the repo root — no external services required.
- E2E tests live under `tests/` (`tests/e2e`, `tests/livesync`, `tests/vault`, `tests/tract_e2e`,
  `tests/gitlab_trigger_e2e`) and are gated behind the `e2e` build tag, so they're excluded from
  `go test ./...` by default.

`make test-e2e` is the standard entrypoint: it brings up the test backends (Postgres, CouchDB, MinIO
on ports 15434/15985/19000, left running afterward so repeated runs don't pay container startup
cost), runs the one-time bootstrap setup (migrations, CouchDB/S3 admin-pool instance registration,
system_settings setup — see `tests/bootstrap`), runs the e2e suites, then always runs bootstrap
cleanup after — even on suite failure — leaving the shared stack empty for the next run:

```bash
make test-e2e
```

To scope to one suite during iteration, bring the stack up once, then run the suite directly
against it (the bootstrap step must have run at least once first so the suites' `harness.GetCouchInstance`/
`GetS3Instance` lookups find a registered pool instance):

```bash
docker compose -f tests/docker-compose.yaml up -d --wait
go test -tags "e2e e2e_bootstrap" -count=1 ./tests/bootstrap/... -run TestEnvSetup
go test -tags e2e ./tests/e2e/...
```

Suites default to `localhost` + the ports above; override via `PG_DSN`, `COUCH_URL`/`COUCH_USER`/`COUCH_PASS`,
and `S3_ENDPOINT` env vars if pointing at different infra.

Each e2e suite constructs its own `svcv1.Services` via `config.EnvironmentConfig{}` — note that
`SubscriptionsEnabled` defaults to `false` there, which routes all quota/feature checks through
`FreeService` (always-allow, no CouchDB/S3 calls). A suite that needs to exercise real plan/quota
enforcement (`PaidService`) must set `cfg.SubscriptionsEnabled = true` explicitly (see
`tests/e2e/quota_test.go`).

`UserContext.UserName` must be non-empty when a test builds one by hand (bypassing the gRPC auth
interceptor) and then calls `Vault.CreateVault` — the CouchDB database name is derived from it, and
an empty value produces a name starting with `-`, which CouchDB rejects. Email/password registration
never populates `domain.User.Username`, so tests stand in with the email's local part.

## MoM (MCP of MCP) — third-party integrations

Artel has a declarative integration framework internally called "MoM" (MCP of MCP — unrelated to
"Mixture of Models"). Despite the name, it does NOT nest the MCP protocol or call out to other MCP
servers — Artel is an MCP server only, never an MCP client. MoM is a DB-stored *tool definition*
layer: each integration ("MoM record") is a JSON document of MCP-shaped tool schemas
(`api_description`) paired with a declarative `action` (`imap` / `smtp` / `http`) that the
backend compiles into a real client call at execution time.

Full spec and worked examples (email via imap/smtp, Trello via http with secret interpolation):
[pkg/mom_examples/README.md](pkg/mom_examples/README.md).

Key pieces:
- `internal/domain/mom.go` — `McpDefinition`, `McpToolDef`, `ToolAction` (imap/smtp/http discriminated union)
- `mcps` table (migration 027) — one row per integration, `tools` JSONB
- `mcp_connectors` table (migration 028) — links an `mcp_key` to an `external_connections` row per MoM; no secrets stored here, secrets always live in `external_connections.credentials_enc`
- `internal/service/v1/mom/mom.go` — `ExecuteToolForKey` dispatches by action type to `internal/service/v1/mcp/executors/` (`EmailExecutor`, `HttpExecutor`)

**Convention: default to MoM for new third-party integrations.** When adding a new external
service (GitLab, Slack, Jira, etc.) that exposes a plain REST API, author it as a MoM `http`-action
tool definition (seeded via migration, following the `gitlab` MoM in migration 030 or `trello` in
the README) instead of writing a bespoke Go SDK client + service package. Use `${{params.*}}` for
caller-supplied input and `__secrets.*` to pull credentials from the linked `external_connections`
row — never hardcode or duplicate credential storage per integration.

Reserve bespoke Go code (new client packages under `internal/clients/`, new transport handlers)
for things that are genuinely **not** an outbound declarative HTTP call:
- inbound webhook ingestion + signature/token verification (e.g. `internal/transport/gitlab_webhook/`)
- multi-step orchestration / output-to-input chaining across tools (the "Workflow" layer — not yet
  implemented; if you find yourself wanting to chain MoM tool calls, that's the intended future
  extension point — extend MoM/add a Workflow layer on top of it, do not introduce a second
  competing tool-execution framework)
- protocols other than plain JSON/REST that don't fit `http`/`imap`/`smtp` (rare — discuss before adding a new action type)

## Go Coding Rules

If editing any `.go` files, read and follow [docs/go-style.md](docs/go-style.md).

## Frontend Coding Rules

If editing any files under `pkg/client/ArtelUI`, read and follow
[pkg/client/ArtelUI/CLAUDE.md](pkg/client/ArtelUI/CLAUDE.md) — it covers the
Feature-Sliced Design layering (pages/widgets/components), when to extract a
component vs. keep it inline, CSS Modules conventions, component structure rules,
and error/confirmation handling primitives.

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

Rules:

- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use
  `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a
  scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough
  context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).
- `graphify-out/` is gitignored and local-only (not committed) — it grew unbounded as daily-dated snapshots and pushed the
  repo over Go's 500MB module-zip limit when imported as a package elsewhere. Each clone/CI run regenerates it via
  `graphify update .`; there is no cross-machine sharing of the graph.
## Client side changes

When change frontend code read pkg/client/ArtelUI/CLAUDE.md before 