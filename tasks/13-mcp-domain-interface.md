# Task 13 — MCP Domain Types + Service Interface

## Goal

Add domain structs for MCP keys and extend the service interface.

## Domain Types

Add file `internal/domain/mcp_key.go`:

```go
package domain

import (
    "time"

    "github.com/google/uuid"
)

type McpKey struct {
    Uuid       uuid.UUID
    VaultUuid  uuid.UUID
    UserUuid   uuid.UUID
    Name       string
    KeyPreview string
    CreatedAt  time.Time
    RevokedAt  *time.Time
}

// McpKeyContext is resolved from a raw bearer token; contains everything
// needed to connect to CouchDB for the associated vault.
type McpKeyContext struct {
    KeyUuid   uuid.UUID
    VaultUuid uuid.UUID
    UserUuid  uuid.UUID
    CouchURL  string
    CouchDb   string
    CouchUser string
    CouchPass string
}
```

## Repository Interface Extension

In `internal/repository/interfaces.go`, add a new interface:

```go
type McpKeyRepository interface {
    CreateMcpKey(ctx context.Context, vaultID, userID uuid.UUID, name string, keyHash []byte, keyPreview string) (domain.McpKey, error)
    ListMcpKeysByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
    GetMcpKeyByID(ctx context.Context, id uuid.UUID) (domain.McpKey, error)
    ListActiveMcpKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
    RevokeMcpKey(ctx context.Context, id uuid.UUID) error
}
```

Also add `McpKeyRepository() McpKeyRepository` method to the existing `Repository` interface.

## Service Interface Extension

In `internal/service/interfaces.go`, add:

```go
type McpService interface {
    // CreateKey generates a new bearer token, stores it hashed, returns the raw token once.
    CreateKey(ctx context.Context, vaultID uuid.UUID, name string) (rawToken string, key domain.McpKey, err error)
    ListKeys(ctx context.Context, vaultID uuid.UUID) ([]domain.McpKey, error)
    RevokeKey(ctx context.Context, keyID uuid.UUID) error
    // ResolveKey validates the raw bearer token and returns vault+couch context.
    ResolveKey(ctx context.Context, rawToken string) (domain.McpKeyContext, error)
}
```

Also add `McpService() McpService` to the existing `Service` interface.

## Coding Rules

- Never check errors inline; always `err := f()` then `if err != nil`.
- Never create struct literals inline in a function call; assign to a variable first.
- Use `rerrors.Wrap(err, "context")` from `go.redsock.ru/rerrors` for error wrapping.
- No all-caps field names: `Uuid` not `UUID`, `Id` not `ID`.

## Verification

- `go build ./...` must pass.
