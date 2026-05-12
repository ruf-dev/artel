# Task 17 — MCP Key Management: gRPC + grpc-gateway

## Goal

Expose key management (Create/List/Revoke) as a gRPC service with an HTTP gateway, following the exact same pattern as `vaults_api`.

## Proto File

Create `api/mcp_keys.proto` (look at existing proto files in `api/` for the package, go_package, and import conventions):

```protobuf
syntax = "proto3";

package artel_api;

import "google/api/annotations.proto";

option go_package = "github.com/ruf-dev/artel/internal/api/server/artel_api";

service McpKeysAPI {
    rpc CreateMcpKey(CreateMcpKey.Request) returns (CreateMcpKey.Response) {
        option (google.api.http) = {
            post: "/api/vaults/{vault_id}/mcp-keys"
            body: "*"
        };
    }

    rpc ListMcpKeys(ListMcpKeys.Request) returns (ListMcpKeys.Response) {
        option (google.api.http) = {
            get: "/api/vaults/{vault_id}/mcp-keys"
        };
    }

    rpc RevokeMcpKey(RevokeMcpKey.Request) returns (RevokeMcpKey.Response) {
        option (google.api.http) = {
            delete: "/api/vaults/{vault_id}/mcp-keys/{key_id}"
        };
    }
}

message McpKeyInfo {
    string id          = 1;
    string vault_id    = 2;
    string name        = 3;
    string key_preview = 4;
    string created_at  = 5;
}

message CreateMcpKey {
    message Request {
        string vault_id = 1;
        string name     = 2;
    }
    message Response {
        McpKeyInfo key       = 1;
        string     raw_token = 2; // returned once only
    }
}

message ListMcpKeys {
    message Request {
        string vault_id = 1;
    }
    message Response {
        repeated McpKeyInfo keys = 1;
    }
}

message RevokeMcpKey {
    message Request {
        string vault_id = 1;
        string key_id   = 2;
    }
    message Response {}
}
```

## Code Generation

Run `moti g` to generate the pb.go files into `internal/api/server/artel_api/`.

## Transport Implementation

Create `internal/transport/mcp_keys_api/mcp_keys_impl.go` following the same structure as `internal/transport/vaults_api/vaults_impl.go`:

```go
package mcp_keys_api

import (
    "context"
    "net/http"

    "github.com/google/uuid"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "github.com/rs/zerolog/log"
    "go.redsock.ru/rerrors"
    "google.golang.org/grpc"

    "github.com/ruf-dev/artel/internal/api/server/artel_api"
    "github.com/ruf-dev/artel/internal/service"
)

type McpKeysImpl struct {
    artel_api.UnimplementedMcpKeysAPIServer
    mcpSvc service.McpService
}

func NewMcpKeysImpl(mcpSvc service.McpService) *McpKeysImpl

func (m *McpKeysImpl) Register(srv grpc.ServiceRegistrar)

func (m *McpKeysImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler)
```

Implement `CreateMcpKey`, `ListMcpKeys`, `RevokeMcpKey` — each calling the corresponding `service.McpService` method.

`Gateway` must mount at `/api/vaults/` (same root as vaults gateway, grpc-gateway routes by full path).

## Coding Rules

- Never check errors inline.
- Never create struct literals inline in a function call.
- Use `rerrors.Wrap(err, "context")` from `go.redsock.ru/rerrors`.
- No all-caps field names.

## Verification

- `moti g` succeeds and generates `internal/api/server/artel_api/mcp_keys.pb.go` etc.
- `go build ./...` must pass.
