package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

const toolConnections = "connections"

func (s *McpServiceImpl) ListTools(ctx context.Context) ([]domain.McpToolDef, error) {
	tools := executors.VaultToolDefinitions()
	tools = append(tools, domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        toolConnections,
			Description: "List the MoMs (Mcp of Mcp tool packages, e.g. email) connected to this key",
			Properties:  map[string]domain.ToolProperty{},
			Required:    []string{},
		},
	})
	return tools, nil
}
