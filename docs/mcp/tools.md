# Tools and When They Appear

`tools/list` returns two kinds of tools. What you see for a given key depends on what's connected
and what that specific key has been granted.

## Built-in tools

Always present for any valid, unrevoked MCP key — no extra setup needed beyond having the key.

- **Vault tools** — `list_files`, `read_file`, `write_file`, `delete_file`, `move_file`,
  `list_folders`, `list_tags`, `get_file_metadata`. Operate on the vault the key is bound to.
  Defined in code (`internal/service/v1/mcp/executors/vault.go`), not stored in the database.
  Vault content is split across two backends by file extension: `.md`/`.markdown` paths live in
  CouchDB, everything else lives in the vault's linked S3-compatible bucket (if any). `write_file`
  takes plain-text `content` for markdown paths and base64-encoded `content` for everything else;
  non-markdown ops fail with a "no linked S3 bucket" error on vaults that have no bucket linked.
  `move_file` cannot cross the markdown/S3 boundary in one call — write to the new path then
  delete the old one instead. `get_file_metadata`'s `rev`/`ctime`/`deleted` fields are
  CouchDB-specific and always zero-valued (`rev` empty, `ctime` == `mtime`, `deleted` false) for
  S3-backed files.
- **`connections`** — always listed, but its *result* depends on what MoMs (see below) are linked
  to this key. Calling it returns the connected MoM packages (name/author/description), e.g.
  `[{"name":"email","author":"Artel","description":"..."}]`, or `[]` if none are linked yet.

## MoM tools (Mcp of Mcp)

MoM tools (e.g. `list_email_folders`, `list_emails`, `read_email`, `send_email` for the `email` MoM)
are dynamic — they only show up in `tools/list` when **both** of these are true:

1. **The user connected the underlying service.** An `external_connections` row exists for that
   provider (e.g. an IMAP/SMTP account added via `AddEmailConnection`, or a Google OAuth connection).
2. **That specific MCP key was granted access to it.** An `mcp_connectors` row links
   `(mcp_key_uuid, mcp_name)` → `external_connection_uuid`. This is what lets a vault owner issue
   multiple keys (e.g. one for a personal assistant, one for a teammate) and decide per key which
   MoMs it can use — connecting an account doesn't automatically expose it to every key.

If only condition 1 holds (service connected, but this key has no connector for it), the MoM's tools
are absent from `tools/list` and calling them by name fails with "tool not found in any connected mcp".

> **Current limitation:** there is no API/UI yet to create the `mcp_connectors` link (step 2) — the
> `mcp_connectors` table and lookup (`ListByKey`) exist, but nothing currently inserts into it.
> Until that's wired up, MoM tools won't appear for any key even if the user has connected the
> service.

## Why dynamic instead of always-on

Vault tools are core product functionality, so they're cheap to keep static. MoM tools represent
optional integrations a user may or may not have set up, and a given key may be scoped to only some
of them — flattening them into `tools/list` (rather than a separate discovery call) keeps the
protocol standard-compliant (plain `tools/list` + `tools/call`, see [examples.md](examples.md)) while
the per-key connector check keeps each key's tool surface scoped to what it's actually allowed to use.
