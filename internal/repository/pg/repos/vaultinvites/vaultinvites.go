package vaultinvites

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type Repo struct {
	q *artel_q.Queries
}

func New(db sqldb.DB) *Repo {
	return &Repo{
		q: artel_q.New(db),
	}
}

func (r *Repo) Create(ctx context.Context, vaultID, createdBy uuid.UUID, role artel_q.VaultRole, token string) (domain.VaultInvite, error) {
	params := artel_q.CreateVaultInviteParams{
		VaultID:   vaultID,
		CreatedBy: createdBy,
		Role:      role,
		Token:     token,
	}

	inv, err := r.q.CreateVaultInvite(ctx, params)
	if err != nil {
		return domain.VaultInvite{}, rerrors.Wrap(err, "create vault invite")
	}

	return mapInvite(inv), nil
}

func (r *Repo) GetByToken(ctx context.Context, token string) (domain.VaultInvite, error) {
	inv, err := r.q.GetVaultInviteByToken(ctx, token)
	if err != nil {
		return domain.VaultInvite{}, rerrors.Wrap(err, "get vault invite by token")
	}

	return mapInvite(inv), nil
}

func (r *Repo) ListByVault(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultInvite, error) {
	rows, err := r.q.ListVaultInvites(ctx, vaultID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list vault invites")
	}

	result := make([]domain.VaultInvite, len(rows))
	for i, row := range rows {
		result[i] = mapInvite(row)
	}
	return result, nil
}

func (r *Repo) Revoke(ctx context.Context, id uuid.UUID) error {
	err := r.q.RevokeVaultInvite(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "revoke vault invite")
	}
	return nil
}

func (r *Repo) WithTx(tx sqldb.DB) repository.VaultInvites {
	return New(tx)
}

func mapInvite(inv artel_q.VaultInvite) domain.VaultInvite {
	d := domain.VaultInvite{
		Uuid:      inv.ID,
		VaultUuid: inv.VaultID,
		CreatedBy: inv.CreatedBy,
		Role:      inv.Role,
		Token:     inv.Token,
		CreatedAt: inv.CreatedAt,
	}
	if inv.RevokedAt.Valid {
		d.RevokedAt = &inv.RevokedAt.Time
	}
	return d
}
