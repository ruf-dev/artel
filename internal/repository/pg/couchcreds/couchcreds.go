package couchcreds

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type CouchCredsRepo struct {
	q             *artel_q.Queries
	encryptionKey []byte
}

func New(q *artel_q.Queries, encryptionKey []byte) *CouchCredsRepo {
	return &CouchCredsRepo{
		q:             q,
		encryptionKey: encryptionKey,
	}
}

func (r *CouchCredsRepo) Store(ctx context.Context, vaultID uuid.UUID, host, username string, passwordPlain []byte) error {
	passwordEnc, err := cryptoutil.Encrypt(r.encryptionKey, passwordPlain)
	if err != nil {
		return rerrors.Wrap(err, "error encrypting password")
	}

	params := artel_q.StoreCouchCredParams{
		VaultID:     vaultID,
		Host:        host,
		Username:    username,
		PasswordEnc: passwordEnc,
	}
	err = r.q.StoreCouchCred(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error storing couch credentials")
	}

	return nil
}

func (r *CouchCredsRepo) Load(ctx context.Context, vaultID uuid.UUID) (domain.CouchCred, error) {
	row, err := r.q.LoadCouchCred(ctx, vaultID)
	if err != nil {
		return domain.CouchCred{}, rerrors.Wrap(err, "error loading couch credentials")
	}

	passwordPlain, err := cryptoutil.Decrypt(r.encryptionKey, row.PasswordEnc)
	if err != nil {
		return domain.CouchCred{}, rerrors.Wrap(err, "error decrypting password")
	}

	c := domain.CouchCred{
		Uuid:      row.ID,
		VaultUuid: row.VaultID,
		Host:      row.Host,
		Username:  row.Username,
		Password:  passwordPlain,
		CreatedAt: row.CreatedAt,
	}
	return c, nil
}

func (r *CouchCredsRepo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.DeleteCouchCred(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting couch credentials")
	}

	return nil
}
