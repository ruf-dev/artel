# Connection Lifecycle: Add, Verify, Detect Vendor, Rotate, Delete

## Proto additions (`api/grpc/external_connections.proto`)

New enum values (append, don't renumber existing ones):

```proto
enum ExternalProvider {
  EXTERNAL_PROVIDER_UNSPECIFIED   = 0;
  EXTERNAL_PROVIDER_GOOGLE_SHEETS = 1;
  EXTERNAL_PROVIDER_TRELLO        = 2;
  EXTERNAL_PROVIDER_MIRO          = 3;
  EXTERNAL_PROVIDER_EMAIL         = 4;
  EXTERNAL_PROVIDER_GITLAB        = 5;
  EXTERNAL_PROVIDER_ANTHROPIC     = 6; // stage 1
  EXTERNAL_PROVIDER_OPENAI        = 7; // reserved — stage 2, not wired to a service yet
}
```

New messages and RPCs, mirroring `AddGitlabConnection`/`CheckGitlabConnection` exactly:

```proto
message LlmKeyConnectionInfo {
  string vendor           = 1; // "anthropic" | "openai"
  string key_preview      = 2;
  string default_model    = 3;
  repeated string available_models = 4;
  string last_verified_at = 5;
}

// Extend ExternalConnectionInfo.details oneof with:
//   LlmKeyConnectionInfo llm_key = 7;

message AddLlmKeyConnection {
  message Request {
    ExternalProvider provider = 1; // ANTHROPIC | OPENAI
    string api_key            = 2;
    string base_url           = 3; // optional override, e.g. Claude Platform on AWS endpoint
    string default_model      = 4; // optional; falls back to the vendor's recommended default
  }
  message Response {
    ExternalConnectionInfo connection = 1;
  }
}

message CheckLlmKeyConnection {
  message Request {
    ExternalProvider provider = 1;
    string api_key            = 2;
    string base_url           = 3;
  }
  message Response {
    string vendor                    = 1; // confirmed/detected vendor label
    repeated string available_models = 2;
    string recommended_default_model = 3;
  }
}
```

`AddLlmKeyConnection`/`CheckLlmKeyConnection` follow the same pair-shape as
`AddGitlabConnection`/`CheckGitlabConnection`: `Check` validates without persisting (backs the
"Test connection" button before the user commits), `Add` validates then persists — internally
`Add` calls the exact same validation function `Check` uses, exactly like
`validateGitlabToken` is shared between `AddGitlabConnection` and `CheckGitlabConnection` today.

A blank `api_key` on `Add` for an *existing* connection (edit flow) means "keep the current key,
just update `base_url`/`default_model`" — same `if password == "" { load stored }` pattern as
`AddEmailConnection`.

## Vendor detection

Two independent signals, used together, never trusted alone:

1. **Tab choice fixes the wire protocol.** The frontend's Add-key form has two tabs (stage 1
   ships only "Claude"; "OpenAI" tab is visibly present but disabled/"coming soon", per the
   request to design Claude now while keeping the other format in mind). Which tab the user is
   on determines which request/response shape the backend uses to talk to `base_url` — Anthropic
   Messages API (`x-api-key` header, `anthropic-version` header) vs. OpenAI Chat Completions
   shape (`Authorization: Bearer`). This is *not* inferred — it's simply which RPC field
   (`provider`) the frontend sends.
2. **Live confirmation, not just a prefix guess.** A key-prefix heuristic (`sk-ant-...` →
   Anthropic, `sk-...`/`sk-proj-...` → OpenAI) is shown as an inline hint in the UI the moment the
   user pastes a key (pure client-side string check, no round-trip) — but it is never the
   authority. The authority is `CheckLlmKeyConnection` actually calling
   `GET {base_url}/v1/models` with the submitted key. Anthropic's Models API
   (`client.models.list()` / raw `GET /v1/models`) needs no beta header, costs zero tokens, and
   both validates the key *and* returns the model catalog (`id`, `display_name`,
   `max_input_tokens`, `max_tokens`) in one call — that catalog becomes `available_models` and
   seeds the Call LLM step's model picker (see [04_tract_llm_step.md](04_tract_llm_step.md)).

   For a custom/self-hosted `base_url` that doesn't identify as either vendor cleanly (a proxy,
   an OpenAI-compatible gateway, Azure OpenAI, etc.) — this is the case the request calls out
   with *"in the end can let user decide the llm behind the api"*: if the models-list call
   succeeds but the response doesn't look like a recognized vendor shape, the UI falls back to
   asking the user to label it themselves (a free-text/short-list "what's actually behind this
   URL" field, stored as `Vendor` in metadata for display purposes only — it never changes which
   request shape the backend sends, since that's still fixed by which tab/RPC field was used).

## `internal/service/v1/externalconnections/external_connections.go` additions

New methods on `Service`, same shape as `AddGitlabConnection`/`CheckGitlabConnection`:

```go
func (s *Service) AddLlmKeyConnection(ctx context.Context, provider, apiKey, baseUrl, defaultModel string) (domain.ExternalConnectionMeta, error)
func (s *Service) CheckLlmKeyConnection(ctx context.Context, provider, apiKey, baseUrl string) (llmKeyCheckResult, error)
```

Internals of `CheckLlmKeyConnection` (shared by `Add`):

1. Resolve `baseUrl` to a default per provider if blank (`https://api.anthropic.com` for
   Anthropic).
2. Build the new `internal/clients/anthropic` client (see
   [04_tract_llm_step.md](04_tract_llm_step.md) for the client itself) with the submitted key +
   base URL.
3. Call its `ListModels(ctx)` — thin wrapper over the SDK's `client.models.list()`.
4. On auth failure, return a validation error (`user_errors.LlmKeyValidationFailed`, new — mirrors
   `GitlabValidationFailed`/`TrelloValidationFailed`).
5. On success, return the model list + a `recommended_default_model` (pick the newest
   Opus-tier/flagship model from the list — same "always default to the most capable" preference
   documented for LLM app builders generally; don't hardcode a specific model ID here, derive it
   from what the list call actually returns so this doesn't rot when models retire).

`AddLlmKeyConnection` calls the same validation, then persists via `s.connections.Upsert` with
`domain.ExternalConnection{Provider: domain.ProviderAnthropic, ProviderType:
artel_q.ExternalProviderTypeApiKey, CredentialsJSON: ..., Metadata: ...}` — identical shape to
`AddGitlabConnection`'s tail end.

## Frontend service layer

`pkg/client/ArtelUI/src/processes/ExternalConnections.ts` (`ExternalConnectionsService`) gets two
new methods, `addLlmKeyConnection`/`checkLlmKeyConnection`, following the existing
`addGitlabConnection`/`checkGitlabConnection` pattern exactly — same file, no new service class.
`useExternalConnections` (`app/hooks/ExternalConnections.ts`) needs no shape change: LLM key
connections are just more `ExternalConnectionInfo` rows in the same `connections` array, filtered
client-side by provider the same way `ContentSegment.tsx` already filters by provider today.

## Delete / disconnect

No new RPC needed — `DisconnectProvider(provider: "anthropic")` already works unchanged, exactly
like GitLab's disconnect flow. `ManageLlmKeyDialog` (see
[03_frontend_connections_page.md](03_frontend_connections_page.md)) reuses the same
`ConfirmDialog` + `disconnect(provider)` pattern as `ManageGitlabDialog`.

## Rotation

"Rotate" is just Add again with a new key on an existing provider — `Upsert` (not `Insert`)
already overwrites the row, matching how GitLab's webhook-secret rotation and email's
edit-password flow work today. No separate rotate endpoint.
