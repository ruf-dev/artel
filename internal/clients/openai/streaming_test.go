package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// sseFixture is a fixed OpenAI-compatible SSE stream: one text delta, one tool call (fully
// accumulated in a single chunk for simplicity), a finish_reason chunk that flips the
// accumulator's state so the tool call is reported "just finished", and a final usage-only
// chunk (no choices), mirroring what stream_options.include_usage produces.
const sseFixture = `data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1,"model":"gpt-4o-mini","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello"},"finish_reason":null}]}

data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1,"model":"gpt-4o-mini","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_1","type":"function","function":{"name":"get_weather","arguments":"{\"loc\":\"NYC\"}"}}]},"finish_reason":null}]}

data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1,"model":"gpt-4o-mini","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}

data: {"id":"chatcmpl-1","object":"chat.completion.chunk","created":1,"model":"gpt-4o-mini","choices":[],"usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}

data: [DONE]

`

func TestStreamComplete_TextAndToolCallAndUsage(t *testing.T) {
	var capturedBody map[string]any

	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		err := json.NewDecoder(r.Body).Decode(&capturedBody)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte(sseFixture))
		require.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	client := New("test-api-key", server.URL)

	toolParams := json.RawMessage(`{"type":"object","properties":{"loc":{"type":"string"}},"required":["loc"]}`)
	req := StreamCompleteRequest{
		Model:  "gpt-4o-mini",
		System: "You are a weather bot.",
		Messages: []Message{
			{Role: "user", Content: "What's the weather in NYC?"},
		},
		Tools: []ToolDefinition{
			{Name: "get_weather", Description: "Get the weather", ParametersJSON: toolParams},
		},
	}

	seq, err := client.StreamComplete(context.Background(), req)
	require.NoError(t, err)

	var chunks []StreamChunk
	for chunk := range seq {
		chunks = append(chunks, chunk)
	}

	require.Len(t, chunks, 4)

	require.Equal(t, "Hello", chunks[0].TextDelta)
	require.Nil(t, chunks[0].ToolCallDelta)
	require.Nil(t, chunks[0].FinishedToolCall)
	require.False(t, chunks[0].Done)

	require.NotNil(t, chunks[1].ToolCallDelta)
	require.Equal(t, "call_1", chunks[1].ToolCallDelta.ID)
	require.Equal(t, "get_weather", chunks[1].ToolCallDelta.Name)
	require.Equal(t, `{"loc":"NYC"}`, chunks[1].ToolCallDelta.ArgumentsDelta)

	require.NotNil(t, chunks[2].FinishedToolCall)
	require.Equal(t, "call_1", chunks[2].FinishedToolCall.ID)
	require.Equal(t, "get_weather", chunks[2].FinishedToolCall.Name)
	require.Equal(t, `{"loc":"NYC"}`, chunks[2].FinishedToolCall.ArgumentsJSON)

	require.True(t, chunks[3].Done)
	require.NoError(t, chunks[3].Err)
	require.NotNil(t, chunks[3].FinalUsage)
	require.Equal(t, int64(10), chunks[3].FinalUsage.InputTokens)
	require.Equal(t, int64(5), chunks[3].FinalUsage.OutputTokens)

	// Assert the request was built correctly: system message, user message, and the tool
	// definition's schema round-tripped as a JSON-Schema object.
	messages, ok := capturedBody["messages"].([]any)
	require.True(t, ok)
	require.Len(t, messages, 2)

	systemMsg, ok := messages[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "system", systemMsg["role"])
	require.Equal(t, "You are a weather bot.", systemMsg["content"])

	userMsg, ok := messages[1].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "user", userMsg["role"])
	require.Equal(t, "What's the weather in NYC?", userMsg["content"])

	tools, ok := capturedBody["tools"].([]any)
	require.True(t, ok)
	require.Len(t, tools, 1)

	tool, ok := tools[0].(map[string]any)
	require.True(t, ok)

	function, ok := tool["function"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "get_weather", function["name"])
	require.Equal(t, "Get the weather", function["description"])

	parameters, ok := function["parameters"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "object", parameters["type"])

	require.Equal(t, true, capturedBody["stream"])
}

func TestStreamComplete_HistoryReplayIncludesToolMessages(t *testing.T) {
	var capturedBody map[string]any

	handler := func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&capturedBody)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("data: [DONE]\n\n"))
		require.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	client := New("test-api-key", server.URL)

	req := StreamCompleteRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{Role: "user", Content: "What's the weather?"},
			{
				Role: "assistant",
				ToolCalls: []ToolCall{
					{ID: "call_1", Name: "get_weather", ArgumentsJSON: `{"loc":"NYC"}`},
				},
			},
			{Role: "tool", Content: `{"tempF":72}`, ToolCallID: "call_1"},
		},
	}

	seq, err := client.StreamComplete(context.Background(), req)
	require.NoError(t, err)

	for range seq {
		// drain
	}

	messages, ok := capturedBody["messages"].([]any)
	require.True(t, ok)
	require.Len(t, messages, 3)

	assistantMsg, ok := messages[1].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "assistant", assistantMsg["role"])

	toolCalls, ok := assistantMsg["tool_calls"].([]any)
	require.True(t, ok)
	require.Len(t, toolCalls, 1)

	toolMsg, ok := messages[2].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "tool", toolMsg["role"])
	require.Equal(t, "call_1", toolMsg["tool_call_id"])
	require.Equal(t, `{"tempF":72}`, toolMsg["content"])
}

func TestStreamComplete_StreamErrorSurfacesOnFinalChunk(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`data: {"error":"boom"}` + "\n\n"))
		require.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	client := New("test-api-key", server.URL)

	req := StreamCompleteRequest{
		Model:    "gpt-4o-mini",
		Messages: []Message{{Role: "user", Content: "hi"}},
	}

	seq, err := client.StreamComplete(context.Background(), req)
	require.NoError(t, err)

	var chunks []StreamChunk
	for chunk := range seq {
		chunks = append(chunks, chunk)
	}

	require.Len(t, chunks, 1)
	require.True(t, chunks[0].Done)
	require.Error(t, chunks[0].Err)
}
