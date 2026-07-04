package executors

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"github.com/ruf-dev/artel/internal/storage"
	"go.redsock.ru/rerrors"
)

const (
	ToolListFiles       = "list_files"
	ToolReadFile        = "read_file"
	ToolWriteFile       = "write_file"
	ToolDeleteFile      = "delete_file"
	ToolMoveFile        = "move_file"
	ToolListFolders     = "list_folders"
	ToolListTags        = "list_tags"
	ToolGetFileMetadata = "get_file_metadata"
)

type VaultExecutor struct{}

func NewVaultExecutor() *VaultExecutor {
	return &VaultExecutor{}
}

// Execute dispatches a vault tool call. client serves the markdown half of the vault (CouchDB
// livesync); bucket serves the non-markdown half (S3-compatible object storage) and may be nil
// if the vault has no linked bucket — handlers that need it then return
// user_errors.NoS3BucketLinked for non-markdown paths.
func (e *VaultExecutor) Execute(ctx context.Context, toolName string, client *couchdb.LiveSyncClient, bucket storage.BinaryStore, params map[string]interface{}) (domain.ToolExecResult, error) {
	switch toolName {
	case ToolListFiles:
		return e.listFiles(ctx, client, bucket)
	case ToolReadFile:
		return e.readFile(ctx, client, bucket, params)
	case ToolWriteFile:
		return e.writeFile(ctx, client, bucket, params)
	case ToolDeleteFile:
		return e.deleteFile(ctx, client, bucket, params)
	case ToolMoveFile:
		return e.moveFile(ctx, client, bucket, params)
	case ToolListFolders:
		return e.listFolders(ctx, client)
	case ToolListTags:
		return e.listTags(ctx, client)
	case ToolGetFileMetadata:
		return e.getFileMetadata(ctx, client, bucket, params)
	default:
		return domain.ToolExecResult{}, rerrors.New(fmt.Sprintf("unknown vault tool: %s", toolName))
	}
}

func (e *VaultExecutor) listFiles(ctx context.Context, client *couchdb.LiveSyncClient, bucket storage.BinaryStore) (domain.ToolExecResult, error) {
	notes, err := client.ListNotes(ctx)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to list notes")
	}

	var entries []map[string]interface{}
	for _, n := range notes {
		entries = append(entries, map[string]interface{}{
			"path":     n.Path,
			"mtime":    n.Mtime,
			"mimeType": couchdb.MimeTypeForPath(n.Path),
		})
	}

	if bucket != nil {
		objects, err := bucket.List(ctx)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to list bucket objects")
		}

		for _, o := range objects {
			// S3's ListObjectsV2 response doesn't carry per-object Content-Type (unlike a
			// HEAD/GetObject request, which Get/Stat use) — derive it from the extension instead
			// of trusting o.MimeType, which is unreliable here.
			entries = append(entries, map[string]interface{}{
				"path":     o.Path,
				"mtime":    o.Mtime,
				"mimeType": couchdb.MimeTypeForPath(o.Path),
			})
		}
	}

	text, err := marshalResult(entries)
	if err != nil {
		return domain.ToolExecResult{}, err
	}

	return domain.ToolExecResult{Text: text}, nil
}

func (e *VaultExecutor) readFile(ctx context.Context, client *couchdb.LiveSyncClient, bucket storage.BinaryStore, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	if storage.IsMarkdown(path) {
		note, err := client.ReadNote(ctx, path)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to read file")
		}

		return domain.ToolExecResult{Text: note.Content}, nil
	}

	if bucket == nil {
		return domain.ToolExecResult{}, user_errors.NoS3BucketLinked
	}

	obj, err := bucket.Get(ctx, path)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to read file")
	}

	return domain.ToolExecResult{Data: obj.RawBytes, MimeType: obj.MimeType, ResourcePath: path}, nil
}

func (e *VaultExecutor) writeFile(ctx context.Context, client *couchdb.LiveSyncClient, bucket storage.BinaryStore, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	content, ok := params["content"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpContentRequired
	}

	if storage.IsMarkdown(path) {
		err := client.WriteNote(ctx, path, content)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to write note")
		}

		return domain.ToolExecResult{Text: "File written successfully"}, nil
	}

	if bucket == nil {
		return domain.ToolExecResult{}, user_errors.NoS3BucketLinked
	}

	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "decode base64 content")
	}

	err = bucket.Put(ctx, path, decoded, couchdb.MimeTypeForPath(path))
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to write file")
	}

	return domain.ToolExecResult{Text: "File written successfully"}, nil
}

func (e *VaultExecutor) deleteFile(ctx context.Context, client *couchdb.LiveSyncClient, bucket storage.BinaryStore, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	if storage.IsMarkdown(path) {
		err := client.DeleteNote(ctx, path)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to delete file")
		}

		return domain.ToolExecResult{Text: "File deleted successfully"}, nil
	}

	if bucket == nil {
		return domain.ToolExecResult{}, user_errors.NoS3BucketLinked
	}

	err := bucket.Delete(ctx, path)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to delete file")
	}

	return domain.ToolExecResult{Text: "File deleted successfully"}, nil
}

func (e *VaultExecutor) moveFile(ctx context.Context, client *couchdb.LiveSyncClient, bucket storage.BinaryStore, params map[string]interface{}) (domain.ToolExecResult, error) {
	oldPath, ok := params["old_path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpOldPathRequired
	}

	newPath, ok := params["new_path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpNewPathRequired
	}

	oldIsMarkdown := storage.IsMarkdown(oldPath)
	newIsMarkdown := storage.IsMarkdown(newPath)

	if oldIsMarkdown != newIsMarkdown {
		return domain.ToolExecResult{}, user_errors.McpCrossBackendMoveNotSupported
	}

	if oldIsMarkdown {
		err := client.MoveNote(ctx, oldPath, newPath)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to move file")
		}

		return domain.ToolExecResult{Text: "File moved successfully"}, nil
	}

	if bucket == nil {
		return domain.ToolExecResult{}, user_errors.NoS3BucketLinked
	}

	err := bucket.Move(ctx, oldPath, newPath)
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

func (e *VaultExecutor) getFileMetadata(ctx context.Context, client *couchdb.LiveSyncClient, bucket storage.BinaryStore, params map[string]interface{}) (domain.ToolExecResult, error) {
	path, ok := params["path"].(string)
	if !ok {
		return domain.ToolExecResult{}, user_errors.McpPathRequired
	}

	if storage.IsMarkdown(path) {
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

	if bucket == nil {
		return domain.ToolExecResult{}, user_errors.NoS3BucketLinked
	}

	entry, err := bucket.Stat(ctx, path)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "failed to get file metadata")
	}

	// S3 has no revision/separate-ctime/tombstone concept: reuse mtime for ctime, leave rev
	// empty, and deleted is always false (a delete simply removes the key, no tombstone row).
	metadataMap := map[string]interface{}{
		"id":      path,
		"rev":     "",
		"mtime":   entry.Mtime,
		"ctime":   entry.Mtime,
		"size":    entry.Size,
		"deleted": false,
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
				Name:        ToolWriteFile,
				Description: "Create or update a file. For .md/.markdown paths, content is plain text stored in CouchDB. For any other path, content must be base64-encoded and is stored in the vault's linked S3 bucket (fails if no bucket is linked).",
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
				Description: "Move or rename a file. Both old_path and new_path must be on the same backend (both .md/.markdown or both non-markdown) — moving a file across the markdown/S3 boundary is not supported in one call.",
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
				Name:        ToolGetFileMetadata,
				Description: "Get metadata for a file. rev/ctime/deleted are CouchDB-specific and zero-valued (rev empty, ctime==mtime, deleted false) for S3-backed (non-markdown) files.",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string"},
				},
				Required: []string{"path"},
			},
			OutputSchema: domain.ToolSchema{
				Properties: map[string]domain.ToolProperty{
					"id":      {Type: "string", Description: "CouchDB document id, or the file path for S3-backed files"},
					"rev":     {Type: "string", Description: "CouchDB document revision; empty for S3-backed files"},
					"mtime":   {Type: "integer", Description: "Last modified time, unix ms"},
					"ctime":   {Type: "integer", Description: "Creation time, unix ms; equals mtime for S3-backed files"},
					"size":    {Type: "integer", Description: "Content size in bytes"},
					"deleted": {Type: "boolean", Description: "Whether the note is tombstoned; always false for S3-backed files"},
				},
				Required: []string{"id", "rev", "mtime", "ctime", "size", "deleted"},
			},
		},
	}
}
