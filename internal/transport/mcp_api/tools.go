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
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
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
	props := make(map[string]interface{}, len(t.ApiDescription.Properties))
	for name, p := range t.ApiDescription.Properties {
		props[name] = map[string]string{"type": p.Type, "description": p.Description}
	}
	schema := map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   t.ApiDescription.Required,
	}
	return ToolDef{
		Name:        t.ApiDescription.Name,
		Description: t.ApiDescription.Description,
		InputSchema: schema,
	}
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
