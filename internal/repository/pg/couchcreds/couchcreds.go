package couchcreds

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
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

func (r *CouchCredsRepo) Load(ctx context.Context, vaultID uuid.UUID) (domain.CouchCred, error) {
	query := `SELECT id, vault_id, host, username, password_enc, created_at FROM couch_credentials WHERE vault_id = $1`

	var c domain.CouchCred
	var passwordEnc []byte
	row := r.db.QueryRowContext(ctx, query, vaultID)
	err := row.Scan(&c.Uuid, &c.VaultUuid, &c.Host, &c.Username, &passwordEnc, &c.CreatedAt)
	if err != nil {
		return domain.CouchCred{}, rerrors.Wrap(err, "error loading couch credentials")
	}

	passwordPlain, err := cryptoutil.Decrypt(r.encryptionKey, passwordEnc)
	if err != nil {
		return domain.CouchCred{}, rerrors.Wrap(err, "error decrypting password")
	}

	c.Password = passwordPlain

	return c, nil
}

func (r *CouchCredsRepo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	query := `DELETE FROM couch_credentials WHERE vault_id = $1`

	_, err := r.db.ExecContext(ctx, query, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting couch credentials")
	}

	return nil
}
