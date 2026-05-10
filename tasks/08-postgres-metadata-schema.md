---
id: "08"
title: "PostgreSQL Metadata Schema"
status: "pending"
model: "qwen2.5-coder:3b"
created: "2026-05-09"
branch: "factory/08-postgres-metadata-schema"
---

# Task 08 — PostgreSQL Metadata Schema

## Goal

Replace the `couchdb_url` environment variable with a PostgreSQL-backed metadata layer: users, vaults, subscriptions,
and per-vault CouchDB credentials encrypted at the application level.

## Context

- `config/config.yaml` currently has `couchdb_url` in the `environment:` section and a `postgres` data source already
  configured.
- `internal/config/environment.go` is **generated** by rscli — never edit it by hand. After changing
  `config/config.yaml`, run `rscli-dev project tidy` to regenerate it.
- `internal/clients/couchdb/config.go` parses `CouchdbURL` from `EnvironmentConfig`. Once the env var is removed, this
  file needs to be deleted and its callers updated.
- `internal/app/custom.go` currently calls `couchdb.NewConfig(a.Cfg.Environment.CouchdbURL)` — this call must be
  removed. The CouchDB client itself stays; it will receive credentials fetched from the DB at runtime.
- `internal/clients/sqldb/` already sets up a `*sql.DB` via `lib/pq`. Use this for repository implementations.
- No `migrations/` directory exists yet.

### Credential Encryption

Store CouchDB credentials (host, username, password) in a `couch_credentials` table. Encrypt only the password using *
*AES-256-GCM** from the Go standard library (`crypto/aes`, `crypto/cipher`). The encryption key comes from a new
`creds_encryption_key` env var (32-byte hex string). Create `internal/cryptoutil/aes.go` with
`Encrypt(key, plaintext []byte) ([]byte, error)` and `Decrypt(key, ciphertext []byte) ([]byte, error)`.

## Acceptance Criteria

- [ ] `couchdb_url` is removed from `config/config.yaml` and `environment.go` is regenerated (no `CouchdbURL` field)
- [ ] `creds_encryption_key` is added to `config/config.yaml` and appears in `EnvironmentConfig`
- [ ] `internal/clients/couchdb/config.go` is deleted; `custom.go` no longer references `CouchdbURL`
- [ ] `migrations/001_initial_schema.sql` exists with all four tables (see Notes)
- [ ] `internal/repository/interfaces.go` defines `Users`, `Vaults`, `Subscriptions`, `CouchCredentials` interfaces
- [ ] `internal/repository/v1/` contains implementations that compile and pass `go test ./...`
- [ ] `internal/cryptoutil/aes.go` compiles; `Encrypt` → `Decrypt` round-trips correctly in a unit test
- [ ] `go test ./...` passes from repo root

## Files to Create / Modify

- `config/config.yaml` — remove `couchdb_url`, add `creds_encryption_key`
- `internal/config/environment.go` — regenerated via `rscli-dev project tidy` (do not edit manually)
- `internal/clients/couchdb/config.go` — **delete**
- `internal/app/custom.go` — remove `couchdb.NewConfig(...)` call; wire postgres `*sql.DB` from
  `a.Cfg.DataSources.Postgres`
- `migrations/001_initial_schema.sql` — new
- `internal/repository/interfaces.go` — new
- `internal/repository/v1/impl.go` — new (Repos struct)
- `internal/repository/v1/users/users.go` — new
- `internal/repository/v1/vaults/vaults.go` — new
- `internal/repository/v1/subscriptions/subscriptions.go` — new
- `internal/repository/v1/couchcreds/couchcreds.go` — new
- `internal/cryptoutil/aes.go` — new
- `internal/cryptoutil/aes_test.go` — new

## Do NOT change

- `internal/app/app.go` — generated, never edit
- `internal/app/config.go` — generated, never edit
- `internal/clients/couchdb/client.go` — keep; CouchDB client is still used
- `internal/clients/sqldb/` — keep as-is; use the existing DB setup

## Notes

### Migration SQL (`migrations/001_initial_schema.sql`)

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    email      TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE subscriptions
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    active     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE vaults
(
    id            UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name          TEXT        NOT NULL,
    couch_db_name TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE couch_credentials
(
    id           UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    vault_id     UUID        NOT NULL REFERENCES vaults (id) ON DELETE CASCADE,
    host         TEXT        NOT NULL,
    username     TEXT        NOT NULL,
    password_enc BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Repository Interface Sketch

```go
// internal/repository/interfaces.go
package repository

type Users interface {
	Create(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id uuid.UUID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

type Vaults interface {
	Create(ctx context.Context, userID uuid.UUID, name, couchDBName string) (Vault, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]Vault, error)
}

type Subscriptions interface {
	Upsert(ctx context.Context, userID uuid.UUID, active bool) (Subscription, error)
	GetByUser(ctx context.Context, userID uuid.UUID) (Subscription, error)
}

type CouchCredentials interface {
	Store(ctx context.Context, vaultID uuid.UUID, host, username string, passwordPlain []byte) error
	Load(ctx context.Context, vaultID uuid.UUID) (CouchCred, error)
}
```

### Encryption Key Config

Add to `config/config.yaml` under `environment:`:

```yaml
- name: creds_encryption_key
  type: string
  value: ""
```

Then run `rscli-dev project tidy`.

### Go Coding Rules (from CLAUDE.md)

- Never check errors inline — always `err := f()` then `if err != nil`.
- Never pass struct literals directly in function calls — assign to a variable first.
- Use `rerrors.Wrap(err, "msg")` for all error wrapping.
