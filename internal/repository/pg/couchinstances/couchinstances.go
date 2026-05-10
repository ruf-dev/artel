package couchinstances

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

type CouchInstancesRepo struct {
	q             *artel_q.Queries
	encryptionKey []byte
}

func New(q *artel_q.Queries, encryptionKey []byte) *CouchInstancesRepo {
	return &CouchInstancesRepo{
		q:             q,
		encryptionKey: encryptionKey,
	}
}

func (r *CouchInstancesRepo) Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error) {
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
		return uuid.UUID{}, rerrors.Wrap(err, "insert couch instance")
	}

	return id, nil
}

func (r *CouchInstancesRepo) Get(ctx context.Context, id uuid.UUID) (domain.CouchInstance, error) {
	row, err := r.q.GetCouchInstance(ctx, id)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "get couch instance")
	}

	instance := domain.CouchInstance{
		Uuid:      row.ID,
		Url:       row.Url,
		Username:  row.Username,
		CreatedAt: row.CreatedAt,
	}
	return instance, nil
}

func (r *CouchInstancesRepo) RandomPick(ctx context.Context) (domain.CouchInstance, error) {
	row, err := r.q.RandomPickCouchInstance(ctx)
	if err != nil {
		return domain.CouchInstance{}, rerrors.Wrap(err, "random pick couch instance")
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

func (r *CouchInstancesRepo) List(ctx context.Context) ([]domain.CouchInstance, error) {
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

func (r *CouchInstancesRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteCouchInstance(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "delete couch instance")
	}
	return nil
}
