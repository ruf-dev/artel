package couchcreds

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/repository"
)

type CouchCredsRepo struct {
	db            sqldb.DB
	encryptionKey []byte
}

func New(db sqldb.DB, encryptionKey []byte) *CouchCredsRepo {
	return &CouchCredsRepo{
		db:            db,
		encryptionKey: encryptionKey,
	}
}

func (r *CouchCredsRepo) Store(ctx context.Context, vaultID uuid.UUID, host, username string, passwordPlain []byte) error {
	passwordEnc, err := cryptoutil.Encrypt(r.encryptionKey, passwordPlain)
	if err != nil {
		return rerrors.Wrap(err, "error encrypting password")
	}

	query := `INSERT INTO couch_credentials (vault_id, host, username, password_enc) VALUES ($1, $2, $3, $4)`

	_, err = r.db.ExecContext(ctx, query, vaultID, host, username, passwordEnc)
	if err != nil {
		return rerrors.Wrap(err, "error storing couch credentials")
	}

	return nil
}

func (r *CouchCredsRepo) Load(ctx context.Context, vaultID uuid.UUID) (repository.CouchCred, error) {
	query := `SELECT id, vault_id, host, username, password_enc, created_at FROM couch_credentials WHERE vault_id = $1`

	var c repository.CouchCred
	row := r.db.QueryRowContext(ctx, query, vaultID)
	err := row.Scan(&c.ID, &c.VaultID, &c.Host, &c.Username, &c.PasswordEnc, &c.CreatedAt)
	if err != nil {
		return repository.CouchCred{}, rerrors.Wrap(err, "error loading couch credentials")
	}

	return c, nil
}
