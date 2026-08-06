package tract

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// TestLlmExecutor_Call_Anthropic verifies Call dispatches an Anthropic connection to the
// Anthropic client with the connection's stored credentials, and maps the response (including
// cache-token usage fields) into LlmCallResult.
func TestLlmExecutor_Call_Anthropic(t *testing.T) {
	body := `{
		"id": "msg_01ABC",
		"type": "message",
		"role": "assistant",
		"model": "claude-opus-4-8",
		"content": [{"type": "text", "text": "hello from anthropic"}],
		"stop_reason": "end_turn",
		"stop_sequence": null,
		"usage": {
			"input_tokens": 10,
			"output_tokens": 5,
			"cache_creation_input_tokens": 2,
			"cache_read_input_tokens": 1
		}
	}`

	var gotPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	creds := domain.AnthropicKeyCredentials{ApiKey: "test-anthropic-key", BaseUrl: server.URL}

	credJSON, err := json.Marshal(creds)
	require.NoError(t, err)

	userUuid := uuid.New()
	connUuid := uuid.New()

	conn := domain.ExternalConnection{
		Uuid:            connUuid,
		UserUuid:        userUuid,
		Provider:        domain.ProviderAnthropic,
		CredentialsJSON: credJSON,
	}

	repo := newFakeExternalConnsRepo()
	repo.conns[connUuid] = conn

	executor := NewLlmExecutor(repo)

	req := LlmCallRequest{
		Model:        "claude-opus-4-8",
		Prompt:       "hi",
		SystemPrompt: "be nice",
		MaxTokens:    100,
	}

	result, err := executor.Call(context.Background(), userUuid, connUuid, req)
	require.NoError(t, err)
	require.Equal(t, "/v1/messages", gotPath)

	want := LlmCallResult{
		Text:  "hello from anthropic",
		Model: "claude-opus-4-8",
		Usage: LlmCallUsage{
			InputTokens:              10,
			OutputTokens:             5,
			CacheCreationInputTokens: 2,
			CacheReadInputTokens:     1,
		},
	}
	require.Equal(t, want, result)
}

// TestLlmExecutor_Call_OpenAI verifies Call dispatches an OpenAI connection to the OpenAI
// client with the connection's stored credentials, and maps the response into LlmCallResult —
// leaving the cache-token usage fields at their zero value (no OpenAI analogue).
func TestLlmExecutor_Call_OpenAI(t *testing.T) {
	body := `{
		"id": "chatcmpl-abc123",
		"object": "chat.completion",
		"created": 1715367049,
		"model": "gpt-4o",
		"choices": [
			{"index": 0, "message": {"role": "assistant", "content": "hello from openai"}, "logprobs": null, "finish_reason": "stop"}
		],
		"usage": {"prompt_tokens": 7, "completion_tokens": 3, "total_tokens": 10}
	}`

	var gotPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	creds := domain.OpenAIKeyCredentials{ApiKey: "test-openai-key", BaseUrl: server.URL}

	credJSON, err := json.Marshal(creds)
	require.NoError(t, err)

	userUuid := uuid.New()
	connUuid := uuid.New()

	conn := domain.ExternalConnection{
		Uuid:            connUuid,
		UserUuid:        userUuid,
		Provider:        domain.ProviderOpenAI,
		CredentialsJSON: credJSON,
	}

	repo := newFakeExternalConnsRepo()
	repo.conns[connUuid] = conn

	executor := NewLlmExecutor(repo)

	req := LlmCallRequest{
		Model:        "gpt-4o",
		Prompt:       "hi",
		SystemPrompt: "be nice",
		MaxTokens:    100,
	}

	result, err := executor.Call(context.Background(), userUuid, connUuid, req)
	require.NoError(t, err)
	require.Equal(t, "/chat/completions", gotPath)

	want := LlmCallResult{
		Text:  "hello from openai",
		Model: "gpt-4o",
		Usage: LlmCallUsage{
			InputTokens:  7,
			OutputTokens: 3,
		},
	}
	require.Equal(t, want, result)
}

// TestLlmExecutor_Call_UnknownProvider verifies a connection whose provider is neither
// anthropic nor openai (e.g. a non-LLM connection reused by mistake, or a future provider not
// yet wired into the tract engine) surfaces as TractLlmConnectionProviderMismatch rather than
// panicking or silently no-op'ing.
func TestLlmExecutor_Call_UnknownProvider(t *testing.T) {
	userUuid := uuid.New()
	connUuid := uuid.New()

	conn := domain.ExternalConnection{
		Uuid:     connUuid,
		UserUuid: userUuid,
		Provider: domain.ProviderGitlab,
	}

	repo := newFakeExternalConnsRepo()
	repo.conns[connUuid] = conn

	executor := NewLlmExecutor(repo)

	req := LlmCallRequest{Model: "whatever", Prompt: "hi"}

	result, err := executor.Call(context.Background(), userUuid, connUuid, req)
	require.Error(t, err)
	require.Equal(t, LlmCallResult{}, result)
	require.ErrorIs(t, err, user_errors.TractLlmConnectionProviderMismatch)
}

// TestLlmExecutor_Call_ConnectionNotOwned verifies a connection owned by a different user
// surfaces as TractConnectionNotOwned rather than leaking another user's credentials.
func TestLlmExecutor_Call_ConnectionNotOwned(t *testing.T) {
	ownerUuid := uuid.New()
	callerUuid := uuid.New()
	connUuid := uuid.New()

	conn := domain.ExternalConnection{
		Uuid:     connUuid,
		UserUuid: ownerUuid,
		Provider: domain.ProviderAnthropic,
	}

	repo := newFakeExternalConnsRepo()
	repo.conns[connUuid] = conn

	executor := NewLlmExecutor(repo)

	req := LlmCallRequest{Model: "whatever", Prompt: "hi"}

	result, err := executor.Call(context.Background(), callerUuid, connUuid, req)
	require.Error(t, err)
	require.Equal(t, LlmCallResult{}, result)
	require.ErrorIs(t, err, user_errors.TractConnectionNotOwned)
}
