package vaults

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type VaultsRepo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *VaultsRepo {
	return &VaultsRepo{q: q}
}

func (r *VaultsRepo) Create(ctx context.Context, userID uuid.UUID, name, couchDBName string) (domain.Vault, error) {
	params := artel_q.CreateVaultParams{
		UserID:      userID,
		Name:        name,
		CouchDbName: couchDBName,
	}
	row, err := r.q.CreateVault(ctx, params)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error creating vault")
	}

	v := domain.Vault{
		Uuid:        row.ID,
		UserUuid:    row.UserID,
		Name:        row.Name,
		CouchDBName: row.CouchDbName,
		CreatedAt:   row.CreatedAt,
	}
	return v, nil
}

func (r *VaultsRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error) {
	rows, err := r.q.ListVaultsByUser(ctx, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing vaults by user")
	}

	vaultList := make([]domain.Vault, 0, len(rows))
	for _, row := range rows {
		v := domain.Vault{
			Uuid:        row.ID,
			UserUuid:    row.UserID,
			Name:        row.Name,
			CouchDBName: row.CouchDbName,
			CreatedAt:   row.CreatedAt,
		}
		vaultList = append(vaultList, v)
	}
	return vaultList, nil
}

func (r *VaultsRepo) GetByName(ctx context.Context, name string) (domain.Vault, error) {
	row, err := r.q.GetVaultByName(ctx, name)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting vault by name")
	}

	v := domain.Vault{
		Uuid:        row.ID,
		UserUuid:    row.UserID,
		Name:        row.Name,
		CouchDBName: row.CouchDbName,
		CreatedAt:   row.CreatedAt,
	}
	return v, nil
}

func (r *VaultsRepo) ListAll(ctx context.Context) ([]domain.Vault, error) {
	rows, err := r.q.ListAllVaults(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing vaults")
	}

	vaultList := make([]domain.Vault, 0, len(rows))
	for _, row := range rows {
		v := domain.Vault{
			Uuid:        row.ID,
			UserUuid:    row.UserID,
			Name:        row.Name,
			CouchDBName: row.CouchDbName,
			CreatedAt:   row.CreatedAt,
		}
		vaultList = append(vaultList, v)
	}
	return vaultList, nil
}

func (r *VaultsRepo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.DeleteVault(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting vault")
	}

	return nil
}
