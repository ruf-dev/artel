package mcp_api

import (
	"context"
	"fmt"
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func dispatchVaultTool(ctx context.Context, toolName string, arguments map[string]interface{}, client *couchdb.LiveSyncClient) (interface{}, error) {
	switch toolName {
	case toolListFiles:
		return handleListFiles(ctx, client)
	case toolReadFile:
		return handleReadFile(ctx, client, arguments)
	case toolWriteNote:
		return handleWriteNote(ctx, client, arguments)
	case toolDeleteFile:
		return handleDeleteFile(ctx, client, arguments)
	case toolMoveFile:
		return handleMoveFile(ctx, client, arguments)
	case toolListFolders:
		return handleListFolders(ctx, client)
	case toolListTags:
		return handleListTags(ctx, client)
	case toolGetNoteMetadata:
		return handleGetNoteMetadata(ctx, client, arguments)
	default:
		return nil, rerrors.New(fmt.Sprintf("unknown vault tool: %s", toolName))
	}
}

func handleListFiles(ctx context.Context, client *couchdb.LiveSyncClient) (ToolResult, error) {
	notes, err := client.ListNotes(ctx)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to list notes")
	}
	binaries, err := client.ListFiles(ctx)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to list files")
	}

	var entries []map[string]interface{}
	for _, n := range notes {
		entries = append(entries, map[string]interface{}{
			"path":     n.Path,
			"mtime":    n.Mtime,
			"mimeType": couchdb.MimeTypeForPath(n.Path),
		})
	}
	for _, f := range binaries {
		entries = append(entries, map[string]interface{}{
			"path":     f.Path,
			"mtime":    f.Mtime,
			"mimeType": f.MimeType,
		})
	}

	return textResult(toJsonString(entries)), nil
}

func handleReadFile(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (ToolResult, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpPathRequired
	}

	if strings.HasPrefix(couchdb.MimeTypeForPath(path), "text/") {
		note, err := client.ReadNote(ctx, path)
		if err != nil {
			return ToolResult{}, rerrors.Wrap(err, "failed to read file")
		}
		return textResult(note.Content), nil
	}

	file, err := client.ReadFile(ctx, path)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to read file")
	}
	block := buildContentBlock(file.RawBytes, file.MimeType, file.Id)
	return ToolResult{Content: []ContentBlock{block}}, nil
}

func handleWriteNote(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (ToolResult, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpPathRequired
	}

	content, ok := arguments["content"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpContentRequired
	}

	err := client.WriteNote(ctx, path, content)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to write note")
	}

	return textResult("Note written successfully"), nil
}

func handleDeleteFile(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (ToolResult, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpPathRequired
	}

	var err error
	if strings.HasPrefix(couchdb.MimeTypeForPath(path), "text/") {
		err = client.DeleteNote(ctx, path)
	} else {
		err = client.DeleteFile(ctx, path)
	}
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to delete file")
	}

	return textResult("File deleted successfully"), nil
}

func handleMoveFile(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (ToolResult, error) {
	oldPath, ok := arguments["old_path"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpOldPathRequired
	}

	newPath, ok := arguments["new_path"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpNewPathRequired
	}

	var err error
	if strings.HasPrefix(couchdb.MimeTypeForPath(oldPath), "text/") {
		err = client.MoveNote(ctx, oldPath, newPath)
	} else {
		err = client.MoveFile(ctx, oldPath, newPath)
	}
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to move file")
	}

	return textResult("File moved successfully"), nil
}

func handleListFolders(ctx context.Context, client *couchdb.LiveSyncClient) (ToolResult, error) {
	folders, err := client.ListFolders(ctx)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to list folders")
	}

	return textResult(toJsonString(folders)), nil
}

func handleListTags(ctx context.Context, client *couchdb.LiveSyncClient) (ToolResult, error) {
	tags, err := client.ListTags(ctx)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to list tags")
	}

	return textResult(toJsonString(tags)), nil
}

func handleGetNoteMetadata(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (ToolResult, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpPathRequired
	}

	metadata, err := client.GetNoteMetadata(ctx, path)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "failed to get note metadata")
	}

	metadataMap := map[string]interface{}{
		"id":      metadata.Id,
		"rev":     metadata.Rev,
		"mtime":   metadata.Mtime,
		"ctime":   metadata.Ctime,
		"size":    metadata.Size,
		"deleted": metadata.Deleted,
	}

	return textResult(toJsonString(metadataMap)), nil
}
