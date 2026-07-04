package tract

// Condition operators (domain.TractCondition.Op values needing string comparison beyond
// simple equality/inequality/numeric operators).
const (
	opContains = "contains"
	opGlob     = "glob"
	opRegex    = "regex"
)

// Field/map-key names and JSON-schema type strings repeated across trigger payload schemas
// and templates in this package.
const (
	schemaTypeObject = "object"
	fieldName        = "name"
	fieldMessage     = "message"
	fieldRef         = "ref"
	fieldBranch      = "branch"
)
