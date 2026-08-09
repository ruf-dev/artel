package vaultpostgresdatabases

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

func (r *Repo) Create(
	ctx context.Context, vaultID, postgresInstanceID uuid.UUID, dbName, roleUsername string, rolePasswordPlain []byte,
) (domain.VaultPostgresDatabase, error) {
	rolePasswordEnc, err := r.encryptor.Encrypt(rolePasswordPlain)
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error encrypting role password")
	}

	params := artel_q.CreateVaultPostgresDatabaseParams{
		VaultID:            vaultID,
		PostgresInstanceID: postgresInstanceID,
		DatabaseName:       dbName,
		RoleUsername:       roleUsername,
		RolePasswordEnc:    rolePasswordEnc,
	}

	row, err := r.q.CreateVaultPostgresDatabase(ctx, params)
	if err != nil {
		return domain.VaultPostgresDatabase{}, pg_err.UnwrapPgErr(err)
	}

	database, err := rowToDomain(row, r.encryptor)
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error decoding vault postgres database")
	}

	return database, nil
}

// GetByVaultID returns Valid: false when vaultID has no Postgres database provisioned — see
// repository.VaultPostgresDatabases.GetByVaultID.
func (r *Repo) GetByVaultID(ctx context.Context, vaultID uuid.UUID) (sql.Null[domain.VaultPostgresDatabase], error) {
	row, err := r.q.GetVaultPostgresDatabaseByVaultID(ctx, vaultID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.VaultPostgresDatabase]{}, nil
		}

		return sql.Null[domain.VaultPostgresDatabase]{}, rerrors.Wrap(err, "error getting vault postgres database")
	}

	database, err := rowToDomain(row, r.encryptor)
	if err != nil {
		return sql.Null[domain.VaultPostgresDatabase]{}, rerrors.Wrap(err, "error decoding vault postgres database")
	}

	result := sql.Null[domain.VaultPostgresDatabase]{V: database, Valid: true}

	return result, nil
}

func (r *Repo) MarkReady(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.MarkVaultPostgresDatabaseReady(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error marking vault postgres database ready")
	}

	return nil
}

func (r *Repo) MarkError(ctx context.Context, vaultID uuid.UUID, message string) error {
	params := artel_q.MarkVaultPostgresDatabaseErrorParams{
		VaultID:      vaultID,
		ErrorMessage: message,
	}

	err := r.q.MarkVaultPostgresDatabaseError(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error marking vault postgres database error")
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.DeleteVaultPostgresDatabase(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting vault postgres database")
	}

	return nil
}

func (r *Repo) WithTx(tx postgres.DB) repository.VaultPostgresDatabases {
	return New(tx, r.encryptor)
}

// rowToDomain decrypts the row's role_password_enc and maps it onto domain.VaultPostgresDatabase.
func rowToDomain(row artel_q.VaultPostgresDatabase, encryptor cryptoutil.Encryptor) (domain.VaultPostgresDatabase, error) {
	decrypted, err := encryptor.Decrypt(row.RolePasswordEnc)
	if err != nil {
		return domain.VaultPostgresDatabase{}, rerrors.Wrap(err, "error decrypting role password")
	}

	database := domain.VaultPostgresDatabase{
		VaultUuid:            row.VaultID,
		PostgresInstanceUuid: row.PostgresInstanceID,
		DatabaseName:         row.DatabaseName,
		RoleUsername:         row.RoleUsername,
		RolePassword:         string(decrypted),
		Status:               domain.VaultPostgresStatus(row.Status),
		ErrorMessage:         row.ErrorMessage,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}

	return database, nil
}
