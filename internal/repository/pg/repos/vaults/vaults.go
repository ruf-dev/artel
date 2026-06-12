package vaults

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
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

func (r *Repo) Upsert(ctx context.Context, userID, couchInstanceID uuid.UUID, name, couchDBName, status, passphrase string) (domain.Vault, error) {
	passphraseEnc, err := cryptoutil.Encrypt(r.encryptionKey, []byte(passphrase))
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error encrypting livesync passphrase")
	}

	params := artel_q.CreateVaultParams{
		UserID:                userID,
		Name:                  name,
		CouchDbName:           couchDBName,
		CouchInstanceID:       uuid.NullUUID{UUID: couchInstanceID, Valid: true},
		Status:                status,
		LivesyncPassphraseEnc: passphraseEnc,
	}
	row, err := r.q.CreateVault(ctx, params)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error creating vault")
	}

	v, err := rowToVault(row.ID, row.UserID, row.Name, row.CouchDbName, row.CouchInstanceID, row.Status, row.LivesyncPassphraseEnc, row.CreatedAt, r.encryptionKey)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error mapping vault row")
	}
	return v, nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (domain.Vault, error) {
	row, err := r.q.GetVaultByID(ctx, id)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting vault by id")
	}

	v, err := rowToVault(row.ID, row.UserID, row.Name, row.CouchDbName, row.CouchInstanceID, row.Status, row.LivesyncPassphraseEnc, row.CreatedAt, r.encryptionKey)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error mapping vault row")
	}
	return v, nil
}

func (r *Repo) GetByNameAndUser(ctx context.Context, userID uuid.UUID, name string) (domain.Vault, error) {
	params := artel_q.GetVaultByNameAndUserParams{
		UserID: userID,
		Name:   name,
	}
	row, err := r.q.GetVaultByNameAndUser(ctx, params)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error getting vault by name and user")
	}

	v, err := rowToVault(row.ID, row.UserID, row.Name, row.CouchDbName, row.CouchInstanceID, row.Status, row.LivesyncPassphraseEnc, row.CreatedAt, r.encryptionKey)
	if err != nil {
		return domain.Vault{}, rerrors.Wrap(err, "error mapping vault row")
	}
	return v, nil
}

func (r *Repo) UpdateStatus(ctx context.Context, vaultID uuid.UUID, status string) error {
	params := artel_q.UpdateVaultStatusParams{
		ID:     vaultID,
		Status: status,
	}
	err := r.q.UpdateVaultStatus(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error updating vault status")
	}
	return nil
}

func (r *Repo) SetLiveSyncPassphrase(ctx context.Context, vaultID uuid.UUID, passphrase string) error {
	passphraseEnc, err := cryptoutil.Encrypt(r.encryptionKey, []byte(passphrase))
	if err != nil {
		return rerrors.Wrap(err, "error encrypting livesync passphrase")
	}

	params := artel_q.SetVaultLiveSyncPassphraseParams{
		ID:                    vaultID,
		LivesyncPassphraseEnc: passphraseEnc,
	}
	err = r.q.SetVaultLiveSyncPassphrase(ctx, params)
	if err != nil {
		return rerrors.Wrap(err, "error setting livesync passphrase")
	}
	return nil
}

func (r *Repo) ListByMembership(ctx context.Context, userID uuid.UUID) ([]domain.Vault, error) {
	rows, err := r.q.ListVaultsByMembership(ctx, userID)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing vaults by membership")
	}

	vaultList := make([]domain.Vault, 0, len(rows))
	for _, row := range rows {
		v, err := rowToVault(row.ID, row.UserID, row.Name, row.CouchDbName, row.CouchInstanceID, row.Status, row.LivesyncPassphraseEnc, row.CreatedAt, r.encryptionKey)
		if err != nil {
			return nil, rerrors.Wrap(err, "error mapping vault row")
		}
		vaultList = append(vaultList, v)
	}
	return vaultList, nil
}

func (r *Repo) Delete(ctx context.Context, vaultID uuid.UUID) error {
	err := r.q.DeleteVault(ctx, vaultID)
	if err != nil {
		return rerrors.Wrap(err, "error deleting vault")
	}

	return nil
}

func (r *Repo) WithTx(tx sqldb.DB) repository.Vaults {
	return New(tx, r.encryptionKey)
}

func rowToVault(id, userID uuid.UUID, name, couchDbName string, couchInstanceID uuid.NullUUID, status string, passphraseEnc []byte, createdAt time.Time, encryptionKey []byte) (domain.Vault, error) {
	v := domain.Vault{
		Uuid:        id,
		UserUuid:    userID,
		Name:        name,
		CouchDBName: couchDbName,
		Status:      status,
		CreatedAt:   createdAt,
	}
	if couchInstanceID.Valid {
		v.CouchInstanceUuid = couchInstanceID.UUID
	}
	if len(passphraseEnc) > 0 {
		decrypted, err := cryptoutil.Decrypt(encryptionKey, passphraseEnc)
		if err != nil {
			return domain.Vault{}, rerrors.Wrap(err, "error decrypting livesync passphrase")
		}
		v.LiveSyncPassphrase = string(decrypted)
	}
	return v, nil
}
