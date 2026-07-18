# Structured/Semi-Structured Data Storage — Discussion Notes

Idea doc, 2026-07-18. Captures a design discussion, not an implementation plan — no code was
written. Useful as a starting point if this direction gets picked up later.

## Problem

Prompted by a concrete case: a user's wife logs cosmetic product test results (ingredients, a
good/ok/bad verdict, notes) as one large free-text Obsidian note. Every time an agent needs to
answer "will this new product suit her," it has to `read_file` the whole note and reason over raw
text client-side. Today's vault MCP surface
(`internal/service/v1/mcp/executors/vault.go`) has no query/aggregation primitive — `list_files`,
`read_file`, `list_tags`, etc. all operate on whole files; there's nothing that lets an agent ask
for a scoped, pre-aggregated answer (e.g. "ingredients grouped by verdict") without pulling
everything into context first.

## Two ways to fix the immediate case

1. **Frontmatter, no new storage.** Split the note into one note per tested item, move the
   structured fields (verdict tag, ingredient list, date) into Obsidian YAML frontmatter, add one
   new builtin tool — `query_notes(folder, filter, group_by)` — that scans frontmatter across notes
   server-side and returns only the aggregate. No new datastore; reuses `list_tags`-style
   vault-wide scanning that already exists. Cheapest option, keeps the data fully Obsidian-native
   (still editable/visible in the Obsidian app). Weakness: CouchDB has no secondary index over
   frontmatter fields, so every query is a full vault scan — fine at hundreds of notes, not at
   thousands.

2. **A real "collections" primitive.** A dedicated feature: `collections` (per-vault, declares
   which fields exist) + `collection_items` (normalized columns for the dimensions actually queried
   on — item name, verdict tag, created_at — plus one JSONB column for the genuinely variable leaf
   data, e.g. the ingredient list), with a GIN index on the JSONB column, and new builtin MCP tools
   (`create_collection`, `add_item`, `query_items` with filter + group-by), analogous to
   `VaultExecutor`. This is a bigger, more general feature (per-vault structured data, not just this
   one use case) and would need its own spec before building.

No decision was made between these; (1) is the cheaper validation step, (2) is the "do it properly"
version if the pattern proves useful beyond one household.

## Should this use MongoDB instead of Postgres?

Raised as a question: since the data is JSON-shaped, wouldn't Mongo be a more natural fit than
Postgres JSONB?

**Conclusion: no, stay on Postgres.** Reasoning:

- Nothing in the repo uses Mongo today — `internal/config/data_sources.go` declares exactly one
  data source (`Postgres`). Adding Mongo means a brand-new engine to provision, secure, back up,
  and integrate a driver for for a workload this small.
- Postgres JSONB (+ GIN index) already gives "loose schema, fast filter" — the actual property
  wanted here — without a second engine. A `data JSONB` column doesn't require predeclaring
  columns any more than a Mongo document does.
- Staying on Postgres means the collections table can `JOIN` against `vaults`/`users` for auth
  scoping the same way notes/tracts already do. Split into Mongo and every ownership check becomes
  an app-level cross-database join instead of SQL — strictly worse.
- Data volume here (one household's test log) is nowhere near where Mongo's horizontal-scale story
  would start paying for itself.
- Mongo's one real edge — its aggregation pipeline reads more naturally than `jsonb_path_query`/
  `jsonb_agg` for deep nested group-bys — is an ergonomics difference, not a capability gap.

### Correction to an "N instances" framing raised mid-discussion

A counter-argument was floated: "no user data lives in Artel's own Postgres anyway — it's already N
instances, same as CouchDB" — implying Mongo would just be another per-vault instance, same
pattern. Checked against the actual code
(`internal/service/v1/couchinstances/couchinstances.go`, `internal/domain/couch_instance.go`): this
is wrong. CouchDB is a small **pool** of admin-registered servers (`CouchInstance`), each hosting
*many* vaults as separate databases (`CouchAccount` = one vault's credentials/database within a
pooled server) — not one instance per vault. And CouchDB isn't a free technology choice at all: it's
forced by the Obsidian Self-hosted LiveSync plugin, which only speaks CouchDB's replication
protocol. Nothing forces a parallel constraint toward Mongo for a new feature.

The underlying instinct was still right, just aimed at the wrong lever: user data (test results)
shouldn't be commingled in the *control-plane* Postgres alongside `users`/`subscriptions`/`vaults`
metadata — it should live in an isolated, poolable *data-plane* store, same as vault content does.
That argues for **mirroring the `CouchInstance`/`CouchAccount` pattern with Postgres** (a small pool
of Postgres hosts, one database/schema per vault, dropped when the vault is deleted) — reusing
pgx/sqlc/migrations/the existing repository pattern — not for introducing a second engine with no
existing tooling in this codebase.

## What else could Mongo be for? (broader brainstorm, not scoped to the notes case)

Three different needs were pulled apart, since "analytics over huge data for agents to scrape over
cases" conflates them:

1. **MoM/tool execution log — the actual best-fit candidate.** Nothing like this exists yet:
   `domain.ToolExecResult` (`internal/domain/tool_exec_result.go`) is a transient in-memory return
   value, never persisted. But if Artel ever started recording every MoM/vault/tract tool call's
   raw input/output/timing/error for audit and "why did this tract fail" debugging, that corpus is
   *genuinely* schema-varying per source — every third-party integration's request/response shape
   is different by design (`domain.HttpAction` is intentionally open-ended), so normalizing it in
   Postgres would mean a table per integration. A single schema-free collection is a legitimate fit
   here, more so than the notes/collections case, because "we don't know the shape and don't want
   to" is actually true for this workload.
2. **Aggregate analytics** ("which tools fail most," "typical tract shapes across users") — Mongo is
   not the right tool despite being document-shaped; it's an OLTP store, not columnar. A partitioned
   Postgres table or an OLAP engine (ClickHouse/DuckDB) would outperform it once volume is real.
3. **"Agent scrapes over cases" as semantic/similarity search** ("find past cases like this one") —
   this is a vector-search problem, not a document-flexibility problem. `pgvector` on Postgres
   covers it without a new engine (Mongo Atlas has vector search too, but that's not a reason to
   adopt Mongo generally).

**If this gets pursued, the recommended scope is (1) only** — a narrowly-framed MoM/tract
execution/audit log, not a general-purpose Mongo adoption.

## Open questions / not decided

- Whether to build the frontmatter/`query_notes` route or the full `collections` primitive first —
  no commitment made either way.
- Whether a per-vault Postgres pool (mirroring `CouchInstance`) is worth building now, or whether
  collections should just live in the existing control-plane Postgres until there's a concrete
  reason to isolate it.
- Whether the MoM/tract execution log is worth building at all yet — no current pain point drives
  it (no user has asked for it), unlike the notes case which has a real, named problem behind it.
