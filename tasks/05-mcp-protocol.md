# Task 05 — MCP Server Protocol

## Goal

Expose artel as an MCP (Model Context Protocol) server so Claude and other MCP clients can invoke vault operations as tool calls over stdio or HTTP/SSE transport.

## Scope

- Add `internal/transport/mcp/` package with an MCP server wrapping the `service.VaultService` interface
- Wire the MCP server in `custom.go` alongside the existing HTTP server
- Expose at minimum two tools: `create_vault` and `delete_vault`
- Use the official Go MCP SDK (`github.com/mark3labs/mcp-go` or equivalent)

## Acceptance Criteria

- `go build ./...` passes
- MCP server starts and registers tools at startup (logged)
- A client can call `create_vault` and receive success/error response
- Transport is stdio by default; HTTP/SSE can be a config flag

## Config

Add `mcp_transport` (string: `stdio` | `sse`) and `mcp_sse_addr` (string, e.g. `:8081`) to `config/config.yaml`, then run `rscli-dev project tidy` to get fields in `EnvironmentConfig`.

## Notes

- Transport handler must accept `service.VaultService` interface, not concrete type
- Follow same Init/Start/Stop pattern as `internal/transport/http/`
- Tool input schemas should be defined as Go structs with JSON tags
