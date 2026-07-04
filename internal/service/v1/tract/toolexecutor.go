package tract

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
)

// ToolExecutor is the narrow slice of tool-execution capability the engine needs. Defined
// here — not imported from internal/service/v1/mcp — to avoid an import cycle: a later phase
// has mcp depend on tract (agents authoring tracts via MCP builtin tools), so tract must not
// depend on mcp.
type ToolExecutor interface {
	IsBuiltinTool(toolName string) bool
	// ExecuteBuiltinTool runs a builtin (vault) tool as userUuid. Builtin tools are
	// vault-scoped but a tract step only carries a user, so the executor resolves the vault.
	ExecuteBuiltinTool(ctx context.Context, userUuid uuid.UUID, toolName string, params map[string]interface{}) (string, error)
	// ExecuteMomTool runs a MoM tool against a specific external connection. Callers (the
	// engine) must verify connection ownership before calling this — it performs no
	// ownership check of its own.
	ExecuteMomTool(ctx context.Context, connectionUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{}) (string, error)
	// ListBuiltinTools returns the builtin (vault ops + connections) tool catalog — the
	// "artel" half of ListTractTools.
	ListBuiltinTools(ctx context.Context) ([]domain.McpToolDef, error)
}

// McpBuiltinExecutor is the slice of the mcp service the tract engine needs for builtin tool
// dispatch. Satisfied structurally by service.McpService — passed to NewToolExecutor without
// tract importing that package.
type McpBuiltinExecutor interface {
	IsBuiltinTool(toolName string) bool
	ExecuteBuiltinToolForUser(ctx context.Context, userUuid uuid.UUID, toolName string, params map[string]interface{}) (string, error)
	ListTools(ctx context.Context) ([]domain.McpToolDef, error)
}

// MomConnectionExecutor is the slice of the mom service the tract engine needs for MoM tool
// dispatch. Satisfied structurally by service.MomService.
type MomConnectionExecutor interface {
	ExecuteToolForConnection(ctx context.Context, exConnUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{}) (string, error)
}

// toolExecutorAdapter composes the mcp and mom services into the single ToolExecutor the
// engine depends on.
type toolExecutorAdapter struct {
	builtin McpBuiltinExecutor
	mom     MomConnectionExecutor
}

// NewToolExecutor wires the mcp and mom services (construct mcp first, then tract — see
// internal/app/custom.go) into a ToolExecutor for the tract engine.
func NewToolExecutor(builtin McpBuiltinExecutor, mom MomConnectionExecutor) ToolExecutor {
	adapter := &toolExecutorAdapter{
		builtin: builtin,
		mom:     mom,
	}

	return adapter
}

func (a *toolExecutorAdapter) IsBuiltinTool(toolName string) bool {
	return a.builtin.IsBuiltinTool(toolName)
}

func (a *toolExecutorAdapter) ExecuteBuiltinTool(ctx context.Context, userUuid uuid.UUID, toolName string, params map[string]interface{}) (string, error) {
	res, err := a.builtin.ExecuteBuiltinToolForUser(ctx, userUuid, toolName, params)
	if err != nil {
		return "", rerrors.Wrap(err, "Error executing builtin tool")
	}

	return res, nil
}

func (a *toolExecutorAdapter) ExecuteMomTool(ctx context.Context, connectionUuid uuid.UUID,
	mcpName string, toolName string, params map[string]interface{}) (string, error) {
	res, err := a.mom.ExecuteToolForConnection(ctx, connectionUuid, mcpName, toolName, params)
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	return res, nil
}

func (a *toolExecutorAdapter) ListBuiltinTools(ctx context.Context) ([]domain.McpToolDef, error) {
	res, err := a.builtin.ListTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "Error listing tools")
	}

	return res, nil
}
