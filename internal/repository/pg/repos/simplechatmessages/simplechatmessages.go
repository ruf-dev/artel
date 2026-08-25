// Package simplechatmessages is the pure-DB layer for simple_chat_messages — see
// domain.SimpleChatMessage and migrations/076_simple_chat_messages.sql.
package simplechatmessages

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/sqlc-dev/pqtype"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(db postgres.DB) *Repo {
	return &Repo{
		q: artel_q.New(db),
	}
}

// Insert persists msg. Uuid/CreatedAt on msg are ignored — the DB assigns both.
func (r *Repo) Insert(ctx context.Context, msg domain.SimpleChatMessage) (domain.SimpleChatMessage, error) {
	params := artel_q.InsertSimpleChatMessageParams{
		ChatID:  msg.ChatUuid,
		Role:    msg.Role,
		Content: msg.Content,
		IsError: msg.IsError,
		Seq:     msg.Seq,
	}

	if msg.ToolCallID != nil {
		params.ToolCallID = sql.NullString{String: *msg.ToolCallID, Valid: true}
	}

	if msg.ToolName != nil {
		params.ToolName = sql.NullString{String: *msg.ToolName, Valid: true}
	}

	if len(msg.ToolInput) > 0 {
		params.ToolInput = pqtype.NullRawMessage{RawMessage: msg.ToolInput, Valid: true}
	}

	if msg.Model != nil {
		params.Model = sql.NullString{String: *msg.Model, Valid: true}
	}

	row, err := r.q.InsertSimpleChatMessage(ctx, params)
	if err != nil {
		return domain.SimpleChatMessage{}, rerrors.Wrap(err, "error inserting simple chat message")
	}

	return rowToSimpleChatMessage(row), nil
}

func (r *Repo) ListByChatID(ctx context.Context, chatUuid uuid.UUID) ([]domain.SimpleChatMessage, error) {
	rows, err := r.q.ListSimpleChatMessagesByChatID(ctx, chatUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing simple chat messages by chat id")
	}

	messages := make([]domain.SimpleChatMessage, 0, len(rows))
	for _, row := range rows {
		messages = append(messages, rowToSimpleChatMessage(row))
	}

	return messages, nil
}

// GetMaxSeq returns the highest seq value stored for chatUuid, or 0 if it has no messages yet
// — the caller adds 1 to compute the next message's seq.
func (r *Repo) GetMaxSeq(ctx context.Context, chatUuid uuid.UUID) (int64, error) {
	maxSeq, err := r.q.GetMaxSeqForChat(ctx, chatUuid)
	if err != nil {
		return 0, rerrors.Wrap(err, "error getting max seq for simple chat")
	}

	return maxSeq, nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.SimpleChatMessages {
	return New(tx)
}

func rowToSimpleChatMessage(row artel_q.SimpleChatMessage) domain.SimpleChatMessage {
	msg := domain.SimpleChatMessage{
		Uuid:      row.ID,
		ChatUuid:  row.ChatID,
		Role:      row.Role,
		Content:   row.Content,
		IsError:   row.IsError,
		Seq:       row.Seq,
		CreatedAt: row.CreatedAt,
	}

	if row.ToolCallID.Valid {
		toolCallID := row.ToolCallID.String
		msg.ToolCallID = &toolCallID
	}

	if row.ToolName.Valid {
		toolName := row.ToolName.String
		msg.ToolName = &toolName
	}

	if row.ToolInput.Valid {
		msg.ToolInput = row.ToolInput.RawMessage
	}

	if row.Model.Valid {
		model := row.Model.String
		msg.Model = &model
	}

	return msg
}
