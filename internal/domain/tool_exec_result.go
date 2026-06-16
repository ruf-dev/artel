package domain

// ToolExecResult is the result of executing a built-in MCP tool (vault tools, connections).
// Text holds marshaled JSON/plain text results; Data/MimeType/ResourcePath are set instead
// for binary file reads.
type ToolExecResult struct {
	Text         string
	Data         []byte
	MimeType     string
	ResourcePath string
}
