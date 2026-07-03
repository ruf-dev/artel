# Tract — Implementation Risks

Companion to the Tract backend plan (v4, 2026-07-03). Each risk: what it is, why it matters,
how it's mitigated, and what to watch for.

## 1. JSONB definition has no referential integrity

**What:** The step tree lives in `tracts.definition` JSONB. Postgres cannot enforce foreign
keys from inside a JSON document, so a definition can reference an `mcp_tools` row, a builtin
tool name, or an `external_connections` row that later disappears or changes shape.

**Why it matters:** A tract that validated at save time can silently rot: a re-seeded tool
rename, a deleted connection, or a dropped MoM record turns a working tract into one that fails
at run time — and for webhook-fired tracts nobody is watching when it fails.

**Mitigation:**
- Validation runs on **every** create/update *and again at run start* (tool existence,
  connection ownership) — a rotted tract fails fast with a precise error in `tract_runs.error`
  instead of half-executing.
- Seed migrations must **upsert** tools by `(mcp_name, name)`, never delete/recreate, so tool
  identity is stable across releases.
- Run history is immune: `tract_run_steps` snapshots `step_id`/`step_name`/`step_type` and
  stores concrete rendered inputs/outputs, so old runs stay readable no matter what happens to
  the definition or the tools.

**Watch for:** any future admin/UI path that deletes `mcp_tools` rows or external connections —
it should warn when tracts reference them (query: `definition::text LIKE '%"tool":"<name>"%'`
is crude but sufficient for a guard).

## 2. Tree model cannot express arbitrary DAGs (canvas mismatch)

**What:** The designer handoff shows a free-form node canvas with drawn connection lines. The
backend models execution as a *tree*: sequences, then/else branches, parallel groups (with
`group` for sequential lanes). Fan-in is only implicit — "the step after a parallel group runs
when all lanes finish."

**Why it matters:** A canvas that lets users draw any edge (e.g. two independent chains merging
mid-flow, diamond joins with partial synchronization) will produce graphs the backend cannot
store or execute. This is a product-level constraint, not a bug.

**Mitigation:**
- The design's own example graph (trigger → lookup → {email ∥ telegram} → update) maps exactly
  to the tree primitives, so the shipped design is fully supported.
- The future canvas UI must *compose* primitives (append step, insert branch, add parallel
  lane) rather than offer free edge drawing; connection lines are then a rendering of the tree,
  and `ui:{x,y}` passthrough metadata preserves node placement.
- If real DAG demand appears later, the JSONB definition can grow a new step type or an edges
  representation without a schema migration — but the engine and validation would need real
  work (topological scheduling, cycle detection). Treat as a separate project.

**Watch for:** frontend tickets that say "let the user connect any two nodes" — push back or
escalate scope before implementation.

## 3. Declared output schemas can drift from real tool outputs

**What:** `mcp_tools.output_schema` (and builtin tools' Go-declared schemas) are *documentation
used for UI hints and agent authoring* — nothing enforces that GitLab/Trello/Telegram actually
return those fields, and upstream APIs change.

**Why it matters:** The whole click-first UX rests on trusting the hint dropdown. If
`create_merge_request` stops returning `url`, every tract binding `{{ create_pr.url }}` renders
an unresolved-reference error at run time, and users authored those bindings in good faith.

**Mitigation:**
- Run history stores **actual** parsed outputs per step; the UI's "Last output" panel (in the
  design) makes truth visible next to the promise, which is how users and agents debug drift.
- Unresolved template refs fail the step with an explicit "reference `create_pr.url` not found
  in output" error rather than passing empty strings downstream.
- Schemas live in migrations/Go, so fixing drift is a normal, reviewable change.

**Watch for:** recurring step failures with reference-not-found errors on a specific tool — the
signal that its schema needs a refresh.

## 4. MoM refactor touches the live MCP wire format

**What:** Part 1 rewrites how tool definitions are stored (`mcps.tools` array → `mcp_tools`
rows), makes `ToolProperty` recursive, and extends the MCP projection (`momToolToToolDef`) with
nested schemas, `enum` preservation, and `outputSchema`. Existing MCP clients (Claude
connections in production) consume this wire format today.

**Why it matters:** A subtle marshalling change (field renamed, `{}` vs `null` schema, enum
emitted where a client expects absence) breaks *current* users of the email/gitlab tools —
before Tract ships anything.

**Mitigation:**
- Phase 1 is explicitly gated on behavior parity: `tools/list` and `tools/call` against
  existing tools must produce byte-equivalent (or semantically identical, verified by diff)
  responses before Phase 2 starts.
- The migration's data move is mechanical (`jsonb_array_elements` lift) and runs in one
  transaction with the `DROP COLUMN`; test it against a copy of real data, not just fixtures.
- `output_schema = '{}'` means "undeclared" and must be omitted from the wire, not sent as an
  empty object — cheap to get wrong, called out in the implementation notes.

**Watch for:** MCP clients erroring on tool listing right after deploy; keep the migration's
inverse (recreate `tools` array from rows) noted in the migration header as the rollback path.

## 5. Multi-trigger tracts have ambiguous input shapes

**What:** One tract can be linked to several triggers whose `payload_schema`s differ. Template
refs like `{{ trigger.branch }}` may exist in one linked trigger's payload and not another's.

**Why it matters:** The tract runs fine from trigger A and fails from trigger B — confusing,
and only discoverable at run time.

**Mitigation:**
- Validation cross-checks `trigger.*` refs against **all** linked triggers' schemas and returns
  warnings (not errors) naming which trigger lacks which field — surfaced at save time and at
  link time.
- At run time a missing field is a hard, well-labeled step failure recorded in run history.
- Per-link filters let users route payload variants: link the same webhook twice with filters
  that guarantee the shape each tract expects.

**Watch for:** users linking `generic`-source triggers (user-declared schemas) — schema honesty
is on them; the warning mechanism is the only guard.

## 6. In-process engine: runs die with the process

**What:** Runs execute in goroutines inside the server process (no queue, no worker pool, no
resume). Deploy/restart/crash kills in-flight runs. This is a deliberate v1 simplification.

**Why it matters:** A run killed between "create MR" and "notify tracker" leaves external
side effects half-applied. Unlike Temporal, v1 does not resume.

**Mitigation (Temporal-inspired, designed in now so resume is cheap later):**
- Persist-before-apply: the `tract_run_steps` row is inserted (`running`, with rendered input)
  *before* the tool call, so post-crash state shows exactly which step was in flight.
- Outputs are recorded append-only keyed by immutable step ids — a future retry-from-failed can
  replay completed steps from recorded outputs and execute only the remainder, with no
  determinism concerns (the definition is data, not code).
- Startup sweep marks stale `running` runs/steps as `failed` so nothing is stuck-forever in the
  UI, and graceful shutdown waits briefly (server-lifecycle context) before hard exit.
- Actions are outbound API calls; most reference-flow tools are safe-ish to re-fire manually.
  True idempotency keys are a v2 concern.

**Watch for:** long-running tools + frequent deploys = frequent half-runs. If that bites,
prioritize retry-from-failed (the schema already supports it) before reaching for a queue.

## 7. Parallel execution shares mutable run state

**What:** Parallel lanes run under an `errgroup`, all resolving templates against a shared
`stepId → output` map and appending run-step rows concurrently.

**Why it matters:** Data races on the map (Go race detector failures at best, corrupted
rendering at worst); interleaved failures must cancel siblings without leaking goroutines.

**Mitigation:**
- Validation guarantees lanes cannot reference each other, so cross-lane reads never occur
  mid-flight; the map is still mutex-guarded because all lanes *write* to it.
- Run-step rows are independent inserts (no shared row updates), ordering reconstructed by
  `started_at`.
- Engine tests run with `-race`, including a failing-lane-cancels-siblings case.

**Watch for:** any future feature letting lanes communicate (shared counters, early-exit
races) — that's where this design stops being sufficient.

## 8. Webhook endpoint is unauthenticated-by-design internet surface

**What:** `POST /tract/hook/{trigger_uuid}` accepts unauthenticated traffic, gated only by an
unguessable UUID + a shared-secret token (`X-Tract-Token` / GitLab's `X-Gitlab-Token`).

**Why it matters:** Forged payloads flow straight into template resolution and out through the
user's connected credentials (create MRs, send comments). Payload spoofing = actions executed
with the victim's tokens.

**Mitigation:**
- Constant-time token comparison against `sha256` at rest (token itself never stored), 1MB body
  cap, immediate 200 with async execution (no timing oracle on filter matching), disabled
  trigger/tract short-circuits.
- Fan-out is bounded by explicit `trigger_tracts` links owned by the same user; a trigger can
  never start another user's tract.
- The template resolver is a data-only path lookup — no code evaluation, no reflection — so a
  hostile payload can at worst parameterize the actions the owner already configured.
- Rate limiting is deliberately deferred; the token gate plus per-user link ownership bounds
  the blast radius to the token holder's own tracts.

**Watch for:** tokens pasted into third-party services (GitLab webhook config) leak via those
services' admin UIs — token regeneration should be an early follow-up feature.

## 9. Template resolver is a new hand-rolled parser

**What:** `{{ trigger.branch }}` / `{{ step[0].field }}` / `{{ $now }}` / `{{ length(x) }}`
gets a custom ~100-line resolver (chosen over `text/template` for TS mirrorability and safety).

**Why it matters:** Hand-rolled parsers attract edge-case bugs: nested braces in literal JSON
bodies, `[N]` on non-arrays, unicode in field names, `{{` inside a resolved string value
(no re-resolution must occur), type coercion of single-token params.

**Mitigation:**
- Single-pass, no re-entry: resolved values are never re-scanned for tokens (prevents payload-
  driven template injection).
- The grammar is deliberately tiny and closed (two namespaces, one function, index access);
  anything else is a parse error at validation time, not silently ignored at run time.
- It's the primary unit-test surface: table-driven tests covering type preservation,
  stringification, missing refs, malformed tokens, `$`-vars, and the MoM `${{params.*}}`
  non-collision (tract renders first; the literal `${{...}}` sequences must survive untouched).

**Watch for:** the future TS mirror drifting from the Go implementation — share the test-case
table (JSON fixtures) between both when the frontend lands.

## 10. Migration 032 is destructive

**What:** `ALTER TABLE mcps DROP COLUMN tools` after the data lift into `mcp_tools`.

**Why it matters:** If the lift is wrong (JSON path typo, tools with missing
`api_description.name` producing NULLs), the source data is gone after deploy.

**Mitigation:** run the INSERT and a row-count + spot-check comparison (`SELECT count(*)` vs
`sum(jsonb_array_length(tools))`) inside the same transaction before the DROP; keep a pre-
migration dump; the migration header documents the reverse transformation.

**Watch for:** environments with hand-edited `mcps.tools` rows that don't match the expected
shape — the `COALESCE` defaults mask absence; a NOT NULL violation on `name` is the desired
loud failure, don't soften it.
