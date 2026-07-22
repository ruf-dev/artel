package anthropic_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/anthropic"
)

// TestListModels_Success verifies a models-list response is mapped into the returned
// []ModelInfo with the expected fields.
func TestListModels_Success(t *testing.T) {
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
		require.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := anthropic.New("test-api-key", server.URL)

	models, err := c.ListModels(context.Background())
	require.NoError(t, err)

	want := []anthropic.ModelInfo{
		{
			Id:             "claude-opus-4-8",
			DisplayName:    "Claude Opus 4.8",
			MaxInputTokens: 1000000,
			MaxTokens:      128000,
		},
		{
			Id:             "claude-haiku-4-5",
			DisplayName:    "Claude Haiku 4.5",
			MaxInputTokens: 200000,
			MaxTokens:      64000,
		},
	}
	require.Equal(t, want, models)
}

// TestListModels_AuthFailure verifies a 401 response from the fake server surfaces as an error.
func TestListModels_AuthFailure(t *testing.T) {
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

	c := anthropic.New("bad-api-key", server.URL)

	models, err := c.ListModels(context.Background())
	require.Error(t, err)
	require.Nil(t, models)
}

// TestComplete_Success verifies a messages-create response is mapped into CompleteResult,
// including the usage object fields.
func TestComplete_Success(t *testing.T) {
	body := `{
		"id": "msg_01ABC",
		"type": "message",
		"role": "assistant",
		"model": "claude-opus-4-8",
		"content": [
			{"type": "text", "text": "Hello, world!"}
		],
		"stop_reason": "end_turn",
		"stop_sequence": null,
		"usage": {
			"input_tokens": 10,
			"output_tokens": 5,
			"cache_creation_input_tokens": 2,
			"cache_read_input_tokens": 1
		}
	}`

	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/messages", r.URL.Path)

		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		gotBody = buf

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := anthropic.New("test-api-key", server.URL)

	req := anthropic.CompleteRequest{
		Model:        "claude-opus-4-8",
		SystemPrompt: "You are a helpful assistant.",
		Prompt:       "Say hello",
		MaxTokens:    1024,
	}

	result, err := c.Complete(context.Background(), req)
	require.NoError(t, err)
	require.NotEmpty(t, gotBody)

	want := anthropic.CompleteResult{
		Text: "Hello, world!",
		Usage: anthropic.Usage{
			InputTokens:              10,
			OutputTokens:             5,
			CacheCreationInputTokens: 2,
			CacheReadInputTokens:     1,
		},
	}
	require.Equal(t, want, result)
}

// TestComplete_DefaultMaxTokens verifies that omitting MaxTokens still produces a successful
// request (the client applies its own default rather than sending max_tokens: 0).
func TestComplete_DefaultMaxTokens(t *testing.T) {
	body := `{
		"id": "msg_01DEF",
		"type": "message",
		"role": "assistant",
		"model": "claude-opus-4-8",
		"content": [{"type": "text", "text": "ok"}],
		"stop_reason": "end_turn",
		"stop_sequence": null,
		"usage": {
			"input_tokens": 1,
			"output_tokens": 1,
			"cache_creation_input_tokens": 0,
			"cache_read_input_tokens": 0
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := anthropic.New("test-api-key", server.URL)

	req := anthropic.CompleteRequest{
		Model:  "claude-opus-4-8",
		Prompt: "Say ok",
	}

	result, err := c.Complete(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "ok", result.Text)
}
