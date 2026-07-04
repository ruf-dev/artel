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

func (s *ServiceImpl) ListTools(_ context.Context) ([]domain.McpToolDef, error) {
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
				fieldName:        {Type: schemaTypeString, Description: "MoM name (each array entry)"},
				fieldAuthor:      {Type: schemaTypeString, Description: "MoM author"},
				fieldDescription: {Type: schemaTypeString, Description: "MoM description"},
			},
			Required: []string{fieldName, fieldAuthor, fieldDescription},
		},
	}
	tools = append(tools, connectionsTool)

	connectionsForTractsTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name: toolConnectionsForTracts,
			Description: "List the caller's external connections (uuid + provider), for discovering " +
				"connection_uuid values to pass into tract step params such as create_tract's connection_uuid",
			Properties: map[string]domain.ToolProperty{},
			Required:   []string{},
		},
		OutputSchema: domain.ToolSchema{
			Properties: map[string]domain.ToolProperty{
				fieldUuid:     {Type: schemaTypeString, Description: "External connection uuid (each array entry)"},
				fieldProvider: {Type: schemaTypeString, Description: "Provider key, e.g. gitlab, trello, email"},
			},
			Required: []string{fieldUuid, fieldProvider},
		},
	}
	tools = append(tools, connectionsForTractsTool)

	return tools, nil
}
