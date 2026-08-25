// Package simplechats is the pure-DB layer for simple_chats — see domain.SimpleChat and
// migrations/075_simple_chats.sql.
package simplechats

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/domain"
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

func (r *Repo) Create(ctx context.Context, vaultUuid, userUuid uuid.UUID, model string, vaultAccess bool) (domain.SimpleChat, error) {
	params := artel_q.CreateSimpleChatParams{
		VaultID:     vaultUuid,
		UserID:      userUuid,
		Model:       model,
		VaultAccess: vaultAccess,
	}

	row, err := r.q.CreateSimpleChat(ctx, params)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(err, "error creating simple chat")
	}

	return rowToSimpleChat(row), nil
}

func (r *Repo) GetByID(ctx context.Context, chatUuid uuid.UUID) (domain.SimpleChat, error) {
	row, err := r.q.GetSimpleChatByID(ctx, chatUuid)
	if err != nil {
		return domain.SimpleChat{}, rerrors.Wrap(err, "error getting simple chat by id")
	}

	return rowToSimpleChat(row), nil
}

func (r *Repo) ListByVaultAndUser(ctx context.Context, vaultUuid, userUuid uuid.UUID) ([]domain.SimpleChat, error) {
	params := artel_q.ListSimpleChatsByVaultAndUserParams{
		VaultID: vaultUuid,
		UserID:  userUuid,
	}

	rows, err := r.q.ListSimpleChatsByVaultAndUser(ctx, params)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing simple chats by vault and user")
	}

	chats := make([]domain.SimpleChat, 0, len(rows))
	for _, row := range rows {
		chats = append(chats, rowToSimpleChat(row))
	}

	return chats, nil
}

func (r *Repo) UpdateLastActivity(ctx context.Context, chatUuid uuid.UUID) error {
	err := r.q.UpdateSimpleChatLastActivity(ctx, chatUuid)
	if err != nil {
		return rerrors.Wrap(err, "error updating simple chat last activity")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, chatUuid uuid.UUID) error {
	err := r.q.DeleteSimpleChat(ctx, chatUuid)
	if err != nil {
		return rerrors.Wrap(err, "error deleting simple chat")
	}

	return nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.SimpleChats {
	return New(tx)
}

func rowToSimpleChat(row artel_q.SimpleChat) domain.SimpleChat {
	chat := domain.SimpleChat{
		Uuid:           row.ID,
		VaultUuid:      row.VaultID,
		UserUuid:       row.UserID,
		Model:          row.Model,
		VaultAccess:    row.VaultAccess,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		LastActivityAt: row.LastActivityAt,
	}

	if row.Title.Valid {
		title := row.Title.String
		chat.Title = &title
	}

	return chat
}
