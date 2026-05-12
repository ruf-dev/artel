# Task 16 — MCP HTTP Transport (JSON-RPC 2.0)

## Goal

Implement a pure `http.Handler` that speaks the MCP (Model Context Protocol) wire protocol: JSON-RPC 2.0 over HTTP POST at `/mcp`.

This is NOT a gRPC service. It registers directly on the HTTP server via `AddHttpHandler`.

## Wire Protocol

Every request body is a JSON-RPC 2.0 object:
```json
{ "jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": { ... } }
```

Every response:
```json
{ "jsonrpc": "2.0", "id": 1, "result": { ... } }
```

Or error:
```json
{ "jsonrpc": "2.0", "id": 1, "error": { "code": -32601, "message": "method not found" } }
```

Notifications (no `id`) must not be responded to.

## Authentication

Extract `Authorization: Bearer <token>` from the request header.

- If missing: respond with HTTP 401, no JSON-RPC body.
- If present: call `mcpSvc.ResolveKey(ctx, token)` to get `domain.McpKeyContext`.
- If ResolveKey fails: respond with JSON-RPC error `{ "code": -32001, "message": "unauthorized" }`.

Exception: the `initialize` method is allowed without a valid key (respond with capabilities before the key is set). Actually — require auth for ALL methods including `initialize` for simplicity.

## Methods to Handle

### `initialize`
Request params:
```json
{ "protocolVersion": "2025-06-18", "capabilities": {}, "clientInfo": { "name": "...", "version": "..." } }
```

Response result:
```json
{
  "protocolVersion": "2025-06-18",
  "capabilities": { "tools": {} },
  "serverInfo": { "name": "artel-mcp", "version": "0.1.0" }
}
```

### `notifications/initialized`
This is a notification (no `id`). Ignore it (return 200 with empty body or `{}`).

### `tools/list`
Response result:
```json
{
  "tools": [
    { "name": "list_notes",       "description": "List all notes in the vault", "inputSchema": { "type": "object", "properties": {}, "required": [] } },
    { "name": "read_note",        "description": "Read a note by path",          "inputSchema": { "type": "object", "properties": { "path": { "type": "string" } }, "required": ["path"] } },
    { "name": "write_note",       "description": "Create or update a note",      "inputSchema": { "type": "object", "properties": { "path": { "type": "string" }, "content": { "type": "string" } }, "required": ["path", "content"] } },
    { "name": "delete_note",      "description": "Delete a note by path",        "inputSchema": { "type": "object", "properties": { "path": { "type": "string" } }, "required": ["path"] } },
    { "name": "move_note",        "description": "Move a note to a new path",    "inputSchema": { "type": "object", "properties": { "old_path": { "type": "string" }, "new_path": { "type": "string" } }, "required": ["old_path", "new_path"] } },
    { "name": "list_folders",     "description": "List all folders in the vault","inputSchema": { "type": "object", "properties": {}, "required": [] } },
    { "name": "list_tags",        "description": "List all tags in the vault",   "inputSchema": { "type": "object", "properties": {}, "required": [] } },
    { "name": "get_note_metadata","description": "Get metadata for a note",      "inputSchema": { "type": "object", "properties": { "path": { "type": "string" } }, "required": ["path"] } }
  ]
}
```

### `tools/call`
Request params:
```json
{ "name": "read_note", "arguments": { "path": "folder/note.md" } }
```

Response result:
```json
{ "content": [ { "type": "text", "text": "<tool output as string>" } ] }
```

Route `name` to the appropriate `LiveSyncClient` method using `domain.McpKeyContext` for connection details.

## Directory Structure

```
internal/transport/mcp_api/
    handler.go      -- main http.Handler; JSON-RPC decode → dispatch → encode
    tools.go        -- tool definitions (tools/list response, tools/call dispatch)
    auth.go         -- Bearer token extraction + ResolveKey call
```

## Key Types

```go
// in handler.go
type rpcRequest struct {
    Jsonrpc string          `json:"jsonrpc"`
    Id      any             `json:"id"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params"`
}

type rpcResponse struct {
    Jsonrpc string `json:"jsonrpc"`
    Id      any    `json:"id,omitempty"`
    Result  any    `json:"result,omitempty"`
    Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

## Handler Constructor

```go
type McpHandler struct {
    mcpSvc service.McpService
}

func NewMcpHandler(mcpSvc service.McpService) *McpHandler

func (h *McpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

## Coding Rules

- Never check errors inline.
- Never create struct literals inline in a function call.
- Use `rerrors.Wrap` from `go.redsock.ru/rerrors`.
- No all-caps field names.
- Set `Content-Type: application/json` on all responses.

## Verification

- `go build ./...` must pass.
