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
	KeyUuid   string
	VaultUuid string
	UserUuid  string
	CouchURL  string
	CouchDb   string
	CouchUser string
	CouchPass string
	HasEmails bool
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
	kc := KeyContext{
		KeyUuid:   keyCtx.KeyUuid.String(),
		VaultUuid: keyCtx.VaultUuid.String(),
		UserUuid:  keyCtx.UserUuid.String(),
		CouchURL:  keyCtx.CouchURL,
		CouchDb:   keyCtx.CouchDb,
		CouchUser: keyCtx.CouchUser,
		CouchPass: keyCtx.CouchPass,
		HasEmails: keyCtx.HasEmails,
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

func getToolDefinitions(hasEmails bool) []ToolDef {
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

	if hasEmails {
		tools = append(tools,
			ToolDef{
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
			ToolDef{
				Name:        toolListEmailAccounts,
				Description: "List the user's configured email accounts",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
			},
			ToolDef{
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
			ToolDef{
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
			ToolDef{
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
		)
	}

	return tools
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

func isKnownTool(name string) bool {
	switch name {
	case toolListFiles, toolReadFile, toolWriteNote, toolDeleteFile, toolMoveFile,
		toolListFolders, toolListTags, toolGetNoteMetadata,
		toolListEmailFolders, toolListEmailAccounts, toolListEmails, toolReadEmail, toolSendEmail:
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
