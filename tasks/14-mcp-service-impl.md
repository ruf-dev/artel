# Task 14 — McpService v1 Implementation

## Goal

Implement `McpService` in `internal/service/v1/mcp/mcp.go`.

## Context

- Existing encryption: `internal/cryptoutil` — `Encrypt(key, plaintext []byte)` / `Decrypt(key, ciphertext []byte)` using AES-GCM.
- The encryption key comes from `EnvironmentConfig.CredsEncryptionKey` (a string, convert to `[]byte`).
- Key hashing for `ResolveKey` lookup: use `golang.org/x/crypto/bcrypt`.
- Raw tokens have the format `artel_vtk_<32 random hex chars>` (use `crypto/rand` + `encoding/hex`).
- `KeyPreview` stores the first 12 characters of the raw token for display (e.g. `artel_vtk_ab`).
- `ResolveKey` must: scan all active keys for the vault — **no**, actually it must look up by trying all active keys (brute-force bcrypt compare per key).
  
  **Better approach:** store a fast lookup token in addition to the bcrypt hash.
  Actually: prepend a non-secret `key_id` prefix to the raw token so we can look up the DB row directly.
  
  Raw token format: `artel_vtk_<uuid_hex>_<32_random_hex>`
  - `uuid_hex` = the `mcp_keys.id` (UUID without dashes), used to look up the row.
  - `32_random_hex` = the secret part that is bcrypt-hashed and stored in `key_hash`.
  
  `ResolveKey` steps:
  1. Parse the token to extract the UUID portion → `keyID`.
  2. Load the row by `GetMcpKeyByID(keyID)`.
  3. If `revoked_at` is not NULL → return error "key revoked".
  4. `bcrypt.CompareHashAndPassword(row.KeyHash, []byte(secretPart))` → error if mismatch.
  5. Load vault + couch instance to build `McpKeyContext` (needs `VaultRepository` and `CouchInstanceRepository`).
  6. Decrypt `CouchInstance.PasswordEnc` using `cryptoutil.Decrypt`.
  7. Return `McpKeyContext`.

## File to Create

`internal/service/v1/mcp/mcp.go`

```go
package mcp

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "go.redsock.ru/rerrors"

    "github.com/ruf-dev/artel/internal/cryptoutil"
    "github.com/ruf-dev/artel/internal/domain"
    "github.com/ruf-dev/artel/internal/repository"
)

const bcryptCost = 12

type McpSvc struct {
    mcpKeys       repository.McpKeyRepository
    vaults        repository.VaultRepository
    couchInstances repository.CouchInstanceRepository
    encKey        []byte
}

func New(
    mcpKeys repository.McpKeyRepository,
    vaults repository.VaultRepository,
    couchInstances repository.CouchInstanceRepository,
    encKey []byte,
) *McpSvc {
    return &McpSvc{
        mcpKeys:        mcpKeys,
        vaults:         vaults,
        couchInstances: couchInstances,
        encKey:         encKey,
    }
}
```

Implement all four methods of `service.McpService`:
- `CreateKey`: generate random bytes, build token string, bcrypt the secret part, store in DB.
- `ListKeys`: delegate to repository.
- `RevokeKey`: delegate to repository.
- `ResolveKey`: parse token → DB lookup → bcrypt verify → load vault + couch instance → decrypt password → return `McpKeyContext`.

## Also Update

`internal/service/v1/impl.go` — add `McpSvc *mcp.McpSvc` field and wire `McpService()` accessor.

## Interfaces to Implement

You must implement the `service.McpService` interface from `internal/service/interfaces.go` exactly.

## Coding Rules

- Never check errors inline.
- Never create struct literals inline in a function call.
- Use `rerrors.Wrap(err, "context")` for error wrapping.
- No all-caps field names.

## Verification

- `go build ./...` must pass.
- `go test ./...` must pass.
