package telegram_webhook

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
)

// fakeExternalConnRepo is a hand-rolled fake of repository.ExternalConnectionRepo — only GetByID
// and Upsert are exercised by this package's tests, the rest panic if called since nothing here
// should reach them.
type fakeExternalConnRepo struct {
	getByIDFunc func(ctx context.Context, id uuid.UUID) (domain.ExternalConnection, error)
	upsertFunc  func(ctx context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error)
}

func (f *fakeExternalConnRepo) Upsert(ctx context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error) {
	if f.upsertFunc != nil {
		return f.upsertFunc(ctx, conn)
	}

	return conn, nil
}

func (f *fakeExternalConnRepo) Insert(ctx context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error) {
	panic("not implemented")
}

func (f *fakeExternalConnRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.ExternalConnection, error) {
	return f.getByIDFunc(ctx, id)
}

func (f *fakeExternalConnRepo) GetByUserAndProvider(
	ctx context.Context, userUuid uuid.UUID, provider string,
) (sql.Null[domain.ExternalConnection], error) {
	panic("not implemented")
}

func (f *fakeExternalConnRepo) ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.ExternalConnection, error) {
	panic("not implemented")
}

func (f *fakeExternalConnRepo) Delete(ctx context.Context, userUuid uuid.UUID, provider string) error {
	panic("not implemented")
}

func (f *fakeExternalConnRepo) DeleteByID(ctx context.Context, userUuid uuid.UUID, id uuid.UUID) error {
	panic("not implemented")
}

// fakeRecentWorkbenchLookup is a hand-rolled fake of recentWorkbenchLookup — the handler_test
// cases below all exercise the "user has no workbench yet" branch, so a fixed empty result is
// enough; other tests in this package exercise the dialing path separately.
type fakeRecentWorkbenchLookup struct{}

func (f *fakeRecentWorkbenchLookup) GetMostRecentByUser(ctx context.Context, userID uuid.UUID) (sql.Null[domain.Workbench], error) {
	return sql.Null[domain.Workbench]{}, nil
}

// fakeBridgeTargetResolver is a hand-rolled fake of bridgeTargetResolver — unused by the auth-path
// tests below (they never get far enough to resolve a workbench target), present only to satisfy
// New's constructor signature.
type fakeBridgeTargetResolver struct{}

func (f *fakeBridgeTargetResolver) ResolveTerminalTarget(ctx context.Context, vaultID, userID uuid.UUID) (string, error) {
	return "", nil
}

// fakeToolExecutor is a hand-rolled fake of toolExecutor, recording every call it receives and
// returning a canned Telegram-shaped success response so callers that parse a message_id out of
// it (sendText/sendKeyboard) don't choke on it.
type fakeToolExecutor struct {
	calls []fakeToolCall
}

type fakeToolCall struct {
	exConnUuid uuid.UUID
	mcpName    string
	toolName   string
	params     map[string]interface{}
}

func (f *fakeToolExecutor) ExecuteToolForConnection(
	ctx context.Context, exConnUuid uuid.UUID, mcpName string, toolName string, params map[string]interface{},
) (string, error) {
	call := fakeToolCall{exConnUuid: exConnUuid, mcpName: mcpName, toolName: toolName, params: params}
	f.calls = append(f.calls, call)

	return `{"ok":true,"result":{"message_id":1}}`, nil
}

// telegramConnection builds a domain.ExternalConnection carrying TelegramCredentials with the
// given webhook secret, for use as the fakeExternalConnRepo.GetByID result.
func telegramConnection(t *testing.T, secret string) domain.ExternalConnection {
	t.Helper()

	creds := domain.TelegramCredentials{
		BotToken:      "test-token",
		WebhookSecret: secret,
	}

	credJSON, err := json.Marshal(creds)
	require.NoError(t, err)

	conn := domain.ExternalConnection{
		Uuid:            uuid.New(),
		UserUuid:        uuid.New(),
		Provider:        domain.ProviderTelegram,
		CredentialsJSON: credJSON,
	}

	return conn
}

func newTestHandler(conn domain.ExternalConnection) (*Handler, *fakeToolExecutor) {
	repo := &fakeExternalConnRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.ExternalConnection, error) {
			return conn, nil
		},
	}

	tools := &fakeToolExecutor{}

	handler := New(context.Background(), repo, &fakeRecentWorkbenchLookup{}, &fakeBridgeTargetResolver{}, tools)

	return handler, tools
}

func webhookRequest(connUuid uuid.UUID, secretHeader string, body []byte) *http.Request {
	target := pathPrefix + connUuid.String()

	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	if secretHeader != "" {
		req.Header.Set("X-Telegram-Bot-Api-Secret-Token", secretHeader)
	}

	return req
}

func TestServeHTTP_SecretAuth(t *testing.T) {
	const realSecret = "correct-horse-battery-staple"

	conn := telegramConnection(t, realSecret)

	tests := []struct {
		name           string
		presentedToken string
		wantStatus     int
		// wantRelayed is true once the secret matches and the handler goes on to actually process
		// the update — here, past the "no workbench" branch (this handler.go has no workbench
		// wired) to send a plain-text explanation back through the telegram MoM tool.
		wantRelayed bool
	}{
		{name: "missing secret header", presentedToken: "", wantStatus: http.StatusUnauthorized, wantRelayed: false},
		{name: "wrong secret", presentedToken: "wrong-secret", wantStatus: http.StatusUnauthorized, wantRelayed: false},
		{name: "correct secret", presentedToken: realSecret, wantStatus: http.StatusOK, wantRelayed: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, tools := newTestHandler(conn)

			body := []byte(`{"message":{"message_id":1,"chat":{"id":42},"from":{"id":7},"text":"hello"}}`)
			req := webhookRequest(conn.Uuid, tt.presentedToken, body)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Equal(t, tt.wantRelayed, len(tools.calls) > 0)
		})
	}
}

func TestServeHTTP_MethodNotAllowed(t *testing.T) {
	conn := telegramConnection(t, "secret")
	handler, _ := newTestHandler(conn)

	req := httptest.NewRequest(http.MethodGet, pathPrefix+conn.Uuid.String(), nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestServeHTTP_InvalidConnectionId(t *testing.T) {
	conn := telegramConnection(t, "secret")
	handler, _ := newTestHandler(conn)

	req := httptest.NewRequest(http.MethodPost, pathPrefix+"not-a-uuid", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestServeHTTP_WrongProvider(t *testing.T) {
	conn := telegramConnection(t, "secret")
	conn.Provider = domain.ProviderGitlab

	handler, _ := newTestHandler(conn)

	req := webhookRequest(conn.Uuid, "secret", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestServeHTTP_ConnectionNotFound(t *testing.T) {
	repo := &fakeExternalConnRepo{
		getByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.ExternalConnection, error) {
			return domain.ExternalConnection{}, sql.ErrNoRows
		},
	}

	handler := New(context.Background(), repo, &fakeRecentWorkbenchLookup{}, &fakeBridgeTargetResolver{}, &fakeToolExecutor{})

	req := webhookRequest(uuid.New(), "secret", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTelegramUpdate_UnmarshalMessage(t *testing.T) {
	raw := []byte(`{"message":{"message_id":10,"chat":{"id":555},"from":{"id":99},"text":"hi bot"}}`)

	var update telegramUpdate

	err := json.Unmarshal(raw, &update)
	require.NoError(t, err)

	require.Nil(t, update.CallbackQuery)
	require.NotNil(t, update.Message)
	require.Equal(t, int64(10), update.Message.MessageId)
	require.Equal(t, int64(555), update.Message.Chat.Id)
	require.Equal(t, int64(99), update.Message.From.Id)
	require.Equal(t, "hi bot", update.Message.Text)
}

func TestTelegramUpdate_UnmarshalCallbackQuery(t *testing.T) {
	raw := []byte(`{"callback_query":{"id":"cbq-1","data":"evt-1:allow_once","message":{"message_id":20,"chat":{"id":777}}}}`)

	var update telegramUpdate

	err := json.Unmarshal(raw, &update)
	require.NoError(t, err)

	require.Nil(t, update.Message)
	require.NotNil(t, update.CallbackQuery)
	require.Equal(t, "cbq-1", update.CallbackQuery.Id)
	require.Equal(t, "evt-1:allow_once", update.CallbackQuery.Data)
	require.NotNil(t, update.CallbackQuery.Message)
	require.Equal(t, int64(20), update.CallbackQuery.Message.MessageId)
	require.Equal(t, int64(777), update.CallbackQuery.Message.Chat.Id)
}

func TestCallbackData_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		eventId  string
		decision permissionDecision
	}{
		{name: "allow once", eventId: uuid.NewString(), decision: decisionAllowOnce},
		{name: "allow always", eventId: uuid.NewString(), decision: decisionAllowAlways},
		{name: "deny", eventId: uuid.NewString(), decision: decisionDeny},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := encodeCallbackData(tt.eventId, tt.decision)

			require.LessOrEqual(t, len(encoded), 64, "callback_data must fit Telegram's 64 byte cap")

			gotEventId, gotDecision, ok := decodeCallbackData(encoded)
			require.True(t, ok)
			require.Equal(t, tt.eventId, gotEventId)
			require.Equal(t, tt.decision, gotDecision)
		})
	}
}

func TestCallbackData_DecodeInvalid(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{name: "no separator", data: "just-an-id"},
		{name: "unknown decision", data: "evt-1:allow_forever"},
		{name: "empty event id", data: ":allow_once"},
		{name: "empty string", data: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, ok := decodeCallbackData(tt.data)
			require.False(t, ok)
		})
	}
}
