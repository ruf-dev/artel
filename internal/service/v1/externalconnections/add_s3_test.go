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

func newS3ConnectionTestService(
	connections *fakeExternalConnsRepo, s3Instances *fakeS3InstancesRepo,
) *Service {
	return New(connections, nil, nil, nil, nil, nil, nil, s3Instances)
}

// s3FakeServer serves a minimal ListAllMyBucketsResult so testS3Connection's bucket-less
// s3client.TestConnection call (ListBuckets) can succeed without a real S3-compatible endpoint.
func s3FakeServer(t *testing.T) *httptest.Server {
	t.Helper()

	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Owner><ID>me</ID><DisplayName>me</DisplayName></Owner>
  <Buckets></Buckets>
</ListAllMyBucketsResult>`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func TestAddS3Connection_Unauthenticated_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	s3Instances := newFakeS3InstancesRepo()
	svc := newS3ConnectionTestService(connections, s3Instances)

	_, err := svc.AddS3Connection(context.Background(), "127.0.0.1:9000", "us-east-1", "access", "secret", false, true)

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.Unauthenticated)
}

// TestAddS3Connection_ConnectivityCheckFails_Errors points at an address nothing listens on
// (loopback, reserved port) so the connection is refused immediately — deterministic, no real
// network required — and verifies neither the connection row nor an owned instance row get
// persisted.
func TestAddS3Connection_ConnectivityCheckFails_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	s3Instances := newFakeS3InstancesRepo()
	svc := newS3ConnectionTestService(connections, s3Instances)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: uuid.New()})

	_, err := svc.AddS3Connection(ctx, "127.0.0.1:1", "us-east-1", "access", "secret", false, true)

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.S3ConnectionValidationFailed)
	assert.Empty(t, connections.conns)
	assert.Zero(t, s3Instances.registerOwnedCalls)
}

func TestAddS3Connection_Success_RegistersOwnedInstance(t *testing.T) {
	server := s3FakeServer(t)
	defer server.Close()

	connections := newFakeExternalConnsRepo()
	s3Instances := newFakeS3InstancesRepo()
	svc := newS3ConnectionTestService(connections, s3Instances)

	userUuid := uuid.New()
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userUuid})

	endpoint := server.Listener.Addr().String()

	meta, err := svc.AddS3Connection(ctx, endpoint, "us-east-1", "access", "secret", false, true)
	require.NoError(t, err)
	assert.Equal(t, domain.ProviderS3, meta.Provider)

	assert.Equal(t, userUuid, connections.lastUpserted.UserUuid)
	assert.Equal(t, domain.ProviderS3, connections.lastUpserted.Provider)

	var storedCreds domain.S3KeyCredentials

	err = json.Unmarshal(connections.lastUpserted.CredentialsJSON, &storedCreds)
	require.NoError(t, err)
	assert.Equal(t, endpoint, storedCreds.Endpoint)
	assert.Equal(t, "access", storedCreds.AccessKey)
	assert.Equal(t, "secret", storedCreds.SecretKey)
	assert.True(t, storedCreds.PathStyle)
	assert.False(t, storedCreds.UseSSL)

	assert.Equal(t, 1, s3Instances.registerOwnedCalls, "expected a new owned instance to be registered")
	assert.Zero(t, s3Instances.updateCalls)
}

// TestAddS3Connection_Success_UpdatesExistingOwnedInstance verifies a second BYOK save for the
// same user updates the existing owned s3_instances row in place instead of registering a
// duplicate — see syncOwnedS3Instance.
func TestAddS3Connection_Success_UpdatesExistingOwnedInstance(t *testing.T) {
	server := s3FakeServer(t)
	defer server.Close()

	connections := newFakeExternalConnsRepo()
	s3Instances := newFakeS3InstancesRepo()

	userUuid := uuid.New()
	existingID := uuid.New()
	s3Instances.owned[userUuid] = domain.S3Instance{Uuid: existingID, Endpoint: "old-endpoint:9000"}

	svc := newS3ConnectionTestService(connections, s3Instances)
	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: userUuid})

	endpoint := server.Listener.Addr().String()

	_, err := svc.AddS3Connection(ctx, endpoint, "us-east-1", "access", "secret", false, true)
	require.NoError(t, err)

	assert.Zero(t, s3Instances.registerOwnedCalls, "expected no new owned instance to be registered")
	assert.Equal(t, 1, s3Instances.updateCalls)
	assert.Equal(t, existingID, s3Instances.lastUpdateID)
	assert.Equal(t, endpoint, s3Instances.lastUpdateEndpoint)
}

func TestCheckS3Connection_Unauthenticated_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	s3Instances := newFakeS3InstancesRepo()
	svc := newS3ConnectionTestService(connections, s3Instances)

	err := svc.CheckS3Connection(context.Background(), "127.0.0.1:9000", "us-east-1", "access", "secret", false, true)

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.Unauthenticated)
}

func TestCheckS3Connection_Success_DoesNotPersist(t *testing.T) {
	server := s3FakeServer(t)
	defer server.Close()

	connections := newFakeExternalConnsRepo()
	s3Instances := newFakeS3InstancesRepo()
	svc := newS3ConnectionTestService(connections, s3Instances)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: uuid.New()})

	endpoint := server.Listener.Addr().String()

	err := svc.CheckS3Connection(ctx, endpoint, "us-east-1", "access", "secret", false, true)

	require.NoError(t, err)
	assert.Empty(t, connections.conns)
	assert.Zero(t, s3Instances.registerOwnedCalls)
}

func TestCheckS3Connection_ConnectivityCheckFails_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	s3Instances := newFakeS3InstancesRepo()
	svc := newS3ConnectionTestService(connections, s3Instances)

	ctx := user_context.WithUserContext(context.Background(), user_context.UserContext{UserUuid: uuid.New()})

	err := svc.CheckS3Connection(ctx, "127.0.0.1:1", "us-east-1", "access", "secret", false, true)

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.S3ConnectionValidationFailed)
}
