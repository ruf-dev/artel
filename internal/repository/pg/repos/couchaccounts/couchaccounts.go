package couchaccounts

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
)

type Repo struct {
	q             *artel_q.Queries
	encryptionKey []byte
}

func New(db sqldb.DB, encryptionKey []byte) *Repo {
	return &Repo{
		q:             artel_q.New(db),
		encryptionKey: encryptionKey,
	}
}

func (r *Repo) Upsert(ctx context.Context, userID, instanceID uuid.UUID, username string, passwordPlain string) (domain.CouchAccount, error) {
	passwordEnc, err := cryptoutil.Encrypt(r.encryptionKey, []byte(passwordPlain))
	if err != nil {
		return domain.CouchAccount{}, rerrors.Wrap(err, "encrypt password")
	}

	params := artel_q.CreateCouchAccountParams{
		UserID:           userID,
		CouchInstanceID:  instanceID,
		CouchUsername:    username,
		CouchPasswordEnc: passwordEnc,
	}
	row, err := r.q.CreateCouchAccount(ctx, params)
	if err != nil {
		return domain.CouchAccount{}, pg_err.UnwrapPgErr(err)
	}

	account := domain.CouchAccount{
		Uuid:              row.ID,
		UserUuid:          row.UserID,
		CouchInstanceUuid: row.CouchInstanceID,
		CouchUsername:     row.CouchUsername,
		CouchPassword:     "",
		CreatedAt:         row.CreatedAt,
	}
	return account, nil
}

func (r *Repo) GetByUserAndInstance(ctx context.Context, userID, instanceID uuid.UUID) (domain.CouchAccount, error) {
	params := artel_q.GetCouchAccountByUserAndInstanceParams{
		UserID:          userID,
		CouchInstanceID: instanceID,
	}
	row, err := r.q.GetCouchAccountByUserAndInstance(ctx, params)
	if err != nil {
		return domain.CouchAccount{}, rerrors.Wrap(err, "get couch account")
	}

	passwordPlain, err := cryptoutil.Decrypt(r.encryptionKey, row.CouchPasswordEnc)
	if err != nil {
		return domain.CouchAccount{}, rerrors.Wrap(err, "decrypt password")
	}

	account := domain.CouchAccount{
		Uuid:              row.ID,
		UserUuid:          row.UserID,
		CouchInstanceUuid: row.CouchInstanceID,
		CouchUsername:     row.CouchUsername,
		CouchPassword:     string(passwordPlain),
		CreatedAt:         row.CreatedAt,
	}
	return account, nil
}

func (r *Repo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.CouchAccount, error) {
	rows, err := r.q.ListCouchAccountsByUser(ctx, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list couch accounts")
	}

	accounts := make([]domain.CouchAccount, len(rows))
	for i, row := range rows {
		passwordPlain, err := cryptoutil.Decrypt(r.encryptionKey, row.CouchPasswordEnc)
		if err != nil {
			return nil, rerrors.Wrap(err, "decrypt password")
		}

		accounts[i] = domain.CouchAccount{
			Uuid:              row.ID,
			UserUuid:          row.UserID,
			CouchInstanceUuid: row.CouchInstanceID,
			CouchUsername:     row.CouchUsername,
			CouchPassword:     string(passwordPlain),
			CreatedAt:         row.CreatedAt,
		}
	}
	return accounts, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteCouchAccount(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "delete couch account")
	}
	return nil
}

func (r *Repo) WithTx(tx sqldb.DB) repository.CouchAccounts {
	return New(tx, r.encryptionKey)
}
