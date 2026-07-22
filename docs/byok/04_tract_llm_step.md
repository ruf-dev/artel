# Tract Step: Call LLM

## Why this is a bespoke engine step, not a MoM tool

`CLAUDE.md`'s MoM convention is for outbound HTTP calls exposed *to an agent as a callable tool*
(email, Trello, GitLab — things an LLM decides to invoke mid-conversation via a declarative
`api_description` + `action`). A "Call LLM" tract step is the inverse: it's the tract *itself*
invoking an LLM as a fixed pipeline stage, with typed inputs (prompt, model, connection) the user
configures at design time in the canvas — structurally identical to `ScriptStep` (also a
hardcoded, typed step kind with its own engine method), not to a MoM tool definition. Forcing it
through the MoM HTTP executor would mean modeling Anthropic's request/response shape as a generic
`${{params.*}}`/`__secrets.*` JSON template, which buys nothing (no reuse across other tools) and
loses type safety on the one thing (`model`, `max_tokens`) that benefits from being real Go
fields validated at save time. Bespoke `internal/clients/anthropic` + a new engine method is the
right call, consistent with the "reserve bespoke Go code for things that are genuinely not a
declarative HTTP tool call" carve-out already in the repo's MoM section.

## Domain (`internal/domain/tract.go`)

`TractStep` gains fields, following the exact comment convention already used for
Script/Action-only fields:

```go
type TractStep struct {
    // ... existing fields unchanged ...

    // LLM fields — only meaningful when Type == "llm_call". ConnectionUuid points at the
    // external_connections row (provider "anthropic"/"openai") supplying the key. Prompt and
    // SystemPrompt are template strings rendered through the same resolver as action Params —
    // {{steps.<id>.output...}} and {{trigger...}} both work here.
    LlmConnectionUuid uuid.UUID `json:"llm_connection_uuid,omitempty"`
    LlmModel          string    `json:"llm_model,omitempty"`
    Prompt            string    `json:"prompt,omitempty"`
    SystemPrompt       string    `json:"system_prompt,omitempty"`
    MaxTokens          int       `json:"max_tokens,omitempty"`
}
```

New step type constant alongside the existing ones defined in `internal/service/v1/tract/`
(`stepTypeAction`, `stepTypeCondition`, etc. — check `field_consts.go`/`dispatch.go` for where
those live): `stepTypeLlmCall = "llm_call"`.

Reuses `ConnectionUuid`'s existing role (a pointer to an `external_connections` row) rather than
introducing a second "connection reference" concept — named `LlmConnectionUuid` only to avoid
colliding with the action-step `ConnectionUuid` field's semantics (MoM connection vs. LLM key
connection are different things even though both are `external_connections` rows).

`Prompt`/`SystemPrompt` are plain strings, not `map[string]string` like `Params` — there's a
single prompt, not a bag of named params, so this doesn't reuse the `Params` field the way script
steps reuse it for input bindings. Variable interpolation still goes through the same
`{{...}}` resolver (`internal/service/v1/tract/template.go`), just applied to one string field
instead of iterated over a map.

## Proto (`api/grpc/tracts.proto`)

```proto
message LlmCallStep {
  string connection_id   = 1; // external_connections.id (anthropic/openai key)
  string model           = 2;
  string prompt          = 3; // template string
  string system_prompt   = 4; // template string, optional
  int32  max_tokens      = 5; // optional, service applies a default if 0
}

message TractStep {
  string id          = 1;
  string name        = 2;
  string description = 3;
  oneof kind {
    ActionStep action       = 4;
    ConditionStep condition = 5;
    ParallelStep parallel   = 6;
    GroupStep group         = 7;
    ScriptStep script        = 8;
    LlmCallStep llm_call     = 9; // new
  }
}
```

`internal/transport/tracts_api/to_proto.go`'s `stepToProto`/`stepFromProto` get a new branch,
mirroring the existing `ScriptStep` branch exactly (field-for-field mapping, no special logic).

## Validation at save time

Mirrors `validateScriptEngines` (script steps validate `Language` has a registered engine before
the tract can be saved). A new `validateLlmConnections` walk over the step tree at
create/update-tract time checks, for every `llm_call` step:

1. `LlmConnectionUuid` resolves to an `external_connections` row owned by the caller, with
   `provider` in `{anthropic, openai}`.
2. `LlmModel` is non-empty and present in that connection's cached `AvailableModels` metadata
   (from the last verification — see [02_connection_lifecycle.md](02_connection_lifecycle.md)).
   A stale cache (models list changed provider-side since last Check) should warn, not hard-block
   — the actual call at run time is the real gate; treat this as a save-time sanity check, same
   spirit as `TractConnectionNotOwned` being re-verified at run time in `executeTool` even though
   the step was validated at save time.

Reuses the existing `TractConnectionNotOwned`-style re-verification at run time in the engine
(below) — save-time validation is a UX nicety, not the security boundary.

## Engine (`internal/service/v1/tract/engine.go`)

New branch in `executeStep`'s switch:

```go
case stepTypeLlmCall:
    return s.executeLlmCall(ctx, state, step)
```

`executeLlmCall` follows the exact same shape as `executeAction`/`executeScript` — snapshot
resolver, render, insert run step, call out, finish run step, set output:

```go
func (s *Service) executeLlmCall(ctx context.Context, state *runState, step domain.TractStep) error {
    rslv := state.snapshotResolver()

    renderedPrompt, err := rslv.render(step.Prompt)
    // ... error handling, same pattern

    renderedSystemPrompt, err := rslv.render(step.SystemPrompt)
    // ...

    inputJSON, _ := json.Marshal(map[string]interface{}{
        "model": step.LlmModel, "prompt": renderedPrompt, "system_prompt": renderedSystemPrompt,
    })

    runStep := domain.TractRunStep{RunUuid: state.run.Uuid, StepId: step.Id, StepName: step.Name,
        StepType: step.Type, Input: inputJSON}
    insertedStep, err := s.tracts.InsertRunStep(ctx, runStep)
    // ...

    result, err := s.llmExecutor.Call(ctx, state.tract.UserUuid, step.LlmConnectionUuid, llm.CallRequest{
        Model: step.LlmModel, Prompt: renderedPrompt, SystemPrompt: renderedSystemPrompt,
        MaxTokens: resolveMaxTokens(step.MaxTokens),
    })
    if err != nil {
        // UpdateRunStepFinish(..., domain.TractRunStepFailed, ...); return wrapped err
    }

    output := map[string]interface{}{
        "text": result.Text, "model": result.Model, "stop_reason": result.StopReason,
        "usage": map[string]interface{}{
            "input_tokens": result.Usage.InputTokens, "output_tokens": result.Usage.OutputTokens,
            "cache_creation_input_tokens": result.Usage.CacheCreationInputTokens,
            "cache_read_input_tokens":    result.Usage.CacheReadInputTokens,
        },
    }
    outputJSON, _ := json.Marshal(output)
    // UpdateRunStepFinish(..., domain.TractRunStepDone, outputJSON, ""); state.setOutput(step.Id, output)
    return nil
}
```

Later steps reference `{{steps.<id>.output.text}}` for the completion text and
`{{steps.<id>.output.usage.output_tokens}}` etc. for token counts — no resolver changes needed,
`output` is just another JSON value in `runState.outputs`, exactly like an action step's parsed
tool result.

### `LlmExecutor` interface (new, `internal/service/v1/tract/llmexecutor.go`)

Same narrow-interface pattern as `ToolExecutor` in `toolexecutor.go` — defined in `tract` to avoid
an import cycle, implemented by an adapter over the connections + LLM client services:

```go
type LlmExecutor interface {
    // Call re-verifies connection ownership itself (userUuid) before making the request —
    // callers must not skip that check, same contract as ExecuteMomTool.
    Call(ctx context.Context, userUuid uuid.UUID, connectionUuid uuid.UUID, req llm.CallRequest) (llm.CallResult, error)
}
```

## `internal/clients/anthropic` (new package)

Thin wrapper over `github.com/anthropics/anthropic-sdk-go`, matching the existing
`internal/clients/{googleapi,imap,smtp}` shape (constructor takes resolved credentials, exposes a
small number of typed methods — no generic "do anything" surface):

```go
package anthropic

type Client struct { sdk *anthropic.Client }

func New(apiKey, baseUrl string) *Client { ... } // option.WithAPIKey + option.WithBaseURL if baseUrl != ""

func (c *Client) ListModels(ctx context.Context) ([]ModelInfo, error) // GET /v1/models, no beta header
func (c *Client) Complete(ctx context.Context, req CompleteRequest) (CompleteResult, error) // POST /v1/messages
```

`Complete` builds a single-turn `client.Messages.New` call per the SDK's Go usage:

```go
resp, err := c.sdk.Messages.New(ctx, anthropic.MessageNewParams{
    Model:     anthropic.Model(req.Model),
    MaxTokens: int64(req.MaxTokens),
    System:    systemBlocks(req.SystemPrompt), // omit if empty
    Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt))},
})
```

No tool use, no thinking, no streaming for stage 1 — a single request/response call matching what
the step actually needs (one prompt in, one completion out). `MaxTokens` defaults to a
conservative value (e.g. 4096) when the step doesn't set one, non-streaming — per the Claude API
guidance, only reach for streaming once a step needs `max_tokens` large enough to risk an HTTP
timeout (~16K+); that's a stage-2 concern if a "long-form generation" use case shows up, not a
stage-1 requirement. Extended thinking/effort tuning is likewise out of scope for stage 1 — expose
`model` and `max_tokens` only; add an "advanced" effort/thinking control later if a real workload
asks for it, rather than speculatively wiring params nobody's configuring yet.

`ListModels` wraps `client.Models.List(ctx)` and maps to the connection's `AvailableModels`
metadata cache — the same call used at Add/Check time in
[02_connection_lifecycle.md](02_connection_lifecycle.md), reused here so there's exactly one
place that talks to Anthropic's Models endpoint.

## Frontend: step type, editor, picker

`pkg/client/ArtelUI/src/processes/tractsTypes.ts`:

```ts
export interface TractStep {
    // ... existing fields ...
    type: "action" | "condition" | "parallel" | "group" | "script" | "llm_call"
    llmConnectionId?: string
    llmModel?: string
    prompt?: string
    systemPrompt?: string
    maxTokens?: number
}
```

- `components/StepPickerDialog/StepPickerDialog.tsx` gets a new option ("Call LLM"), alongside
  `ConnectionStep.tsx`/`ToolStep.tsx` — likely a new sibling `LlmCallStep.tsx` in the same
  `components/StepPickerDialog/components/` folder, following that folder's existing per-kind
  split.
- `pages/tract-canvas/components/ActionBody/ActionBody.tsx` is the reference for how a step's
  body renders/edits inline in the canvas — a new colocated sibling (`LlmCallBody.tsx` under
  `pages/tract-canvas/components/`) renders:
  - A connection picker filtered to `useExternalConnections()` rows where
    `provider ∈ {ANTHROPIC, OPENAI}` (reuses the hook, no new fetch).
  - A model `<select>` populated from that connection's `metadata.available_models` (already
    cached, no extra network call from the canvas).
  - Prompt/system-prompt fields using the existing `components/TemplateInput/TemplateInput.tsx` —
    it already solves "insert a `{{steps.X.output...}}` reference into a text field," which is
    exactly what this step needs and is the component other step kinds already use for template
    strings (confirm by checking who else imports `TemplateInput` — if `ActionBody`'s `Params`
    fields already use it, this is a drop-in reuse, not a new capability).
  - `MaxTokens` as a plain number input, optional, defaulting to blank (service-side default
    applies).
- `components/TractStepTree/components/ActionCard.tsx` (the step-tree row renderer) needs a
  `llm_call` branch for its icon/label, same shape as its existing per-type branches.

## What Call LLM step deliberately does NOT do in stage 1

- No tool use / function calling from within the LLM call (the step is "prompt in, text out," not
  an agent loop). If a future need arises for the LLM to call other tract tools mid-step, that's
  a much bigger design (effectively a nested agent step) — flag it as an idea, don't build it now.
- No streaming output into the canvas run log. `TractRunStepItem.output` is written once, on
  completion, like every other step type — the whole tract engine is fire-and-collect, not
  streaming, and this step shouldn't be the one exception.
- No multi-turn/conversation state across steps. Each `llm_call` step is one independent request;
  chaining several is done by referencing a prior step's `output.text` in a later step's
  `Prompt`, the same composition model already used for action steps' outputs.
