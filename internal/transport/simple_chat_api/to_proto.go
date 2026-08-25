package simple_chat_api

import (
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
)

// timestampLayout matches the RFC3339-in-UTC rendering every other to_proto converter in this
// repo uses for timestamp strings.
const timestampLayout = "2006-01-02T15:04:05Z"

func chatToProto(chat domain.SimpleChat) *artel_api.SimpleChat {
	out := &artel_api.SimpleChat{
		Id:             chat.Uuid.String(),
		VaultId:        chat.VaultUuid.String(),
		Model:          chat.Model,
		VaultAccess:    chat.VaultAccess,
		CreatedAt:      chat.CreatedAt.UTC().Format(timestampLayout),
		UpdatedAt:      chat.UpdatedAt.UTC().Format(timestampLayout),
		LastActivityAt: chat.LastActivityAt.UTC().Format(timestampLayout),
	}

	if chat.Title != nil {
		out.Title = *chat.Title
	}

	return out
}

func chatsToProto(chats []domain.SimpleChat) []*artel_api.SimpleChat {
	out := make([]*artel_api.SimpleChat, 0, len(chats))
	for _, chat := range chats {
		out = append(out, chatToProto(chat))
	}

	return out
}

// messageToProto flattens a transcript row for the wire. tool_call_id and tool_input are
// deliberately not exposed — they're internal plumbing for replaying history back to the model.
func messageToProto(msg domain.SimpleChatMessage) *artel_api.SimpleChatMessage {
	out := &artel_api.SimpleChatMessage{
		Id:        msg.Uuid.String(),
		Role:      msg.Role,
		Content:   msg.Content,
		IsError:   msg.IsError,
		Seq:       msg.Seq,
		CreatedAt: msg.CreatedAt.UTC().Format(timestampLayout),
	}

	if msg.ToolName != nil {
		out.ToolName = *msg.ToolName
	}

	if msg.Model != nil {
		out.Model = *msg.Model
	}

	return out
}

func messagesToProto(messages []domain.SimpleChatMessage) []*artel_api.SimpleChatMessage {
	out := make([]*artel_api.SimpleChatMessage, 0, len(messages))
	for _, msg := range messages {
		out = append(out, messageToProto(msg))
	}

	return out
}
