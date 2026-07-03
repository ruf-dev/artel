package executors

import (
	"context"
	"fmt"
	"strings"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

const (
	ToolListFiles       = "list_files"
	ToolReadFile        = "read_file"
	ToolWriteNote       = "write_note"
	ToolDeleteFile      = "delete_file"
	ToolMoveFile        = "move_file"
	ToolListFolders     = "list_folders"
	ToolListTags        = "list_tags"
	ToolGetNoteMetadata = "get_note_metadata"
)

type VaultExecutor struct{}

func NewVaultExecutor() *VaultExecutor {
	return &VaultExecutor{}
}

func (e *VaultExecutor) Execute(ctx context.Context, toolName string, client *couchdb.LiveSyncClient, params map[string]interface{}) (domain.ToolExecResult, error) {
	switch toolName {
	case ToolListFiles:
		return e.listFiles(ctx, client)
	case ToolReadFile:
		return e.readFile(ctx, client, params)
	case ToolWriteNote:
		return e.writeNote(ctx, client, params)
	case ToolDeleteFile:
		return e.deleteFile(ctx, client, params)
	case ToolMoveFile:
		return e.moveFile(ctx, client, params)
	case ToolListFolders:
		return e.listFolders(ctx, client)
	case ToolListTags:
		return e.listTags(ctx, client)
	case ToolGetNoteMetadata:
		return e.getNoteMetadata(ctx, client, params)
	default:
		return domain.ToolExecResult{}, rerrors.New(fmt.Sprintf("unknown vault tool: %s", toolName))
	}
}

func (e *VaultExecutor) listFiles(ctx context.Context, client *couchdb.LiveSyncClient) (domain.ToolExecResult, error) {
	notes, err := client.ListNotes(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to list notes")
	}
	binaries, err := client.ListFiles(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to list files")
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

	text, err := marshalResult(entries)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

func (e *VaultExecutor) readFile(ctx context.Context, client *couchdb.LiveSyncClient, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	if strings.HasPrefix(couchdb.MimeTypeForPath(path), "text/") {
		note, err := client.ReadNote(ctx, path)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to read file")
		}
		return domain.ToolExecResult{Text: note.Content}, nil
	}

	file, err := client.ReadFile(ctx, path)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to read file")
	}
	return domain.ToolExecResult{Data: file.RawBytes, MimeType: file.MimeType, ResourcePath: file.Id}, nil
}

func (e *VaultExecutor) writeNote(ctx context.Context, client *couchdb.LiveSyncClient, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	content, ok := params["content"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpContentRequired
	}

	err := client.WriteNote(ctx, path, content)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to write note")
	}

	return domain.ToolExecResult{Text: "Note written successfully"}, nil
}

func (e *VaultExecutor) deleteFile(ctx context.Context, client *couchdb.LiveSyncClient, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	var err error
	if strings.HasPrefix(couchdb.MimeTypeForPath(path), "text/") {
		err = client.DeleteNote(ctx, path)
	} else {
		err = client.DeleteFile(ctx, path)
	}
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to delete file")
	}

	return domain.ToolExecResult{Text: "File deleted successfully"}, nil
}

func (e *VaultExecutor) moveFile(ctx context.Context, client *couchdb.LiveSyncClient, params map[string]interface{}) (domain.ToolExecResult, error) {
	oldPath, ok := params["old_path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpOldPathRequired
	}

	newPath, ok := params["new_path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpNewPathRequired
	}

	var err error
	if strings.HasPrefix(couchdb.MimeTypeForPath(oldPath), "text/") {
		err = client.MoveNote(ctx, oldPath, newPath)
	} else {
		err = client.MoveFile(ctx, oldPath, newPath)
	}
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to move file")
	}

	return domain.ToolExecResult{Text: "File moved successfully"}, nil
}

func (e *VaultExecutor) listFolders(ctx context.Context, client *couchdb.LiveSyncClient) (domain.ToolExecResult, error) {
	folders, err := client.ListFolders(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to list folders")
	}

	text, err := marshalResult(folders)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

func (e *VaultExecutor) listTags(ctx context.Context, client *couchdb.LiveSyncClient) (domain.ToolExecResult, error) {
	tags, err := client.ListTags(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to list tags")
	}

	text, err := marshalResult(tags)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

func (e *VaultExecutor) getNoteMetadata(ctx context.Context, client *couchdb.LiveSyncClient, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	metadata, err := client.GetNoteMetadata(ctx, path)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to get note metadata")
	}

	metadataMap := map[string]interface{}{
		"id":      metadata.Id,
		"rev":     metadata.Rev,
		"mtime":   metadata.Mtime,
		"ctime":   metadata.Ctime,
		"size":    metadata.Size,
		"deleted": metadata.Deleted,
	}

	text, err := marshalResult(metadataMap)
	if err != nil {
		return domain.ToolExecResult{}, err
	}
	return domain.ToolExecResult{Text: text}, nil
}

// Output schemas below are hints for tool-picker/port-chip UIs, not enforced at runtime.
// list_files/list_folders/list_tags/connections marshal a JSON array as their text result;
// since ToolSchema (unlike ToolProperty) has no array/root type of its own, Properties there
// describes the shape of each array element rather than the response envelope itself.

func VaultToolDefinitions() []domain.McpToolDef {
	return []domain.McpToolDef{
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolListFiles,
				Description: "List all files in the vault (notes and binary files)",
				Properties:  map[string]domain.ToolProperty{},
				Required:    []string{},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"path":     {Type: "string", Description: "Vault-relative file path (each array entry)"},
					"mtime":    {Type: "integer", Description: "Last modified time, unix ms"},
					"mimeType": {Type: "string", Description: "MIME type of the file"},
				},
				Required: []string{"path", "mtime", "mimeType"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolReadFile,
				Description: "Read any file by path. Returns text content for notes/text files, base64-encoded binary for images/PDFs.",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string"},
				},
				Required: []string{"path"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"content": {Type: "string", Description: "File content as text; binary files are returned as a resource/image content block instead"},
				},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolWriteNote,
				Description: "Create or update a note",
				Properties: map[string]domain.ToolProperty{
					"path":    {Type: "string"},
					"content": {Type: "string"},
				},
				Required: []string{"path", "content"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"message": {Type: "string", Description: "Confirmation message"},
				},
				Required: []string{"message"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolDeleteFile,
				Description: "Delete any file by path",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string"},
				},
				Required: []string{"path"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"message": {Type: "string", Description: "Confirmation message"},
				},
				Required: []string{"message"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolMoveFile,
				Description: "Move or rename any file. Not supported for large chunked binary files.",
				Properties: map[string]domain.ToolProperty{
					"old_path": {Type: "string"},
					"new_path": {Type: "string"},
				},
				Required: []string{"old_path", "new_path"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"message": {Type: "string", Description: "Confirmation message"},
				},
				Required: []string{"message"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolListFolders,
				Description: "List all folders in the vault",
				Properties:  map[string]domain.ToolProperty{},
				Required:    []string{},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string", Description: "Folder path (each array entry)"},
				},
				Required: []string{"path"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolListTags,
				Description: "List all tags in the vault",
				Properties:  map[string]domain.ToolProperty{},
				Required:    []string{},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"tag": {Type: "string", Description: "Tag name (each array entry)"},
				},
				Required: []string{"tag"},
			},
		},
		{
			ApiDescription: domain.ToolApiDescription{
				Name:        ToolGetNoteMetadata,
				Description: "Get metadata for a note",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string"},
				},
				Required: []string{"path"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"id":      {Type: "string", Description: "CouchDB document id"},
					"rev":     {Type: "string", Description: "CouchDB document revision"},
					"mtime":   {Type: "integer", Description: "Last modified time, unix ms"},
					"ctime":   {Type: "integer", Description: "Creation time, unix ms"},
					"size":    {Type: "integer", Description: "Content size in bytes"},
					"deleted": {Type: "boolean", Description: "Whether the note is tombstoned"},
				},
				Required: []string{"id", "rev", "mtime", "ctime", "size", "deleted"},
			},
		},
	}
}
