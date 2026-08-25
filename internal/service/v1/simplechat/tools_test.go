package simplechat

import (
	"encoding/json"
	"testing"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/v1/mcp/executors"
	"github.com/stretchr/testify/require"
)

func TestMcpToolToOpenAIToolDef_FlatProperties(t *testing.T) {
	props := map[string]domain.ToolProperty{
		"path": {
			Type:        "string",
			Description: "Path to the note",
		},
		"mode": {
			Type: "string",
			Enum: []string{"append", "overwrite"},
		},
	}
	tool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:        "write_file",
			Description: "Write a note",
			Properties:  props,
			Required:    []string{"path"},
		},
	}

	def, err := mcpToolToOpenAIToolDef(tool)
	require.NoError(t, err)

	require.Equal(t, "write_file", def.Name)
	require.Equal(t, "Write a note", def.Description)

	var schema map[string]any
	err = json.Unmarshal(def.ParametersJSON, &schema)
	require.NoError(t, err)

	require.Equal(t, "object", schema["type"])
	require.Equal(t, []any{"path"}, schema["required"])

	properties, ok := schema["properties"].(map[string]any)
	require.True(t, ok)

	pathProp, ok := properties["path"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "string", pathProp["type"])
	require.Equal(t, "Path to the note", pathProp["description"])

	modeProp, ok := properties["mode"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, []any{"append", "overwrite"}, modeProp["enum"])
	// Description was empty, so the key must be omitted rather than emitted as "".
	_, hasDescription := modeProp["description"]
	require.False(t, hasDescription)
}

// A tool taking no parameters must still declare an empty "properties" object — some providers
// reject a function schema that omits it entirely.
func TestMcpToolToOpenAIToolDef_NoParametersStillEmitsProperties(t *testing.T) {
	tool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:       "list_files",
			Properties: map[string]domain.ToolProperty{},
			Required:   []string{},
		},
	}

	def, err := mcpToolToOpenAIToolDef(tool)
	require.NoError(t, err)

	var schema map[string]any
	err = json.Unmarshal(def.ParametersJSON, &schema)
	require.NoError(t, err)

	require.Equal(t, "object", schema["type"])

	properties, ok := schema["properties"].(map[string]any)
	require.True(t, ok)
	require.Empty(t, properties)

	// Required was empty, so it must be omitted entirely.
	_, hasRequired := schema["required"]
	require.False(t, hasRequired)
}

func TestMcpToolToOpenAIToolDef_NestedObjectAndArray(t *testing.T) {
	itemProp := domain.ToolProperty{Type: "string"}
	filterProps := map[string]domain.ToolProperty{
		"field": {Type: "string", Description: "Column"},
		"tags": {
			Type:  "array",
			Items: &itemProp,
		},
	}
	props := map[string]domain.ToolProperty{
		"filter": {
			Type:       "object",
			Properties: filterProps,
			Required:   []string{"field"},
		},
	}
	tool := domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:       "query",
			Properties: props,
			Required:   []string{"filter"},
		},
	}

	def, err := mcpToolToOpenAIToolDef(tool)
	require.NoError(t, err)

	var schema map[string]any
	err = json.Unmarshal(def.ParametersJSON, &schema)
	require.NoError(t, err)

	properties, ok := schema["properties"].(map[string]any)
	require.True(t, ok)

	filter, ok := properties["filter"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "object", filter["type"])
	require.Equal(t, []any{"field"}, filter["required"])

	nested, ok := filter["properties"].(map[string]any)
	require.True(t, ok)

	fieldProp, ok := nested["field"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "string", fieldProp["type"])

	tagsProp, ok := nested["tags"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "array", tagsProp["type"])

	items, ok := tagsProp["items"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "string", items["type"])
}

func TestFilterToolsForChat_VaultAccessKeepsEverything(t *testing.T) {
	tools := []domain.McpToolDef{
		newToolDef(executors.ToolReadFile),
		newToolDef("tract_run"),
	}

	filtered := filterToolsForChat(tools, true)

	require.Len(t, filtered, 2)
}

func TestFilterToolsForChat_NoVaultAccessDropsOnlyVaultTools(t *testing.T) {
	tools := []domain.McpToolDef{
		newToolDef(executors.ToolReadFile),
		newToolDef(executors.ToolWriteFile),
		newToolDef("tract_run"),
		newToolDef("postgres_query"),
	}

	filtered := filterToolsForChat(tools, false)

	names := make([]string, 0, len(filtered))
	for _, tool := range filtered {
		names = append(names, tool.ApiDescription.Name)
	}

	require.ElementsMatch(t, []string{"tract_run", "postgres_query"}, names)
}

// Every tool executors.VaultToolDefinitions declares must be excluded when vault access is off
// — this pins the filter to the real catalogue rather than a hardcoded copy of it.
func TestFilterToolsForChat_ExcludesEveryVaultToolDefinition(t *testing.T) {
	vaultDefs := executors.VaultToolDefinitions()
	require.NotEmpty(t, vaultDefs)

	filtered := filterToolsForChat(vaultDefs, false)

	require.Empty(t, filtered)
}

func TestBuildToolDefinitions_ConvertsFilteredSet(t *testing.T) {
	tools := []domain.McpToolDef{
		newToolDef(executors.ToolReadFile),
		newToolDef("tract_run"),
	}

	defs, err := buildToolDefinitions(tools, false)
	require.NoError(t, err)

	require.Len(t, defs, 1)
	require.Equal(t, "tract_run", defs[0].Name)
	require.NotEmpty(t, defs[0].ParametersJSON)
}

func newToolDef(name string) domain.McpToolDef {
	return domain.McpToolDef{
		ApiDescription: domain.ToolApiDescription{
			Name:       name,
			Properties: map[string]domain.ToolProperty{},
		},
	}
}
