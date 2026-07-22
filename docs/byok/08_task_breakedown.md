# Stage 1 Execution Progress

Tracks actual implementation progress against [06_task_breakdown.md](06_task_breakdown.md).
Scope for this pass: Tasks 1–4 only ("most basic feature" — connection plumbing + BYOK tab
shell). Usage/metrics (Task 8) and the Call LLM tract step (Tasks 6–7) are deliberately deferred
to a later pass, per explicit instruction to ship the basic feature first.

Each task is executed by one sub-agent, one at a time (later tasks depend on earlier ones —
proto/domain changes from Task 1 are needed before Task 3's service methods and Task 4's
frontend codegen). Backend tasks must ship with `go test ./...` coverage before being marked done.

## Task 1 — Domain, proto, migration for LLM key connections
Status: **done**
Agent: a092bb4ba60317ea0
Notes: `ProviderAnthropic` const + `AnthropicKeyCredentials` (domain), `anthropicConnectionMeta`
(service, unwired), `AddAnthropicConnection`/`CheckAnthropicConnection` proto + RPCs,
`EXTERNAL_PROVIDER_ANTHROPIC` enum value. No new oneof branch on `ExternalConnectionInfo.details`
— anthropic falls through to `genericDetails` like gitlab/trello. No migration needed (partial
unique index from 045 already covers it). `moti g` + `bun gen` ran clean. `go build/vet/test` and
`bun run build` all pass.

## Task 2 — `internal/clients/anthropic` package
Status: **done**
Agent: a1325f13d5ef8bbef
Notes: Added `github.com/anthropics/anthropic-sdk-go v1.58.0`. `Client{New, ListModels, Complete}`,
`ModelInfo`, `CompleteRequest`/`CompleteResult`/`Usage` (real SDK field names only). Tests use
`httptest.NewServer` (no live calls): ListModels success/auth-failure, Complete success/default
max-tokens. `go build/vet/test` + lint + gofmt all clean.

## Task 3 — Connection service: Add/Check for LLM keys
Status: **done**
Agent: ab7823d955834d35f
Notes: `Service.AddAnthropicConnection`/`CheckAnthropicConnection` + `validateAnthropicKey` (using
Task 2's client), `storedAnthropicApiKey` (blank-key-keeps-current, mirrors
`storedEmailPassword`), `recommendedDefaultAnthropicModel` (first/newest model — no tier field
on `ModelInfo` to do better yet, noted as TODO). New `user_errors.LlmKeyValidationFailed` /
`LlmKeyRequired`. Transport handlers `add_anthropic.go`/`check_anthropic.go` mirror the gitlab
pair. `to_proto.go` provider map extended (still falls through to `genericDetails`, no new oneof
branch). Tests cover `validateAnthropicKey` success/401 + the helper functions via httptest, no
repo mock needed. `go build/vet/test` + `make lint` all clean.

## Task 4 — Frontend: BYOK tab shell
Status: **done**
Agent: aca530b7f3c9c48f7
Notes: New `components/atoms/Tabs` (tier-2, genuine chures gap). `ConnectionsPage.tsx` now owns
`?tab=` state and renders `Tabs` under `HeroSegment`, switching between unchanged `ContentSegment`
and new `pages/connections/components/BYOKSection`. `BYOKSection` renders the Anthropic
`ProviderCard` (no `onClick` yet — `ManageLlmKeyDialog` is a follow-up, not yet started) plus 3
`ComingSoonCard` placeholders (CouchDB/S3/WebDAV) using the *existing* global `react-tooltip`
instance (`data-tooltip-id="root-tooltip"`) — corrected a stale design-doc assumption that a new
Tooltip atom was needed; it wasn't. `processes/ExternalConnections.ts` +
`app/hooks/ExternalConnections.ts` gained `addAnthropicConnection`/`checkAnthropicConnection`,
unused until the manage dialog lands. Post-review fix: the agent's OAuth-callback `?status=`
cleanup regressed to setting `status: ""` instead of deleting it (would've left a stray query
param and dropped `?tab=`) — corrected to a functional `setSearchParams` that deletes just
`status`. `bun run build` + `bun run lint` clean after the fix.

## Overall status

Tasks 1–4 done. End-to-end result: a user can now add/verify an Anthropic API key via
`AddAnthropicConnection`/`CheckAnthropicConnection` (backend-complete, no UI form wired to it
yet), and the Connections page has a working BYOK tab showing the connection's state. **Not yet
usable end-to-end from the UI** — there's no dialog to actually type in a key (Task 5,
`ManageLlmKeyDialog`, not started). Tasks 5–9 (manage dialog, tract `llm_call` step, usage
capture, docs) remain per [06_task_breakdown.md](06_task_breakdown.md), not requested yet.
