// Package claudecli drives the `claude` CLI in headless (`--print`) mode and translates its
// stream-json output into chatprotocol events.
//
// The shapes handled here were confirmed by capturing real output of
//
//	claude -p "..." --output-format stream-json --verbose --include-partial-messages
//
// from CLI version 2.1.234. Every line is a self-contained JSON object; the ones that matter are:
//
//	{"type":"system","subtype":"init",...}                          first line of a turn
//	{"type":"stream_event","event":{"type":"message_start",...}}    carries the assistant message id
//	{"type":"stream_event","event":{"type":"content_block_delta","delta":{"type":"text_delta",...}}}
//	{"type":"assistant","message":{"content":[...]}}                complete blocks (text and tool_use)
//	{"type":"user","message":{"content":[{"type":"tool_result",...}]}}
//	{"type":"result","session_id":"...","total_cost_usd":...}       last line of a turn
//
// Text is streamed from the `content_block_delta` lines and *completed* from the following
// `assistant` line, rather than reconstructed by concatenating deltas: the assistant line carries
// the authoritative final text, so a dropped or malformed delta can never leave a consumer with a
// silently truncated message.
//
// Tool calls are taken from the `assistant` line rather than from `content_block_start` +
// `input_json_delta`, because only the assistant line carries the tool's fully-assembled input;
// the partial-json deltas would have to be reassembled by hand for no benefit, since a tool call
// isn't actionable until its input is complete anyway.
package claudecli

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"workbenchbridge/internal/chatprotocol"
)

// inputSummaryLimit / outputLimit bound what a single event carries. Tool inputs are summarised
// for display (a consumer that wants the exact input reads Event.Input, which is never
// truncated); tool *output* has no separate untruncated field, so its limit is set generously —
// enough for a realistic file read or command output, short of letting one runaway tool result
// blow out the hub's replay backlog inside a memory-limited container.
const (
	inputSummaryLimit = 400
	outputLimit       = 16384
)

// Translator converts stream-json lines into chatprotocol events. It is stateful only in the
// small amount needed to correlate deltas with their completed block (the current assistant
// message id), so one Translator is used per turn.
//
// Not safe for concurrent use; the bridge runs exactly one turn at a time.
type Translator struct {
	// messageId is the id of the assistant message currently streaming, from the most recent
	// message_start. Combined with a content block's index it forms the correlation id shared by
	// assistant_text_delta and assistant_text_done, which is what lets a consumer replace an
	// accumulating bubble with the final text rather than appending it twice.
	messageId string
}

// NewTranslator returns a Translator for one turn.
func NewTranslator() *Translator {
	return &Translator{}
}

// streamLine is the outer envelope every stream-json line shares.
type streamLine struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`

	Event   json.RawMessage `json:"event"`
	Message json.RawMessage `json:"message"`

	SessionId    string  `json:"session_id"`
	IsError      bool    `json:"is_error"`
	TotalCostUsd float64 `json:"total_cost_usd"`
	Result       string  `json:"result"`
}

// streamEvent is the raw Messages-API streaming event nested under a "stream_event" line.
type streamEvent struct {
	Type  string `json:"type"`
	Index int    `json:"index"`

	Message *struct {
		Id string `json:"id"`
	} `json:"message"`

	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

// apiMessage is the complete assistant/user message carried by an "assistant"/"user" line.
type apiMessage struct {
	Id      string        `json:"id"`
	Content []contentPart `json:"content"`
}

// contentPart is one content block. The four block types that matter carry disjoint fields, so
// they share one struct rather than being modelled as a discriminated union — Type says which
// fields are populated.
type contentPart struct {
	Type string `json:"type"`

	// text
	Text string `json:"text"`

	// tool_use
	Id    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input"`

	// tool_result
	ToolUseId string          `json:"tool_use_id"`
	Content   json.RawMessage `json:"content"`
	IsError   bool            `json:"is_error"`
}

// Translate parses one stream-json line and returns the chatprotocol events it produces — zero
// for the many lines (rate_limit_event, system/status, message_stop, ...) that carry nothing a
// chat consumer needs.
//
// A line that isn't valid JSON produces a single error event rather than an error return: the
// `claude` CLI occasionally interleaves plain-text warnings into stdout, and one such line must
// not abort a turn that is otherwise streaming fine.
func (t *Translator) Translate(line []byte) []chatprotocol.Event {
	trimmed := strings.TrimSpace(string(line))
	if trimmed == "" {
		return nil
	}

	var parsed streamLine

	err := json.Unmarshal([]byte(trimmed), &parsed)
	if err != nil {
		event := chatprotocol.Event{
			Type: chatprotocol.EventError,
			Text: "unparseable claude output: " + truncate(trimmed, inputSummaryLimit),
		}

		return []chatprotocol.Event{event}
	}

	switch parsed.Type {
	case "system":
		return t.translateSystem(parsed)
	case "stream_event":
		return t.translateStreamEvent(parsed)
	case "assistant":
		return t.translateAssistant(parsed)
	case "user":
		return t.translateUser(parsed)
	case "result":
		return t.translateResult(parsed)
	default:
		return nil
	}
}

// translateSystem emits system_init for the init line that opens every turn. Other system
// subtypes ("status", ...) are progress noise with no chat-level meaning.
func (t *Translator) translateSystem(parsed streamLine) []chatprotocol.Event {
	if parsed.Subtype != "init" {
		return nil
	}

	event := chatprotocol.Event{
		Type: chatprotocol.EventSystemInit,
	}

	return []chatprotocol.Event{event}
}

// translateStreamEvent handles the raw streaming events: it latches the current message id from
// message_start and turns each text_delta into an assistant_text_delta. input_json_delta (the
// streaming form of a tool call's input) is deliberately ignored — see the package doc.
func (t *Translator) translateStreamEvent(parsed streamLine) []chatprotocol.Event {
	var event streamEvent

	err := json.Unmarshal(parsed.Event, &event)
	if err != nil {
		return nil
	}

	if event.Type == "message_start" && event.Message != nil {
		t.messageId = event.Message.Id

		return nil
	}

	if event.Type != "content_block_delta" || event.Delta == nil {
		return nil
	}

	if event.Delta.Type != "text_delta" || event.Delta.Text == "" {
		return nil
	}

	delta := chatprotocol.Event{
		Type: chatprotocol.EventAssistantTextDelta,
		ID:   blockId(t.messageId, event.Index),
		Text: event.Delta.Text,
	}

	return []chatprotocol.Event{delta}
}

// translateAssistant turns a completed assistant message into one assistant_text_done per text
// block and one tool_call_started per tool_use block, in the blocks' own order.
func (t *Translator) translateAssistant(parsed streamLine) []chatprotocol.Event {
	var message apiMessage

	err := json.Unmarshal(parsed.Message, &message)
	if err != nil {
		return nil
	}

	if message.Id != "" {
		t.messageId = message.Id
	}

	events := make([]chatprotocol.Event, 0, len(message.Content))

	for index, part := range message.Content {
		switch part.Type {
		case "text":
			if part.Text == "" {
				continue
			}

			done := chatprotocol.Event{
				Type: chatprotocol.EventAssistantTextDone,
				ID:   blockId(message.Id, index),
				Text: part.Text,
			}

			events = append(events, done)
		case "tool_use":
			started := chatprotocol.Event{
				Type:         chatprotocol.EventToolCallStarted,
				ID:           part.Id,
				ToolName:     part.Name,
				Input:        part.Input,
				InputSummary: summariseInput(part.Input),
			}

			events = append(events, started)
		}
	}

	return events
}

// translateUser turns the tool_result blocks of the synthetic user message Claude Code emits
// after running a tool into tool_call_result events, correlated back to tool_call_started by the
// tool_use id.
//
// A user line with no tool_result blocks is the *actual* user's prompt being echoed back into the
// transcript, which the consumer already knows about — it sent it — so it produces nothing.
func (t *Translator) translateUser(parsed streamLine) []chatprotocol.Event {
	var message apiMessage

	err := json.Unmarshal(parsed.Message, &message)
	if err != nil {
		return nil
	}

	var events []chatprotocol.Event

	for _, part := range message.Content {
		if part.Type != "tool_result" {
			continue
		}

		result := chatprotocol.Event{
			Type:    chatprotocol.EventToolCallResult,
			ID:      part.ToolUseId,
			Output:  truncate(toolResultText(part.Content), outputLimit),
			IsError: part.IsError,
		}

		events = append(events, result)
	}

	return events
}

// translateResult closes a turn: turn_done always, preceded by an error event when the CLI itself
// reported the turn as failed.
func (t *Translator) translateResult(parsed streamLine) []chatprotocol.Event {
	var events []chatprotocol.Event

	if parsed.IsError {
		text := parsed.Result
		if text == "" {
			text = "claude reported the turn as failed (" + parsed.Subtype + ")"
		}

		failure := chatprotocol.Event{
			Type: chatprotocol.EventError,
			Text: truncate(text, outputLimit),
		}

		events = append(events, failure)
	}

	done := chatprotocol.Event{
		Type:      chatprotocol.EventTurnDone,
		SessionID: parsed.SessionId,
		CostUSD:   parsed.TotalCostUsd,
	}

	return append(events, done)
}

// blockId builds the correlation id shared by a text block's deltas and its completion. An
// unknown message id (a delta arriving before any message_start, which shouldn't happen but
// costs nothing to tolerate) still yields a stable per-index id.
func blockId(messageId string, index int) string {
	return messageId + ":" + strconv.Itoa(index)
}

// toolResultText renders a tool_result block's content, which the API models either as a plain
// string or as an array of content blocks (a tool returning images alongside text, for instance).
// Anything that is neither is passed through as its raw JSON, which is strictly more useful to a
// human reading the chat than an empty string.
func toolResultText(content json.RawMessage) string {
	if len(content) == 0 {
		return ""
	}

	var asString string

	err := json.Unmarshal(content, &asString)
	if err == nil {
		return asString
	}

	var asParts []contentPart

	err = json.Unmarshal(content, &asParts)
	if err != nil {
		return string(content)
	}

	var builder strings.Builder

	for _, part := range asParts {
		if part.Text == "" {
			continue
		}

		if builder.Len() > 0 {
			builder.WriteString("\n")
		}

		builder.WriteString(part.Text)
	}

	return builder.String()
}

// summariseInput renders a tool call's input as a single compact line for display next to the
// tool name. The untruncated input travels alongside it in Event.Input, so this is free to be
// aggressive about length.
func summariseInput(input json.RawMessage) string {
	if len(input) == 0 {
		return ""
	}

	var compacted bytes.Buffer

	err := json.Compact(&compacted, input)
	if err != nil {
		return truncate(string(input), inputSummaryLimit)
	}

	return truncate(compacted.String(), inputSummaryLimit)
}

// truncate shortens s to at most limit runes, marking that it did so. Rune-based rather than
// byte-based so a cut never lands mid-UTF-8-sequence and produces invalid JSON on the wire.
func truncate(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}

	return string(runes[:limit]) + "… (truncated)"
}
