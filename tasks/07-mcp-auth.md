# Task 07 — MCP Auth

## Goal

Protect MCP tool calls with per-user authentication so each caller operates only on their own vaults.

## Scope

- Add bearer token validation middleware to the MCP HTTP/SSE transport
- Store user→token mappings in CouchDB (a `_users`-style database or a dedicated `artel_auth` db)
- Add MCP tools: `register_user`, `issue_token` (admin-only)
- Propagate user identity into service calls so vault/note operations are scoped to the caller

## Acceptance Criteria

- Unauthenticated MCP requests over SSE are rejected with a 401
- Stdio transport (local/trusted) can bypass auth via config flag `mcp_auth_required: false`
- `register_user` + `issue_token` work end-to-end with a CouchDB-backed store
- `create_vault` uses the authenticated user as the vault owner; other users cannot access it

## Config

Add `mcp_auth_required` (bool) to `config/config.yaml` + `rscli-dev project tidy`.

## Notes

- Token validation should be a middleware/interceptor, not inline in each tool handler
- Tokens are opaque random strings (UUID v4); no JWT needed at this stage
- Admin credentials come from `EnvironmentConfig` (add `admin_token` env var)
