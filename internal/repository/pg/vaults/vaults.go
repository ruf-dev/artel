package vaults

import (
	"context"
	"time"

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

func (r *VaultsRepo) Create(ctx context.Context, userID, couchInstanceID uuid.UUID, name, couchDBName string) (domain.Vault, error) {
	params := artel_q.CreateVaultParams{
		UserID:          userID,
		Name:            name,
		CouchDbName:     couchDBName,
		CouchInstanceID: uuid.NullUUID{UUID: couchInstanceID, Valid: true},
	}
	row, err := r.q.CreateVault(ctx, params)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error creating vault")
	}

	v := rowToVault(row.ID, row.UserID, row.Name, row.CouchDbName, row.CouchInstanceID, row.CreatedAt)
	return v, nil
}

func (r *VaultsRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
	row, err := r.q.GetVaultByID(ctx, id)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting vault by id")
	}

	v := rowToVault(row.ID, row.UserID, row.Name, row.CouchDbName, row.CouchInstanceID, row.CreatedAt)
	return v, nil
}

func (r *VaultsRepo) ListByMembership(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error) {
	rows, err := r.q.ListVaultsByMembership(ctx, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing vaults by membership")
	}

	vaultList := make([]domain.Vault, 0, len(rows))
	for _, row := range rows {
		v := rowToVault(row.ID, row.UserID, row.Name, row.CouchDbName, row.CouchInstanceID, row.CreatedAt)
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

func rowToVault(id, userID uuid.UUID, name, couchDbName string, couchInstanceID uuid.NullUUID, createdAt time.Time) domain.Vault {
	v := domain.Vault{
		Uuid:        id,
		UserUuid:    userID,
		Name:        name,
		CouchDBName: couchDbName,
		CreatedAt:   createdAt,
	}
	if couchInstanceID.Valid {
		v.CouchInstanceUuid = couchInstanceID.UUID
	}
	return v
}
