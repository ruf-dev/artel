package externalconnections

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/ruf-dev/artel/internal/clients/postgres"
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

// fakeCouchInstancesRepo is a minimal in-memory repository.CouchInstances used to verify
// AddCouchDBConnection's owned-instance sync (RegisterOwned vs Update, see syncOwnedCouchInstance)
// without a real database.
type fakeCouchInstancesRepo struct {
	owned map[uuid.UUID]domain.CouchInstance

	registerOwnedCalls int
	updateCalls        int
	lastUpdateID       uuid.UUID
	lastUpdateURL      string
}

func newFakeCouchInstancesRepo() *fakeCouchInstancesRepo {
	return &fakeCouchInstancesRepo{owned: map[uuid.UUID]domain.CouchInstance{}}
}

func (f *fakeCouchInstancesRepo) Register(context.Context, string, string, []byte) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (f *fakeCouchInstancesRepo) Get(context.Context, uuid.UUID) (domain.CouchInstance, error) {
	return domain.CouchInstance{}, nil
}

func (f *fakeCouchInstancesRepo) RandomPick(context.Context) (domain.CouchInstanceWithAccount, error) {
	return domain.CouchInstanceWithAccount{}, nil
}

func (f *fakeCouchInstancesRepo) PickForUser(context.Context, uuid.UUID) (domain.CouchInstanceWithAccount, error) {
	return domain.CouchInstanceWithAccount{}, nil
}

func (f *fakeCouchInstancesRepo) GetOwned(_ context.Context, userID uuid.UUID) (sql.Null[domain.CouchInstance], error) {
	instance, ok := f.owned[userID]
	if !ok {
		return sql.Null[domain.CouchInstance]{}, nil
	}

	return sql.Null[domain.CouchInstance]{V: instance, Valid: true}, nil
}

func (f *fakeCouchInstancesRepo) List(context.Context) ([]domain.CouchInstance, error) {
	return nil, nil
}

func (f *fakeCouchInstancesRepo) Update(_ context.Context, id uuid.UUID, url, _ string, _ []byte) error {
	f.updateCalls++
	f.lastUpdateID = id
	f.lastUpdateURL = url

	return nil
}

func (f *fakeCouchInstancesRepo) RegisterOwned(
	_ context.Context, ownerUserID uuid.UUID, url, username string, _ []byte,
) (uuid.UUID, error) {
	f.registerOwnedCalls++

	id := uuid.New()
	f.owned[ownerUserID] = domain.CouchInstance{Uuid: id, Url: url, Username: username}

	return id, nil
}

func (f *fakeCouchInstancesRepo) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeCouchInstancesRepo) DeleteOwnedIfUnreferenced(_ context.Context, userID uuid.UUID) error {
	delete(f.owned, userID)

	return nil
}

func (f *fakeCouchInstancesRepo) Exists(context.Context) (bool, error) { return false, nil }

func (f *fakeCouchInstancesRepo) WithTx(postgres.DB) repository.CouchInstances { return f }

var _ repository.CouchInstances = (*fakeCouchInstancesRepo)(nil)

// fakeS3InstancesRepo is a minimal in-memory repository.S3Instances used to verify
// AddS3Connection's owned-instance sync (RegisterOwned vs Update, see syncOwnedS3Instance)
// without a real database.
type fakeS3InstancesRepo struct {
	owned map[uuid.UUID]domain.S3Instance

	registerOwnedCalls int
	updateCalls        int
	lastUpdateID       uuid.UUID
	lastUpdateEndpoint string
}

func newFakeS3InstancesRepo() *fakeS3InstancesRepo {
	return &fakeS3InstancesRepo{owned: map[uuid.UUID]domain.S3Instance{}}
}

func (f *fakeS3InstancesRepo) Register(
	context.Context, string, string, bool, bool, string, []byte,
) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (f *fakeS3InstancesRepo) Get(context.Context, uuid.UUID) (domain.S3Instance, error) {
	return domain.S3Instance{}, nil
}

func (f *fakeS3InstancesRepo) PickForUser(context.Context, uuid.UUID) (domain.S3Instance, error) {
	return domain.S3Instance{}, nil
}

func (f *fakeS3InstancesRepo) GetOwned(_ context.Context, userID uuid.UUID) (sql.Null[domain.S3Instance], error) {
	instance, ok := f.owned[userID]
	if !ok {
		return sql.Null[domain.S3Instance]{}, nil
	}

	return sql.Null[domain.S3Instance]{V: instance, Valid: true}, nil
}

func (f *fakeS3InstancesRepo) List(context.Context) ([]domain.S3Instance, error) { return nil, nil }

func (f *fakeS3InstancesRepo) Update(
	_ context.Context, id uuid.UUID, endpoint, _ string, _, _ bool, _ string, _ []byte,
) error {
	f.updateCalls++
	f.lastUpdateID = id
	f.lastUpdateEndpoint = endpoint

	return nil
}

func (f *fakeS3InstancesRepo) RegisterOwned(
	_ context.Context, ownerUserID uuid.UUID, endpoint, region string, useSSL, pathStyle bool, accessKey string, _ []byte,
) (uuid.UUID, error) {
	f.registerOwnedCalls++

	id := uuid.New()
	f.owned[ownerUserID] = domain.S3Instance{
		Uuid:      id,
		Endpoint:  endpoint,
		Region:    region,
		UseSSL:    useSSL,
		PathStyle: pathStyle,
		AccessKey: accessKey,
	}

	return id, nil
}

func (f *fakeS3InstancesRepo) Delete(context.Context, uuid.UUID) error { return nil }

func (f *fakeS3InstancesRepo) DeleteOwnedIfUnreferenced(_ context.Context, userID uuid.UUID) error {
	delete(f.owned, userID)

	return nil
}

func (f *fakeS3InstancesRepo) Exists(context.Context) (bool, error) { return false, nil }

func (f *fakeS3InstancesRepo) WithTx(postgres.DB) repository.S3Instances { return f }

var _ repository.S3Instances = (*fakeS3InstancesRepo)(nil)
