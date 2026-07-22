# Task Breakdown

Ordered so each task is independently reviewable/mergeable and later tasks build on earlier ones.
File paths reference the current codebase layout. None of this is implemented yet — this is the
execution order for whoever picks it up next.

## Stage 1 — Anthropic BYOK + Call LLM step

### Task 1 — Domain, proto, migration for LLM key connections
_Ref: [01_data_model_and_security.md](01_data_model_and_security.md), [02_connection_lifecycle.md](02_connection_lifecycle.md)_

- [ ] `internal/domain/external_connection.go`: add `ProviderAnthropic`/`ProviderOpenAI`
      constants, `AnthropicKeyCredentials` struct.
- [ ] `internal/domain/external_connection.go` (or wherever `gitlabConnectionMeta`-style structs
      live — check current file, likely `internal/service/v1/externalconnections/`):
      `llmKeyConnectionMeta` struct.
- [ ] `api/grpc/external_connections.proto`: `EXTERNAL_PROVIDER_ANTHROPIC`/`_OPENAI` enum values,
      `LlmKeyConnectionInfo` message + `details` oneof branch, `AddLlmKeyConnection`/
      `CheckLlmKeyConnection` request/response messages + RPCs.
- [ ] `moti g` (regen Go server stubs) + `bun gen` from `pkg/client/ArtelUI` (regen TS clients) —
      per repo convention, never hand-edit `.pb.go`/`.pb.ts` files.
- [ ] No DB migration needed for the connection itself (reuses `external_connections` as-is —
      confirm the partial unique index from migration `045` already covers the new providers
      correctly, i.e. they fall under the "single connection" branch).

### Task 2 — `internal/clients/anthropic` package
_Ref: [04_tract_llm_step.md](04_tract_llm_step.md)_

- [ ] Add `github.com/anthropics/anthropic-sdk-go` dependency.
- [ ] `internal/clients/anthropic/client.go`: `New(apiKey, baseUrl string) *Client`,
      `ListModels(ctx) ([]ModelInfo, error)`, `Complete(ctx, CompleteRequest) (CompleteResult,
      error)`. Use `claude-opus-4-8` as the recommended-default model when a connection's
      verification response doesn't obviously suggest otherwise — don't hardcode it as the *only*
      option, the model list comes from the live API.
- [ ] Unit tests against a fake HTTP transport (same style as `mcp/executors/http_test.go`) — no
      live API calls in `go test ./...`.

### Task 3 — Connection service: Add/Check for LLM keys
_Ref: [02_connection_lifecycle.md](02_connection_lifecycle.md)_

- [ ] `internal/service/v1/externalconnections/external_connections.go`:
      `AddLlmKeyConnection`/`CheckLlmKeyConnection`, using the Task 2 client.
- [ ] New `user_errors.LlmKeyValidationFailed`.
- [ ] `internal/transport/external_connections_api/`: wire the two new RPCs (gRPC handler +
      gateway route), mirroring the GitLab handler pair.
- [ ] Update `pkg/mom_examples/README.md`/`docs/mcp/*` **only if** they document the full
      provider list somewhere stale — check before editing, don't assume.

### Task 4 — Frontend: BYOK tab shell
_Ref: [03_frontend_connections_page.md](03_frontend_connections_page.md)_

- [ ] `components/atoms/Tabs/Tabs.tsx` (new tier-2 atom).
- [ ] `components/atoms/Tooltip/Tooltip.tsx` (new tier-2 atom, portaled per the no-`z-index` rule).
- [ ] `pages/connections/ConnectionsPage.tsx`: own tab state, render `Tabs` under `HeroSegment`.
- [ ] `pages/connections/components/BYOKSection/BYOKSection.tsx`: LLM key cards (reusing
      `widgets/ProviderCard/ProviderCard.tsx`) + `ComingSoonCard`s for CouchDB/S3/WebDAV.
- [ ] `components/ComingSoonCard/ComingSoonCard.tsx` (new tier-3 component).
- [ ] `app/api/artel/external_connections.pb.ts`, `processes/ExternalConnections.ts`,
      `app/hooks/ExternalConnections.ts`: add `addLlmKeyConnection`/`checkLlmKeyConnection` (the
      first two are generated/thin-wrapper, not hand-designed).

### Task 5 — Frontend: `ManageLlmKeyDialog`
_Ref: [03_frontend_connections_page.md](03_frontend_connections_page.md)_

- [ ] `dialogs/ManageLlmKeyDialog/ManageLlmKeyDialog.tsx` + colocated `components/ConnectForm/`,
      `components/ConnectedContent/` — mirrors `ManageGitlabDialog`'s folder shape exactly.
- [ ] Vendor tabs (Claude active, OpenAI disabled+tooltipped) inside `ConnectForm`.
- [ ] Wire "Test connection" → `checkLlmKeyConnection`, "Save" → `addLlmKeyConnection`,
      "Disconnect" → existing `disconnect(provider)` + `ConfirmDialog`.

### Task 6 — Tract domain/proto/engine for `llm_call` step
_Ref: [04_tract_llm_step.md](04_tract_llm_step.md)_

- [ ] `internal/domain/tract.go`: new `TractStep` fields (`LlmConnectionUuid`, `LlmModel`,
      `Prompt`, `SystemPrompt`, `MaxTokens`).
- [ ] `api/grpc/tracts.proto`: `LlmCallStep` message + `TractStep.kind` oneof branch.
- [ ] `moti g` + `bun gen`.
- [ ] `internal/transport/tracts_api/to_proto.go`: `stepToProto`/`stepFromProto` branch.
- [ ] `internal/service/v1/tract/field_consts.go` (or wherever step-type string constants live):
      `stepTypeLlmCall = "llm_call"`.
- [ ] `internal/service/v1/tract/llmexecutor.go`: `LlmExecutor` interface + adapter, following
      `toolexecutor.go`'s shape.
- [ ] `internal/service/v1/tract/engine.go`: `executeLlmCall`, dispatch branch in `executeStep`.
- [ ] Save-time validation (`validateLlmConnections`, mirroring `validateScriptEngines`) — find
      where tract create/update calls its step-tree validators and add this alongside.
- [ ] `internal/service/v1/tract/engine_test.go`: new tests for `executeLlmCall` (success,
      connection-not-owned, provider validation failure), following the file's existing test
      shape (check for a fake `ToolExecutor`/`LlmExecutor` test double pattern already in
      `mocks_test.go`).

### Task 7 — Frontend: `llm_call` step editing
_Ref: [04_tract_llm_step.md](04_tract_llm_step.md)_

- [ ] `processes/tractsTypes.ts`: `TractStep.type` union + new fields.
- [ ] `processes/tractSteps.ts`/`tractCanvasLayout.ts`: confirm the generic step-tree walk code
      doesn't special-case existing types in a way that silently drops `llm_call` (check for an
      exhaustive switch that would need a new arm, vs. code that's already type-agnostic).
- [ ] `components/StepPickerDialog/components/LlmCallStep.tsx` (new, alongside
      `ConnectionStep.tsx`/`ToolStep.tsx`) + wiring into `StepPickerDialog.tsx`.
- [ ] `pages/tract-canvas/components/LlmCallBody/LlmCallBody.tsx` (new, colocated, mirrors
      `ActionBody.tsx`'s structure): connection picker, model select, prompt/system-prompt
      `TemplateInput` fields, max-tokens input.
- [ ] `components/TractStepTree/components/ActionCard.tsx`: icon/label branch for `llm_call`.

### Task 8 — Usage capture
_Ref: [05_metrics_and_usage.md](05_metrics_and_usage.md)_

- [ ] Migration: `llm_usage_events` table.
- [ ] `internal/repository/pg/queries/llm_usage.sql` + `sqlc generate`.
- [ ] `internal/repository/pg/repos/llmusage/` repo (insert, sum-by-connection-since).
- [ ] `executeLlmCall` writes a usage event after a successful call (best-effort, logged on
      failure, never fails the run).
- [ ] `TractCanvasLogPanel.tsx`: render `output.usage` on `llm_call` run-step log entries.
- [ ] `GetLlmUsageSummary` RPC + service method + `ManageLlmKeyDialog`'s connected view renders it.

### Task 9 — Docs & cleanup

- [ ] Update `docs/architecture.md` if it enumerates services/tables (check before editing).
- [ ] `CLAUDE.md` gets a short pointer to this folder if the BYOK/Call-LLM surface becomes a
      recurring area of work (per the "durable conventions belong in CLAUDE.md" rule) — don't
      duplicate this folder's content there, just link it.

## Stage 2 — deferred, not scoped in detail here

- OpenAI-format BYOK connections (`ProviderOpenAI` wired to a real `internal/clients/openai`
  client + Chat Completions request shape). The proto/enum/UI groundwork from Stage 1 already
  reserves the slot; this is "implement the second client + enable the disabled tab."
- CouchDB instance / S3 bucket / WebDAV connections under the BYOK tab (today: `ComingSoonCard`
  placeholders only). Each is its own connection type with its own credential shape — worth a
  separate design pass per type when prioritized, not bundled into this plan.
- Cost estimation (dollar figures from token counts) — see the note in
  [05_metrics_and_usage.md](05_metrics_and_usage.md) for why this is deliberately deferred.
- Streaming / multi-turn / tool-use inside the Call LLM step, if a real workload demands it.
