package executors

import (
	"net/http"
	"time"
)

// ToolIdent identifies which MoM tool an outbound HTTP request belongs to. It is an ordinary
// parameter of HttpExecutor.Execute (never passed through context) and is used at request time
// only to pick the pre-wound http.Client for that tool.
type ToolIdent struct {
	McpName  string
	ToolName string
}

// HttpMiddleware wraps an http.RoundTripper, returning a RoundTripper that runs custom logic
// around the wrapped one.
type HttpMiddleware func(next http.RoundTripper) http.RoundTripper

// middlewareReg is one WithToolHttpMiddleware / WithMcpHttpMiddleware call, captured so the
// executor can wind per-tool and per-mcp http.Clients once, at construction.
type middlewareReg struct {
	tool McpTool // "" when the registration is mcp-scoped
	mcp  McpName // "" when the registration is tool-scoped
	wrap HttpMiddleware
}

// HttpExecutorOption configures a HttpExecutor at construction time.
type HttpExecutorOption func(*httpExecutorConfig)

type httpExecutorConfig struct {
	regs []middlewareReg
}

// WithToolHttpMiddleware registers mw around every outbound request of exactly one tool
// ("<mcp>.<tool>"). Registrations stack LIFO: the last one registered for a given scope is the
// outermost wrapper and runs first.
func WithToolHttpMiddleware(tool McpTool, mw HttpMiddleware) HttpExecutorOption {
	return func(cfg *httpExecutorConfig) {
		reg := middlewareReg{tool: tool, wrap: mw}

		cfg.regs = append(cfg.regs, reg)
	}
}

// WithMcpHttpMiddleware registers mw around every outbound request of every tool of one mcp. It
// also applies to a tool that has its own WithToolHttpMiddleware registration — the mcp-scoped
// wrappers sit inside (closer to the base transport than) that tool's own. mcp-scoped
// registrations stack LIFO too.
func WithMcpHttpMiddleware(mcp McpName, mw HttpMiddleware) HttpExecutorOption {
	return func(cfg *httpExecutorConfig) {
		reg := middlewareReg{mcp: mcp, wrap: mw}

		cfg.regs = append(cfg.regs, reg)
	}
}

// buildHttpClients winds the per-tool and per-mcp http.Clients once, from base, applying every
// registration in cfg. Tools and mcps with no registration are served by the shared base client
// (with no registration at all it is byte-identical to using base directly). timeout is copied
// onto every client so behaviour matches the pre-middleware single-client executor.
func buildHttpClients(
	cfg *httpExecutorConfig, base http.RoundTripper, timeout time.Duration,
) (baseClient *http.Client, byTool map[McpTool]*http.Client, byMcp map[McpName]*http.Client) {
	baseClient = &http.Client{Timeout: timeout, Transport: base}

	toolKeys := map[McpTool]struct{}{}
	mcpKeys := map[McpName]struct{}{}

	for _, reg := range cfg.regs {
		if reg.tool != "" {
			toolKeys[reg.tool] = struct{}{}
		}

		if reg.mcp != "" {
			mcpKeys[reg.mcp] = struct{}{}
		}
	}

	byMcp = make(map[McpName]*http.Client, len(mcpKeys))

	for mcp := range mcpKeys {
		rt := windMcp(cfg.regs, base, mcp)

		byMcp[mcp] = &http.Client{Timeout: timeout, Transport: rt}
	}

	byTool = make(map[McpTool]*http.Client, len(toolKeys))

	for tool := range toolKeys {
		rt := windTool(cfg.regs, base, tool)

		byTool[tool] = &http.Client{Timeout: timeout, Transport: rt}
	}

	return baseClient, byTool, byMcp
}

// windMcp stacks every mcp-scoped registration for mcp onto base, in registration order, so the
// last-registered one ends up outermost.
func windMcp(regs []middlewareReg, base http.RoundTripper, mcp McpName) http.RoundTripper {
	rt := base

	for _, reg := range regs {
		if reg.mcp == mcp {
			rt = reg.wrap(rt)
		}
	}

	return rt
}

// windTool stacks the mcp-scoped registrations for the tool's mcp (in registration order) then
// the tool-scoped registrations for that exact tool (in registration order) onto base, so the
// last tool-scoped registration is outermost and the mcp-scoped wrappers run innermost.
func windTool(regs []middlewareReg, base http.RoundTripper, tool McpTool) http.RoundTripper {
	mcp := mcpOf(tool)

	rt := base

	for _, reg := range regs {
		if reg.mcp != "" && reg.mcp == mcp {
			rt = reg.wrap(rt)
		}
	}

	for _, reg := range regs {
		if reg.tool == tool {
			rt = reg.wrap(rt)
		}
	}

	return rt
}
