package externalconnections

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/service"
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

// fakeMomService is a minimal in-memory service.MomService used to verify what
// ExecuteToolWithSecrets is called with, without a real MoM/http executor. Only
// ExecuteToolWithSecrets is exercised by the telegram connection tests, so the remaining
// interface methods are unused stubs.
type fakeMomService struct {
	result string
	err    error

	lastMcpName  string
	lastToolName string
	lastSecrets  map[string]interface{}
}

func (f *fakeMomService) ListToolsForKey(_ context.Context, _ uuid.UUID) ([]domain.McpToolDef, error) {
	return nil, nil
}

func (f *fakeMomService) ExecuteToolForKey(
	_ context.Context, _ uuid.UUID, _ string, _ map[string]interface{},
) (string, error) {
	return "", nil
}

func (f *fakeMomService) ExecuteToolForUserConnection(
	_ context.Context, _ uuid.UUID, _ string, _ string, _ map[string]interface{},
) (string, error) {
	return "", nil
}

func (f *fakeMomService) ExecuteToolForConnection(
	_ context.Context, _ uuid.UUID, _ string, _ string, _ map[string]interface{},
) (string, error) {
	return "", nil
}

func (f *fakeMomService) ExecuteToolWithSecrets(
	_ context.Context, mcpName, toolName string, secrets, _ map[string]interface{},
) (string, error) {
	f.lastMcpName = mcpName
	f.lastToolName = toolName
	f.lastSecrets = secrets

	return f.result, f.err
}

var _ service.MomService = (*fakeMomService)(nil)
