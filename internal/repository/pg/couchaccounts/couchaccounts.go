package couchaccounts

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type CouchAccountsRepo struct {
	q             *artel_q.Queries
	encryptionKey []byte
}

func New(q *artel_q.Queries, encryptionKey []byte) *CouchAccountsRepo {
	return &CouchAccountsRepo{
		q:             q,
		encryptionKey: encryptionKey,
	}
}

func (r *CouchAccountsRepo) Create(ctx context.Context, userID, instanceID uuid.UUID, username string, passwordPlain []byte) (domain.CouchAccount, error) {
	passwordEnc, err := cryptoutil.Encrypt(r.encryptionKey, passwordPlain)
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
		return domain.CouchAccount{}, rerrors.Wrap(err, "create couch account")
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

func (r *CouchAccountsRepo) GetByUserAndInstance(ctx context.Context, userID, instanceID uuid.UUID) (domain.CouchAccount, error) {
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

func (r *CouchAccountsRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.CouchAccount, error) {
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

func (r *CouchAccountsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteCouchAccount(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "delete couch account")
	}
	return nil
}
