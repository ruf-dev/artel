package domain

import (
	"bufio"
	"bytes"
	"encoding/json"
	"sort"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"
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

// ChatHistoryFolderPrefix is the vault-path prefix under which every SimpleChat's JSONL
// transcript lives — hidden from the regular Notes/Files listings the same way
// SkillsFolderPrefix hides .skills/ (see internal/clients/couchdb/livesync.go).
const ChatHistoryFolderPrefix = ".chat_history/"

// ChatHistoryPath returns the vault-relative CouchDB path a SimpleChat's JSONL transcript lives
// at — one doc per chat, namespaced under the owning user so a per-user listing (ListChats) never
// has to read another member's threads.
func ChatHistoryPath(userUuid, chatUuid uuid.UUID) string {
	return ChatHistoryFolderPrefix + userUuid.String() + "/" + chatUuid.String() + ".jsonl"
}

// simpleChatLineType discriminates the JSON lines making up a chat's JSONL transcript.
type simpleChatLineType string

const (
	simpleChatLineHeader  simpleChatLineType = "header"
	simpleChatLineMessage simpleChatLineType = "message"
)

// SimpleChatHeader is line 0 of a chat's JSONL transcript: the thread's own metadata plus its
// remembered "allow always" tool-permission decisions — the CouchDB-era replacement for
// migrations/075_simple_chats.sql's row shape and migrations/077_simple_chat_tool_allowances.sql.
type SimpleChatHeader struct {
	Type      simpleChatLineType `json:"type"`
	Uuid      uuid.UUID          `json:"id"`
	VaultUuid uuid.UUID          `json:"vault_id"`
	UserUuid  uuid.UUID          `json:"user_id"`
	// Title is nil until the engine (or the user) names the chat — see SimpleChat.Title.
	Title          *string   `json:"title,omitempty"`
	Model          string    `json:"model"`
	VaultAccess    bool      `json:"vault_access"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
	// ToolAllowances remembers a chat's "allow always" answers, keyed by tool name.
	ToolAllowances map[string]string `json:"tool_allowances,omitempty"`
}

// ToSimpleChat projects a decoded header into the SimpleChat wire/service shape.
func (h SimpleChatHeader) ToSimpleChat() SimpleChat {
	chat := SimpleChat{
		Uuid:           h.Uuid,
		VaultUuid:      h.VaultUuid,
		UserUuid:       h.UserUuid,
		Title:          h.Title,
		Model:          h.Model,
		VaultAccess:    h.VaultAccess,
		CreatedAt:      h.CreatedAt,
		UpdatedAt:      h.UpdatedAt,
		LastActivityAt: h.LastActivityAt,
	}

	return chat
}

// simpleChatMessageLine is one JSONL line (lines 1..N) of a chat's transcript — the on-disk shape
// SimpleChatMessage marshals to/from.
type simpleChatMessageLine struct {
	Type       simpleChatLineType `json:"type"`
	Seq        int64              `json:"seq"`
	Role       string             `json:"role"`
	Content    string             `json:"content"`
	ToolCallID *string            `json:"tool_call_id,omitempty"`
	ToolName   *string            `json:"tool_name,omitempty"`
	ToolInput  json.RawMessage    `json:"tool_input,omitempty"`
	IsError    bool               `json:"is_error,omitempty"`
	Model      *string            `json:"model,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

// SimpleChatFile is a chat's JSONL transcript decoded into memory: the header line plus every
// message line, in file order. simplechat.Service reads, mutates, and rewrites this whole on
// every change — see ChatHistoryPath.
type SimpleChatFile struct {
	Header   SimpleChatHeader
	Messages []SimpleChatMessage
}

// HasMessages reports whether f has at least one message line — the CouchDB-era replacement for
// the Postgres EXISTS(...) check ListSimpleChatsByVaultAndUser used to filter on: a thread
// created but never sent a first message stays hidden from ListChats until this turns true.
func (f SimpleChatFile) HasMessages() bool {
	return len(f.Messages) > 0
}

// NextSeq returns the seq value the next appended message should carry — 1-based position among
// message lines, so it needs no separately stored counter.
func (f SimpleChatFile) NextSeq() int64 {
	return int64(len(f.Messages)) + 1
}

// EncodeSimpleChatFile renders f as JSONL bytes: the header line, then one line per message, in
// f.Messages order.
func EncodeSimpleChatFile(file SimpleChatFile) ([]byte, error) {
	var buf bytes.Buffer

	header := file.Header
	header.Type = simpleChatLineHeader

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return nil, rerrors.Wrap(err, "error marshaling simple chat header")
	}

	buf.Write(headerJSON)
	buf.WriteByte('\n')

	for _, msg := range file.Messages {
		line := simpleChatMessageLine{
			Type:       simpleChatLineMessage,
			Seq:        msg.Seq,
			Role:       msg.Role,
			Content:    msg.Content,
			ToolCallID: msg.ToolCallID,
			ToolName:   msg.ToolName,
			ToolInput:  msg.ToolInput,
			IsError:    msg.IsError,
			Model:      msg.Model,
			CreatedAt:  msg.CreatedAt,
		}

		lineJSON, marshalErr := json.Marshal(line)
		if marshalErr != nil {
			return nil, rerrors.Wrap(marshalErr, "error marshaling simple chat message")
		}

		buf.Write(lineJSON)
		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

// DecodeSimpleChatFile parses JSONL bytes written by EncodeSimpleChatFile back into a
// SimpleChatFile. Each decoded message's ChatUuid is stamped from the header's id, since the
// per-line JSON carries no chat id of its own — it's implied by the file itself.
func DecodeSimpleChatFile(content []byte) (SimpleChatFile, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	var file SimpleChatFile

	headerSeen := false

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var kind struct {
			Type simpleChatLineType `json:"type"`
		}

		err := json.Unmarshal(line, &kind)
		if err != nil {
			return SimpleChatFile{}, rerrors.Wrap(err, "error unmarshaling simple chat line type")
		}

		switch kind.Type {
		case simpleChatLineHeader:
			var header SimpleChatHeader

			err = json.Unmarshal(line, &header)
			if err != nil {
				return SimpleChatFile{}, rerrors.Wrap(err, "error unmarshaling simple chat header")
			}

			file.Header = header
			headerSeen = true
		case simpleChatLineMessage:
			var msgLine simpleChatMessageLine

			err = json.Unmarshal(line, &msgLine)
			if err != nil {
				return SimpleChatFile{}, rerrors.Wrap(err, "error unmarshaling simple chat message")
			}

			msg := SimpleChatMessage{
				ChatUuid:   file.Header.Uuid,
				Role:       msgLine.Role,
				Content:    msgLine.Content,
				ToolCallID: msgLine.ToolCallID,
				ToolName:   msgLine.ToolName,
				ToolInput:  msgLine.ToolInput,
				IsError:    msgLine.IsError,
				Model:      msgLine.Model,
				Seq:        msgLine.Seq,
				CreatedAt:  msgLine.CreatedAt,
			}

			file.Messages = append(file.Messages, msg)
		}
	}

	err := scanner.Err()
	if err != nil {
		return SimpleChatFile{}, rerrors.Wrap(err, "error scanning simple chat jsonl")
	}

	if !headerSeen {
		return SimpleChatFile{}, rerrors.New("simple chat jsonl file has no header line")
	}

	return file, nil
}

// SortSimpleChatsByLastActivityDesc sorts chats by LastActivityAt, most recent first — the order
// ListChats returns threads in.
func SortSimpleChatsByLastActivityDesc(chats []SimpleChat) {
	sort.Slice(chats, func(i, j int) bool {
		return chats[i].LastActivityAt.After(chats[j].LastActivityAt)
	})
}
