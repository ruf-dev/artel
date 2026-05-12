package mcp_api

import (
	"context"
	"encoding/json"
	"fmt"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
)

type KeyContext struct {
	VaultUuid string
	UserUuid  string
	CouchURL  string
	CouchDb   string
	CouchUser string
	CouchPass string
}

type contextKey string

const keyContextKey contextKey = "mcpKeyContext"

func contextWithKeyCtx(ctx context.Context, keyCtx domain.McpKeyContext) context.Context {
	return context.WithValue(ctx, keyContextKey, KeyContext{
		VaultUuid: keyCtx.VaultUuid.String(),
		UserUuid:  keyCtx.UserUuid.String(),
		CouchURL:  keyCtx.CouchURL,
		CouchDb:   keyCtx.CouchDb,
		CouchUser: keyCtx.CouchUser,
		CouchPass: keyCtx.CouchPass,
	})
}

type ToolDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

func getToolDefinitions() []ToolDef {
	return []ToolDef{
		{
			Name:        "list_notes",
			Description: "List all notes in the vault",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        "read_note",
			Description: "Read a note by path",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "write_note",
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
			Name:        "delete_note",
			Description: "Delete a note by path",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "move_note",
			Description: "Move a note to a new path",
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
			Name:        "list_folders",
			Description: "List all folders in the vault",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        "list_tags",
			Description: "List all tags in the vault",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        "get_note_metadata",
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
}

func dispatchToolCall(ctx context.Context, toolName string, arguments map[string]interface{}, keyCtx KeyContext) (interface{}, error) {
	client := couchdb.NewLiveSyncClient(keyCtx.CouchURL, keyCtx.CouchDb, keyCtx.CouchUser, keyCtx.CouchPass)

	switch toolName {
	case "list_notes":
		return handleListNotes(ctx, client)
	case "read_note":
		return handleReadNote(ctx, client, arguments)
	case "write_note":
		return handleWriteNote(ctx, client, arguments)
	case "delete_note":
		return handleDeleteNote(ctx, client, arguments)
	case "move_note":
		return handleMoveNote(ctx, client, arguments)
	case "list_folders":
		return handleListFolders(ctx, client)
	case "list_tags":
		return handleListTags(ctx, client)
	case "get_note_metadata":
		return handleGetNoteMetadata(ctx, client, arguments)
	default:
		return nil, rerrors.New(fmt.Sprintf("unknown tool: %s", toolName))
	}
}

func handleListNotes(ctx context.Context, client *couchdb.LiveSyncClient) (interface{}, error) {
	notes, err := client.ListNotes(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to list notes")
	}

	var notesInfo []map[string]interface{}
	for _, note := range notes {
		notesInfo = append(notesInfo, map[string]interface{}{
			"path":  note.Path,
			"mtime": note.Mtime,
		})
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(notesInfo),
			},
		},
	}, nil
}

func handleReadNote(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return nil, rerrors.New("path is required and must be a string")
	}

	note, err := client.ReadNote(ctx, path)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to read note")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": note.Content,
			},
		},
	}, nil
}

func handleWriteNote(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return nil, rerrors.New("path is required and must be a string")
	}

	content, ok := arguments["content"].(string)
	if !ok {
		return nil, rerrors.New("content is required and must be a string")
	}

	err := client.WriteNote(ctx, path, content)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to write note")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Note written successfully",
			},
		},
	}, nil
}

func handleDeleteNote(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return nil, rerrors.New("path is required and must be a string")
	}

	err := client.DeleteNote(ctx, path)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to delete note")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Note deleted successfully",
			},
		},
	}, nil
}

func handleMoveNote(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	oldPath, ok := arguments["old_path"].(string)
	if !ok {
		return nil, rerrors.New("old_path is required and must be a string")
	}

	newPath, ok := arguments["new_path"].(string)
	if !ok {
		return nil, rerrors.New("new_path is required and must be a string")
	}

	err := client.MoveNote(ctx, oldPath, newPath)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to move note")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Note moved successfully",
			},
		},
	}, nil
}

func handleListFolders(ctx context.Context, client *couchdb.LiveSyncClient) (interface{}, error) {
	folders, err := client.ListFolders(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to list folders")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(folders),
			},
		},
	}, nil
}

func handleListTags(ctx context.Context, client *couchdb.LiveSyncClient) (interface{}, error) {
	tags, err := client.ListTags(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to list tags")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(tags),
			},
		},
	}, nil
}

func handleGetNoteMetadata(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return nil, rerrors.New("path is required and must be a string")
	}

	metadata, err := client.GetNoteMetadata(ctx, path)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to get note metadata")
	}

	metadataMap := map[string]interface{}{
		"id":      metadata.Id,
		"rev":     metadata.Rev,
		"mtime":   metadata.Mtime,
		"ctime":   metadata.Ctime,
		"size":    metadata.Size,
		"deleted": metadata.Deleted,
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(metadataMap),
			},
		},
	}, nil
}

func toJsonString(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
