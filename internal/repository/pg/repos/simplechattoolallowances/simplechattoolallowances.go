// Package simplechattoolallowances is the pure-DB layer for simple_chat_tool_allowances — the
// per-chat "allow always" memory for a tool permission prompt (see
// migrations/077_simple_chat_tool_allowances.sql). decision is stored as the raw string of
// whatever internal/chatprotocol.PermissionDecision value the caller chose to persist (in
// practice only DecisionAllowAlways is ever expected to be upserted here — allow_once/deny
// don't need to be remembered past the current turn).
package simplechattoolallowances

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
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

func (r *Repo) Upsert(ctx context.Context, chatUuid uuid.UUID, toolName, decision string) error {
	params := artel_q.UpsertSimpleChatToolAllowanceParams{
		ChatID:   chatUuid,
		ToolName: toolName,
		Decision: decision,
	}

	err := r.q.UpsertSimpleChatToolAllowance(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error upserting simple chat tool allowance")
	}

	return nil
}

// Get returns Valid: false when chatUuid has no stored decision for toolName yet.
func (r *Repo) Get(ctx context.Context, chatUuid uuid.UUID, toolName string) (sql.Null[string], error) {
	params := artel_q.GetSimpleChatToolAllowanceParams{
		ChatID:   chatUuid,
		ToolName: toolName,
	}

	decision, err := r.q.GetSimpleChatToolAllowance(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[string]{}, nil
		}

		return sql.Null[string]{}, rerrors.Wrap(err, "error getting simple chat tool allowance")
	}

	result := sql.Null[string]{V: decision, Valid: true}

	return result, nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.SimpleChatToolAllowances {
	return New(tx)
}
