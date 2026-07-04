package s3instances

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
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

func (r *Repo) Register(
	ctx context.Context,
	endpoint, region string,
	useSSL, pathStyle bool,
	accessKey string,
	secretKeyPlain []byte,
) (uuid.UUID, error) {
	secretKeyEnc, err := cryptoutil.Encrypt(r.encryptionKey, secretKeyPlain)
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

	decrypted, err := cryptoutil.Decrypt(r.encryptionKey, row.SecretKeyEnc)
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
	secretKeyEnc, err := cryptoutil.Encrypt(r.encryptionKey, secretKeyPlain)
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

func (r *Repo) WithTx(tx sqldb.DB) repository.S3Instances {
	return New(tx, r.encryptionKey)
}
