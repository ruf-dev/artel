package s3instances

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
	ctx context.Context,
	endpoint, region string,
	useSSL, pathStyle bool,
	accessKey string,
	secretKeyPlain []byte,
) (uuid.UUID, error) {
	secretKeyEnc, err := r.encryptor.Encrypt(secretKeyPlain)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "error encrypting secret key")
	}

	params := artel_q.RegisterS3InstanceParams{
		Endpoint:     endpoint,
		Region:       region,
		AccessKey:    accessKey,
		SecretKeyEnc: secretKeyEnc,
		UseSsl:       useSSL,
		PathStyle:    pathStyle,
	}

	id, err := r.q.RegisterS3Instance(ctx, params)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id uuid.UUID) (domain.S3Instance, error) {
	row, err := r.q.GetS3InstanceWithCreds(ctx, id)
	if err != nil {
		return domain.S3Instance{}, rerrors.Wrap(err, "error getting s3 instance")
	}

	decrypted, err := r.encryptor.Decrypt(row.SecretKeyEnc)
	if err != nil {
		return domain.S3Instance{}, rerrors.Wrap(err, "error decrypting secret key")
	}

	instance := domain.S3Instance{
		Uuid:      row.ID,
		Endpoint:  row.Endpoint,
		Region:    row.Region,
		AccessKey: row.AccessKey,
		SecretKey: string(decrypted),
		UseSSL:    row.UseSsl,
		PathStyle: row.PathStyle,
		CreatedAt: row.CreatedAt,
	}

	return instance, nil
}

func (r *Repo) List(ctx context.Context) ([]domain.S3Instance, error) {
	rows, err := r.q.ListS3Instances(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing s3 instances")
	}

	instances := make([]domain.S3Instance, len(rows))
	for i, row := range rows {
		instances[i] = domain.S3Instance{
			Uuid:      row.ID,
			Endpoint:  row.Endpoint,
			Region:    row.Region,
			AccessKey: row.AccessKey,
			UseSSL:    row.UseSsl,
			PathStyle: row.PathStyle,
			CreatedAt: row.CreatedAt,
		}
	}

	return instances, nil
}

func (r *Repo) Update(
	ctx context.Context,
	id uuid.UUID,
	endpoint, region string,
	useSSL, pathStyle bool,
	accessKey string,
	secretKeyPlain []byte,
) error {
	secretKeyEnc, err := r.encryptor.Encrypt(secretKeyPlain)
	if err != nil {
		return rerrors.Wrap(err, "error encrypting secret key")
	}

	params := artel_q.UpdateS3InstanceParams{
		ID:           id,
		Endpoint:     endpoint,
		Region:       region,
		AccessKey:    accessKey,
		SecretKeyEnc: secretKeyEnc,
		UseSsl:       useSSL,
		PathStyle:    pathStyle,
	}

	err = r.q.UpdateS3Instance(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error updating s3 instance")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteS3Instance(ctx, id)
	if err != nil {
		return rerrors.Wrap(err, "error deleting s3 instance")
	}

	return nil
}

// PickForUser tries userID's owned instance first, falling back to a random pick from the
// shared admin pool when they have none — see repository.S3Instances.PickForUser.
func (r *Repo) PickForUser(ctx context.Context, userID uuid.UUID) (domain.S3Instance, error) {
	owned, err := r.GetOwned(ctx, userID)
	if err != nil {
		return domain.S3Instance{}, rerrors.Wrap(err, "get owned s3 instance")
	}

	if owned.Valid {
		return owned.V, nil
	}

	row, err := r.q.RandomPickS3Instance(ctx)
	if err != nil {
		return domain.S3Instance{}, pg_err.UnwrapPgErr(err)
	}

	decrypted, err := r.encryptor.Decrypt(row.SecretKeyEnc)
	if err != nil {
		return domain.S3Instance{}, rerrors.Wrap(err, "error decrypting secret key")
	}

	instance := domain.S3Instance{
		Uuid:      row.ID,
		Endpoint:  row.Endpoint,
		Region:    row.Region,
		AccessKey: row.AccessKey,
		SecretKey: string(decrypted),
		UseSSL:    row.UseSsl,
		PathStyle: row.PathStyle,
		CreatedAt: row.CreatedAt,
	}

	return instance, nil
}

// GetOwned returns the instance owned by userID, if any — see repository.S3Instances.GetOwned.
func (r *Repo) GetOwned(ctx context.Context, userID uuid.UUID) (sql.Null[domain.S3Instance], error) {
	ownerUserID := uuid.NullUUID{UUID: userID, Valid: true}

	row, err := r.q.PickOwnedS3Instance(ctx, ownerUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.S3Instance]{}, nil
		}

		return sql.Null[domain.S3Instance]{}, rerrors.Wrap(err, "error picking owned s3 instance")
	}

	decrypted, err := r.encryptor.Decrypt(row.SecretKeyEnc)
	if err != nil {
		return sql.Null[domain.S3Instance]{}, rerrors.Wrap(err, "error decrypting secret key")
	}

	instance := domain.S3Instance{
		Uuid:      row.ID,
		Endpoint:  row.Endpoint,
		Region:    row.Region,
		AccessKey: row.AccessKey,
		SecretKey: string(decrypted),
		UseSSL:    row.UseSsl,
		PathStyle: row.PathStyle,
		CreatedAt: row.CreatedAt,
	}

	result := sql.Null[domain.S3Instance]{V: instance, Valid: true}

	return result, nil
}

// RegisterOwned is like Register but stamps owner_user_id — see repository.S3Instances.RegisterOwned.
func (r *Repo) RegisterOwned(
	ctx context.Context,
	ownerUserID uuid.UUID,
	endpoint, region string,
	useSSL, pathStyle bool,
	accessKey string,
	secretKeyPlain []byte,
) (uuid.UUID, error) {
	secretKeyEnc, err := r.encryptor.Encrypt(secretKeyPlain)
	if err != nil {
		return uuid.UUID{}, rerrors.Wrap(err, "error encrypting secret key")
	}

	params := artel_q.RegisterOwnedS3InstanceParams{
		Endpoint:     endpoint,
		Region:       region,
		AccessKey:    accessKey,
		SecretKeyEnc: secretKeyEnc,
		UseSsl:       useSSL,
		PathStyle:    pathStyle,
		OwnerUserID:  uuid.NullUUID{UUID: ownerUserID, Valid: true},
	}

	id, err := r.q.RegisterOwnedS3Instance(ctx, params)
	if err != nil {
		return uuid.UUID{}, pg_err.UnwrapPgErr(err)
	}

	return id, nil
}

// DeleteOwnedIfUnreferenced removes ownerUserID's owned instance row if no vault references it
// — see repository.S3Instances.DeleteOwnedIfUnreferenced.
func (r *Repo) DeleteOwnedIfUnreferenced(ctx context.Context, ownerUserID uuid.UUID) error {
	err := r.q.DeleteOwnedS3InstanceIfUnreferenced(ctx, uuid.NullUUID{UUID: ownerUserID, Valid: true})
	if err != nil {
		return rerrors.Wrap(err, "error deleting owned s3 instance if unreferenced")
	}

	return nil
}

func (r *Repo) Exists(ctx context.Context) (bool, error) {
	exists, err := r.q.S3InstanceExists(ctx)
	if err != nil {
		return false, rerrors.Wrap(err, "error checking s3 instance existence")
	}

	return exists, nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.S3Instances {
	return New(tx, r.encryptor)
}
