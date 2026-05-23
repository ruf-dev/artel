package mcp_api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service"
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
		{
			Name:        "list_email_folders",
			Description: "List IMAP folders available in an email account",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"account_id": map[string]string{"type": "string", "description": "UUID of the email account"},
				},
				"required": []string{"account_id"},
			},
		},
		{
			Name:        "list_email_accounts",
			Description: "List the user's configured email accounts",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        "list_emails",
			Description: "List recent emails from an account's inbox",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"account_id": map[string]string{"type": "string", "description": "UUID of the email account"},
					"limit":      map[string]interface{}{"type": "integer", "description": "Max number of emails to return (default 20)"},
				},
				"required": []string{"account_id"},
			},
		},
		{
			Name:        "read_email",
			Description: "Read the full content of an email by its ID",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"account_id": map[string]string{"type": "string", "description": "UUID of the email account"},
					"id":         map[string]string{"type": "string", "description": "Email UID from list_emails"},
				},
				"required": []string{"account_id", "id"},
			},
		},
		{
			Name:        "send_email",
			Description: "Send an email from one of the user's configured accounts",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"account_id": map[string]string{"type": "string", "description": "UUID of the email account"},
					"to":         map[string]string{"type": "string", "description": "Recipient email address"},
					"subject":    map[string]string{"type": "string", "description": "Email subject"},
					"body":       map[string]string{"type": "string", "description": "Email body (plain text)"},
				},
				"required": []string{"account_id", "to", "subject", "body"},
			},
		},
		{
			Name:        "list_files",
			Description: "List all binary files (images, PDFs, etc.) in the vault. Does not include text or markdown notes — use list_notes for those.",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        "read_file",
			Description: "Read a binary file (image, PDF, etc.) by path. Returns base64-encoded content with MIME type. Use read_note for text/markdown files.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "delete_file",
			Description: "Delete a binary file by path. Use delete_note for text/markdown files.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]string{"type": "string"},
				},
				"required": []string{"path"},
			},
		},
		{
			Name:        "move_file",
			Description: "Move or rename a binary file. Not supported for large chunked files — use Obsidian directly for those. Use move_note for text/markdown files.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"old_path": map[string]string{"type": "string"},
					"new_path": map[string]string{"type": "string"},
				},
				"required": []string{"old_path", "new_path"},
			},
		},
	}
}

func dispatchToolCall(ctx context.Context, toolName string, arguments map[string]interface{}, keyCtx KeyContext, emailSvc service.EmailService) (interface{}, error) {
	switch toolName {
	case "list_notes", "read_note", "write_note", "delete_note", "move_note", "list_folders", "list_tags", "get_note_metadata",
		"list_files", "read_file", "delete_file", "move_file":
		client := couchdb.NewLiveSyncClient(keyCtx.CouchURL, keyCtx.CouchDb, keyCtx.CouchUser, keyCtx.CouchPass)
		return dispatchVaultTool(ctx, toolName, arguments, client)
	case "list_email_folders":
		return handleListEmailFolders(ctx, arguments, emailSvc)
	case "list_email_accounts":
		return handleListEmailAccounts(ctx, keyCtx, emailSvc)
	case "list_emails":
		return handleListEmails(ctx, arguments, emailSvc)
	case "read_email":
		return handleReadEmail(ctx, arguments, emailSvc)
	case "send_email":
		return handleSendEmail(ctx, arguments, emailSvc)
	default:
		return nil, rerrors.New(fmt.Sprintf("unknown tool: %s", toolName))
	}
}

func dispatchVaultTool(ctx context.Context, toolName string, arguments map[string]interface{}, client *couchdb.LiveSyncClient) (interface{}, error) {
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
	case "list_files":
		return handleListFiles(ctx, client)
	case "read_file":
		return handleReadFile(ctx, client, arguments)
	case "delete_file":
		return handleDeleteFile(ctx, client, arguments)
	case "move_file":
		return handleMoveFile(ctx, client, arguments)
	default:
		return nil, rerrors.New(fmt.Sprintf("unknown vault tool: %s", toolName))
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

func handleListEmailFolders(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (interface{}, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return nil, rerrors.New("account_id is required and must be a string")
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return nil, rerrors.Wrap(err, "invalid account_id")
	}
	folders, err := emailSvc.ListFolders(ctx, accountUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list email folders")
	}
	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": toJsonString(folders)},
		},
	}, nil
}

func handleListEmailAccounts(ctx context.Context, keyCtx KeyContext, emailSvc service.EmailService) (interface{}, error) {
	userUuid, err := uuid.Parse(keyCtx.UserUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "invalid user uuid")
	}

	accounts, err := emailSvc.ListAccounts(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "list email accounts")
	}

	var result []map[string]interface{}
	for _, a := range accounts {
		result = append(result, map[string]interface{}{
			"id":         a.Uuid.String(),
			"email":      a.Email,
			"imap_host":  a.ImapHost,
			"smtp_host":  a.SmtpHost,
			"created_at": a.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		})
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(result),
			},
		},
	}, nil
}

func handleListEmails(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (interface{}, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return nil, rerrors.New("account_id is required and must be a string")
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return nil, rerrors.Wrap(err, "invalid account_id")
	}

	limit := 20
	if l, ok := arguments["limit"]; ok {
		switch v := l.(type) {
		case float64:
			limit = int(v)
		case int:
			limit = v
		}
	}

	emails, err := emailSvc.ListEmails(ctx, accountUuid, limit)
	if err != nil {
		return nil, rerrors.Wrap(err, "list emails")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(emails),
			},
		},
	}, nil
}

func handleReadEmail(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (interface{}, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return nil, rerrors.New("account_id is required and must be a string")
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return nil, rerrors.Wrap(err, "invalid account_id")
	}

	id, ok := arguments["id"].(string)
	if !ok {
		return nil, rerrors.New("id is required and must be a string")
	}

	msg, err := emailSvc.ReadEmail(ctx, accountUuid, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "read email")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(msg),
			},
		},
	}, nil
}

func handleSendEmail(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (interface{}, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return nil, rerrors.New("account_id is required and must be a string")
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return nil, rerrors.Wrap(err, "invalid account_id")
	}

	to, ok := arguments["to"].(string)
	if !ok {
		return nil, rerrors.New("to is required and must be a string")
	}
	subject, ok := arguments["subject"].(string)
	if !ok {
		return nil, rerrors.New("subject is required and must be a string")
	}
	body, ok := arguments["body"].(string)
	if !ok {
		return nil, rerrors.New("body is required and must be a string")
	}

	if err := emailSvc.SendEmail(ctx, accountUuid, to, subject, body); err != nil {
		return nil, rerrors.Wrap(err, "send email")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Email sent successfully",
			},
		},
	}, nil
}

func handleListFiles(ctx context.Context, client *couchdb.LiveSyncClient) (interface{}, error) {
	files, err := client.ListFiles(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to list files")
	}

	var filesInfo []map[string]interface{}
	for _, f := range files {
		filesInfo = append(filesInfo, map[string]interface{}{
			"path":     f.Path,
			"mtime":    f.Mtime,
			"mimeType": f.MimeType,
		})
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJsonString(filesInfo),
			},
		},
	}, nil
}

func handleReadFile(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return nil, rerrors.New("path is required and must be a string")
	}

	file, err := client.ReadFile(ctx, path)
	if err != nil {
		return nil, rerrors.Wrap(err, "failed to read file")
	}

	contentBlock := buildContentBlock(file.RawBytes, file.MimeType, file.Id)

	return map[string]interface{}{
		"content": []interface{}{contentBlock},
	}, nil
}

func handleDeleteFile(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	path, ok := arguments["path"].(string)
	if !ok {
		return nil, rerrors.New("path is required and must be a string")
	}

	if err := client.DeleteFile(ctx, path); err != nil {
		return nil, rerrors.Wrap(err, "failed to delete file")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": "File deleted successfully"},
		},
	}, nil
}

func handleMoveFile(ctx context.Context, client *couchdb.LiveSyncClient, arguments map[string]interface{}) (interface{}, error) {
	oldPath, ok := arguments["old_path"].(string)
	if !ok {
		return nil, rerrors.New("old_path is required and must be a string")
	}
	newPath, ok := arguments["new_path"].(string)
	if !ok {
		return nil, rerrors.New("new_path is required and must be a string")
	}

	if err := client.MoveFile(ctx, oldPath, newPath); err != nil {
		return nil, rerrors.Wrap(err, "failed to move file")
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": "File moved successfully"},
		},
	}, nil
}

func buildContentBlock(data []byte, mimeType, path string) map[string]interface{} {
	b64 := base64.StdEncoding.EncodeToString(data)

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return map[string]interface{}{
			"type":     "image",
			"data":     b64,
			"mimeType": mimeType,
		}
	case mimeType == "application/pdf":
		return map[string]interface{}{
			"type":     "document",
			"data":     b64,
			"mimeType": mimeType,
		}
	default:
		return map[string]interface{}{
			"type": "resource",
			"resource": map[string]interface{}{
				"uri":      "file://" + path,
				"mimeType": mimeType,
				"blob":     b64,
			},
		}
	}
}

func toJsonString(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
