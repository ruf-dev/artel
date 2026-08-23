package telegram_webhook

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/chatprotocol"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

// bridgeWSPath is the workbench bridge's WebSocket endpoint, dialed relative to the
// "http://<host>:<port>" base URL bridgeTargetResolver.ResolveTerminalTarget returns.
//
// ASSUMPTION, flagged for the parallel track building deploy/workbench/bridge: this mirrors ttyd's
// own "/ws" endpoint off its root (see internal/transport/vaults_api/workbench_terminal.go's
// comment on ttyd's path layout), since as of writing deploy/workbench/bridge contained only a
// hand-mirrored copy of internal/chatprotocol/events.go and no server file yet — there was nothing
// to read the real path off. Confirm/update this constant once the bridge lands.
const bridgeWSPath = "/ws"

// deltaEditInterval bounds how often an in-flight assistant_text_delta message gets re-edited on
// Telegram — coalescing many small token deltas into occasional edits rather than one Bot API
// call per token, per Telegram's message-edit rate limits.
const deltaEditInterval = time.Second

// dialTimeout bounds how long dialing a workbench bridge WebSocket may take before giving up —
// this runs synchronously inside webhook request handling (see handler.go's dispatch), so a
// hung dial must not hold the request (and the sender goroutine) open indefinitely.
const dialTimeout = 10 * time.Second

// bridgeSession is one relay's worth of state for a single external_connections row: a workbench
// bridge WebSocket connection plus enough per-inflight-event bookkeeping to coalesce
// assistant_text_delta edits and to finalize a permission_request message once its decision comes
// back over the callback_query. One session per connection (not per chat — v1 scope is one
// telegram connection per user, one chat per connection).
type bridgeSession struct {
	handler   *Handler
	connUuid  uuid.UUID
	chatId    int64
	targetURL string

	writeMu sync.Mutex
	conn    *websocket.Conn

	stateMu          sync.Mutex
	closed           bool
	awaitingAuthCode bool
	deltaMsgId       map[string]int64
	deltaText        map[string]string
	deltaLastEdit    map[string]time.Time
	permMsgId        map[string]int64
	permText         map[string]string
}

// getOrDialSession returns the connection's live bridge session, dialing (or re-dialing, if the
// previous connection dropped) one on demand. chatId is threaded through explicitly since it may
// have just been captured by handleMessage moments ago and not yet reflected in any stored state
// this method reads.
func (h *Handler) getOrDialSession(ctx context.Context, exConnUuid uuid.UUID, userUuid uuid.UUID, chatId int64) (*bridgeSession, error) {
	if v, ok := h.sessions.Load(exConnUuid); ok {
		sess, castOk := v.(*bridgeSession)
		if castOk {
			sess.stateMu.Lock()
			closed := sess.closed
			sess.stateMu.Unlock()

			if !closed {
				return sess, nil
			}
		}

		h.sessions.Delete(exConnUuid)
	}

	wbResult, err := h.workbenches.GetMostRecentByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting most recent workbench for telegram relay")
	}

	if !wbResult.Valid {
		return nil, user_errors.TelegramNoWorkbench
	}

	wb := wbResult.V

	target, err := h.workbenchSvc.ResolveTerminalTarget(ctx, wb.VaultUuid, wb.UserUuid)
	if err != nil {
		return nil, err
	}

	dialCtx, cancel := context.WithTimeout(ctx, dialTimeout)
	defer cancel()

	conn, err := dialBridge(dialCtx, target)
	if err != nil {
		return nil, rerrors.Wrap(err, "error dialing workbench chat bridge")
	}

	sess := newBridgeSession(h, exConnUuid, chatId, target, conn)
	h.sessions.Store(exConnUuid, sess)

	go sess.readLoop()

	return sess, nil
}

func newBridgeSession(h *Handler, connUuid uuid.UUID, chatId int64, targetURL string, conn *websocket.Conn) *bridgeSession {
	return &bridgeSession{
		handler:       h,
		connUuid:      connUuid,
		chatId:        chatId,
		targetURL:     targetURL,
		conn:          conn,
		deltaMsgId:    make(map[string]int64),
		deltaText:     make(map[string]string),
		deltaLastEdit: make(map[string]time.Time),
		permMsgId:     make(map[string]int64),
		permText:      make(map[string]string),
	}
}

func dialBridge(ctx context.Context, targetBaseURL string) (*websocket.Conn, error) {
	parsed, err := url.Parse(targetBaseURL)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing workbench bridge target url")
	}

	scheme := "ws"
	if parsed.Scheme == "https" {
		scheme = "wss"
	}

	parsed.Scheme = scheme
	parsed.Path = bridgeWSPath

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, parsed.String(), nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error dialing workbench bridge websocket")
	}

	return conn, nil
}

func (s *bridgeSession) markClosed() {
	s.stateMu.Lock()
	s.closed = true
	s.stateMu.Unlock()
}

// takeAwaitingAuthCode reports whether the bridge is waiting on an auth_code_submit, clearing the
// flag in the same step — the next plain-text message consumed as a code is a one-shot.
func (s *bridgeSession) takeAwaitingAuthCode() bool {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	awaiting := s.awaitingAuthCode
	s.awaitingAuthCode = false

	return awaiting
}

func (s *bridgeSession) setAwaitingAuthCode() {
	s.stateMu.Lock()
	s.awaitingAuthCode = true
	s.stateMu.Unlock()
}

func (s *bridgeSession) writeEvent(event chatprotocol.Event) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	err := s.conn.WriteJSON(event)
	if err != nil {
		return rerrors.Wrap(err, "error writing event to workbench bridge")
	}

	return nil
}

func (s *bridgeSession) sendUserMessage(text string) {
	event := chatprotocol.Event{
		Type: chatprotocol.EventUserMessage,
		Text: text,
	}

	err := s.writeEvent(event)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send user_message to bridge")
	}
}

func (s *bridgeSession) sendAuthCodeSubmit(code string) {
	event := chatprotocol.Event{
		Type: chatprotocol.EventAuthCodeSubmit,
		Code: code,
	}

	err := s.writeEvent(event)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send auth_code_submit to bridge")
	}
}

func (s *bridgeSession) sendPermissionDecision(eventId string, decision permissionDecision) {
	event := chatprotocol.Event{
		Type:     chatprotocol.EventPermissionDecision,
		ID:       eventId,
		Decision: decision,
	}

	err := s.writeEvent(event)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send permission_decision to bridge")
	}
}

// finalizePermissionMessage edits the Telegram message that carried a permission_request's
// buttons to show which decision was taken, and clears its keyboard so the buttons stop looking
// clickable. Falls back to the callback_query's own message id/chat id (chatMsgId/chatId) when
// this session never recorded the permission (e.g. the session was re-dialed after a restart, so
// permText/permMsgId lost their in-memory state) — the edit still succeeds since Telegram
// addresses messages by (chat_id, message_id), not by this session's bookkeeping.
func (s *bridgeSession) finalizePermissionMessage(ctx context.Context, eventId string, chatId, messageId int64, label string) {
	s.stateMu.Lock()
	origText, hasText := s.permText[eventId]
	delete(s.permMsgId, eventId)
	delete(s.permText, eventId)
	s.stateMu.Unlock()

	text := label
	if hasText {
		text = origText + "\n\n→ " + label
	}

	err := s.handler.editText(ctx, s.connUuid, chatId, messageId, text, true)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to finalize permission message")
	}
}

// readLoop pumps events off the bridge WebSocket until it errors or closes, translating each into
// Telegram API calls. Runs for the lifetime of the session; markClosed on exit lets
// getOrDialSession know to redial next time this connection needs to send something.
func (s *bridgeSession) readLoop() {
	defer s.markClosed()

	ctx := s.handler.baseCtx

	for {
		var event chatprotocol.Event

		err := s.conn.ReadJSON(&event)
		if err != nil {
			log.Warn().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: bridge websocket read loop ended")

			return
		}

		s.handleEvent(ctx, event)
	}
}

func (s *bridgeSession) handleEvent(ctx context.Context, event chatprotocol.Event) {
	switch event.Type {
	case chatprotocol.EventSystemInit:
		// No payload, nothing to relay.
	case chatprotocol.EventAssistantTextDelta:
		s.handleAssistantDelta(ctx, event)
	case chatprotocol.EventAssistantTextDone:
		s.handleAssistantDone(ctx, event)
	case chatprotocol.EventToolCallStarted:
		s.handleToolCallStarted(ctx, event)
	case chatprotocol.EventToolCallResult:
		s.handleToolCallResult(ctx, event)
	case chatprotocol.EventPermissionRequest:
		s.handlePermissionRequest(ctx, event)
	case chatprotocol.EventTurnDone:
		// No Telegram message — assistant_text_done already delivered the final text.
	case chatprotocol.EventError:
		s.handleError(ctx, event)
	case chatprotocol.EventAuthLink:
		s.handleAuthLink(ctx, event)
	case chatprotocol.EventAuthCodeNeeded:
		s.handleAuthCodeNeeded(ctx, event)
	case chatprotocol.EventAuthComplete:
		// No Telegram message — nothing user-facing to relay.
	case chatprotocol.EventUserMessage:
		// No Telegram message — this consumer already knows what it (or another consumer) sent.
	case chatprotocol.EventNewChat:
		// No Telegram message — nothing user-facing to relay; a fresh conversation just means the
		// next reply isn't a continuation of the old one.
	default:
		log.Warn().Str("external_connection_id", s.connUuid.String()).Str("event_type", string(event.Type)).
			Msg("telegram webhook: unrecognized bridge event type")
	}
}

func (s *bridgeSession) handleAssistantDelta(ctx context.Context, event chatprotocol.Event) {
	s.stateMu.Lock()
	text := s.deltaText[event.ID] + event.Text
	s.deltaText[event.ID] = text
	msgId, hasMsg := s.deltaMsgId[event.ID]
	lastEdit := s.deltaLastEdit[event.ID]
	s.stateMu.Unlock()

	if !hasMsg {
		newId, err := s.handler.sendText(ctx, s.connUuid, s.chatId, text)
		if err != nil {
			log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send assistant text delta")

			return
		}

		s.stateMu.Lock()
		s.deltaMsgId[event.ID] = newId
		s.deltaLastEdit[event.ID] = time.Now()
		s.stateMu.Unlock()

		return
	}

	if time.Since(lastEdit) < deltaEditInterval {
		return
	}

	err := s.handler.editText(ctx, s.connUuid, s.chatId, msgId, text, false)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to edit assistant text delta")

		return
	}

	s.stateMu.Lock()
	s.deltaLastEdit[event.ID] = time.Now()
	s.stateMu.Unlock()
}

func (s *bridgeSession) handleAssistantDone(ctx context.Context, event chatprotocol.Event) {
	s.stateMu.Lock()
	msgId, hasMsg := s.deltaMsgId[event.ID]
	delete(s.deltaMsgId, event.ID)
	delete(s.deltaText, event.ID)
	delete(s.deltaLastEdit, event.ID)
	s.stateMu.Unlock()

	if event.Text == "" {
		return
	}

	if !hasMsg {
		_, err := s.handler.sendText(ctx, s.connUuid, s.chatId, event.Text)
		if err != nil {
			log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send assistant text done")
		}

		return
	}

	err := s.handler.editText(ctx, s.connUuid, s.chatId, msgId, event.Text, false)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to edit assistant text done")
	}
}

func (s *bridgeSession) handleToolCallStarted(ctx context.Context, event chatprotocol.Event) {
	text := fmt.Sprintf("🔧 Running: %s", event.ToolName)
	if event.InputSummary != "" {
		text += "\n" + event.InputSummary
	}

	_, err := s.handler.sendText(ctx, s.connUuid, s.chatId, text)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send tool_call_started")
	}
}

func (s *bridgeSession) handleToolCallResult(ctx context.Context, event chatprotocol.Event) {
	status := "✅ Done"
	if event.IsError {
		status = "❌ Failed"
	}

	text := fmt.Sprintf("%s: %s", status, event.ToolName)

	_, err := s.handler.sendText(ctx, s.connUuid, s.chatId, text)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send tool_call_result")
	}
}

func (s *bridgeSession) handlePermissionRequest(ctx context.Context, event chatprotocol.Event) {
	var textBuilder strings.Builder

	textBuilder.WriteString("🔐 Permission requested\nTool: ")
	textBuilder.WriteString(event.ToolName)

	if event.InputSummary != "" {
		textBuilder.WriteString("\n")
		textBuilder.WriteString(event.InputSummary)
	}

	text := textBuilder.String()

	buttons := []inlineKeyboardButton{
		{Text: "Allow", CallbackData: encodeCallbackData(event.ID, decisionAllowOnce)},
		{Text: "Allow always", CallbackData: encodeCallbackData(event.ID, decisionAllowAlways)},
		{Text: "Deny", CallbackData: encodeCallbackData(event.ID, decisionDeny)},
	}

	keyboard := inlineKeyboardMarkup{InlineKeyboard: [][]inlineKeyboardButton{buttons}}

	msgId, err := s.handler.sendKeyboard(ctx, s.connUuid, s.chatId, text, keyboard)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send permission_request")

		return
	}

	s.stateMu.Lock()
	s.permMsgId[event.ID] = msgId
	s.permText[event.ID] = text
	s.stateMu.Unlock()
}

func (s *bridgeSession) handleError(ctx context.Context, event chatprotocol.Event) {
	_, err := s.handler.sendText(ctx, s.connUuid, s.chatId, "⚠️ "+event.Text)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send error event")
	}
}

func (s *bridgeSession) handleAuthLink(ctx context.Context, event chatprotocol.Event) {
	button := inlineKeyboardButton{Text: "Open link", Url: event.URL}
	keyboard := inlineKeyboardMarkup{InlineKeyboard: [][]inlineKeyboardButton{{button}}}

	_, err := s.handler.sendKeyboard(ctx, s.connUuid, s.chatId, "Authentication required", keyboard)
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send auth_link")
	}
}

// handleAuthCodeNeeded ignores its event parameter — auth_code_needed carries no payload fields
// (see internal/chatprotocol/events.go) — but keeps it for a uniform handleX(ctx, event)
// signature across handleEvent's dispatch table.
func (s *bridgeSession) handleAuthCodeNeeded(ctx context.Context, _ chatprotocol.Event) {
	s.setAwaitingAuthCode()

	_, err := s.handler.sendText(ctx, s.connUuid, s.chatId, "Please reply with your authentication code.")
	if err != nil {
		log.Error().Err(err).Str("external_connection_id", s.connUuid.String()).Msg("telegram webhook: failed to send auth_code_needed")
	}
}
