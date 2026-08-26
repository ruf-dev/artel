package externalconnections

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
)

// fakeExternalConnectionRepo is a hand-rolled fake of repository.ExternalConnectionRepo, scoped
// to the GetByUserAndProvider/Upsert methods AddOpenAIConnection uses.
type fakeExternalConnectionRepo struct {
	getByUserAndProviderFunc func(
		ctx context.Context, userUuid uuid.UUID, provider string,
	) (sql.Null[domain.ExternalConnection], error)
	upsertFunc func(ctx context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error)
}

func (f *fakeExternalConnectionRepo) Upsert(
	ctx context.Context, conn domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	return f.upsertFunc(ctx, conn)
}

func (f *fakeExternalConnectionRepo) Insert(
	_ context.Context, _ domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	panic("not implemented")
}

func (f *fakeExternalConnectionRepo) GetByID(_ context.Context, _ uuid.UUID) (domain.ExternalConnection, error) {
	panic("not implemented")
}

func (f *fakeExternalConnectionRepo) GetByUserAndProvider(
	ctx context.Context, userUuid uuid.UUID, provider string,
) (sql.Null[domain.ExternalConnection], error) {
	return f.getByUserAndProviderFunc(ctx, userUuid, provider)
}

func (f *fakeExternalConnectionRepo) ListByUser(_ context.Context, _ uuid.UUID) ([]domain.ExternalConnection, error) {
	panic("not implemented")
}

func (f *fakeExternalConnectionRepo) Delete(_ context.Context, _ uuid.UUID, _ string) error {
	panic("not implemented")
}

func (f *fakeExternalConnectionRepo) DeleteByID(_ context.Context, _, _ uuid.UUID) error {
	panic("not implemented")
}

// hostRedirectTransport rewrites any request targeting targetHost to dest instead, so a test can
// intercept calls the code under test makes to a real, hardcoded external hostname (here,
// openrouter.ai) and serve them from a local httptest.Server without touching the network.
type hostRedirectTransport struct {
	targetHost string
	dest       *url.URL
	underlying http.RoundTripper
}

func (t *hostRedirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == t.targetHost {
		req = req.Clone(req.Context())
		req.URL.Scheme = t.dest.Scheme
		req.URL.Host = t.dest.Host
		req.Host = t.dest.Host
	}

	return t.underlying.RoundTrip(req)
}

// withRedirectedHost points requests to targetHost at server for the duration of the test,
// restoring the original http.DefaultTransport on cleanup.
func withRedirectedHost(t *testing.T, targetHost string, server *httptest.Server) {
	t.Helper()

	dest, err := url.Parse(server.URL)
	require.NoError(t, err)

	original := http.DefaultTransport
	http.DefaultTransport = &hostRedirectTransport{targetHost: targetHost, dest: dest, underlying: original}

	t.Cleanup(func() {
		http.DefaultTransport = original
	})
}

// newFakeModelsListServer serves a minimal, valid GET /v1/models response — enough for
// openaiClient.ListModels to succeed without asserting on the request path, since the only
// caller in this test always hits "/models" relative to whatever base URL it was given.
func newFakeModelsListServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(`{"object":"list","data":[{"id":"test-model","object":"model","owned_by":"test","created":0}]}`))
		require.NoError(t, err)
	}))
}

// TestAddOpenAIConnection_OpenRouterBlankBaseUrlPersistsDefault guards against the bug where a
// blank baseUrl on an OpenRouter BYOK connection was persisted as-is (empty) instead of the
// resolved OpenRouter default, which made the OpenAI SDK client fall back to its own hardcoded
// api.openai.com default at chat time — sending an OpenRouter key to OpenAI's endpoint.
func TestAddOpenAIConnection_OpenRouterBlankBaseUrlPersistsDefault(t *testing.T) {
	server := newFakeModelsListServer(t)
	defer server.Close()

	withRedirectedHost(t, "openrouter.ai", server)

	userUuid := uuid.New()
	uc := user_context.UserContext{UserUuid: userUuid}
	ctx := user_context.WithUserContext(context.Background(), uc)

	var savedConn domain.ExternalConnection

	repo := &fakeExternalConnectionRepo{
		upsertFunc: func(_ context.Context, conn domain.ExternalConnection) (domain.ExternalConnection, error) {
			savedConn = conn
			conn.Uuid = uuid.New()

			return conn, nil
		},
	}

	s := &Service{connections: repo}

	_, err := s.AddOpenAIConnection(ctx, "sk-or-test", "", "", domain.ProviderOpenRouter)
	require.NoError(t, err)

	require.Equal(t, artel_q.ExternalProviderTypeApiKey, savedConn.ProviderType)

	var savedCreds domain.OpenAIKeyCredentials

	err = json.Unmarshal(savedConn.CredentialsJSON, &savedCreds)
	require.NoError(t, err)

	require.Equal(t, "sk-or-test", savedCreds.ApiKey)
	require.Equal(t, openrouterDefaultBaseUrl, savedCreds.BaseUrl)
	require.NotEmpty(t, savedCreds.BaseUrl)
}
