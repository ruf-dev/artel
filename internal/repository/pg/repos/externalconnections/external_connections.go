package externalconnections

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
)

type Repo struct {
	q             *artel_q.Queries
	encryptionKey []byte
}

func New(q *artel_q.Queries, encryptionKey []byte) *Repo {
	return &Repo{
		q:             q,
		encryptionKey: encryptionKey,
	}
}

func (r *Repo) Upsert(ctx context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error) {
	credEnc, err := cryptoutil.Encrypt(r.encryptionKey, conn.CredentialsJSON)
	if err != nil {
		return domain.ExternalConnection{}, rerrors.Wrap(err, "error encrypting credentials")
	}

	params := artel_q.UpsertExternalConnectionParams{
		UserID:         conn.UserUuid,
		Provider:       conn.Provider,
		ProviderType:   conn.ProviderType,
		CredentialsEnc: credEnc,
		Metadata:       toNullRawMessage(conn.Metadata),
	}

	row, err := r.q.UpsertExternalConnection(ctx, params)
	if err != nil {
		return domain.ExternalConnection{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error upserting external connection")
	}

	credJSON, err := cryptoutil.Decrypt(r.encryptionKey, row.CredentialsEnc)
	if err != nil {
		return domain.ExternalConnection{}, rerrors.Wrap(err, "error decrypting credentials")
	}

	return toDomain(row, credJSON), nil
}

func (r *Repo) GetByUserAndProvider(ctx context.Context, userUuid uuid.UUID, provider string) (sql.Null[domain.ExternalConnection], error) {
	params := artel_q.GetExternalConnectionByUserAndProviderParams{
		UserID:   userUuid,
		Provider: provider,
	}

	row, err := r.q.GetExternalConnectionByUserAndProvider(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.ExternalConnection]{}, nil
		}
		return sql.Null[domain.ExternalConnection]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting external connection")
	}

	credJSON, err := cryptoutil.Decrypt(r.encryptionKey, row.CredentialsEnc)
	if err != nil {
		return sql.Null[domain.ExternalConnection]{}, rerrors.Wrap(err, "error decrypting credentials")
	}

	result := sql.Null[domain.ExternalConnection]{
		V:     toDomain(row, credJSON),
		Valid: true,
	}
	return result, nil
}

func (r *Repo) ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.ExternalConnection, error) {
	rows, err := r.q.ListExternalConnectionsByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing external connections")
	}

	conns := make([]domain.ExternalConnection, len(rows))
	for i, row := range rows {
		credJSON, err := cryptoutil.Decrypt(r.encryptionKey, row.CredentialsEnc)
		if err != nil {
			return nil, rerrors.Wrap(err, "error decrypting credentials")
		}
		conns[i] = toDomain(row, credJSON)
	}
	return conns, nil
}

func (r *Repo) Delete(ctx context.Context, userUuid uuid.UUID, provider string) error {
	params := artel_q.DeleteExternalConnectionParams{
		UserID:   userUuid,
		Provider: provider,
	}
	err := r.q.DeleteExternalConnection(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting external connection")
	}
	return nil
}

func toDomain(row artel_q.ExternalConnection, credJSON []byte) domain.ExternalConnection {
	return domain.ExternalConnection{
		Uuid:            row.ID,
		UserUuid:        row.UserID,
		Provider:        row.Provider,
		ProviderType:    row.ProviderType,
		CredentialsJSON: json.RawMessage(credJSON),
		Metadata:        fromNullRawMessage(row.Metadata),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func toNullRawMessage(m json.RawMessage) pqtype.NullRawMessage {
	return pqtype.NullRawMessage{
		RawMessage: json.RawMessage(m),
		Valid:      m != nil,
	}
}

func fromNullRawMessage(m pqtype.NullRawMessage) json.RawMessage {
	if !m.Valid {
		return nil
	}
	return json.RawMessage(m.RawMessage)
}
