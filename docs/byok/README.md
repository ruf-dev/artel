# Bring Your Own Key (BYOK) — Design Plan

Status: **planning only, nothing in this folder is implemented yet.**

## Goal

Let a user hand Artel their own LLM API key (Claude format first, OpenAI format later) so that
Tracts can call an LLM directly — without Artel operating a shared/metered LLM account. Reuse
the existing `external_connections` secret-storage machinery rather than inventing a parallel
credentials system.

Two other things ride along because they were requested in the same breath:

1. A restructuring of the Connections page: today it's one flat grid of provider cards. It
   becomes two tabs — **External Connections** (today's Google Sheets/Trello/Miro/Email/GitLab
   grid, unchanged) and **BYOK** (new: LLM keys now, mock/"coming soon" cards for CouchDB/S3/WebDAV
   connections later).
2. A new Tract step, **Call LLM**, that is the actual consumer of a BYOK connection: prompt +
   system prompt + model, with access to earlier steps' outputs via the existing template
   resolver, executed through the engine the same way `action`/`script`/`condition` steps are.

## Documents in this folder

| File | Covers |
|---|---|
| [01_data_model_and_security.md](01_data_model_and_security.md) | How/where the key is stored, encryption, what's plaintext vs encrypted, provider enum additions |
| [02_connection_lifecycle.md](02_connection_lifecycle.md) | Add / verify / list / rotate / delete flows, vendor detection, RPC contracts |
| [03_frontend_connections_page.md](03_frontend_connections_page.md) | Tabs, BYOK section, mock cards, dialogs, component tiering |
| [04_tract_llm_step.md](04_tract_llm_step.md) | New `llm_call` step: domain, engine, resolver integration, frontend step editor |
| [05_metrics_and_usage.md](05_metrics_and_usage.md) | What the Anthropic API actually reports, how usage is persisted and surfaced |
| [06_task_breakdown.md](06_task_breakdown.md) | Ordered, file-level implementation tasks referencing the above |

## Key architectural decisions (summary — see linked docs for the "why")

- **Storage**: LLM keys are `external_connections` rows like every other provider. No new table
  for the keys themselves. `provider = "anthropic"` (stage 1), `"openai"` (stage 2, enum reserved
  now). Same `credentials_enc` AES encryption via `cryptoutil`, same repo/service shape as
  GitLab/Trello.
- **Categorization is a frontend concern.** The backend doesn't need an `is_byok` flag — the
  Connections page simply renders two hardcoded provider lists (today's `PROVIDERS` array plus a
  new `BYOK_PROVIDERS` array), exactly like `ContentSegment.tsx` already does for the External
  Connections tab.
- **Vendor detection** is a UX affordance, not a security boundary: the tab the user picks fixes
  the wire protocol (Anthropic Messages API shape vs. OpenAI Chat Completions shape); the actual
  vendor behind a given `base_url` is confirmed by a live, zero-token `GET /v1/models` call at
  Add/Check time, which doubles as key validation and populates the model picker.
- **"Call LLM" is a bespoke engine step, not a MoM tool.** MoM is for outbound HTTP calls exposed
  to an *agent* as a callable tool; the LLM-call step is a first-class Tract primitive analogous
  to `ScriptStep` — hardcoded step kind, typed engine logic, executed via the official
  `anthropic-sdk-go` under a new `internal/clients/anthropic` package (matches the existing
  `internal/clients/{googleapi,imap,smtp}` convention).
- **Metrics come from `response.usage`, never estimated.** Anthropic's Messages API returns
  `input_tokens` / `output_tokens` / `cache_creation_input_tokens` / `cache_read_input_tokens` on
  every response — that's the only source of truth recorded. No client-side token counting.
