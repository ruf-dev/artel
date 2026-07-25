package externalconnections

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
)

// fakeExternalConnsRepo is a minimal in-memory repository.ExternalConnectionRepo used to verify
// what AddGenericConnection persists without a real database.
type fakeExternalConnsRepo struct {
	conns        map[uuid.UUID]domain.ExternalConnection
	lastUpserted domain.ExternalConnection
}

func newFakeExternalConnsRepo() *fakeExternalConnsRepo {
	repo := &fakeExternalConnsRepo{conns: map[uuid.UUID]domain.ExternalConnection{}}

	return repo
}

func (f *fakeExternalConnsRepo) Upsert(
	_ context.Context,
	conn domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	if conn.Uuid == uuid.Nil {
		conn.Uuid = uuid.New()
	}

	f.conns[conn.Uuid] = conn
	f.lastUpserted = conn

	return conn, nil
}

func (f *fakeExternalConnsRepo) Insert(
	_ context.Context,
	conn domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	f.conns[conn.Uuid] = conn

	return conn, nil
}

func (f *fakeExternalConnsRepo) GetByID(_ context.Context, id uuid.UUID) (domain.ExternalConnection, error) {
	conn, ok := f.conns[id]
	if !ok {
		return domain.ExternalConnection{}, sql.ErrNoRows
	}

	return conn, nil
}

func (f *fakeExternalConnsRepo) GetByUserAndProvider(
	_ context.Context,
	userUuid uuid.UUID,
	provider string,
) (sql.Null[domain.ExternalConnection], error) {
	for _, conn := range f.conns {
		if conn.UserUuid == userUuid && conn.Provider == provider {
			result := sql.Null[domain.ExternalConnection]{V: conn, Valid: true}

			return result, nil
		}
	}

	return sql.Null[domain.ExternalConnection]{}, nil
}

func (f *fakeExternalConnsRepo) ListByUser(_ context.Context, _ uuid.UUID) ([]domain.ExternalConnection, error) {
	return nil, nil
}

func (f *fakeExternalConnsRepo) Delete(_ context.Context, _ uuid.UUID, _ string) error {
	return nil
}

func (f *fakeExternalConnsRepo) DeleteByID(_ context.Context, _ uuid.UUID, id uuid.UUID) error {
	delete(f.conns, id)

	return nil
}

var _ repository.ExternalConnectionRepo = (*fakeExternalConnsRepo)(nil)
