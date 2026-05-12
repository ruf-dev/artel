# Task 18 — Wire MCP into the Application

## Goal

Connect all MCP components into the running app by updating `custom.go`, `pg/impl.go`, and `service/v1/impl.go`.

## Files to Modify

### `internal/repository/pg/impl.go`

1. Add `mcpKeys repository.McpKeyRepository` field to `Repos` struct.
2. Add `McpKeys() repository.McpKeyRepository` accessor method.
3. In `New(...)`, construct the mcp keys repo:
   - Create `internal/repository/pg/repos/mcpkeys/` package with a `New(q *artel_q.Queries) *Repo` constructor.
   - The repo must implement `repository.McpKeyRepository`.
   - Each method wraps the corresponding sqlc-generated query (from task 12).
   - Wire it in `New(...)`: `mcpKeys: mcpkeys.New(q)`.

### `internal/service/v1/impl.go`

Look at how existing services are wired. Add:
- Import `mcpsvc "github.com/ruf-dev/artel/internal/service/v1/mcp"`.
- Add `McpSvc *mcpsvc.McpSvc` field to `Services` struct.
- Add `McpService() service.McpService { return s.McpSvc }` accessor.
- In `New(repo repository.Repository, env config.EnvironmentConfig)`:
  - Construct `mcpsvc.New(repo.McpKeys(), repo.Vaults(), repo.CouchInstances(), encKey)`.
  - The `encKey` is already derived in `New` (look at how existing services use it).
  - Assign to `McpSvc`.

### `internal/app/custom.go`

Add after existing implementations:

```go
mcpKeysImpl := mcp_keys_api.NewMcpKeysImpl(services.McpService())
mcpHandler := mcp_api.NewMcpHandler(services.McpService())

c.Transport.AddImplementation(mcpKeysImpl)
c.Transport.AddHttpHandler("/mcp", mcpHandler)
```

Add imports for the new packages:
- `"github.com/ruf-dev/artel/internal/transport/mcp_keys_api"`
- `"github.com/ruf-dev/artel/internal/transport/mcp_api"`

Also add `pb.McpKeysAPI_CreateMcpKey_FullMethodName` etc. to the ignored auth paths IF the MCP key endpoint should be public — actually no, key management requires session auth so do NOT add to ignored paths.

The `/mcp` HTTP endpoint does its own auth via bearer token (not the gRPC session interceptor), so it must be excluded from the gRPC auth interceptor. The interceptor only applies to gRPC methods, not plain HTTP handlers, so no change needed there.

## Repository Package to Create

`internal/repository/pg/repos/mcpkeys/mcpkeys.go`:

Implement `repository.McpKeyRepository` using sqlc-generated queries. Map `artel_q.MpcKey` rows to `domain.McpKey`. Handle `sql.NullTime` for `revoked_at`: if `Valid`, set `domain.McpKey.RevokedAt` to `&t.Time`, else leave nil.

## Coding Rules

- Never check errors inline.
- Never create struct literals inline in a function call.
- Use `rerrors.Wrap(err, "context")` from `go.redsock.ru/rerrors`.
- No all-caps field names.

## Verification

- `go build ./...` must pass.
- `go test ./...` must pass.
- Service starts without errors when run locally.
