package couchinstances

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

func (r *Repo) Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error) {
	passwordEnc, err := cryptoutil.Encrypt(r.encryptionKey, passwordPlain)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "encrypt password")
	}

	params := artel_q.RegisterCouchInstanceParams{
		Url:         url,
		Username:    username,
		PasswordEnc: passwordEnc,
	}
	id, err := r.q.RegisterCouchInstance(ctx, params)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error) {
	row, err := r.q.GetCouchInstanceWithCreds(ctx, id)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "get couch instance")
	}

	decrypted, err := cryptoutil.Decrypt(r.encryptionKey, row.PasswordEnc)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "decrypt password")
	}

	instance := domain.CouchInstance{
		Uuid:      row.ID,
		Url:       row.Url,
		Username:  row.Username,
		Password:  string(decrypted),
		CreatedAt: row.CreatedAt,
	}
	return instance, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.CouchInstance, error) {
	rows, err := r.q.ListCouchInstances(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "list couch instances")
	}

	instances := make([]domain.CouchInstance, len(rows))
	for i, row := range rows {
		instances[i] = domain.CouchInstance{
			Uuid:      row.ID,
			Url:       row.Url,
			Username:  row.Username,
			CreatedAt: row.CreatedAt,
		}
	}
	return instances, nil
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, url, username string, passwordPlain []byte) error {
	passwordEnc, err := cryptoutil.Encrypt(r.encryptionKey, passwordPlain)
	if err != nil {
		return rerrors.Wrap(err, "encrypt password")
	}

	params := artel_q.UpdateCouchInstanceParams{
		ID:          id,
		Url:         url,
		Username:    username,
		PasswordEnc: passwordEnc,
	}
	err = r.q.UpdateCouchInstance(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "update couch instance")
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteCouchInstance(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "delete couch instance")
	}
	return nil
}

func (r *Repo) RandomPick(ctx context.Context) (domain.CouchInstanceWithAccount, error) {
	row, err := r.q.RandomPickCouchInstance(ctx)
	if err != nil {
		return domain.CouchInstanceWithAccount{}, pg_err.UnwrapPgErr(err)
	}

	decrypted, err := cryptoutil.Decrypt(r.encryptionKey, row.PasswordEnc)
	if err != nil {
		return domain.CouchInstanceWithAccount{}, rerrors.Wrap(err, "decrypt password")
	}

	result := domain.CouchInstanceWithAccount{
		Instance: domain.CouchInstance{
			Uuid:      row.ID,
			Url:       row.Url,
			Username:  row.Username,
			Password:  string(decrypted),
			CreatedAt: row.CreatedAt,
		},
		Account: nil,
	}

	return result, nil
}

func (r *Repo) WithTx(tx sqldb.DB) repository.CouchInstances {
	return New(tx, r.encryptionKey)
}
