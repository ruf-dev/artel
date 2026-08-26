package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestChatHistoryPath(t *testing.T) {
	userUuid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	chatUuid := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	got := ChatHistoryPath(userUuid, chatUuid)

	want := ".chat_history/11111111-1111-1111-1111-111111111111/22222222-2222-2222-2222-222222222222.jsonl"
	require.Equal(t, want, got)
}

// TestEncodeDecodeSimpleChatFile_RoundTrip guards the JSONL wire format: a header plus several
// messages must decode back byte-for-byte equivalent, including the tool-call fields that are
// only ever set on "tool" role messages.
func TestEncodeDecodeSimpleChatFile_RoundTrip(t *testing.T) {
	chatUuid := uuid.New()
	vaultUuid := uuid.New()
	userUuid := uuid.New()

	title := "My chat"
	now := time.Now().UTC().Truncate(time.Second)

	toolCallID := "call_1"
	toolName := "list_notes"
	toolModel := "openrouter/some-model"

	header := SimpleChatHeader{
		// Type is set explicitly (rather than left zero) since EncodeSimpleChatFile stamps it on
		// its own copy of the header, and this struct is compared byte-for-byte against the
		// decoded result below.
		Type:           simpleChatLineHeader,
		Uuid:           chatUuid,
		VaultUuid:      vaultUuid,
		UserUuid:       userUuid,
		Title:          &title,
		Model:          "openrouter/some-model",
		VaultAccess:    true,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastActivityAt: now,
		ToolAllowances: map[string]string{"list_notes": "allow_always"},
	}

	messages := []SimpleChatMessage{
		{
			ChatUuid:  chatUuid,
			Role:      string(SimpleChatRoleUser),
			Content:   "hello",
			Seq:       1,
			CreatedAt: now,
		},
		{
			ChatUuid:   chatUuid,
			Role:       string(SimpleChatRoleTool),
			Content:    `{"ok":true}`,
			ToolCallID: &toolCallID,
			ToolName:   &toolName,
			ToolInput:  json.RawMessage(`{"path":"a.md"}`),
			IsError:    false,
			Seq:        2,
			CreatedAt:  now,
		},
		{
			ChatUuid:  chatUuid,
			Role:      string(SimpleChatRoleAssistant),
			Content:   "done",
			Model:     &toolModel,
			Seq:       3,
			CreatedAt: now,
		},
	}

	file := SimpleChatFile{Header: header, Messages: messages}

	content, err := EncodeSimpleChatFile(file)
	require.NoError(t, err)

	decoded, err := DecodeSimpleChatFile(content)
	require.NoError(t, err)

	require.Equal(t, header, decoded.Header)
	require.Equal(t, messages, decoded.Messages)
}

func TestSimpleChatFile_HasMessages(t *testing.T) {
	tests := []struct {
		name string
		file SimpleChatFile
		want bool
	}{
		{name: "no messages", file: SimpleChatFile{}, want: false},
		{
			name: "one message",
			file: SimpleChatFile{Messages: []SimpleChatMessage{{Seq: 1}}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.file.HasMessages())
		})
	}
}

func TestSimpleChatFile_NextSeq(t *testing.T) {
	tests := []struct {
		name string
		file SimpleChatFile
		want int64
	}{
		{name: "empty file", file: SimpleChatFile{}, want: 1},
		{
			name: "two messages already",
			file: SimpleChatFile{Messages: []SimpleChatMessage{{Seq: 1}, {Seq: 2}}},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.file.NextSeq())
		})
	}
}

func TestDecodeSimpleChatFile_MissingHeaderErrors(t *testing.T) {
	_, err := DecodeSimpleChatFile([]byte(`{"type":"message","seq":1,"role":"user","content":"hi"}` + "\n"))
	require.Error(t, err)
}

func TestSortSimpleChatsByLastActivityDesc(t *testing.T) {
	now := time.Now()

	oldest := SimpleChat{Uuid: uuid.New(), LastActivityAt: now.Add(-2 * time.Hour)}
	middle := SimpleChat{Uuid: uuid.New(), LastActivityAt: now.Add(-1 * time.Hour)}
	newest := SimpleChat{Uuid: uuid.New(), LastActivityAt: now}

	chats := []SimpleChat{oldest, newest, middle}

	SortSimpleChatsByLastActivityDesc(chats)

	require.Equal(t, []SimpleChat{newest, middle, oldest}, chats)
}
