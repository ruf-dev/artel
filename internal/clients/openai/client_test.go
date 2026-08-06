package openai_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ruf-dev/artel/internal/clients/openai"
)

// TestListModels_Success verifies a models-list response is mapped into the returned
// []ModelInfo with the expected fields.
func TestListModels_Success(t *testing.T) {
	body := `{
		"object": "list",
		"data": [
			{
				"id": "gpt-4o",
				"object": "model",
				"created": 1715367049,
				"owned_by": "system"
			},
			{
				"id": "gpt-4o-mini",
				"object": "model",
				"created": 1721172741,
				"owned_by": "system"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := openai.New("test-api-key", server.URL)

	models, err := c.ListModels(context.Background())
	require.NoError(t, err)

	want := []openai.ModelInfo{
		{Id: "gpt-4o", OwnedBy: "system"},
		{Id: "gpt-4o-mini", OwnedBy: "system"},
	}
	require.Equal(t, want, models)
}

// TestListModels_AuthFailure verifies a 401 response from the fake server surfaces as an error.
func TestListModels_AuthFailure(t *testing.T) {
	body := `{
		"error": {
			"message": "Incorrect API key provided",
			"type": "invalid_request_error",
			"param": null,
			"code": "invalid_api_key"
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := openai.New("bad-api-key", server.URL)

	models, err := c.ListModels(context.Background())
	require.Error(t, err)
	require.Nil(t, models)

	code, ok := openai.StatusCode(err)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, code)
}

// TestComplete_Success verifies a chat-completion response is mapped into CompleteResult,
// including the usage object fields.
func TestComplete_Success(t *testing.T) {
	body := `{
		"id": "chatcmpl-abc123",
		"object": "chat.completion",
		"created": 1715367049,
		"model": "gpt-4o",
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "Hello, world!"
				},
				"logprobs": null,
				"finish_reason": "stop"
			}
		],
		"usage": {
			"prompt_tokens": 10,
			"completion_tokens": 5,
			"total_tokens": 15
		}
	}`

	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/chat/completions", r.URL.Path)

		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		gotBody = buf

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := openai.New("test-api-key", server.URL)

	req := openai.CompleteRequest{
		Model:        "gpt-4o",
		SystemPrompt: "You are a helpful assistant.",
		Prompt:       "Say hello",
		MaxTokens:    1024,
	}

	result, err := c.Complete(context.Background(), req)
	require.NoError(t, err)
	require.NotEmpty(t, gotBody)

	want := openai.CompleteResult{
		Text: "Hello, world!",
		Usage: openai.Usage{
			InputTokens:  10,
			OutputTokens: 5,
		},
	}
	require.Equal(t, want, result)
}

// TestComplete_DefaultMaxTokens verifies that omitting MaxTokens still produces a successful
// request (the client applies its own default rather than sending max_completion_tokens: 0).
func TestComplete_DefaultMaxTokens(t *testing.T) {
	body := `{
		"id": "chatcmpl-def456",
		"object": "chat.completion",
		"created": 1715367049,
		"model": "gpt-4o",
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "ok"
				},
				"logprobs": null,
				"finish_reason": "stop"
			}
		],
		"usage": {
			"prompt_tokens": 1,
			"completion_tokens": 1,
			"total_tokens": 2
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := openai.New("test-api-key", server.URL)

	req := openai.CompleteRequest{
		Model:  "gpt-4o",
		Prompt: "Say ok",
	}

	result, err := c.Complete(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "ok", result.Text)
}

// TestComplete_Failure verifies an error response from the fake server surfaces as an error,
// and that the status code is recoverable via StatusCode.
func TestComplete_Failure(t *testing.T) {
	body := `{
		"error": {
			"message": "The model ` + "`gpt-bogus`" + ` does not exist",
			"type": "invalid_request_error",
			"param": null,
			"code": "model_not_found"
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	defer server.Close()

	c := openai.New("test-api-key", server.URL)

	req := openai.CompleteRequest{
		Model:  "gpt-bogus",
		Prompt: "hi",
	}

	result, err := c.Complete(context.Background(), req)
	require.Error(t, err)
	require.Equal(t, openai.CompleteResult{}, result)

	code, ok := openai.StatusCode(err)
	require.True(t, ok)
	require.Equal(t, http.StatusNotFound, code)
}
