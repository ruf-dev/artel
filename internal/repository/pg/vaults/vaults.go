package vaults

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
)

type VaultsRepo struct {
	db sqldb.DB
}

func New(db sqldb.DB) *VaultsRepo {
	return &VaultsRepo{db: db}
}

func (r *VaultsRepo) Create(ctx context.Context, userID uuid.UUID, name, couchDBName string) (domain.Vault, error) {
	query := `INSERT INTO vaults (user_id, name, couch_db_name) VALUES ($1, $2, $3) RETURNING id, user_id, name, couch_db_name, created_at`

	var v domain.Vault
	row := r.db.QueryRowContext(ctx, query, userID, name, couchDBName)
	err := row.Scan(&v.Uuid, &v.UserUuid, &v.Name, &v.CouchDBName, &v.CreatedAt)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error creating vault")
	}

	return v, nil
}

func (r *VaultsRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error) {
	query := `SELECT id, user_id, name, couch_db_name, created_at FROM vaults WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing vaults")
	}
	defer rows.Close()

	var vaultList []domain.Vault
	for rows.Next() {
		var v domain.Vault
		err = rows.Scan(&v.Uuid, &v.UserUuid, &v.Name, &v.CouchDBName, &v.CreatedAt)
		if err != nil {
			return nil, rerrors.Wrap(err, "error scanning vault")
		}
		vaultList = append(vaultList, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, rerrors.Wrap(err, "error iterating vaults")
	}

	return vaultList, nil
}
