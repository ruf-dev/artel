package externalconnections

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func newCouchDBConnectionTestService(
	connections *fakeExternalConnsRepo, couchInstances *fakeCouchInstancesRepo,
) *Service {
	return New(connections, nil, nil, nil, nil, nil, couchInstances, nil)
}

// couchDBFakeServer serves the minimal read-only surface Client.GetSetupStatus touches
// (GET /_cluster_setup, HEAD/GET /_users, /_replicator) so testCouchDBConnection can succeed
// without a real CouchDB instance.
func couchDBFakeServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/_cluster_setup":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"state":"cluster_finished"}`))
		case "/_users", "/_replicator":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestAddCouchDBConnection_Unauthenticated_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	couchInstances := newFakeCouchInstancesRepo()
	svc := newCouchDBConnectionTestService(connections, couchInstances)

	_, err := svc.AddCouchDBConnection(context.Background(), "http://localhost:5984", "admin", "password")

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.Unauthenticated)
}

// TestAddCouchDBConnection_ConnectivityCheckFails_Errors points at an address nothing listens on
// (loopback, reserved port) so the connection is refused immediately — deterministic, no real
// network required — and verifies neither the connection row nor an owned instance row get
// persisted.
func TestAddCouchDBConnection_ConnectivityCheckFails_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	couchInstances := newFakeCouchInstancesRepo()
	svc := newCouchDBConnectionTestService(connections, couchInstances)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: uuid.New()})

	_, err := svc.AddCouchDBConnection(ctx, "http://127.0.0.1:1", "admin", "password")

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.CouchDBConnectionValidationFailed)
	assert.Empty(t, connections.conns)
	assert.Zero(t, couchInstances.registerOwnedCalls)
}

func TestAddCouchDBConnection_Success_RegistersOwnedInstance(t *testing.T) {
	server := couchDBFakeServer(t)
	defer server.Close()

	connections := newFakeExternalConnsRepo()
	couchInstances := newFakeCouchInstancesRepo()
	svc := newCouchDBConnectionTestService(connections, couchInstances)

	userUuid := uuid.New()
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userUuid})

	meta, err := svc.AddCouchDBConnection(ctx, server.URL, "admin", "password")
	require.NoError(t, err)
	assert.Equal(t, domain.ProviderCouchDB, meta.Provider)

	assert.Equal(t, userUuid, connections.lastUpserted.UserUuid)
	assert.Equal(t, domain.ProviderCouchDB, connections.lastUpserted.Provider)

	var storedCreds domain.CouchDBKeyCredentials

	err = json.Unmarshal(connections.lastUpserted.CredentialsJSON, &storedCreds)
	require.NoError(t, err)
	assert.Equal(t, server.URL, storedCreds.URL)
	assert.Equal(t, "admin", storedCreds.Username)
	assert.Equal(t, "password", storedCreds.Password)

	assert.Equal(t, 1, couchInstances.registerOwnedCalls, "expected a new owned instance to be registered")
	assert.Zero(t, couchInstances.updateCalls)
}

// TestAddCouchDBConnection_Success_UpdatesExistingOwnedInstance verifies a second BYOK save for
// the same user updates the existing owned couch_instances row in place instead of registering a
// duplicate — see syncOwnedCouchInstance.
func TestAddCouchDBConnection_Success_UpdatesExistingOwnedInstance(t *testing.T) {
	server := couchDBFakeServer(t)
	defer server.Close()

	connections := newFakeExternalConnsRepo()
	couchInstances := newFakeCouchInstancesRepo()

	userUuid := uuid.New()
	existingID := uuid.New()
	couchInstances.owned[userUuid] = domain.CouchInstance{Uuid: existingID, Url: "http://old-url"}

	svc := newCouchDBConnectionTestService(connections, couchInstances)
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userUuid})

	_, err := svc.AddCouchDBConnection(ctx, server.URL, "admin", "password")
	require.NoError(t, err)

	assert.Zero(t, couchInstances.registerOwnedCalls, "expected no new owned instance to be registered")
	assert.Equal(t, 1, couchInstances.updateCalls)
	assert.Equal(t, existingID, couchInstances.lastUpdateID)
	assert.Equal(t, server.URL, couchInstances.lastUpdateURL)
}

func TestCheckCouchDBConnection_Unauthenticated_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	couchInstances := newFakeCouchInstancesRepo()
	svc := newCouchDBConnectionTestService(connections, couchInstances)

	err := svc.CheckCouchDBConnection(context.Background(), "http://localhost:5984", "admin", "password")

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.Unauthenticated)
}

func TestCheckCouchDBConnection_Success_DoesNotPersist(t *testing.T) {
	server := couchDBFakeServer(t)
	defer server.Close()

	connections := newFakeExternalConnsRepo()
	couchInstances := newFakeCouchInstancesRepo()
	svc := newCouchDBConnectionTestService(connections, couchInstances)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: uuid.New()})

	err := svc.CheckCouchDBConnection(ctx, server.URL, "admin", "password")

	require.NoError(t, err)
	assert.Empty(t, connections.conns)
	assert.Zero(t, couchInstances.registerOwnedCalls)
}

func TestCheckCouchDBConnection_ConnectivityCheckFails_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	couchInstances := newFakeCouchInstancesRepo()
	svc := newCouchDBConnectionTestService(connections, couchInstances)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: uuid.New()})

	err := svc.CheckCouchDBConnection(ctx, "http://127.0.0.1:1", "admin", "password")

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.CouchDBConnectionValidationFailed)
}
