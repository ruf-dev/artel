package simplechat

import (
	"encoding/json"

	"github.com/ruf-dev/artel/internal/clients/openai"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
	"go.redsock.ru/rerrors"
)

// mcpToolToOpenAIToolDef converts one MCP tool definition into the OpenAI function-tool shape.
// domain.ToolProperty carries no json tags, so the JSON-Schema object is assembled by hand
// field-by-field rather than marshalled straight off the domain struct.
func mcpToolToOpenAIToolDef(tool domain.McpToolDef) (openai.ToolDefinition, error) {
	schema := objectSchema(tool.ApiDescription.Properties, tool.ApiDescription.Required)

	raw, err := json.Marshal(schema)
	if err != nil {
		return openai.ToolDefinition{}, rerrors.Wrap(err, "error marshalling tool parameters schema")
	}

	def := openai.ToolDefinition{
		Name:           tool.ApiDescription.Name,
		Description:    tool.ApiDescription.Description,
		ParametersJSON: raw,
	}

	return def, nil
}

// objectSchema builds a JSON-Schema "object" node from a property map plus its required list.
// "properties" is always emitted (as an empty object when the tool takes no parameters) since
// a function tool whose schema omits it is rejected by some providers.
func objectSchema(properties map[string]domain.ToolProperty, required []string) map[string]any {
	props := make(map[string]any, len(properties))
	for name, prop := range properties {
		props[name] = toolPropertyToJSONSchema(prop)
	}

	schema := map[string]any{
		"type":       "object",
		"properties": props,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// toolPropertyToJSONSchema recursively renders one ToolProperty as a JSON-Schema node: nested
// object fields via Properties (+ their own Required), array elements via Items.
func toolPropertyToJSONSchema(prop domain.ToolProperty) map[string]any {
	node := map[string]any{}

	if prop.Type != "" {
		node["type"] = prop.Type
	}

	if prop.Description != "" {
		node["description"] = prop.Description
	}

	if len(prop.Enum) > 0 {
		node["enum"] = prop.Enum
	}

	if len(prop.Properties) > 0 {
		nested := make(map[string]any, len(prop.Properties))
		for name, nestedProp := range prop.Properties {
			nested[name] = toolPropertyToJSONSchema(nestedProp)
		}

		node["properties"] = nested
	}

	if len(prop.Required) > 0 {
		node["required"] = prop.Required
	}

	if prop.Items != nil {
		node["items"] = toolPropertyToJSONSchema(*prop.Items)
	}

	return node
}

// vaultToolNames is the set of tool names that read or write vault notes — the tools a chat
// created with vault_access=false must not be offered. Derived from
// executors.VaultToolDefinitions so the two can never drift.
func vaultToolNames() map[string]struct{} {
	defs := executors.VaultToolDefinitions()

	names := make(map[string]struct{}, len(defs))
	for _, def := range defs {
		names[def.ApiDescription.Name] = struct{}{}
	}

	return names
}

// filterToolsForChat drops the vault-note tools from tools when the chat was created without
// vault access. The full catalogue (tract, postgres, skills, connectors, ...) is otherwise
// passed through untouched — this never reimplements or duplicates the tool list, it only
// filters the one ListTools returned.
func filterToolsForChat(tools []domain.McpToolDef, vaultAccess bool) []domain.McpToolDef {
	if vaultAccess {
		return tools
	}

	excluded := vaultToolNames()

	filtered := make([]domain.McpToolDef, 0, len(tools))
	for _, tool := range tools {
		_, isVaultTool := excluded[tool.ApiDescription.Name]
		if isVaultTool {
			continue
		}

		filtered = append(filtered, tool)
	}

	return filtered
}

// buildToolDefinitions converts the chat's permitted MCP tools into OpenAI tool definitions.
func buildToolDefinitions(tools []domain.McpToolDef, vaultAccess bool) ([]openai.ToolDefinition, error) {
	permitted := filterToolsForChat(tools, vaultAccess)

	defs := make([]openai.ToolDefinition, 0, len(permitted))
	for _, tool := range permitted {
		def, err := mcpToolToOpenAIToolDef(tool)
		if err != nil {
			return nil, rerrors.Wrap(err, "error converting mcp tool to openai tool definition")
		}

		defs = append(defs, def)
	}

	return defs, nil
}
