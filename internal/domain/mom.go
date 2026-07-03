package domain

import (
	"encoding/json"
	"time"
)

type McpDefinition struct {
	Name        string
	Author      string
	Description string
	Tools       []McpToolDef
	CreatedAt   time.Time
}

// MomCandidate is an McpDefinition paired with the caller's external connections
// that satisfy the providers its tools require.
type MomCandidate struct {
	Name        string
	Author      string
	Description string
	Connections []ExternalConnectionMeta
	Tools       []McpToolDef
}

type McpToolDef struct {
	ApiDescription ToolApiDescription
	OutputSchema   ToolSchema
	Action         ToolAction
}

// McpToolRef pairs a tool definition with the name of the MoM it belongs to — used by
// cross-MoM catalogs (e.g. McpDefinitionsRepo.ListAllTools) where the owning MoM isn't
// otherwise recoverable from McpToolDef alone.
type McpToolRef struct {
	McpName string
	Tool    McpToolDef
}

// ToolApiDescription is the LLM-visible metadata for a tool.
type ToolApiDescription struct {
	Name        string
	Description string
	Properties  map[string]ToolProperty
	Required    []string
}

// ToolSchema is a standalone {properties, required} object schema — the same shape as
// ToolApiDescription's input fields, reused to describe a tool's output. Empty (no properties,
// no required) means the output shape is undeclared.
type ToolSchema struct {
	Properties map[string]ToolProperty
	Required   []string
}

// ToolProperty describes one input or output parameter. Recursive for object/array types:
// object properties nest via Properties (+ their own Required), array elements via Items.
type ToolProperty struct {
	Type        string // "string" | "integer" | "number" | "boolean" | "array" | "object"
	Description string
	Enum        []string                // non-empty when the param is constrained to a fixed set of values
	Properties  map[string]ToolProperty // nested fields when Type == "object"
	Items       *ToolProperty           // element schema when Type == "array"
	Required    []string                // required nested field names when Type == "object"
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
	// Body is a JSON template sent as the request body (Content-Type: application/json).
	// String leaves are interpolated the same way as Headers/Query values: whole-value
	// "__secrets.field" replacement, otherwise ${{params.*}} substitution.
	Body json.RawMessage
	// TODO: replace with a typed enum once external_connections.provider is constrained in the DB schema.
	Credentials string // external_connections.provider to resolve __secrets.*
}
