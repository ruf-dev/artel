# Data Model & Security

## Where the key lives

Reuse `external_connections` (migration `024_external_connections.sql`) exactly as GitLab/Trello
do today — no new table for the key itself. Fields already on the table:

```sql
CREATE TABLE external_connections (
    id               UUID                   PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID                   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider         TEXT                   NOT NULL,
    provider_type    external_provider_type NOT NULL,   -- 'google_oauth' | 'api_key' | 'password'
    credentials_enc  BYTEA                  NOT NULL,    -- AES via cryptoutil, keyed by app secret
    metadata         JSONB,                              -- plaintext, non-sensitive display data
    created_at       TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, provider)  -- (superseded by a partial unique index; trello is excluded — see below)
);
```

`provider_type = 'api_key'` is reused as-is (same value GitLab/Trello already use) — an LLM key
is, mechanically, an API key. No new `external_provider_type` enum value needed.

### New `provider` values

`internal/domain/external_connection.go` gets two new constants, following the existing
lower-snake style (`ProviderGitlab = "gitlab"`, etc.):

```go
ProviderAnthropic = "anthropic"
ProviderOpenAI    = "openai" // reserved now, wired in a later stage — see 06_task_breakdown.md
```

One connection per provider per user, same as every provider except Trello — the existing
partial unique index from migration `045_trello_multi_connection.sql` already expresses "single
connection except trello":

```sql
CREATE UNIQUE INDEX external_connections_user_provider_single_uidx
    ON external_connections (user_id, provider)
    WHERE provider <> 'trello';
```

No migration needed for this constraint — `anthropic`/`openai` fall under the default "single
connection" branch automatically. This matches "each key type is a separate connection" from the
request: a user can hold one Anthropic-format connection and one OpenAI-format connection
simultaneously, each independently addressable, but not two Anthropic ones (nothing in the
request asks for multiple keys of the same vendor — if that's wanted later, follow the Trello
precedent: exclude `anthropic`/`openai` from the partial index too).

### `CredentialsJSON` shape (encrypted)

```go
// domain.AnthropicKeyCredentials — encrypted in credentials_enc
type AnthropicKeyCredentials struct {
    ApiKey  string `json:"api_key"`
    BaseUrl string `json:"base_url,omitempty"` // override for Claude Platform on AWS / regional endpoints; defaults to https://api.anthropic.com
}
```

Never round-tripped back to the frontend after creation. Editing follows the email-connection
precedent (`AddEmailConnection`/`storedEmailPassword`): the edit form's key field is left blank to
mean "keep the existing key" — a blank submission causes the service to load the stored key
before re-validating and re-saving, rather than the frontend ever receiving the plaintext key
back.

### `Metadata` shape (plaintext, non-sensitive — mirrors `gitlabConnectionMeta`/`trelloConnectionMeta`)

```go
type llmKeyConnectionMeta struct {
    Vendor          string   `json:"vendor"`           // "anthropic" | "openai" | "" (unconfirmed)
    KeyPreview      string   `json:"key_preview"`       // e.g. "sk-ant-...ab12", computed once at Add time
    DefaultModel    string   `json:"default_model"`     // last-selected default for new Call LLM steps
    AvailableModels []string `json:"available_models"`  // cached from the verification call's model list
    LastVerifiedAt  string   `json:"last_verified_at"`  // RFC3339, set on every successful Check/Add
}
```

`KeyPreview` is computed once, at write time, from the plaintext key the user submitted — never
by decrypting `credentials_enc` on every `ListConnections` call. This mirrors how GitLab stores
`Username`/`InstanceUrl` in plaintext metadata instead of re-decrypting the PAT to display it.

`AvailableModels` is a cache of the last verification response, not a live value — refreshed
every time the user hits "Test connection" or edits the key. It exists so the Call LLM step's
model picker doesn't need its own round-trip to the provider; see
[02_connection_lifecycle.md](02_connection_lifecycle.md).

## Encryption — nothing new required

`internal/repository/pg/repos/externalconnections/external_connections.go` already encrypts
`CredentialsJSON` via `cryptoutil.Encrypt(r.encryptionKey, ...)` before every `Insert`/`Upsert`,
and decrypts on every read. An LLM key connection goes through the exact same repo methods — the
repo has no per-provider branching today and needs none added. This is the same encryption key
(`ARTEL_ENCRYPTION_KEY` / whatever config wires `encryptionKey` into `externalconnections.New`)
already protecting Gmail OAuth tokens, GitLab PATs, and Trello tokens. No new secret-management
surface, no new KMS integration, no new envelope encryption — deliberately, to avoid a bespoke
security model for "the dangerous one." The trust story for the user is: *"your key is encrypted
the same way your GitLab token already is — we don't invent a special case that's weaker or
stronger."*

## Threat model notes (for whoever implements this)

- **In transit to Artel**: standard TLS to the API endpoint, same as every other Add-connection
  RPC. No special handling needed beyond what gRPC-gateway already provides.
- **At rest**: `credentials_enc` (AES, existing key). Compromise of the DB alone does not yield
  the key; compromise of DB + app secret does — same blast radius as every other stored
  credential today. Do not build BYOK-specific extra hardening (e.g. per-user KMS keys) unless
  the user explicitly asks for it later — that's a much bigger, separate initiative and the task
  didn't ask for it.
- **In the tract engine at call time**: decrypted only inside `executeLLMCall`'s call to the new
  `internal/clients/anthropic` client, immediately before the outbound HTTPS request — same
  lifetime as how `executeTool` resolves MoM secrets today (see `mcp/executors/http.go`'s
  `secrets` map, which is decrypted per-call and never logged). The new LLM client must not log
  the key; `log_transport.go`'s pattern (used by the MoM HTTP executor for request/response
  logging) needs an explicit redaction rule if the LLM client reuses that transport — flag this
  as a review item in the task breakdown, don't silently inherit logging behavior tuned for
  non-secret MoM headers.
- **Never surfaced back to the frontend**: `ListConnections`/`ExternalConnectionInfo` for an LLM
  key connection carries only `KeyPreview`, `Vendor`, `DefaultModel`, `AvailableModels`,
  `LastVerifiedAt` — the same "write-only secret" contract GitLab/Trello already follow (compare
  `GenerateGitlabWebhookSecret`'s explicit "one-time reveal; never retrievable again" comment).
