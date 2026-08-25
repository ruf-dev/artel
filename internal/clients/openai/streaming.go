package openai

import (
	"context"
	"encoding/json"
	"iter"

	sdk "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"

	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

// Message is one turn of conversation history for StreamComplete. The caller (the Simple Chat
// engine) owns the full message history across turns — StreamComplete is a single, synchronous
// per-turn call, not a long-lived session.
type Message struct {
	// Role is one of "system" | "user" | "assistant" | "tool".
	Role string
	// Content is the message text. For Role == "assistant" it may be empty when the turn was
	// pure tool calls.
	Content string
	// ToolCallID is set only for Role == "tool" — the id of the ToolCall this message answers.
	ToolCallID string
	// ToolCalls is set only for an assistant message that made tool calls, so it can be replayed
	// back to the model as history on a later turn.
	ToolCalls []ToolCall
}

// ToolCall is one function call the model requested — either a finished call replayed as
// history (see Message.ToolCalls) or the result of a StreamChunk.FinishedToolCall.
type ToolCall struct {
	ID            string
	Name          string
	ArgumentsJSON string
}

// ToolDefinition describes one tool the model may call. ParametersJSON is a pre-built
// JSON-Schema object; converting domain.ToolApiDescription/domain.ToolProperty into this shape
// is the caller's responsibility, not this package's.
type ToolDefinition struct {
	Name           string
	Description    string
	ParametersJSON json.RawMessage
}

// ToolCallDelta is an incremental fragment of a tool call still being streamed — the raw,
// index-keyed accumulation openai-go exposes per chunk, translated to the tool call it belongs
// to. ID and Name are only set on the first delta for a given tool call; ArgumentsDelta is the
// incremental JSON-arguments fragment carried by this specific chunk. Most callers only need
// StreamChunk.FinishedToolCall; ToolCallDelta exists for a caller that wants to reflect
// in-progress tool calls live (e.g. "assistant is calling web_search...") before it completes.
type ToolCallDelta struct {
	Index          int64
	ID             string
	Name           string
	ArgumentsDelta string
}

// StreamChunk is one event yielded by StreamComplete's iterator. Only the fields relevant to
// what happened in this chunk are populated:
//   - TextDelta: new assistant text — forward directly to a websocket as it arrives.
//   - ToolCallDelta: an incremental fragment of a tool call still being accumulated.
//   - FinishedToolCall: a tool call that just finished accumulating — the caller learns about a
//     completed call as soon as this fires, without tracking openai-go's index-keyed state
//     itself.
//   - Done: true on the final chunk of the stream (success or failure — check Err).
//   - FinalUsage: token usage for the whole turn (see the Usage type in client.go), populated
//     alongside Done when the provider returns it.
//   - Err: set on the final chunk when the stream failed (network error, API error, or ctx
//     cancellation). The iterator always stops after a chunk with Err set.
type StreamChunk struct {
	TextDelta        string
	ToolCallDelta    *ToolCallDelta
	FinishedToolCall *ToolCall
	Done             bool
	FinalUsage       *Usage
	Err              error
}

// StreamCompleteRequest is a single-turn, streaming, tool-calling-capable completion request.
type StreamCompleteRequest struct {
	Model     string
	System    string
	Messages  []Message
	Tools     []ToolDefinition
	MaxTokens int64
}

// StreamComplete performs one streaming turn against the Chat Completions API and returns a
// range-over-func iterator (iter.Seq[StreamChunk], go1.23+) of the chunks produced. This
// package has no other established async idiom to match (no channel-based streaming precedent
// elsewhere in this codebase, and the module already targets go1.25 — see go.mod), so iter.Seq
// was chosen as the more idiomatic, allocation-light fit for a synchronous "range over the
// stream, forward to a websocket" consumption pattern: the caller does
//
//	for chunk := range seq { ... }
//
// and a `break` mid-loop cleanly tears down the underlying SSE stream (the yield function
// returning false stops iteration and triggers cleanup), which a raw channel would need extra
// plumbing (a done channel, or draining) to achieve safely.
//
// Iteration never panics; a mid-stream failure is surfaced as a final chunk with Err set
// (Done is also true on that chunk) rather than as StreamComplete's own return error — that
// return error is reserved for request-construction failures that happen before any network
// call (there are none today, but the signature leaves room for future validation).
func (c *Client) StreamComplete(ctx context.Context, req StreamCompleteRequest) (iter.Seq[StreamChunk], error) {
	params := buildStreamParams(req)

	seq := func(yield func(StreamChunk) bool) {
		stream := c.sdk.Chat.Completions.NewStreaming(ctx, params)
		defer utils.CloseWithLog(stream, "openai chat completion stream")

		var acc sdk.ChatCompletionAccumulator

		for stream.Next() {
			chunk := stream.Current()

			acc.AddChunk(chunk)

			if !emitChunkEvents(chunk, &acc, yield) {
				return
			}
		}

		err := stream.Err()
		if err != nil {
			final := StreamChunk{
				Done: true,
				Err:  rerrors.Wrap(err, "error streaming openai chat completion"),
			}
			yield(final)

			return
		}

		final := StreamChunk{Done: true}

		if acc.Usage.PromptTokens != 0 || acc.Usage.CompletionTokens != 0 {
			usage := Usage{
				InputTokens:  acc.Usage.PromptTokens,
				OutputTokens: acc.Usage.CompletionTokens,
			}
			final.FinalUsage = &usage
		}

		yield(final)
	}

	return seq, nil
}

// emitChunkEvents translates one raw ChatCompletionChunk (plus the accumulator's just-updated
// state) into zero or more StreamChunk yields. Returns false as soon as yield asks to stop.
func emitChunkEvents(chunk sdk.ChatCompletionChunk, acc *sdk.ChatCompletionAccumulator, yield func(StreamChunk) bool) bool {
	if len(chunk.Choices) == 0 {
		return true
	}

	delta := chunk.Choices[0].Delta

	if delta.Content != "" {
		textChunk := StreamChunk{TextDelta: delta.Content}
		if !yield(textChunk) {
			return false
		}
	}

	for _, deltaTool := range delta.ToolCalls {
		toolDelta := ToolCallDelta{
			Index:          deltaTool.Index,
			ID:             deltaTool.ID,
			Name:           deltaTool.Function.Name,
			ArgumentsDelta: deltaTool.Function.Arguments,
		}
		toolChunk := StreamChunk{ToolCallDelta: &toolDelta}

		if !yield(toolChunk) {
			return false
		}
	}

	finished, ok := acc.JustFinishedToolCall()
	if ok {
		toolCall := ToolCall{
			ID:            finished.ID,
			Name:          finished.Name,
			ArgumentsJSON: finished.Arguments,
		}
		finishedChunk := StreamChunk{FinishedToolCall: &toolCall}

		if !yield(finishedChunk) {
			return false
		}
	}

	return true
}

// buildStreamParams maps a StreamCompleteRequest onto the SDK's request params, including
// enabling the usage-accounting final chunk (see ChatCompletionStreamOptionsParam) so
// StreamChunk.FinalUsage can be populated.
func buildStreamParams(req StreamCompleteRequest) sdk.ChatCompletionNewParams {
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	messages := make([]sdk.ChatCompletionMessageParamUnion, 0, len(req.Messages)+1)

	if req.System != "" {
		messages = append(messages, sdk.SystemMessage(req.System))
	}

	for _, msg := range req.Messages {
		messages = append(messages, toSDKMessage(msg))
	}

	tools := make([]sdk.ChatCompletionToolParam, 0, len(req.Tools))
	for _, toolDef := range req.Tools {
		tools = append(tools, toSDKTool(toolDef))
	}

	streamOptions := sdk.ChatCompletionStreamOptionsParam{
		IncludeUsage: sdk.Bool(true),
	}

	params := sdk.ChatCompletionNewParams{
		Model:               req.Model,
		Messages:            messages,
		MaxCompletionTokens: sdk.Int(maxTokens),
		StreamOptions:       streamOptions,
	}

	if len(tools) > 0 {
		params.Tools = tools
	}

	return params
}

// toSDKMessage maps a Message onto the SDK's discriminated message-union constructors.
func toSDKMessage(msg Message) sdk.ChatCompletionMessageParamUnion {
	switch msg.Role {
	case "user":
		return sdk.UserMessage(msg.Content)
	case "system":
		return sdk.SystemMessage(msg.Content)
	case "tool":
		return sdk.ToolMessage(msg.Content, msg.ToolCallID)
	case "assistant":
		return toSDKAssistantMessage(msg)
	default:
		return sdk.UserMessage(msg.Content)
	}
}

// toSDKAssistantMessage builds an assistant history message, attaching ToolCalls (if any) so a
// later turn can replay them back to the model per the Chat Completions API's requirement that
// a tool result message be preceded by the assistant message that requested it.
func toSDKAssistantMessage(msg Message) sdk.ChatCompletionMessageParamUnion {
	assistantMsg := sdk.AssistantMessage(msg.Content)

	if len(msg.ToolCalls) == 0 {
		return assistantMsg
	}

	toolCallParams := make([]sdk.ChatCompletionMessageToolCallParam, 0, len(msg.ToolCalls))
	for _, tc := range msg.ToolCalls {
		function := sdk.ChatCompletionMessageToolCallFunctionParam{
			Name:      tc.Name,
			Arguments: tc.ArgumentsJSON,
		}
		toolCallParam := sdk.ChatCompletionMessageToolCallParam{
			ID:       tc.ID,
			Function: function,
		}
		toolCallParams = append(toolCallParams, toolCallParam)
	}

	assistantMsg.OfAssistant.ToolCalls = toolCallParams

	return assistantMsg
}

// toSDKTool maps a ToolDefinition onto the SDK's function-tool param, parsing ParametersJSON
// into the map[string]any shape shared.FunctionParameters requires.
func toSDKTool(def ToolDefinition) sdk.ChatCompletionToolParam {
	var parameters shared.FunctionParameters

	if len(def.ParametersJSON) > 0 {
		_ = json.Unmarshal(def.ParametersJSON, &parameters)
	}

	function := shared.FunctionDefinitionParam{
		Name:        def.Name,
		Description: sdk.String(def.Description),
		Parameters:  parameters,
	}

	return sdk.ChatCompletionToolParam{
		Function: function,
	}
}
