---
id: "09"
title: "gRPC CouchDB Management API"
status: "pending"
model: "qwen2.5-coder:3b"
created: "2026-05-10"
branch: "factory/09-grpc-couchdb-api"
---

# Task 09 — gRPC CouchDB Management API

## Goal

Add two new proto files (`vaults.proto`, `users.proto`) with separate gRPC services, run `moti g` to generate Go code, implement both handlers, and wire them into the app — replacing the current broken HTTP handler wiring.

## Context

- `api/grpc/artel_api.proto` shows the required message pattern: each RPC uses nested `Request`/`Response` sub-messages under a top-level message named after the operation (e.g. `message Version { message Request {}; message Response {} }`). Follow this pattern for every RPC in the new files.
- `moti.yaml` input is `directory: "api/grpc"` — it picks up all `.proto` files automatically. Adding new files there is enough; do not edit `moti.yaml`.
- After editing protos, run `moti g` to regenerate `internal/api/server/artel_api/`. Never hand-edit generated files.
- `internal/api/server/artel_api_impl/` does not exist yet — create implementation files there.
- `internal/transport/grpc.go` — `GrpcImpl` interface: `Register(grpc.ServiceRegistrar)`. `GrpcWithGateway`: `Gateway(ctx, endpoint, opts...) (rootRoute string, handler http.Handler)`. Each new service impl must satisfy `GrpcImpl`; implement `GrpcWithGateway` too so the HTTP gateway registers automatically.
- `internal/transport/manager.go` — `AddImplementation(GrpcImpl...)` wires a handler and its gateway into the cmux.
- `internal/app/custom.go` — currently broken: references `c.httpServer.Start()/Stop()` (field does not exist on `Custom`; manager's methods are `c.Transport.Start()/Stop()`) and calls `c.Transport.AddHttpHandler()` with no args. Fix this wiring: remove the broken calls, add `c.Transport.AddImplementation(vaultsImpl, usersImpl)`, and call `c.Transport.Start()/Stop()`.
- `internal/service/interfaces.go` defines `VaultService` with `CreateVault`/`DeleteVault`. Extend with `GetVault`, `ListVaults`, and add a new `UserService` interface.
- `internal/clients/couchdb/client.go` — the CouchDB client used by service impls; check its existing methods before adding service logic.

### New proto files

**`api/grpc/vaults.proto`** — package `artel_vaults`, service `VaultsAPI`:

```
service VaultsAPI {
  rpc CreateVault(CreateVault.Request) returns (CreateVault.Response) { GET /api/vaults }
  rpc GetVault(GetVault.Request)       returns (GetVault.Response)    { GET /api/vaults/{name} }
  rpc ListVaults(ListVaults.Request)   returns (ListVaults.Response)  { GET /api/vaults }
  rpc DeleteVault(DeleteVault.Request) returns (DeleteVault.Response) { DELETE /api/vaults/{name} }
}

message CreateVault {
  message Request  { string name = 1; }
  message Response { string name = 1; string db_url = 2; }
}
message GetVault {
  message Request  { string name = 1; }
  message Response { string name = 1; string db_url = 2; }
}
message ListVaults {
  message Request  {}
  message Response { repeated VaultItem vaults = 1; }
}
message DeleteVault {
  message Request  { string name = 1; }
  message Response { bool ok = 1; }
}
message VaultItem { string name = 1; string db_url = 2; }
```

**`api/grpc/users.proto`** — package `artel_users`, service `UsersAPI`:

```
service UsersAPI {
  rpc CreateUser(CreateUser.Request) returns (CreateUser.Response) { POST /api/users }
  rpc GetUser(GetUser.Request)       returns (GetUser.Response)    { GET /api/users/{username} }
  rpc UpdateUser(UpdateUser.Request) returns (UpdateUser.Response) { PUT /api/users/{username} }
  rpc DeleteUser(DeleteUser.Request) returns (DeleteUser.Response) { DELETE /api/users/{username} }
}

message CreateUser {
  message Request  { string username = 1; string password = 2; repeated string roles = 3; }
  message Response { string username = 1; repeated string roles = 2; }
}
message GetUser {
  message Request  { string username = 1; }
  message Response { string username = 1; repeated string roles = 2; }
}
message UpdateUser {
  message Request  { string username = 1; string password = 2; repeated string roles = 3; }
  message Response { string username = 1; repeated string roles = 2; }
}
message DeleteUser {
  message Request  { string username = 1; }
  message Response { bool ok = 1; }
}
```

Both files must import `google/api/annotations.proto` and set `go_package` matching the existing convention.

### Service interface extension

```go
// internal/service/interfaces.go
type VaultService interface {
    CreateVault(ctx context.Context, name string) error
    GetVault(ctx context.Context, name string) (Vault, error)
    ListVaults(ctx context.Context) ([]Vault, error)
    DeleteVault(ctx context.Context, name string) error
}

type UserService interface {
    CreateUser(ctx context.Context, username, password string, roles []string) error
    GetUser(ctx context.Context, username string) (User, error)
    UpdateUser(ctx context.Context, username, password string, roles []string) error
    DeleteUser(ctx context.Context, username string) error
}

type Vault struct{ Name, DBURL string }
type User  struct{ Username string; Roles []string }
```

### Go Coding Rules (from CLAUDE.md)

- Never check errors inline — always `err := f()` then `if err != nil`.
- Never pass struct literals directly in function calls — assign to a named variable first.
- Use `rerrors.Wrap(err, "msg")` for all error wrapping.

## Acceptance Criteria

- [ ] `api/grpc/vaults.proto` defines `VaultsAPI` service with 4 RPCs; all messages use nested `Request`/`Response` sub-message pattern
- [ ] `api/grpc/users.proto` defines `UsersAPI` service with 4 RPCs; same nested pattern
- [ ] `moti g` runs without error and regenerates `internal/api/server/artel_api/`
- [ ] `internal/api/server/artel_api_impl/vaults_impl.go` implements `VaultsAPIServer` (all 4 methods); satisfies `GrpcImpl` and `GrpcWithGateway`
- [ ] `internal/api/server/artel_api_impl/users_impl.go` implements `UsersAPIServer` (all 4 methods); satisfies `GrpcImpl` and `GrpcWithGateway`
- [ ] `internal/service/interfaces.go` has updated `VaultService`, new `UserService`, and `Vault`/`User` structs
- [ ] `internal/service/v1/vault/vault.go` implements `GetVault` and `ListVaults` (new methods)
- [ ] `internal/service/v1/users/users.go` implements all `UserService` methods
- [ ] `internal/service/v1/impl.go` exposes `UserService` field
- [ ] `internal/app/custom.go` wiring fixed: `c.Transport.Start()`/`Stop()` used; `AddImplementation(vaultsImpl, usersImpl)` called
- [ ] `go build ./...` passes from repo root
- [ ] `go test ./...` passes from repo root

## Files to Create / Modify

- `api/grpc/vaults.proto` — new
- `api/grpc/users.proto` — new
- `internal/api/server/artel_api/` — **regenerated by `moti g`**, do not edit by hand
- `internal/api/server/artel_api_impl/vaults_impl.go` — new
- `internal/api/server/artel_api_impl/users_impl.go` — new
- `internal/service/interfaces.go` — extend `VaultService`, add `UserService`, add `Vault`/`User` structs
- `internal/service/v1/impl.go` — add `UserService` field
- `internal/service/v1/vault/vault.go` — add `GetVault`, `ListVaults`
- `internal/service/v1/users/users.go` — new
- `internal/app/custom.go` — fix broken wiring; wire both gRPC impls

## Do NOT change

- `internal/app/app.go` — generated, never edit
- `internal/app/config.go` — generated, never edit
- `api/grpc/artel_api.proto` — keep as-is; `Version` RPC stays
- `internal/transport/grpc.go` — keep as-is
- `internal/transport/http.go` — keep as-is
- `internal/transport/manager.go` — keep as-is
- `internal/clients/couchdb/client.go` — keep as-is

## Notes

- `moti g` must be run after every proto edit; the generated `artel_api/` directory is the source of truth for the gRPC interface.
- grpc-gateway registers HTTP JSON routes automatically via `GrpcWithGateway`; no separate HTTP handler code is needed.
- The broken `c.httpServer` references in `custom.go` must be removed entirely.