package externalconnections

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	anthropicClient "github.com/ruf-dev/artel/internal/clients/anthropic"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// newValidationTestService builds a *Service with no live dependencies. The pieces under test
// here (validateAnthropicKey, recommendedDefaultAnthropicModel, anthropicModelIds,
// anthropicKeyPreview) never touch the repository, pending-codes store, spreadsheets repo, mail
// suggestions repo, oauth config, or mom service, so nil stand-ins are safe.
func newValidationTestService() *Service {
	return New(nil, nil, nil, nil, nil, nil)
}

// TestValidateAnthropicKey_Success verifies a successful models-list response against a fake
// Anthropic endpoint is returned as-is (mirrors internal/clients/anthropic/client_test.go's
// TestListModels_Success fixture).
func TestValidateAnthropicKey_Success(t *testing.T) {
	body := `{
		"data": [
			{
				"id": "claude-opus-4-8",
				"display_name": "Claude Opus 4.8",
				"max_input_tokens": 1000000,
				"max_tokens": 128000,
				"type": "model",
				"created_at": "2026-01-01T00:00:00Z"
			},
			{
				"id": "claude-haiku-4-5",
				"display_name": "Claude Haiku 4.5",
				"max_input_tokens": 200000,
				"max_tokens": 64000,
				"type": "model",
				"created_at": "2026-01-01T00:00:00Z"
			}
		],
		"has_more": false,
		"first_id": "claude-opus-4-8",
		"last_id": "claude-haiku-4-5"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	svc := newValidationTestService()

	models, err := svc.validateAnthropicKey(context.Background(), "test-api-key", server.URL, "")
	require.NoError(t, err)

	want := []anthropicClient.ModelInfo{
		{Id: "claude-opus-4-8", DisplayName: "Claude Opus 4.8", MaxInputTokens: 1000000, MaxTokens: 128000},
		{Id: "claude-haiku-4-5", DisplayName: "Claude Haiku 4.5", MaxInputTokens: 200000, MaxTokens: 64000},
	}
	require.Equal(t, want, models)
}

// TestValidateAnthropicKey_AuthFailure verifies a 401 from the provider surfaces as
// user_errors.LlmKeyRejected rather than a raw transport error.
func TestValidateAnthropicKey_AuthFailure(t *testing.T) {
	body := `{
		"type": "error",
		"error": {
			"type": "authentication_error",
			"message": "invalid x-api-key"
		},
		"request_id": "req_01"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	svc := newValidationTestService()

	models, err := svc.validateAnthropicKey(context.Background(), "bad-api-key", server.URL, "")
	require.Error(t, err)
	require.Nil(t, models)
	require.ErrorIs(t, err, user_errors.LlmKeyRejected)
}

// TestValidateAnthropicKey_ModelListingUnsupported verifies a 404 from the provider (a
// non-official Anthropic-compatible endpoint that doesn't implement GET /v1/models) surfaces as
// user_errors.LlmKeyModelListingUnsupported when the caller supplied no fallback model to ping.
func TestValidateAnthropicKey_ModelListingUnsupported(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	svc := newValidationTestService()

	models, err := svc.validateAnthropicKey(context.Background(), "test-api-key", server.URL, "")
	require.Error(t, err)
	require.Nil(t, models)
	require.ErrorIs(t, err, user_errors.LlmKeyModelListingUnsupported)
}

// TestValidateAnthropicKey_FallbackPingSucceeds verifies that when GET /v1/models 404s but the
// caller supplied a defaultModel, validation falls back to a minimal completion against that
// model and succeeds if the provider accepts it — the case that unblocks non-official
// Anthropic-compatible providers that only implement the Messages API.
func TestValidateAnthropicKey_FallbackPingSucceeds(t *testing.T) {
	messageBody := `{
		"id": "msg_01ABC",
		"type": "message",
		"role": "assistant",
		"model": "custom-model",
		"content": [{"type": "text", "text": "hi"}],
		"stop_reason": "end_turn",
		"stop_sequence": null,
		"usage": {"input_tokens": 1, "output_tokens": 1, "cache_creation_input_tokens": 0, "cache_read_input_tokens": 0}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/messages" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, err := w.Write([]byte(messageBody))
			require.NoError(t, err)

			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	svc := newValidationTestService()

	models, err := svc.validateAnthropicKey(context.Background(), "test-api-key", server.URL, "custom-model")
	require.NoError(t, err)
	require.Empty(t, models)
}

// TestValidateAnthropicKey_FallbackPingFails verifies that when both GET /v1/models and the
// fallback ping fail, validation reports the key as rejected rather than silently succeeding.
func TestValidateAnthropicKey_FallbackPingFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	svc := newValidationTestService()

	models, err := svc.validateAnthropicKey(context.Background(), "test-api-key", server.URL, "custom-model")
	require.Error(t, err)
	require.Nil(t, models)
	require.ErrorIs(t, err, user_errors.LlmKeyRejected)
}

// TestValidateAnthropicKey_ProviderUnreachable verifies a network-level failure (no server
// listening at all, as opposed to an API error response) surfaces as
// user_errors.LlmKeyProviderUnreachable.
func TestValidateAnthropicKey_ProviderUnreachable(t *testing.T) {
	svc := newValidationTestService()

	models, err := svc.validateAnthropicKey(context.Background(), "test-api-key", "http://127.0.0.1:1", "")
	require.Error(t, err)
	require.Nil(t, models)
	require.ErrorIs(t, err, user_errors.LlmKeyProviderUnreachable)
}

// TestValidateAnthropicKey_DefaultBaseUrl verifies a blank baseUrl is resolved to
// anthropicDefaultBaseUrl rather than being sent to the SDK as-is (which would fall back to the
// SDK's own default host instead of this package's documented default).
func TestValidateAnthropicKey_DefaultBaseUrl(t *testing.T) {
	svc := newValidationTestService()

	// A blank baseUrl resolves to the real Anthropic API host, which this test must not call.
	// Instead, confirm the resolution happens by checking a non-blank baseUrl is passed through
	// unchanged to the client (exercised by TestValidateAnthropicKey_Success above) and that
	// resolveAnthropicBaseUrl itself returns the constant for a blank input.
	require.Equal(t, anthropicDefaultBaseUrl, resolveAnthropicBaseUrl(""))
	require.Equal(t, "https://example.com", resolveAnthropicBaseUrl("https://example.com"))

	_ = svc
}

// TestRecommendedDefaultAnthropicModel_UsesProvidedDefault verifies an explicit defaultModel wins
// over the model list.
func TestRecommendedDefaultAnthropicModel_UsesProvidedDefault(t *testing.T) {
	models := []anthropicClient.ModelInfo{
		{Id: "claude-opus-4-8"},
		{Id: "claude-haiku-4-5"},
	}

	got := recommendedDefaultAnthropicModel("claude-haiku-4-5", models)
	require.Equal(t, "claude-haiku-4-5", got)
}

// TestRecommendedDefaultAnthropicModel_FallsBackToNewest verifies that, absent an explicit
// default, the first model in the list (newest, per Anthropic's documented ordering) is picked.
func TestRecommendedDefaultAnthropicModel_FallsBackToNewest(t *testing.T) {
	models := []anthropicClient.ModelInfo{
		{Id: "claude-opus-4-8"},
		{Id: "claude-haiku-4-5"},
	}

	got := recommendedDefaultAnthropicModel("", models)
	require.Equal(t, "claude-opus-4-8", got)
}

// TestRecommendedDefaultAnthropicModel_EmptyModels verifies an empty model list with no explicit
// default doesn't panic and returns an empty string.
func TestRecommendedDefaultAnthropicModel_EmptyModels(t *testing.T) {
	got := recommendedDefaultAnthropicModel("", nil)
	require.Equal(t, "", got)
}

// TestAnthropicModelIds verifies model IDs are extracted in order.
func TestAnthropicModelIds(t *testing.T) {
	models := []anthropicClient.ModelInfo{
		{Id: "claude-opus-4-8"},
		{Id: "claude-haiku-4-5"},
	}

	require.Equal(t, []string{"claude-opus-4-8", "claude-haiku-4-5"}, anthropicModelIds(models))
	require.Equal(t, []string{}, anthropicModelIds(nil))
}

// TestAnthropicKeyPreview verifies the preview never exposes more than the trailing 4 characters
// of the key, and doesn't panic on a short key.
func TestAnthropicKeyPreview(t *testing.T) {
	require.Equal(t, "...ab12", anthropicKeyPreview("sk-ant-api03-verylongkeyab12"))
	require.Equal(t, "...ab", anthropicKeyPreview("ab"))
	require.Equal(t, "...", anthropicKeyPreview(""))
}
