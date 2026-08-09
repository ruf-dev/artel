package executors

// JSON-schema type strings, repeated across every ToolProperty literal in this package.
const (
	schemaTypeString  = "string"
	schemaTypeObject  = "object"
	schemaTypeArray   = "array"
	schemaTypeBoolean = "boolean"
	schemaTypeInteger = "integer"
)

// Field/map-key names shared across tool schema definitions and result maps in this package.
const (
	fieldName         = "name"
	fieldDescription  = "description"
	fieldDefinition   = "definition"
	fieldUuid         = "uuid"
	fieldSource       = "source"
	fieldFilters      = "filters"
	fieldStatus       = "status"
	fieldTract        = "tract"
	fieldWarnings     = "warnings"
	fieldMessage      = "message"
	fieldMcp          = "mcp"
	fieldCreatedAt    = "created_at"
	fieldUpdatedAt    = "updated_at"
	fieldTractUuid    = "tract_uuid"
	fieldTriggerUuid  = "trigger_uuid"
	fieldWebhookUrl   = "webhook_url"
	fieldWebhookToken = "webhook_token"
	fieldPath         = "path"
	fieldMtime        = "mtime"
	fieldMimeType     = "mimeType"
	fieldRev          = "rev"
	fieldCtime        = "ctime"
	fieldSize         = "size"
	fieldDeleted      = "deleted"
	fieldContent      = "content"
	fieldEnabled      = "enabled"
	fieldTableName    = "table_name"
	fieldSql          = "sql"
	fieldRows         = "rows"
	fieldTruncated    = "truncated"
)

const msgConfirmation = "Confirmation message"
