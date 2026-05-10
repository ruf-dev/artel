package vaults

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/repository"
)

type VaultsRepo struct {
	db sqldb.DB
}

func New(db sqldb.DB) *VaultsRepo {
	return &VaultsRepo{db: db}
}

func (r *VaultsRepo) Create(ctx context.Context, userID uuid.UUID, name, couchDBName string) (repository.Vault, error) {
	query := `INSERT INTO vaults (user_id, name, couch_db_name) VALUES ($1, $2, $3) RETURNING id, user_id, name, couch_db_name, created_at`

	var v repository.Vault
	row := r.db.QueryRowContext(ctx, query, userID, name, couchDBName)
	err := row.Scan(&v.ID, &v.UserID, &v.Name, &v.CouchDBName, &v.CreatedAt)
	if err != nil {
		return repository.Vault{}, rerrors.Wrap(err, "error creating vault")
	}

	return v, nil
}

func (r *VaultsRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]repository.Vault, error) {
	query := `SELECT id, user_id, name, couch_db_name, created_at FROM vaults WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing vaults")
	}
	defer rows.Close()

	var vaultList []repository.Vault
	for rows.Next() {
		var v repository.Vault
		err = rows.Scan(&v.ID, &v.UserID, &v.Name, &v.CouchDBName, &v.CreatedAt)
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
