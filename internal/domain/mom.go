package domain

import "time"

type McpDefinition struct {
	Name        string
	Author      string
	Description string
	Tools       []McpToolDef
	CreatedAt   time.Time
}

type McpToolDef struct {
	ApiDescription ToolApiDescription
	Action         ToolAction
}

// ToolApiDescription is the LLM-visible metadata for a tool.
type ToolApiDescription struct {
	Name        string
	Description string
	Properties  map[string]ToolProperty
	Required    []string
}

// ToolProperty describes one input parameter.
type ToolProperty struct {
	Type        string // "string" | "integer" | "number" | "boolean" | "array" | "object"
	Description string
}

// ToolAction is a discriminated union — exactly one field must be non-nil.
type ToolAction struct {
	Imap *ImapAction
	Smtp *SmtpAction
	Http *HttpAction
}

type ImapOperation string

const (
	IMAP_OP_LIST_FOLDERS  ImapOperation = "list_folders"
	IMAP_OP_LIST_MESSAGES ImapOperation = "list_messages"
	IMAP_OP_FETCH_MESSAGE ImapOperation = "fetch_message"
)

type ImapAction struct {
	Operation ImapOperation
}

type SmtpOperation string

const (
	SMTP_OP_SEND SmtpOperation = "send"
)

type SmtpAction struct {
	Operation SmtpOperation
}

// HttpAction drives the HTTP executor with ${{params.*}} / __secrets.* interpolation.
type HttpAction struct {
	Method  string            // GET | POST | PUT | PATCH | DELETE
	Url     string            // may contain ${{params.name}} placeholders
	Headers map[string]string // values may be static or __secrets.*
	Query   map[string]string // values may be static, ${{params.*}}, or __secrets.*
	// TODO: replace with a typed enum once external_connections.provider is constrained in the DB schema.
	Credentials string // external_connections.provider to resolve __secrets.*
}
