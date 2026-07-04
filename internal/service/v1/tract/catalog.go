package tract

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// ListTractTools returns the tract action picker's tool catalog: builtins (paired with the
// "artel" sentinel Mcp name) plus every mcp_tools row across every MoM.
func (s *Service) ListTractTools(ctx context.Context) ([]domain.McpToolRef, error) {
	_, ok := user_context.GetUserContext(ctx)
	if !ok {
		return nil, rerrors.Wrap(user_errors.Unauthenticated)
	}

	builtins, err := s.toolExecutor.ListBuiltinTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing builtin tools")
	}

	momTools, err := s.mcpDefs.ListAllTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing mom tools")
	}

	catalog := make([]domain.McpToolRef, 0, len(builtins)+len(momTools))

	for _, tool := range builtins {
		ref := domain.McpToolRef{McpName: builtinMcpName, Tool: tool}
		catalog = append(catalog, ref)
	}

	catalog = append(catalog, momTools...)

	return catalog, nil
}
