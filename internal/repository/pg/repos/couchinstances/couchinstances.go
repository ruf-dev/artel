package couchinstances

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/postgres"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q         *artel_q.Queries
	encryptor cryptoutil.Encryptor
}

func New(db postgres.DB, encryptor cryptoutil.Encryptor) *Repo {
	return &Repo{
		q:         artel_q.New(db),
		encryptor: encryptor,
	}
}

func (r *Repo) Register(ctx context.Context, url, username string, passwordPlain []byte) (uuid.UUID, error) {
	passwordEnc, err := r.encryptor.Encrypt(passwordPlain)
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

	decrypted, err := r.encryptor.Decrypt(row.PasswordEnc)
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
	passwordEnc, err := r.encryptor.Encrypt(passwordPlain)
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

	decrypted, err := r.encryptor.Decrypt(row.PasswordEnc)
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

// PickForUser tries userID's owned instance first, falling back to RandomPick's shared admin
// pool when they have none — see repository.CouchInstances.PickForUser.
func (r *Repo) PickForUser(ctx context.Context, userID uuid.UUID) (domain.CouchInstanceWithAccount, error) {
	owned, err := r.GetOwned(ctx, userID)
	if err != nil {
		return domain.CouchInstanceWithAccount{}, rerrors.Wrap(err, "get owned couch instance")
	}

	if owned.Valid {
		result := domain.CouchInstanceWithAccount{
			Instance: owned.V,
			Account:  nil,
		}

		return result, nil
	}

	result, err := r.RandomPick(ctx)
	if err != nil {
		return domain.CouchInstanceWithAccount{}, rerrors.Wrap(err, "random pick couch instance")
	}

	return result, nil
}

// GetOwned returns the instance owned by userID, if any — see repository.CouchInstances.GetOwned.
func (r *Repo) GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.CouchInstance], error) {
	ownerUserID := uuid.NullUUID{UUID: userID, Valid: true}

	row, err := r.q.PickOwnedCouchInstance(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.CouchInstance]{}, nil
		}

		return sql.Null[domain.CouchInstance]{}, rerrors.Wrap(err, "pick owned couch instance")
	}

	decrypted, err := r.encryptor.Decrypt(row.PasswordEnc)
	if err != nil {
		return sql.Null[domain.CouchInstance]{}, rerrors.Wrap(err, "decrypt password")
	}

	instance := domain.CouchInstance{
		Uuid:      row.ID,
		Url:       row.Url,
		Username:  row.Username,
		Password:  string(decrypted),
		CreatedAt: row.CreatedAt,
	}

	result := sql.Null[domain.CouchInstance]{V: instance, Valid: true}

	return result, nil
}

// RegisterOwned is like Register but stamps owner_user_id — see
// repository.CouchInstances.RegisterOwned.
func (r *Repo) RegisterOwned(
	ctx context.Context, ownerUserID uuid.UUID, url, username string, passwordPlain []byte,
) (uuid.UUID, error) {
	passwordEnc, err := r.encryptor.Encrypt(passwordPlain)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "encrypt password")
	}

	params := artel_q.RegisterOwnedCouchInstanceParams{
		Url:         url,
		Username:    username,
		PasswordEnc: passwordEnc,
		OwnerUserID: uuid.NullUUID{UUID: ownerUserID, Valid: true},
	}

	id, err := r.q.RegisterOwnedCouchInstance(ctx, params)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

// DeleteOwnedIfUnreferenced removes ownerUserID's owned instance row if no vault references it
// — see repository.CouchInstances.DeleteOwnedIfUnreferenced.
func (r *Repo) DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error {
	err := r.q.DeleteOwnedCouchInstanceIfUnreferenced(ctx, uuid.NullUUID{UUID: ownerUserID, Valid: true})
	if err != nil {
		return rerrors.Wrap(err, "delete owned couch instance if unreferenced")
	}

	return nil
}

func (r *Repo) Exists(ctx context.Context) (bool, error) {
	exists, err := r.q.CouchInstanceExists(ctx)
	if err != nil {
		return false, rerrors.Wrap(err, "check couch instance existence")
	}

	return exists, nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.CouchInstances {
	return New(tx, r.encryptor)
}
