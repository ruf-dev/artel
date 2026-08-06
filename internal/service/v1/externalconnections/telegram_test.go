package externalconnections

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func newTelegramConnectionTestService(connections *fakeExternalConnsRepo, mom *fakeMomService) *Service {
	return New(connections, nil, nil, nil, nil, mom)
}

func TestAddTelegramConnection_Unauthenticated_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	mom := &fakeMomService{result: `{"ok":true,"result":{"username":"testbot"}}`}
	svc := newTelegramConnectionTestService(connections, mom)

	_, err := svc.AddTelegramConnection(context.Background(), "bot-token")

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.Unauthenticated)
}

func TestAddTelegramConnection_ValidToken_UpsertsCredentialsAndReturnsMeta(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	mom := &fakeMomService{result: `{"ok":true,"result":{"username":"testbot"}}`}
	svc := newTelegramConnectionTestService(connections, mom)

	userUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: userUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	meta, err := svc.AddTelegramConnection(ctx, "bot-token")

	require.NoError(t, err)
	assert.Equal(t, "testbot", meta.DisplayName)

	assert.Equal(t, "telegram", mom.lastMcpName)
	assert.Equal(t, "get_me", mom.lastToolName)
	assert.Equal(t, "bot-token", mom.lastSecrets["bot_token"])

	assert.Equal(t, userUuid, connections.lastUpserted.UserUuid)
	assert.Equal(t, domain.ProviderTelegram, connections.lastUpserted.Provider)

	var storedCreds domain.TelegramCredentials

	err = json.Unmarshal(connections.lastUpserted.CredentialsJSON, &storedCreds)
	require.NoError(t, err)
	assert.Equal(t, "bot-token", storedCreds.BotToken)

	var storedMeta telegramConnectionMeta

	err = json.Unmarshal(connections.lastUpserted.Metadata, &storedMeta)
	require.NoError(t, err)
	assert.Equal(t, "testbot", storedMeta.BotUsername)
}

func TestAddTelegramConnection_InvalidToken_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	mom := &fakeMomService{result: `{"ok":false}`}
	svc := newTelegramConnectionTestService(connections, mom)

	userUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: userUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	_, err := svc.AddTelegramConnection(ctx, "bad-token")

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.TelegramValidationFailed)
	assert.Empty(t, connections.conns)
}

func TestCheckTelegramConnection_ValidToken_ReturnsUsername(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	mom := &fakeMomService{result: `{"ok":true,"result":{"username":"testbot"}}`}
	svc := newTelegramConnectionTestService(connections, mom)

	userUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: userUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	botUsername, err := svc.CheckTelegramConnection(ctx, "bot-token")

	require.NoError(t, err)
	assert.Equal(t, "testbot", botUsername)
	assert.Empty(t, connections.conns)
}

func TestCheckTelegramConnection_Unauthenticated_Errors(t *testing.T) {
	connections := newFakeExternalConnsRepo()
	mom := &fakeMomService{result: `{"ok":true,"result":{"username":"testbot"}}`}
	svc := newTelegramConnectionTestService(connections, mom)

	_, err := svc.CheckTelegramConnection(context.Background(), "bot-token")

	require.Error(t, err)
	assert.ErrorIs(t, err, user_errors.Unauthenticated)
}
