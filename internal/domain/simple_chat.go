package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SimpleChatRole is the role of a SimpleChatMessage — see that struct's doc comment.
type SimpleChatRole string

const (
	SimpleChatRoleUser      SimpleChatRole = "user"
	SimpleChatRoleAssistant SimpleChatRole = "assistant"
	SimpleChatRoleTool      SimpleChatRole = "tool"
)

// SimpleChat is one saved chat thread of the in-process, container-free "Simple Chat" agent —
// a vault member picks an OpenRouter BYOK model and chats with an agent that can call Artel's
// existing MCP tools. A (vault, user) pair may have several SimpleChat rows (multiple saved
// threads); see migrations/075_simple_chats.sql.
type SimpleChat struct {
	Uuid      uuid.UUID
	VaultUuid uuid.UUID
	UserUuid  uuid.UUID
	// Title is nil until the engine (or the user) names the chat — a nil title is rendered as
	// e.g. "New chat" / the first user message by the frontend, not stored as such here.
	Title          *string
	Model          string
	VaultAccess    bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastActivityAt time.Time
}

// SimpleChatMessage is one turn (or tool step) in a SimpleChat's history — see
// migrations/076_simple_chat_messages.sql. Role is one of the SimpleChatRole constants above:
//   - user: ToolCallID/ToolName/ToolInput/Model unset, Content is the member's message text.
//   - assistant: Content is the model's reply text (possibly empty when the turn was pure tool
//     calls); Model is the model that produced it. ToolCallID/ToolName/ToolInput are unset here
//     — a requested tool call is recorded as its own row (see below), mirroring how
//     internal/chatprotocol.Event separates assistant text from tool_call_started.
//   - tool: one row per tool invocation the agent made — ToolCallID/ToolName/ToolInput carry the
//     call, Content carries the tool's result (or the error message when IsError is true).
type SimpleChatMessage struct {
	Uuid       uuid.UUID
	ChatUuid   uuid.UUID
	Role       string
	Content    string
	ToolCallID *string
	ToolName   *string
	ToolInput  json.RawMessage
	IsError    bool
	Model      *string
	Seq        int64
	CreatedAt  time.Time
}
