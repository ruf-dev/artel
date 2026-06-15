package mcp_api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
)

type KeyContext struct {
	KeyUuid   string
	VaultUuid string
	UserUuid  string
	CouchURL  string
	CouchDb   string
	CouchUser string
	CouchPass string
}

type contextKey string

const keyContextKey contextKey = "mcpKeyContext"

const (
	toolListFiles       = "list_files"
	toolReadFile        = "read_file"
	toolWriteNote       = "write_note"
	toolDeleteFile      = "delete_file"
	toolMoveFile        = "move_file"
	toolListFolders     = "list_folders"
	toolListTags        = "list_tags"
	toolGetNoteMetadata = "get_note_metadata"
)

func contextWithKeyCtx(ctx context.Context, keyCtx domain.McpKeyContext) context.Context {
	kc := KeyContext{
		KeyUuid:   keyCtx.KeyUuid.String(),
		VaultUuid: keyCtx.VaultUuid.String(),
		UserUuid:  keyCtx.UserUuid.String(),
		CouchURL:  keyCtx.CouchURL,
		CouchDb:   keyCtx.CouchDb,
		CouchUser: keyCtx.CouchUser,
		CouchPass: keyCtx.CouchPass,
	}
	return context.WithValue(ctx, keyContextKey, kc)
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

func getToolDefinitions() []ToolDef {
	tools := []ToolDef{
		{
			Name:        toolListFiles,
			Description: "List all files in the vault (notes and binary files)",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        toolReadFile,
			Description: "Read any file by path. Returns text content for notes/text files, base64-encoded binary for images/PDFs.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        toolWriteNote,
			Description: "Create or update a note",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path":    map[string]string{"type": "string"},
					"content": map[string]string{"type": "string"},
				},
				"required": []string{"path", "content"},
			},
		},
		{
			Name:        toolDeleteFile,
			Description: "Delete any file by path",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        toolMoveFile,
			Description: "Move or rename any file. Not supported for large chunked binary files.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"old_path": map[string]string{"type": "string"},
					"new_path": map[string]string{"type": "string"},
				},
				"required": []string{"old_path", "new_path"},
			},
		},
		{
			Name:        toolListFolders,
			Description: "List all folders in the vault",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        toolListTags,
			Description: "List all tags in the vault",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        toolGetNoteMetadata,
			Description: "Get metadata for a note",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
	}

	return tools
}

func dispatchToolCall(ctx context.Context, toolName string, arguments map[string]interface{}, keyCtx KeyContext) (interface{}, error) {
	switch toolName {
	case toolListFiles, toolReadFile, toolWriteNote, toolDeleteFile, toolMoveFile,
		toolListFolders, toolListTags, toolGetNoteMetadata:
		client := couchdb.NewLiveSyncClient(keyCtx.CouchURL, keyCtx.CouchDb, keyCtx.CouchUser, keyCtx.CouchPass)
		return dispatchVaultTool(ctx, toolName, arguments, client)
	default:
		return nil, rerrors.New(fmt.Sprintf("unknown tool: %s", toolName))
	}
}

func isKnownTool(name string) bool {
	switch name {
	case toolListFiles, toolReadFile, toolWriteNote, toolDeleteFile, toolMoveFile,
		toolListFolders, toolListTags, toolGetNoteMetadata:
		return true
	}
	return false
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

func toJsonString(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
