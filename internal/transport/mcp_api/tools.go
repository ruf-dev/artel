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
	"github.com/ruf-dev/artel/internal/service/user_errors"
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

const (
	toolListFiles       = "list_files"
	toolReadFile        = "read_file"
	toolWriteNote       = "write_note"
	toolDeleteFile      = "delete_file"
	toolMoveFile        = "move_file"
	toolListFolders     = "list_folders"
	toolListTags        = "list_tags"
	toolGetNoteMetadata = "get_note_metadata"

	toolListEmailFolders  = "list_email_folders"
	toolListEmailAccounts = "list_email_accounts"
	toolListEmails        = "list_emails"
	toolReadEmail         = "read_email"
	toolSendEmail         = "send_email"
)

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

type ToolResult struct {
	Content []ContentBlock `json:"content"`
}

type ContentBlock struct {
	Type     string           `json:"type"`
	Text     string           `json:"text,omitempty"`
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
	return []ToolDef{
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
		{
			Name:        toolListEmailFolders,
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
			Name:        toolListEmailAccounts,
			Description: "List the user's configured email accounts",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
		{
			Name:        toolListEmails,
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
			Name:        toolReadEmail,
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
			Name:        toolSendEmail,
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
	}
}

func dispatchToolCall(ctx context.Context, toolName string, arguments map[string]interface{}, keyCtx KeyContext, emailSvc service.EmailService) (interface{}, error) {
	switch toolName {
	case toolListFiles, toolReadFile, toolWriteNote, toolDeleteFile, toolMoveFile,
		toolListFolders, toolListTags, toolGetNoteMetadata:
		client := couchdb.NewLiveSyncClient(keyCtx.CouchURL, keyCtx.CouchDb, keyCtx.CouchUser, keyCtx.CouchPass)
		return dispatchVaultTool(ctx, toolName, arguments, client)
	case toolListEmailFolders:
		return handleListEmailFolders(ctx, arguments, emailSvc)
	case toolListEmailAccounts:
		return handleListEmailAccounts(ctx, keyCtx, emailSvc)
	case toolListEmails:
		return handleListEmails(ctx, arguments, emailSvc)
	case toolReadEmail:
		return handleReadEmail(ctx, arguments, emailSvc)
	case toolSendEmail:
		return handleSendEmail(ctx, arguments, emailSvc)
	default:
		return nil, rerrors.New(fmt.Sprintf("unknown tool: %s", toolName))
	}
}

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

func handleListEmailFolders(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (ToolResult, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpAccountIdRequired
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "invalid account_id")
	}
	folders, err := emailSvc.ListFolders(ctx, accountUuid)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "list email folders")
	}
	return textResult(toJsonString(folders)), nil
}

func handleListEmailAccounts(ctx context.Context, keyCtx KeyContext, emailSvc service.EmailService) (ToolResult, error) {
	userUuid, err := uuid.Parse(keyCtx.UserUuid)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "invalid user uuid")
	}

	accounts, err := emailSvc.ListAccounts(ctx, userUuid)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "list email accounts")
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

	return textResult(toJsonString(result)), nil
}

func handleListEmails(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (ToolResult, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpAccountIdRequired
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "invalid account_id")
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
		return ToolResult{}, rerrors.Wrap(err, "list emails")
	}

	return textResult(toJsonString(emails)), nil
}

func handleReadEmail(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (ToolResult, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpAccountIdRequired
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "invalid account_id")
	}

	id, ok := arguments["id"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpIdRequired
	}

	msg, err := emailSvc.ReadEmail(ctx, accountUuid, id)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "read email")
	}

	return textResult(toJsonString(msg)), nil
}

func handleSendEmail(ctx context.Context, arguments map[string]interface{}, emailSvc service.EmailService) (ToolResult, error) {
	accountIdStr, ok := arguments["account_id"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpAccountIdRequired
	}
	accountUuid, err := uuid.Parse(accountIdStr)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "invalid account_id")
	}

	to, ok := arguments["to"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpToRequired
	}
	subject, ok := arguments["subject"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpSubjectRequired
	}
	body, ok := arguments["body"].(string)
	if !ok {
		return ToolResult{}, user_errors.McpBodyRequired
	}

	err = emailSvc.SendEmail(ctx, accountUuid, to, subject, body)
	if err != nil {
		return ToolResult{}, rerrors.Wrap(err, "send email")
	}

	return textResult("Email sent successfully"), nil
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
