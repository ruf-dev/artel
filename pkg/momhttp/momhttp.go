// Package momhttp is the public surface for injecting per-mcp / per-tool HTTP middleware into
// artel's MoM http executor from outside the module. Pair it with pkg/app.WithToolHttpMiddleware
// / pkg/app.WithMcpHttpMiddleware and pkg/app.New.
//
// pkg/ importing internal/ within the same module is allowed (pkg/app already does it) — Go's
// internal rule only blocks importers outside the module.
package momhttp

import (
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

// Re-exported identifiers, so embedders never import internal/ directly.
//
// McpName / McpToolName are the sqlc-generated Go types for the Postgres mcp_name / mcp_tool_name
// enums (migration 080). Scope a middleware to one tool with McpToolName("<mcp>.<tool>") — e.g.
// McpToolName("trello.create_card") — or to a whole mcp with one of the McpName* constants.
type (
	McpName     = artel_q.McpName
	McpToolName = artel_q.McpToolName
	Middleware  = executors.HttpMiddleware
	ToolIdent   = executors.ToolIdent
)

// Re-exported MoM tool-set identifiers (the Postgres mcp_name enum labels).
const (
	McpNameEmail    = artel_q.McpNameEmail
	McpNameGitlab   = artel_q.McpNameGitlab
	McpNameTelegram = artel_q.McpNameTelegram
	McpNameTrello   = artel_q.McpNameTrello
)
