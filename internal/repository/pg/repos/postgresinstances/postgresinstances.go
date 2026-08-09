package postgresinstances

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

func (r *Repo) Register(
	ctx context.Context, host string, port int, adminDatabase, username string, passwordPlain []byte, sslMode string,
) (uuid.UUID, error) {
	passwordEnc, err := r.encryptor.Encrypt(passwordPlain)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "error encrypting password")
	}

	params := artel_q.RegisterPostgresInstanceParams{
		Host:          host,
		Port:          int32(port),
		AdminDatabase: adminDatabase,
		Username:      username,
		PasswordEnc:   passwordEnc,
		SslMode:       sslMode,
	}

	id, err := r.q.RegisterPostgresInstance(ctx, params)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (domain.PostgresInstance, error) {
	row, err := r.q.GetPostgresInstance(ctx, id)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "error getting postgres instance")
	}

	instance, err := rowToDomain(row, r.encryptor)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "error decoding postgres instance")
	}

	return instance, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.PostgresInstance, error) {
	rows, err := r.q.ListPostgresInstances(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing postgres instances")
	}

	instances := make([]domain.PostgresInstance, len(rows))
	for i, row := range rows {
		instance, err := rowToDomain(row, r.encryptor)
		if err != nil {
			return nil, rerrors.Wrap(err, "error decoding postgres instance")
		}

		instances[i] = instance
	}

	return instances, nil
}

func (r *Repo) Update(
	ctx context.Context, id uuid.UUID, host string, port int, adminDatabase, username string,
	passwordPlain []byte, sslMode string,
) error {
	passwordEnc, err := r.encryptor.Encrypt(passwordPlain)
	if err != nil {
		return rerrors.Wrap(err, "error encrypting password")
	}

	params := artel_q.UpdatePostgresInstanceParams{
		ID:            id,
		Host:          host,
		Port:          int32(port),
		AdminDatabase: adminDatabase,
		Username:      username,
		PasswordEnc:   passwordEnc,
		SslMode:       sslMode,
	}

	err = r.q.UpdatePostgresInstance(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error updating postgres instance")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeletePostgresInstance(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "error deleting postgres instance")
	}

	return nil
}

func (r *Repo) RandomPick(ctx context.Context) (domain.PostgresInstance, error) {
	row, err := r.q.RandomPickPostgresInstance(ctx)
	if err != nil {
		return domain.PostgresInstance{}, pg_err.UnwrapPgErr(err)
	}

	instance, err := rowToDomain(row, r.encryptor)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "error decoding postgres instance")
	}

	return instance, nil
}

// PickForUser tries userID's owned instance first, falling back to RandomPick's shared admin
// pool when they have none — see repository.PostgresInstances.PickForUser.
func (r *Repo) PickForUser(ctx context.Context, userID uuid.UUID) (domain.PostgresInstance, error) {
	owned, err := r.GetOwned(ctx, userID)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "error getting owned postgres instance")
	}

	if owned.Valid {
		return owned.V, nil
	}

	instance, err := r.RandomPick(ctx)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "error random picking postgres instance")
	}

	return instance, nil
}

// GetOwned returns the instance owned by userID, if any — see repository.PostgresInstances.GetOwned.
func (r *Repo) GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.PostgresInstance], error) {
	ownerUserID := uuid.NullUUID{UUID: userID, Valid: true}

	row, err := r.q.PickOwnedPostgresInstance(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.PostgresInstance]{}, nil
		}

		return sql.Null[domain.PostgresInstance]{}, rerrors.Wrap(err, "error picking owned postgres instance")
	}

	instance, err := rowToDomain(row, r.encryptor)
	if err != nil {
		return sql.Null[domain.PostgresInstance]{}, rerrors.Wrap(err, "error decoding postgres instance")
	}

	result := sql.Null[domain.PostgresInstance]{V: instance, Valid: true}

	return result, nil
}

// RegisterOwned is like Register but stamps owner_user_id — see
// repository.PostgresInstances.RegisterOwned.
func (r *Repo) RegisterOwned(
	ctx context.Context, ownerUserID uuid.UUID, host string, port int, adminDatabase, username string,
	passwordPlain []byte, sslMode string,
) (uuid.UUID, error) {
	passwordEnc, err := r.encryptor.Encrypt(passwordPlain)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "error encrypting password")
	}

	params := artel_q.RegisterOwnedPostgresInstanceParams{
		Host:          host,
		Port:          int32(port),
		AdminDatabase: adminDatabase,
		Username:      username,
		PasswordEnc:   passwordEnc,
		SslMode:       sslMode,
		OwnerUserID:   uuid.NullUUID{UUID: ownerUserID, Valid: true},
	}

	id, err := r.q.RegisterOwnedPostgresInstance(ctx, params)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

// DeleteOwnedIfUnreferenced removes ownerUserID's owned instance row if no vault postgres
// database references it — see repository.PostgresInstances.DeleteOwnedIfUnreferenced.
func (r *Repo) DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error {
	err := r.q.DeleteOwnedPostgresInstanceIfUnreferenced(ctx, uuid.NullUUID{UUID: ownerUserID, Valid: true})
	if err != nil {
		return rerrors.Wrap(err, "error deleting owned postgres instance if unreferenced")
	}

	return nil
}

func (r *Repo) Exists(ctx context.Context) (bool, error) {
	exists, err := r.q.PostgresInstanceExists(ctx)
	if err != nil {
		return false, rerrors.Wrap(err, "error checking postgres instance existence")
	}

	return exists, nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.PostgresInstances {
	return New(tx, r.encryptor)
}

// rowToDomain decrypts the row's password_enc and maps it onto domain.PostgresInstance,
// populating OwnerUserUuid — unlike CouchInstance/S3Instance, PostgresInstance surfaces ownership
// directly on the domain struct (see internal/domain/postgres_instance.go doc comment).
func rowToDomain(row artel_q.PostgresInstance, encryptor cryptoutil.Encryptor) (domain.PostgresInstance, error) {
	decrypted, err := encryptor.Decrypt(row.PasswordEnc)
	if err != nil {
		return domain.PostgresInstance{}, rerrors.Wrap(err, "error decrypting password")
	}

	var ownerUserUuid *uuid.UUID
	if row.OwnerUserID.Valid {
		ownerUserUuid = &row.OwnerUserID.UUID
	}

	instance := domain.PostgresInstance{
		Uuid:          row.ID,
		Host:          row.Host,
		Port:          int(row.Port),
		AdminDatabase: row.AdminDatabase,
		Username:      row.Username,
		Password:      string(decrypted),
		SSLMode:       row.SslMode,
		OwnerUserUuid: ownerUserUuid,
		CreatedAt:     row.CreatedAt,
	}

	return instance, nil
}
