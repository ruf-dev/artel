package mcp_api

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/ruf-dev/artel/internal/domain"
)

type contextKey string

const keyContextKey contextKey = "mcpKeyContext"

func contextWithKeyCtx(ctx context.Context, keyCtx domain.McpKeyContext) context.Context {
	return context.WithValue(ctx, keyContextKey, keyCtx)
}

type ToolDef struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	InputSchema  map[string]interface{} `json:"inputSchema"`
	OutputSchema map[string]interface{} `json:"outputSchema,omitempty"`
}

type ToolResult struct {
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Type     string           `json:"type"`
	Text     string           `json:"text"`
	Data     string           `json:"data,omitempty"`
	MimeType string           `json:"mimeType,omitempty"`
	Resource *ResourceContent `json:"resource,omitempty"`
}

type ResourceContent struct {
	Uri      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Blob     string `json:"blob"`
}

func textResult(text string) ToolResult {
	block := ContentBlock{Type: "text", Text: text}
	return ToolResult{Content: []ContentBlock{block}}
}

func toolResultFromExec(res domain.ToolExecResult) ToolResult {
	if res.Data != nil {
		return ToolResult{Content: []ContentBlock{buildContentBlock(res.Data, res.MimeType, res.ResourcePath)}}
	}
	return textResult(res.Text)
}

func momToolToToolDef(t domain.McpToolDef) ToolDef {
	inputSchema := toolSchemaToWire(t.ApiDescription.Properties, t.ApiDescription.Required)

	def := ToolDef{
		Name:        t.ApiDescription.Name,
		Description: t.ApiDescription.Description,
		InputSchema: inputSchema,
	}

	// {} (no properties, no required) means the output shape is undeclared — omit the key
	// entirely rather than sending an empty object.
	if len(t.OutputSchema.Properties) > 0 || len(t.OutputSchema.Required) > 0 {
		def.OutputSchema = toolSchemaToWire(t.OutputSchema.Properties, t.OutputSchema.Required)
	}

	return def
}

// toolSchemaToWire renders a {properties, required} pair into MCP's JSON-Schema wire shape,
// recursing into nested object/array properties.
func toolSchemaToWire(properties map[string]domain.ToolProperty, required []string) map[string]interface{} {
	props := make(map[string]interface{}, len(properties))
	for name, p := range properties {
		props[name] = toolPropertyToWire(p)
	}
	schema := map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
	return schema
}

func toolPropertyToWire(p domain.ToolProperty) map[string]interface{} {
	wire := map[string]interface{}{"type": p.Type, "description": p.Description}

	if len(p.Enum) > 0 {
		wire["enum"] = p.Enum
	}

	if len(p.Properties) > 0 {
		nested := make(map[string]interface{}, len(p.Properties))
		for name, sub := range p.Properties {
			nested[name] = toolPropertyToWire(sub)
		}
		wire["properties"] = nested
		wire["required"] = p.Required
	}

	if p.Items != nil {
		wire["items"] = toolPropertyToWire(*p.Items)
	}

	return wire
}

func buildContentBlock(data []byte, mimeType, path string) ContentBlock {
	b64 := base64.StdEncoding.EncodeToString(data)

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return ContentBlock{Type: "image", Data: b64, MimeType: mimeType}
	case mimeType == "application/pdf":
		return ContentBlock{Type: "document", Data: b64, MimeType: mimeType}
	default:
		res := &ResourceContent{Uri: "file://" + path, MimeType: mimeType, Blob: b64}
		return ContentBlock{Type: "resource", Resource: res}
	}
}
