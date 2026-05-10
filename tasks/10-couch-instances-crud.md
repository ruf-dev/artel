---
id: "10"
title: "CouchDB Instances CRUD API"
status: "done"
---

## Goal

Expand the CouchDB instances API from a single RegisterCouchInstance endpoint to full CRUD. Change the register route from `/api/couch-instances` to `/api/couch/add`. Add Get, List, Delete endpoints.

## Context

The proto is at `api/grpc/couch_instances.proto`. Generated Go lives in `internal/api/server/artel_api/`. The repo layer is in `internal/repository/pg/couchinstances/couchinstances.go`. Service layer in `internal/service/v1/couchinstances/couchinstances.go`. Transport in `internal/transport/couch_instances_api/couch_instances_impl.go`. SQL queries in `internal/repository/pg/queries/couch_instances.sql`. Repo interface in `internal/repository/interfaces.go`. Service interface in `internal/service/interfaces.go`.

The `couch_instances` table has columns: `id UUID`, `url TEXT`, `username TEXT`, `password_enc BYTEA`, `created_at TIMESTAMPTZ`.

IMPORTANT: Do NOT create git branches. Do NOT commit. Do NOT run moti g or sqlc generate — the generated files must be written by hand based on existing generated file patterns.

## Files to Modify / Create

### 1. `api/grpc/couch_instances.proto`
Change RegisterCouchInstance route from `post: "/api/couch-instances"` to `post: "/api/couch/add"`. Add three new RPCs:
- `GetCouchInstance`: `GET /api/couch/{id}`
- `ListCouchInstances`: `GET /api/couch`
- `DeleteCouchInstance`: `DELETE /api/couch/{id}`

New messages needed:
```
message GetCouchInstance {
  message Request { string id = 1; }
  message Response { string id = 1; string url = 2; string username = 3; string created_at = 4; }
}
message ListCouchInstances {
  message Request {}
  message Response { repeated GetCouchInstance.Response instances = 1; }
}
message DeleteCouchInstance {
  message Request { string id = 1; }
  message Response {}
}
```

### 2. `internal/repository/pg/queries/couch_instances.sql`
Add three queries:
```sql
-- name: GetCouchInstance :one
SELECT id, url, username, created_at FROM couch_instances WHERE id = $1;

-- name: ListCouchInstances :many
SELECT id, url, username, created_at FROM couch_instances ORDER BY created_at DESC;

-- name: DeleteCouchInstance :exec
DELETE FROM couch_instances WHERE id = $1;
```

### 3. `internal/repository/pg/generated/couch_instances.sql.go`
Add the generated Go for the three new queries. Look at existing generated files (e.g. `internal/repository/pg/generated/vaults.sql.go`) as the pattern. The new rows type should be:
```go
type GetCouchInstanceRow struct {
    Id        uuid.UUID
    Url       string
    Username  string
    CreatedAt time.Time
}
```
Methods to add to `*Queries`:
- `GetCouchInstance(ctx, id uuid.UUID) (GetCouchInstanceRow, error)`
- `ListCouchInstances(ctx) ([]GetCouchInstanceRow, error)`
- `DeleteCouchInstance(ctx, id uuid.UUID) error`

### 4. `internal/repository/interfaces.go`
Expand `CouchInstances` interface:
```go
type CouchInstances interface {
    Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error)
    Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error)
    List(ctx context.Context) ([]domain.CouchInstance, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### 5. `internal/repository/pg/couchinstances/couchinstances.go`
Add three methods:
- `Get(ctx, id uuid.UUID) (domain.CouchInstance, error)` — calls `q.GetCouchInstance`, maps to domain
- `List(ctx) ([]domain.CouchInstance, error)` — calls `q.ListCouchInstances`, maps slice to domain
- `Delete(ctx, id uuid.UUID) error` — calls `q.DeleteCouchInstance`

### 6. `internal/service/interfaces.go`
Expand `CouchInstanceService`:
```go
type CouchInstanceService interface {
    RegisterCouchInstance(ctx context.Context, url, username, password string) (string, error)
    GetCouchInstance(ctx context.Context, id string) (domain.CouchInstance, error)
    ListCouchInstances(ctx context.Context) ([]domain.CouchInstance, error)
    DeleteCouchInstance(ctx context.Context, id string) error
}
```

### 7. `internal/service/v1/couchinstances/couchinstances.go`
Add three methods:
- `GetCouchInstance(ctx, id string) (domain.CouchInstance, error)` — parse uuid from id string, call repo.Get
- `ListCouchInstances(ctx) ([]domain.CouchInstance, error)` — call repo.List
- `DeleteCouchInstance(ctx, id string) error` — parse uuid from id string, call repo.Delete

### 8. `internal/api/server/artel_api/couch_instances.pb.go` (and related generated files)
Add the new message structs for the three new RPCs. Look at existing structs in the file as pattern. Add:
- `GetCouchInstance_Request`, `GetCouchInstance_Response`
- `ListCouchInstances_Request`, `ListCouchInstances_Response`
- `DeleteCouchInstance_Request`, `DeleteCouchInstance_Response`

### 9. `internal/api/server/artel_api/couch_instances_grpc.pb.go`
Add the new RPC methods to `CouchInstancesAPIServer` interface and `UnimplementedCouchInstancesAPIServer`. Add handler registrations.

### 10. `internal/api/server/artel_api/couch_instances.pb.gw.go`
Add gateway handler registrations for the three new routes.

### 11. `internal/transport/couch_instances_api/couch_instances_impl.go`
Add handler methods:
- `GetCouchInstance(ctx, req *artel_api.GetCouchInstance_Request) (*artel_api.GetCouchInstance_Response, error)`
- `ListCouchInstances(ctx, req *artel_api.ListCouchInstances_Request) (*artel_api.ListCouchInstances_Response, error)`
- `DeleteCouchInstance(ctx, req *artel_api.DeleteCouchInstance_Request) (*artel_api.DeleteCouchInstance_Response, error)`

### 12. `examples/artel_api.http`
Change route to `/api/couch/add`. Add examples:
```http
### Get CouchDB instance
GET {{artel_url}}/api/couch/{{couch_instance_id}}

### List CouchDB instances
GET {{artel_url}}/api/couch

### Delete CouchDB instance
DELETE {{artel_url}}/api/couch/{{couch_instance_id}}
```

### 13. `examples/http-client.env.json`
Add `"couch_instance_id": "00000000-0000-0000-0000-000000000000"` to the dev environment.

## Acceptance Criteria

- [x] Proto has four RPCs: RegisterCouchInstance (POST /api/couch/add), GetCouchInstance (GET /api/couch/{id}), ListCouchInstances (GET /api/couch), DeleteCouchInstance (DELETE /api/couch/{id})
- [x] SQL queries file has all four queries
- [x] Generated repo code has all four methods
- [x] `CouchInstances` repo interface has Register, Get, List, Delete
- [x] `CouchInstanceService` service interface has all four methods
- [x] Repo impl and service impl have all four methods
- [x] Transport handler has all four RPC methods
- [x] Generated proto Go files have structs and handlers for all four RPCs
- [x] `go build ./...` passes with no errors
- [x] `examples/artel_api.http` has all four example calls
