package simplechat

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
)

// fakeExternalConnectionRepo is a hand-rolled fake of repository.ExternalConnectionRepo, scoped
// to the GetByUserAndProvider method resolveOpenRouterCredentials uses.
type fakeExternalConnectionRepo struct {
	getByUserAndProviderFunc func(
		ctx context.Context, userUuid uuid.UUID, provider string,
	) (sql.Null[domain.ExternalConnection], error)
}

func (f *fakeExternalConnectionRepo) Upsert(
	_ context.Context, _ domain.ExternalConnection,
) (domain.ExternalConnection, error) {
	panic("not implemented")
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

// TestResolveOpenRouterCredentials_DefaultsEmptyBaseUrl guards against a regression of the bug
// where a connection persisted with an empty BaseUrl (either from before the AddOpenAIConnection
// fix, or any other write path that omits it) made the OpenAI SDK client silently fall back to
// its own hardcoded api.openai.com default — sending an OpenRouter key to OpenAI's endpoint.
func TestResolveOpenRouterCredentials_DefaultsEmptyBaseUrl(t *testing.T) {
	userUuid := uuid.New()

	storedCreds := domain.OpenAIKeyCredentials{
		ApiKey:  "sk-or-test",
		BaseUrl: "",
	}

	credsJSON, err := json.Marshal(storedCreds)
	require.NoError(t, err)

	repo := &fakeExternalConnectionRepo{
		getByUserAndProviderFunc: func(
			_ context.Context, gotUserUuid uuid.UUID, gotProvider string,
		) (sql.Null[domain.ExternalConnection], error) {
			require.Equal(t, userUuid, gotUserUuid)
			require.Equal(t, domain.ProviderOpenRouter, gotProvider)

			conn := domain.ExternalConnection{
				UserUuid:        userUuid,
				Provider:        domain.ProviderOpenRouter,
				CredentialsJSON: json.RawMessage(credsJSON),
			}

			return sql.Null[domain.ExternalConnection]{V: conn, Valid: true}, nil
		},
	}

	s := &Service{connections: repo}

	creds, err := s.resolveOpenRouterCredentials(context.Background(), userUuid)
	require.NoError(t, err)
	require.Equal(t, "sk-or-test", creds.ApiKey)
	require.Equal(t, openrouterDefaultBaseUrl, creds.BaseUrl)
}

// TestResolveOpenRouterCredentials_PreservesExplicitBaseUrl confirms the defaulting is only
// applied when BaseUrl is empty — a caller-supplied override (e.g. a proxy endpoint) must pass
// through unchanged.
func TestResolveOpenRouterCredentials_PreservesExplicitBaseUrl(t *testing.T) {
	userUuid := uuid.New()

	storedCreds := domain.OpenAIKeyCredentials{
		ApiKey:  "sk-or-test",
		BaseUrl: "https://my-proxy.example.com/v1",
	}

	credsJSON, err := json.Marshal(storedCreds)
	require.NoError(t, err)

	repo := &fakeExternalConnectionRepo{
		getByUserAndProviderFunc: func(
			_ context.Context, _ uuid.UUID, _ string,
		) (sql.Null[domain.ExternalConnection], error) {
			conn := domain.ExternalConnection{
				UserUuid:        userUuid,
				Provider:        domain.ProviderOpenRouter,
				CredentialsJSON: json.RawMessage(credsJSON),
			}

			return sql.Null[domain.ExternalConnection]{V: conn, Valid: true}, nil
		},
	}

	s := &Service{connections: repo}

	creds, err := s.resolveOpenRouterCredentials(context.Background(), userUuid)
	require.NoError(t, err)
	require.Equal(t, "https://my-proxy.example.com/v1", creds.BaseUrl)
}

// TestResolveOpenRouterCredentials_MissingConnection confirms the absent-connection error still
// surfaces correctly and isn't masked by the new defaulting logic.
func TestResolveOpenRouterCredentials_MissingConnection(t *testing.T) {
	userUuid := uuid.New()

	repo := &fakeExternalConnectionRepo{
		getByUserAndProviderFunc: func(
			_ context.Context, _ uuid.UUID, _ string,
		) (sql.Null[domain.ExternalConnection], error) {
			return sql.Null[domain.ExternalConnection]{Valid: false}, nil
		},
	}

	s := &Service{connections: repo}

	_, err := s.resolveOpenRouterCredentials(context.Background(), userUuid)
	require.Error(t, err)
}
