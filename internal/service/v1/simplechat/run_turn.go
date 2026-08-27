package simplechat

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ruf-dev/artel/internal/chatprotocol"
	"github.com/ruf-dev/artel/internal/clients/openai"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

// maxToolIterations caps how many consecutive tool-calling rounds one turn may run before it is
// abandoned — see user_errors.SimpleChatToolLoopExceeded.
const maxToolIterations = 25

// inputSummaryLimit is how much of a tool call's argument JSON is echoed into an event's
// InputSummary, which the UI renders as a one-line preview beside the tool name.
const inputSummaryLimit = 160

// deniedToolResult is fed back to the model as the tool's result when the user denies a call, so
// the model can acknowledge the refusal and choose another route instead of silently stalling.
const deniedToolResult = "User denied this tool call."

// streamResult is one completed StreamComplete pass: the assistant text it produced plus any
// tool calls it asked for. StreamComplete exposes no finish_reason, so "this turn is finished"
// is derived here from an empty toolCalls slice on a stream that ended without error.
type streamResult struct {
	text      string
	toolCalls []openai.ToolCall
}

// RunTurn drives one full agent turn: it persists the user's message, then loops
// StreamComplete → tool calls → tool results until the model answers with text and no further
// tool calls. It is synchronous and single-goroutine; the websocket handler runs it on its own
// goroutine per turn so its read loop stays free to deliver permission decisions through sink.
//
// model overrides the thread's stored default for this turn only (the "switchable mid-
// conversation" path); an empty model falls back to chat.Model.
//
// Cancellation: ctx is cancelled when the client disconnects, which tears down the in-flight
// stream. Persistence deliberately runs on a context.WithoutCancel copy of ctx, so a turn cut
// off mid-stream still keeps the messages it had already completed rather than losing them to
// the cancelled context. This is best-effort rather than transactional: a turn interrupted
// between an assistant row and its tool-result rows leaves the partial prefix persisted, which
// replays as valid history because tool rows always carry their own tool_call_id.
func (s *Service) RunTurn(
	ctx context.Context, chat domain.SimpleChat, userText, model string, sink chatprotocol.EventSink,
) error {
	// Writes outlive ctx so a disconnect mid-turn doesn't discard completed work.
	persistCtx := context.WithoutCancel(ctx)

	if model == "" {
		model = chat.Model
	}

	creds, err := s.resolveOpenRouterCredentials(ctx, chat.UserUuid)
	if err != nil {
		return failTurn(sink, rerrors.Wrap(err, "error resolving openrouter credentials"))
	}

	client := openai.New(creds.ApiKey, creds.BaseUrl)

	file, err := s.currentChatFile(ctx, chat)
	if err != nil {
		return failTurn(sink, rerrors.Wrap(err, "error reading simple chat file"))
	}

	rows := file.Messages

	userRow := domain.SimpleChatMessage{
		Role:    string(domain.SimpleChatRoleUser),
		Content: userText,
	}

	err = s.appendMessage(persistCtx, chat, userRow)
	if err != nil {
		return failTurn(sink, rerrors.Wrap(err, "error inserting user message"))
	}

	tools, err := s.resolveTools(ctx, chat)
	if err != nil {
		return failTurn(sink, rerrors.Wrap(err, "error resolving tools"))
	}

	systemSettings, err := s.systemSettings.Get(ctx)
	if err != nil {
		return failTurn(sink, rerrors.Wrap(err, "error getting system prompt"))
	}
	userSettings, err := s.userSettings.Get(ctx, chat.UserUuid)
	if err != nil {
		return failTurn(sink, rerrors.Wrap(err, "error getting user prompt"))
	}
	vault, err := s.vaults.GetByID(ctx, chat.VaultUuid)
	if err != nil {
		return failTurn(sink, rerrors.Wrap(err, "error getting vault prompt"))
	}
	system := mergePrompts(systemSettings.SystemPrompt, userSettings.UserPrompt, vault.Prompt, vault.UseSystemPrompt)

	messages := historyToMessages(rows)
	userMsg := openai.Message{Role: "user", Content: userText}
	messages = append(messages, userMsg)

	for iteration := 0; iteration < maxToolIterations; iteration++ {
		result, streamErr := s.streamOnce(ctx, client, model, system, messages, tools, sink)
		if streamErr != nil {
			return failTurn(sink, streamErr)
		}

		if result.text != "" {
			assistantRow := domain.SimpleChatMessage{
				Role:    string(domain.SimpleChatRoleAssistant),
				Content: result.text,
				Model:   &model,
			}

			err = s.appendMessage(persistCtx, chat, assistantRow)
			if err != nil {
				return failTurn(sink, rerrors.Wrap(err, "error inserting assistant message"))
			}

			textMsg := openai.Message{Role: "assistant", Content: result.text}
			messages = append(messages, textMsg)
		}

		if len(result.toolCalls) == 0 {
			return s.finishTurn(persistCtx, chat, model, sink)
		}

		// The Chat Completions API requires the assistant message that requested the calls to
		// precede their results; historyToMessages rebuilds this same shape on a later turn.
		callMsg := openai.Message{Role: "assistant", ToolCalls: result.toolCalls}
		messages = append(messages, callMsg)

		for _, call := range result.toolCalls {
			toolMsg, toolErr := s.runToolCall(ctx, persistCtx, chat, call, sink)
			if toolErr != nil {
				return failTurn(sink, toolErr)
			}

			messages = append(messages, toolMsg)
		}
	}

	return failTurn(sink, rerrors.Wrap(user_errors.SimpleChatToolLoopExceeded, chat.Uuid.String()))
}

// finishTurn bumps the thread's activity timestamp and emits turn_done.
func (s *Service) finishTurn(
	persistCtx context.Context, chat domain.SimpleChat, model string, sink chatprotocol.EventSink,
) error {
	bumpFn := func(file *domain.SimpleChatFile) {
		now := time.Now().UTC()
		file.Header.LastActivityAt = now
		file.Header.UpdatedAt = now
	}

	err := s.mutateChatFile(persistCtx, chat, bumpFn)
	if err != nil {
		log.Warn().Err(err).Str("chat_id", chat.Uuid.String()).
			Msg("error updating simple chat last activity")
	}

	doneEvent := chatprotocol.Event{
		Type:      chatprotocol.EventTurnDone,
		SessionID: chat.Uuid.String(),
		Model:     model,
	}

	sendErr := sink.Send(doneEvent)
	if sendErr != nil {
		return rerrors.Wrap(sendErr, "error sending turn done")
	}

	return nil
}

// streamOnce runs a single StreamComplete pass, forwarding text deltas to the client as they
// arrive and collecting the tool calls the model asked for.
//
// Every delta and the closing assistant_text_done carry the same generated message id: the web
// client folds assistant text by event id, appending deltas and replacing the accumulated text
// on done — without a stable id each delta would render as its own message bubble.
func (s *Service) streamOnce(
	ctx context.Context,
	client *openai.Client,
	model string,
	system string,
	messages []openai.Message,
	tools []openai.ToolDefinition,
	sink chatprotocol.EventSink,
) (streamResult, error) {
	req := openai.StreamCompleteRequest{
		Model:    model,
		System:   system,
		Messages: messages,
		Tools:    tools,
	}

	chunks, err := client.StreamComplete(ctx, req)
	if err != nil {
		return streamResult{}, rerrors.Wrap(err, "error starting openai stream")
	}

	messageId := uuid.NewString()

	var builder strings.Builder
	var toolCalls []openai.ToolCall
	var loopErr error

	for chunk := range chunks {
		if chunk.Err != nil {
			loopErr = rerrors.Wrap(chunk.Err, "error streaming openai completion")

			break
		}

		if chunk.TextDelta != "" {
			builder.WriteString(chunk.TextDelta)

			deltaEvent := chatprotocol.Event{
				Type:  chatprotocol.EventAssistantTextDelta,
				ID:    messageId,
				Text:  chunk.TextDelta,
				Model: model,
			}

			sendErr := sink.Send(deltaEvent)
			if sendErr != nil {
				loopErr = rerrors.Wrap(sendErr, "error sending assistant text delta")

				break
			}
		}

		if chunk.FinishedToolCall != nil {
			toolCalls = append(toolCalls, *chunk.FinishedToolCall)
		}
	}

	if loopErr != nil {
		return streamResult{}, loopErr
	}

	text := builder.String()

	if text != "" {
		doneEvent := chatprotocol.Event{
			Type:  chatprotocol.EventAssistantTextDone,
			ID:    messageId,
			Text:  text,
			Model: model,
		}

		err = sink.Send(doneEvent)
		if err != nil {
			return streamResult{}, rerrors.Wrap(err, "error sending assistant text done")
		}
	}

	result := streamResult{
		text:      text,
		toolCalls: toolCalls,
	}

	return result, nil
}

// runToolCall announces one tool call, resolves permission for it, executes it (unless denied)
// and persists the outcome. It returns the tool-result message to feed back to the model.
//
// A denied call is still persisted and still answered — the model is told the user refused,
// rather than being left waiting on a result that never comes.
func (s *Service) runToolCall(
	ctx context.Context,
	persistCtx context.Context,
	chat domain.SimpleChat,
	call openai.ToolCall,
	sink chatprotocol.EventSink,
) (openai.Message, error) {
	input := rawArguments(call.ArgumentsJSON)

	startedEvent := chatprotocol.Event{
		Type:         chatprotocol.EventToolCallStarted,
		ID:           call.ID,
		ToolName:     call.Name,
		Input:        input,
		InputSummary: summarize(call.ArgumentsJSON),
	}

	err := sink.Send(startedEvent)
	if err != nil {
		return openai.Message{}, rerrors.Wrap(err, "error sending tool call started")
	}

	decision, err := s.resolvePermission(ctx, persistCtx, chat, call, input, sink)
	if err != nil {
		return openai.Message{}, rerrors.Wrap(err, "error resolving tool permission")
	}

	output := deniedToolResult
	isError := true

	if decision != chatprotocol.DecisionDeny {
		output, isError = s.executeTool(ctx, chat, call)
	}

	resultEvent := chatprotocol.Event{
		Type:     chatprotocol.EventToolCallResult,
		ID:       call.ID,
		ToolName: call.Name,
		Output:   output,
		IsError:  isError,
	}

	err = sink.Send(resultEvent)
	if err != nil {
		return openai.Message{}, rerrors.Wrap(err, "error sending tool call result")
	}

	callId := call.ID
	toolName := call.Name

	toolRow := domain.SimpleChatMessage{
		Role:       string(domain.SimpleChatRoleTool),
		Content:    output,
		ToolCallID: &callId,
		ToolName:   &toolName,
		ToolInput:  input,
		IsError:    isError,
	}

	err = s.appendMessage(persistCtx, chat, toolRow)
	if err != nil {
		return openai.Message{}, rerrors.Wrap(err, "error inserting tool message")
	}

	toolMsg := openai.Message{
		Role:       "tool",
		Content:    output,
		ToolCallID: call.ID,
	}

	return toolMsg, nil
}

// executeTool runs the tool through the MCP service as the chat's owner. A failure is reported
// back to the model as the tool's (error) result rather than aborting the turn — the model can
// often recover, e.g. by correcting an argument.
func (s *Service) executeTool(
	ctx context.Context, chat domain.SimpleChat, call openai.ToolCall,
) (string, bool) {
	params := map[string]interface{}{}

	if call.ArgumentsJSON != "" {
		err := json.Unmarshal([]byte(call.ArgumentsJSON), &params)
		if err != nil {
			return rerrors.Wrap(err, "error parsing tool arguments").Error(), true
		}
	}

	keyCtx, err := s.mcp.BuildKeyContext(ctx, chat.VaultUuid, chat.UserUuid)
	if err != nil {
		return rerrors.Wrap(err, "error building mcp key context").Error(), true
	}

	result, err := s.mcp.ExecuteTool(ctx, keyCtx, call.Name, params)
	if err != nil {
		return err.Error(), true
	}

	return result.Text, false
}

// resolvePermission decides whether call may run: a remembered allow_always for this
// (chat, tool) skips the prompt entirely, otherwise the user is asked and the answer awaited.
// An allow_always answer is remembered for the rest of the thread.
func (s *Service) resolvePermission(
	ctx context.Context,
	persistCtx context.Context,
	chat domain.SimpleChat,
	call openai.ToolCall,
	input json.RawMessage,
	sink chatprotocol.EventSink,
) (chatprotocol.PermissionDecision, error) {
	storedDecision, ok, err := s.getToolAllowance(ctx, chat, call.Name)
	if err != nil {
		return "", rerrors.Wrap(err, "error getting simple chat tool allowance")
	}

	if ok && storedDecision == string(chatprotocol.DecisionAllowAlways) {
		return chatprotocol.DecisionAllowAlways, nil
	}

	options := []chatprotocol.PermissionDecision{
		chatprotocol.DecisionAllowOnce,
		chatprotocol.DecisionAllowAlways,
		chatprotocol.DecisionDeny,
	}

	requestEvent := chatprotocol.Event{
		Type:         chatprotocol.EventPermissionRequest,
		ID:           call.ID,
		ToolName:     call.Name,
		Input:        input,
		InputSummary: summarize(call.ArgumentsJSON),
		Options:      options,
	}

	err = sink.Send(requestEvent)
	if err != nil {
		return "", rerrors.Wrap(err, "error sending permission request")
	}

	decision, err := sink.AwaitPermissionDecision(ctx, call.ID)
	if err != nil {
		return "", rerrors.Wrap(err, "error awaiting permission decision")
	}

	if decision == chatprotocol.DecisionAllowAlways {
		err = s.upsertToolAllowance(persistCtx, chat, call.Name, string(decision))
		if err != nil {
			return "", rerrors.Wrap(err, "error upserting simple chat tool allowance")
		}
	}

	return decision, nil
}

// resolveTools lists the MCP tool catalogue once and converts the portion this chat may use.
// A chat created without vault access never sees the vault-note tools.
func (s *Service) resolveTools(
	ctx context.Context, chat domain.SimpleChat,
) ([]openai.ToolDefinition, error) {
	tools, err := s.mcp.ListTools(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing mcp tools")
	}

	defs, err := buildToolDefinitions(tools, chat.VaultAccess)
	if err != nil {
		return nil, rerrors.Wrap(err, "error building tool definitions")
	}

	return defs, nil
}

// historyToMessages replays persisted rows as OpenAI chat history. A consecutive run of tool
// rows is emitted as one assistant message carrying all of their tool calls, followed by one
// tool message per result — the order the Chat Completions API requires, and the same shape
// RunTurn builds live.
func historyToMessages(rows []domain.SimpleChatMessage) []openai.Message {
	messages := make([]openai.Message, 0, len(rows))

	idx := 0
	for idx < len(rows) {
		row := rows[idx]

		switch domain.SimpleChatRole(row.Role) {
		case domain.SimpleChatRoleUser:
			msg := openai.Message{Role: "user", Content: row.Content}
			messages = append(messages, msg)
			idx++
		case domain.SimpleChatRoleAssistant:
			msg := openai.Message{Role: "assistant", Content: row.Content}
			messages = append(messages, msg)
			idx++
		case domain.SimpleChatRoleTool:
			end := idx
			for end < len(rows) && domain.SimpleChatRole(rows[end].Role) == domain.SimpleChatRoleTool {
				end++
			}

			calls := make([]openai.ToolCall, 0, end-idx)
			for _, toolRow := range rows[idx:end] {
				call := openai.ToolCall{
					ID:            derefString(toolRow.ToolCallID),
					Name:          derefString(toolRow.ToolName),
					ArgumentsJSON: string(toolRow.ToolInput),
				}
				calls = append(calls, call)
			}

			callMsg := openai.Message{Role: "assistant", ToolCalls: calls}
			messages = append(messages, callMsg)

			for _, toolRow := range rows[idx:end] {
				resultMsg := openai.Message{
					Role:       "tool",
					Content:    toolRow.Content,
					ToolCallID: derefString(toolRow.ToolCallID),
				}
				messages = append(messages, resultMsg)
			}

			idx = end
		default:
			idx++
		}
	}

	return messages
}

// failTurn reports err to the client as an error event and returns it unchanged.
func failTurn(sink chatprotocol.EventSink, err error) error {
	errorEvent := chatprotocol.Event{
		Type: chatprotocol.EventError,
		Text: err.Error(),
	}

	sendErr := sink.Send(errorEvent)
	if sendErr != nil {
		log.Warn().Err(sendErr).Msg("error sending simple chat error event")
	}

	return err
}

// rawArguments returns the tool arguments as a JSON payload for the wire, or nil when the model
// sent nothing parseable — Event.Input is omitempty, and emitting invalid JSON there would break
// the consumer's unmarshal of the whole envelope.
func rawArguments(argumentsJSON string) json.RawMessage {
	if argumentsJSON == "" {
		return nil
	}

	if !json.Valid([]byte(argumentsJSON)) {
		return nil
	}

	return json.RawMessage(argumentsJSON)
}

// summarize renders a short one-line preview of a tool call's arguments for the UI.
func summarize(argumentsJSON string) string {
	compact := strings.Join(strings.Fields(argumentsJSON), " ")

	if len(compact) <= inputSummaryLimit {
		return compact
	}

	return compact[:inputSummaryLimit] + "..."
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
