# Metrics & Usage

## What the Anthropic API actually gives us

Every `POST /v1/messages` response carries a `usage` object — this is the *only* source of token
counts; nothing is estimated client-side (no `tiktoken`-style guessing, which is wrong for
Claude's tokenizer anyway):

| Field | Meaning |
|---|---|
| `input_tokens` | Tokens processed at full price (prompt + system prompt, minus anything served from cache) |
| `output_tokens` | Tokens generated in the completion |
| `cache_creation_input_tokens` | Tokens written to prompt cache this request (not relevant to stage 1 — Call LLM doesn't use caching, see below — included for completeness/forward-compat) |
| `cache_read_input_tokens` | Tokens served from cache (same caveat) |

"Words" isn't a unit the API reports — token counts are the only thing worth persisting.
Anything word-based (e.g. "average words per response" in a dashboard) should be derived
client-side from `output_tokens`/response length only as a rough display heuristic, never stored
as if it were an API-reported metric.

Stage 1's `Complete` call (single prompt/response, no multi-turn, no explicit `cache_control`
breakpoints — see [04_tract_llm_step.md](04_tract_llm_step.md)) won't produce meaningful cache
figures; they'll read zero. That's fine — the schema below still has the columns so a future
stage that adds prompt caching to repeated Call LLM invocations (e.g. a shared system prompt
across steps) doesn't need a migration.

## Persisting usage

Two things already capture *part* of this for free, and one new table is needed for anything
beyond a single run.

### Already covered: per-step output

`executeLlmCall` (see [04_tract_llm_step.md](04_tract_llm_step.md)) writes the full usage object
into the step's `output` JSON, which lands in `tract_run_steps.output` via the existing
`UpdateRunStepFinish` call — no extra work. This is enough to show token counts for one run in
`TractCanvasLogPanel.tsx` (`TractRunStepItem.output` already flows to the frontend; the log panel
just needs to render the `usage` sub-object it wasn't expecting to see before, e.g. as a small
"1,204 in · 340 out" chip on the step's log entry).

This is **not** enough for cross-run aggregation ("how many tokens has this BYOK key used this
month") — `tract_run_steps` rows aren't indexed by connection, and nothing guarantees they're
retained indefinitely (a future run-history retention/pruning policy would silently lose the
data). Aggregation needs its own table.

### New: `llm_usage_events` table

```sql
-- migration N_llm_usage_events.sql
CREATE TABLE llm_usage_events (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    external_connection_id      UUID NOT NULL REFERENCES external_connections(id) ON DELETE CASCADE,
    tract_id                    UUID REFERENCES tracts(id) ON DELETE SET NULL, -- NULL-able: tract may be deleted later
    run_step_id                 UUID, -- best-effort pointer into tract_run_steps, not FK'd (that table's lifecycle may differ)
    provider                    TEXT NOT NULL, -- "anthropic" | "openai"
    model                       TEXT NOT NULL,
    input_tokens                BIGINT NOT NULL DEFAULT 0,
    output_tokens                BIGINT NOT NULL DEFAULT 0,
    cache_creation_input_tokens  BIGINT NOT NULL DEFAULT 0,
    cache_read_input_tokens      BIGINT NOT NULL DEFAULT 0,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON llm_usage_events(external_connection_id, created_at);
CREATE INDEX ON llm_usage_events(user_id, created_at);
```

Written by `executeLlmCall` right after a successful call, via a new narrow repo
(`internal/repository/pg/repos/llmusage/`) — **best-effort, never fails the tract run**: if the
insert errors, log it (same `log.Error().Err(err)...` pattern used elsewhere in `engine.go` for
non-fatal bookkeeping failures like `UpdateRunStepFinish` errors on the failure path) and
continue; a lost usage row is a minor analytics gap, not a reason to fail a user's automation.

### Cost estimation — explicitly a "nice to have," not stage 1

Turning tokens into a dollar estimate needs a price table per model, and pricing changes over
time independent of code deploys. Two options, worth deciding at implementation time rather than
now:

- **Static Go map** (`internal/service/v1/tract/llm_pricing.go`), reviewed/updated by hand when
  prices change. Simple, but drifts silently if nobody remembers to update it.
- **DB table** (`llm_model_pricing(model TEXT PRIMARY KEY, input_price_per_mtok NUMERIC,
  output_price_per_mtok NUMERIC, updated_at)`), seeded via migration, editable without a deploy.
  More durable, small amount of extra plumbing (an admin-only update path).

Recommendation: **skip cost estimation in stage 1 entirely** — surface raw token counts only.
Add pricing once real usage data shows people actually want a dollar figure, and prefer the DB
table then (it's the option that doesn't need a redeploy every time Anthropic changes pricing).

## Surfacing usage

- **Per-connection summary** on the BYOK connection card / `ManageLlmKeyDialog`'s connected view
  (see [03_frontend_connections_page.md](03_frontend_connections_page.md)): "12,400 tokens · 34
  calls in the last 30 days," backed by a new `GetLlmUsageSummary` RPC —
  `SELECT count(*), sum(input_tokens), sum(output_tokens) FROM llm_usage_events WHERE
  external_connection_id = $1 AND created_at > now() - interval '30 days'`. One query, no
  aggregation service needed beyond that.
- **Per-run-step detail** in `TractCanvasLogPanel.tsx`, from the existing `output.usage` JSON —
  no new RPC, just a rendering change.
- **Aggregate dashboard across all keys/tracts** — explicitly out of scope for stage 1. If wanted
  later, it's a straightforward extension of the same table (group by `provider`/`model`/`tract_id`
  instead of just `external_connection_id`), not a redesign.
