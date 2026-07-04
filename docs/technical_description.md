# Artel — Technical Infrastructure

## Application Architecture

```
cmd/service/main.go
  └── internal/app/custom.go        (Init / Start / Stop wiring)
        ├── internal/transport/     (HTTP + gRPC handlers, MCP, OAuth)
        ├── internal/service/v1/    (business logic per domain)
        ├── internal/repository/    (PostgreSQL queries via SQLC)
        └── internal/clients/       (CouchDB, IMAP, SMTP)
              ├── PostgreSQL         (metadata, sessions, credentials)
              └── CouchDB            (vault document storage, per-vault isolated)

pkg/client/ArtelUI/                 (React + Bun frontend, served at /)
```

The app is scaffolded by [rscli](https://github.com/Red-Sock/rscli). `app.go` and `config.go` are generated and never edited — all wiring goes in `custom.go`.

---

## Infrastructure Components

### PostgreSQL — Primary Metadata Store

Manages all relational data. Queries are generated with [sqlc](https://sqlc.dev/); migrations run via [Goose](https://github.com/pressly/goose).

**Tables:**

| Table | Purpose |
|-------|---------|
| `users` | Accounts: email, username, photo, identity type (Artel/Telegram) |
| `subscriptions` | Feature gate — whether a user has an active subscription |
| `sessions` | Auth tokens with expiry |
| `user_permissions` | Per-user flags: `is_administrator`, `has_emails` |
| `vaults` | Vault metadata: name, status, owning CouchDB instance |
| `vault_members` | Many-to-many: vault ↔ user with role |
| `couch_instances` | Admin-registered CouchDB servers (URL, credentials encrypted) |
| `couch_accounts` | Per-user CouchDB credentials for a given instance |
| `mcp_keys` | MCP bearer tokens: bcrypt-hashed secret, vault scope, preview string |
| `email_accounts` | User IMAP/SMTP credentials (AES-256 encrypted) |
| `prompts` | Pre-defined setup prompts served to the UI (e.g., vault setup for Claude) |

Schema evolves through 15 SQL migrations (`migrations/`).

### CouchDB — Per-Vault Document Store

Each vault maps to a dedicated CouchDB database on an admin-registered instance. Users get individual CouchDB accounts — no credential sharing.

**Client:** `internal/clients/couchdb/livesync.go` (`LiveSyncClient`)

**Operations exposed:**
- `WriteNote` / `ReadNote` / `DeleteFile`
- `ListFiles` / `ListFolders` / `ListTags`
- `MoveFile`

**Initialization:** System databases (`_users`, `_replicator`) are auto-created when a new CouchDB instance is registered.

**Vault creation flow** (atomic across both stores):
1. Insert vault record in PostgreSQL
2. Create CouchDB user + database via CouchDB HTTP API
3. Insert user_permissions record
4. All three steps run inside a single PostgreSQL transaction manager (`tx_manager`)

### Email — IMAP + SMTP Clients

- `internal/clients/imap/` — list folders, list messages, fetch by UID
- `internal/clients/smtp/` — send email via PlainAuth

Clients are instantiated **per-request** (lazy, no persistent connections). Credentials are decrypted in-memory from the `email_accounts` table at call time.

---

## API Surface

### gRPC + REST (grpc-gateway)

Single binary serves both binary gRPC (HTTP/2) and JSON REST on **port 1551**. REST is translated at `/api` by grpc-gateway.

**Proto services → REST prefixes:**

| Service | REST prefix |
|---------|-------------|
| `AuthAPI` | `/api/auth/` |
| `VaultsAPI` | `/api/vaults/` |
| `CouchInstancesAPI` | `/api/couch/` |
| `EmailAccountsAPI` | `/api/email-accounts/` |
| `McpKeysAPI` | `/api/mcp/{vault_id}/` |
| `PromptsAPI` | `/api/prompts/` |

Proto definitions: `api/grpc/`. Generated Go stubs: `internal/api/server/artel_api/`.

### MCP JSON-RPC Endpoint

Exposes vault content and email to external AI tools (Claude, etc.) using the [Model Context Protocol](https://modelcontextprotocol.io/).

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/mcp` | `POST` | JSON-RPC 2.0 request |
| `/mcp` | `GET` | Server-Sent Events stream |
| `/.well-known/oauth-authorization-server` | `GET` | OAuth discovery |
| `/.well-known/oauth-protected-resource` | `GET` | Resource metadata |
| `/oauth/login` | `POST` | MCP session login |
| `/oauth/vaults` | `POST` | List vaults for MCP session |
| `/oauth/vault` | `POST` | Authorize MCP client for a vault |
| `/token` | `POST` | OAuth token exchange |

**Auth:** Bearer token (`artel_vtk_{uuid}_{hex_secret}`) resolved to vault + CouchDB credentials on each request.

**MCP tools exposed:** `list_files`, `read_file`, `write_file`, `delete_file`, `move_file`, `list_folders`, `list_tags`, `list_email_accounts`, `list_email_folders`, `list_emails`, `read_email`, `send_email`

Vault file content is split across two backends by extension: `.md`/`.markdown` paths are stored in CouchDB; every other path is stored in the vault's linked S3-compatible bucket (optional, linked separately from the mandatory CouchDB instance). `write_file` requires base64-encoded `content` for non-markdown paths.

### OAuth 2.0 PKCE Flow

Used by Claude and other MCP clients to authorize access to a specific vault:
1. Client redirects to Artel with `code_challenge`
2. User selects vault in the Artel UI
3. Artel issues authorization code → redirect back to client
4. Client exchanges code for bearer token at `/token`

---

## Security Model

| Concern | Mechanism |
|---------|-----------|
| User passwords | bcrypt, cost=12 |
| Session tokens | Cryptographically random, expiring, stored in PostgreSQL |
| CouchDB / IMAP / SMTP credentials | AES-256 encrypted before storage; key from `creds_encryption_key` env var |
| MCP token format | `artel_vtk_{uuid}_{hex_secret}` — secret bcrypt-hashed, raw token shown once |
| Admin endpoints | gRPC interceptor checks `is_administrator` flag before allowing `CouchInstancesAPI` calls |
| Email feature | Gated by `has_emails` permission flag per user |

The `creds_encryption_key` environment variable must be set at runtime (hex-encoded AES key). It is never committed.

---

## Key Data Flows

### Vault Creation
```
Client → VaultsAPI.CreateVault
  → service/vault: generate CouchDB password
  → tx_manager (PostgreSQL tx):
      INSERT vaults
      INSERT couch_accounts
      INSERT user_permissions
  → CouchDB HTTP API:
      PUT /_users/{username}
      PUT /{db_name}
  → return vault + db_url
```

### MCP Request Auth
```
POST /mcp (Bearer: artel_vtk_abc_def)
  → parse uuid from token
  → SELECT mcp_key WHERE uuid = ?
  → bcrypt.Compare(def, stored_hash)
  → load vault + couch_account
  → instantiate LiveSyncClient
  → execute tool, return result
```

### Email Credential Lifecycle
```
AddEmailAccount → AES-256 encrypt(password) → INSERT email_accounts
ListEmails      → SELECT email_accounts → AES-256 decrypt → new imap.Client → fetch
SendEmail       → SELECT email_accounts → AES-256 decrypt → new smtp.Client → send
```

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go |
| Scaffolding | rscli (RedSock CLI) |
| Config | matreshka (YAML + env vars) |
| API framework | gRPC + grpc-gateway (REST bridge) |
| Database (metadata) | PostgreSQL |
| Query generation | sqlc |
| Migrations | Goose |
| Database (vault content) | CouchDB (LiveSync protocol) |
| Email | go-imap (IMAP), net/smtp (SMTP) |
| Auth (external) | Telegram OAuth (JWKS validation via golang-jwt) |
| Password hashing | bcrypt (golang.org/x/crypto) |
| Logging | zerolog (rs/zerolog) |
| Error handling | rerrors (go.redsock.ru/rerrors) |
| Graceful shutdown | toolbox/closer (go.redsock.ru/toolbox) |
| Frontend | React 18 + TypeScript |
| Frontend build | Bun |
| Proto codegen (TS) | `bun gen` (from pkg/client/ArtelUI) |
| Proto codegen (Go) | `moti g` |
