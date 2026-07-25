package externalconnections

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/middleware/user_context"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func newGenericConnectionTestService(connections *fakeExternalConnsRepo) *Service {
	return New(connections, nil, nil, nil, nil, nil)
}

func TestAddGenericConnection_Unauthenticated_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	svc := newGenericConnectionTestService(connections)

	credentials := map[string]string{"token": "secret"}

	_, err := svc.AddGenericConnection(context.Background(), "acme_connector", credentials)

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.Unauthenticated)
}

func TestAddGenericConnection_EmptyProvider_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	svc := newGenericConnectionTestService(connections)

	uc := user_context.UserContext{UserUuid: uuid.New()}
	ctx := user_context.WithUserContext(context.Background(), uc)

	credentials := map[string]string{"token": "secret"}

	_, err := svc.AddGenericConnection(ctx, "", credentials)

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.GenericProviderRequired)
}

func TestAddGenericConnection_Valid_UpsertsMarshaledCredentials(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	svc := newGenericConnectionTestService(connections)

	userUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: userUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	credentials := map[string]string{"api_key": "abc123", "region": "eu"}

	meta, err := svc.AddGenericConnection(ctx, "acme_connector", credentials)

	require.NoError(t, err)
	assert.Equal(t, "acme_connector", meta.Provider)
	assert.Equal(t, artel_q.ExternalProviderTypeApiKey, meta.ProviderType)

	assert.Equal(t, userUuid, connections.lastUpserted.UserUuid)
	assert.Equal(t, "acme_connector", connections.lastUpserted.Provider)
	assert.Equal(t, artel_q.ExternalProviderTypeApiKey, connections.lastUpserted.ProviderType)

	var storedCreds map[string]string

	err = json.Unmarshal(connections.lastUpserted.CredentialsJSON, &storedCreds)
	require.NoError(t, err)
	assert.Equal(t, credentials, storedCreds)
}
