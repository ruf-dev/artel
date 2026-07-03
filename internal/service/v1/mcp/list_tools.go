package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

const (
	toolConnections          = "connections"
	toolConnectionsForTracts = "list_connections_for_tracts"
)

func (s *McpServiceImpl) ListTools(ctx context.Context) ([]domain.McpToolDef, error) {
	tools := executors.VaultToolDefinitions()
	tools = append(tools, executors.TractToolDefinitions()...)
	connectionsTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        toolConnections,
			Description: "List the MoMs (Mcp of Mcp tool packages, e.g. email) connected to this key",
			Properties:  map[string]domain.ToolProperty{},
			Required:    []string{},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				"name":        {Type: "string", Description: "MoM name (each array entry)"},
				"author":      {Type: "string", Description: "MoM author"},
				"description": {Type: "string", Description: "MoM description"},
			},
			Required: []string{"name", "author", "description"},
		},
	}
	tools = append(tools, connectionsTool)

	connectionsForTractsTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        toolConnectionsForTracts,
			Description: "List the caller's external connections (uuid + provider), for discovering connection_uuid values to pass into tract step params such as create_tract's connection_uuid",
			Properties:  map[string]domain.ToolProperty{},
			Required:    []string{},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				"uuid":     {Type: "string", Description: "External connection uuid (each array entry)"},
				"provider": {Type: "string", Description: "Provider key, e.g. gitlab, trello, email"},
			},
			Required: []string{"uuid", "provider"},
		},
	}
	tools = append(tools, connectionsForTractsTool)

	return tools, nil
}
