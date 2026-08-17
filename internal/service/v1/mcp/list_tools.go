package mcp

import (
	"context"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
)

const (
	toolConnections              = "connections"
	toolConnectionsForTracts     = "list_connections_for_tracts"
	toolCreateCommunityConnector = "create_community_connector"
)

func (s *ServiceImpl) ListTools(_ context.Context) ([]domain.McpToolDef, error) {
	tools := executors.VaultToolDefinitions()
	tools = append(tools, executors.TractToolDefinitions()...)
	// Postgres tools are always advertised, regardless of whether the vault actually has a
	// Postgres database enabled — mirrors write_file being listed even without a linked S3
	// bucket; ExecuteTool is where keyCtx.Postgres == nil is checked and fails at call time.
	tools = append(tools, executors.PostgresToolDefinitions()...)
	connectionsTool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        toolConnections,
			Description: "List the MoMs (Mcp of Mcp tool packages, e.g. email) connected to this key",
			Properties:  map[string]domain.ToolProperty{},
			Required:    []string{},
		},
		OutputSchema: domain.ToolSchema{
			IsArray: true,
			Items: &domain.ToolProperty{
				Type: schemaTypeObject,
				Properties: map[string]domain.ToolProperty{
					fieldName:        {Type: schemaTypeString, Description: "MoM name"},
					fieldAuthor:      {Type: schemaTypeString, Description: "MoM author"},
					fieldDescription: {Type: schemaTypeString, Description: "MoM description"},
				},
				Required: []string{fieldName, fieldAuthor, fieldDescription},
			},
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
			IsArray: true,
			Items: &domain.ToolProperty{
				Type: schemaTypeObject,
				Properties: map[string]domain.ToolProperty{
					fieldUuid:     {Type: schemaTypeString, Description: "External connection uuid"},
					fieldProvider: {Type: schemaTypeString, Description: "Provider key, e.g. gitlab, trello, email"},
				},
				Required: []string{fieldUuid, fieldProvider},
			},
		},
	}
	tools = append(tools, connectionsForTractsTool)

	tools = append(tools, createCommunityConnectorTool())

	tools = append(tools, skillsToolDefinitions()...)

	return tools, nil
}

// createCommunityConnectorTool describes create_community_connector's input contract, which
// mirrors the MoM JSON shape documented in pkg/mom_examples/README.md exactly: a {name, author,
// description, tools[]} record where each tool is {api_description, action}.
func createCommunityConnectorTool() domain.McpToolDef {
	httpActionProps := map[string]domain.ToolProperty{
		"method": {Type: schemaTypeString, Description: "HTTP method, e.g. GET, POST, PUT, PATCH, DELETE"},
		"url": {
			Type: schemaTypeString,
			Description: "Request URL; may contain ${{params.name}} placeholders interpolated from the " +
				"tool call's input arguments",
		},
		"headers": {
			Type: schemaTypeObject,
			Description: "Static string headers; values may be __secrets.field to inject a value from " +
				"the linked external connection's decrypted credentials",
		},
		"query": {
			Type: schemaTypeObject,
			Description: "URL query parameters; values may be static, ${{params.name}}, or " +
				"__secrets.field",
		},
		"body": {
			Type:        schemaTypeObject,
			Description: "JSON request body template; string leaves support the same ${{params.*}}/__secrets.* interpolation as headers/query",
		},
		"credentials": {
			Type: schemaTypeString,
			Description: "external_connections.provider to resolve __secrets.* against — normally set " +
				"to this connector's own name, so a caller's AddGenericConnection(provider=<name>, ...) " +
				"call resolves these secrets",
		},
	}

	actionProps := map[string]domain.ToolProperty{
		"http": {
			Type:        schemaTypeObject,
			Description: "HTTP executor definition — set this for a plain REST call",
			Properties:  httpActionProps,
			Required:    []string{"method", "url"},
		},
		"imap": {
			Type:        schemaTypeObject,
			Description: "IMAP executor definition — set this to reuse Artel's built-in email client",
			Properties: map[string]domain.ToolProperty{
				"operation": {Type: schemaTypeString, Description: "list_folders | list_messages | fetch_message"},
			},
			Required: []string{"operation"},
		},
		"smtp": {
			Type:        schemaTypeObject,
			Description: "SMTP executor definition — set this to reuse Artel's built-in email client",
			Properties: map[string]domain.ToolProperty{
				"operation": {Type: schemaTypeString, Description: "send"},
			},
			Required: []string{"operation"},
		},
	}

	apiDescriptionProps := map[string]domain.ToolProperty{
		fieldName:        {Type: schemaTypeString, Description: "Tool name, unique within this connector"},
		fieldDescription: {Type: schemaTypeString, Description: "What the LLM sees describing this tool"},
		"properties": {
			Type: schemaTypeObject,
			Description: "Input schema properties, keyed by param name, each a JSON-Schema-shaped " +
				"{type, description, ...} object — the standard MCP tool input schema",
		},
		"required": {
			Type:        "array",
			Description: "Names of required input params",
			Items:       &domain.ToolProperty{Type: schemaTypeString},
		},
	}

	toolItemSchema := domain.ToolProperty{
		Type: schemaTypeObject,
		Properties: map[string]domain.ToolProperty{
			"api_description": {
				Type:        schemaTypeObject,
				Description: "LLM-visible metadata for this tool",
				Properties:  apiDescriptionProps,
				Required:    []string{fieldName, fieldDescription, "properties", "required"},
			},
			"action": {
				Type:        schemaTypeObject,
				Description: "Executor definition — must contain exactly one of http, imap, or smtp",
				Properties:  actionProps,
				Required:    []string{},
			},
		},
		Required: []string{"api_description", "action"},
	}

	description := "Admin-only: persist a community MoM connector (a set of MoM-shaped tool " +
		"definitions produced from a user's Swagger/OpenAPI spec) so it becomes available to " +
		"connect. Only administrators may call this successfully — CheckIsAdmin is enforced " +
		"server-side. Calling again with the same name re-defines (upserts) your own existing " +
		"connector, but is rejected if that name already belongs to a system MoM or to a " +
		"different admin's connector. For http-action tools, set action.http.credentials to " +
		"this connector's own name so that a caller's separate AddGenericConnection(provider=<name>, " +
		"...) call resolves __secrets.* placeholders against the right external connection."

	return domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        toolCreateCommunityConnector,
			Description: description,
			Properties: map[string]domain.ToolProperty{
				fieldName:        {Type: schemaTypeString, Description: "Unique name/primary key for this MoM"},
				fieldDescription: {Type: schemaTypeString, Description: "Human-readable description of this MoM"},
				"tools": {
					Type:        "array",
					Description: "The MoM's tool definitions",
					Items:       &toolItemSchema,
				},
			},
			Required: []string{fieldName, fieldDescription, "tools"},
		},
	}
}
