package simple_chat_api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/ruf-dev/artel/internal/chatprotocol"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	"github.com/ruf-dev/artel/internal/service"
	"go.redsock.ru/rerrors"
)

const (
	// ChatWsRoutePattern is the http.ServeMux pattern this handler is registered under in
	// internal/app/custom.go.
	//
	// It sits under "/api/" on purpose: the session cookie this handler authenticates with is
	// scoped to middleware.CookiePath ("/api/"), so a route outside that prefix would never
	// receive it. Being more specific than the grpc-gateway's "/api/simple-chats/" subtree
	// pattern, http.ServeMux routes a matching request here rather than to the gateway.
	ChatWsRoutePattern = "/api/simple-chats/{chatId}/ws"

	// chatIdPathValue is the wildcard name in ChatWsRoutePattern, read back via
	// http.Request.PathValue.
	chatIdPathValue = "chatId"
)

// chatAuthService is the narrow subset of service.AuthService this handler depends on — the same
// token validation the gRPC auth interceptor performs, reused rather than reimplemented.
type chatAuthService interface {
	ValidateToken(ctx context.Context, token string) (domain.User, error)
}

// chatWsUpgrader upgrades the client side of a Simple Chat connection.
var chatWsUpgrader = websocket.Upgrader{
	CheckOrigin: allowChatWsOrigin,
}

// allowChatWsOrigin accepts every origin: ServeHTTP has already authenticated the caller and
// verified they own the chat before any upgrade happens, and the browser client always
// same-origins here through the artel UI — nothing extra to check at the WS layer.
func allowChatWsOrigin(_ *http.Request) bool {
	return true
}

// ChatWsHandler serves the live Simple Chat turn exchange. It speaks the same
// chatprotocol.Event envelope as the Docker workbench's in-container bridge, so the web client
// reuses its existing chat components unchanged — but the agent runs in this process rather than
// in a container.
//
// Unlike the workbench bridge there is no in-memory backlog or JSONL replay: a reconnecting
// client rebuilds history from Postgres via the GetSimpleChat RPC, so this handler only ever
// streams events for turns it is currently running.
type ChatWsHandler struct {
	authSvc       chatAuthService
	simpleChatSvc service.SimpleChatService
}

func NewChatWsHandler(authSvc chatAuthService, simpleChatSvc service.SimpleChatService) *ChatWsHandler {
	return &ChatWsHandler{
		authSvc:       authSvc,
		simpleChatSvc: simpleChatSvc,
	}
}

func (h *ChatWsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	chatUuid, err := uuid.Parse(r.PathValue(chatIdPathValue))
	if err != nil {
		http.Error(w, "invalid chat id", http.StatusBadRequest)

		return
	}

	user, err := h.authenticate(ctx, r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)

		return
	}

	// Ownership, not vault membership: a Simple Chat thread is personal to its creator.
	chat, err := h.simpleChatSvc.OwnedChat(ctx, chatUuid, user.Uuid)
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)

		return
	}

	conn, err := chatWsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn().Err(err).Msg("error upgrading simple chat websocket")

		return
	}

	h.serve(ctx, conn, chat)
}

// authenticate resolves the caller from the access_token cookie, the same credential the
// grpc-gateway's CookieToMetadataAnnotator forwards as "authorization" metadata for regular
// RPCs. An absent cookie is deliberately passed through to ValidateToken as an empty token
// rather than short-circuited: it fails the session lookup exactly like an invalid one, while
// still honouring the no-auth local-dev mode, where ValidateToken returns the fixed dev user
// regardless of what was presented.
func (h *ChatWsHandler) authenticate(ctx context.Context, r *http.Request) (domain.User, error) {
	token := ""

	cookie, err := r.Cookie(middleware.AccessTokenCookieName)
	if err == nil {
		token = cookie.Value
	}

	user, err := h.authSvc.ValidateToken(ctx, token)
	if err != nil {
		return domain.User{}, rerrors.Wrap(err, "error validating simple chat access token")
	}

	return user, nil
}

// serve runs the connection's read loop until the client goes away. Each user_message starts a
// turn on its own goroutine so this loop keeps reading — that is what lets a permission_decision
// arrive while a turn is parked inside EventSink.AwaitPermissionDecision.
func (h *ChatWsHandler) serve(ctx context.Context, conn *websocket.Conn, chat domain.SimpleChat) {
	connCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sink := newConnSink(conn)

	var turns sync.WaitGroup

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var event chatprotocol.Event

		err = json.Unmarshal(payload, &event)
		if err != nil {
			log.Warn().Err(err).Msg("error decoding simple chat inbound event")

			continue
		}

		switch event.Type {
		case chatprotocol.EventUserMessage:
			h.startTurn(connCtx, &turns, sink, chat, event)
		case chatprotocol.EventPermissionDecision:
			sink.resolvePermission(event.ID, event.Decision)
		default:
			log.Debug().Str("type", string(event.Type)).
				Msg("simple chat: ignoring unexpected inbound event type")
		}
	}

	// Cancel any in-flight turn, then wait for it to unwind so nothing writes to a closed
	// connection after this returns.
	cancel()
	turns.Wait()

	closeErr := conn.Close()
	if closeErr != nil {
		log.Debug().Err(closeErr).Msg("error closing simple chat websocket")
	}
}

// startTurn echoes the user's message back and dispatches the turn.
//
// The echo mirrors the workbench bridge (see deploy/workbench/bridge/main.go's runMainLoop,
// which Broadcasts an inbound user_message before running the turn): the web client tags its
// optimistic message with a client-generated id and recognizes the echo by that id, so echoing
// keeps this handler interchangeable with the bridge for the shared chat components.
func (h *ChatWsHandler) startTurn(
	ctx context.Context,
	turns *sync.WaitGroup,
	sink *connSink,
	chat domain.SimpleChat,
	event chatprotocol.Event,
) {
	err := sink.Send(event)
	if err != nil {
		log.Warn().Err(err).Msg("error echoing simple chat user message")

		return
	}

	// The per-turn model override: the client stamps Model on its user_message so the model can
	// be switched mid-conversation. Empty falls back to the thread's stored default.
	model := event.Model
	text := event.Text

	turns.Add(1)

	go h.runTurn(ctx, turns, sink, chat, text, model)
}

func (h *ChatWsHandler) runTurn(
	ctx context.Context,
	turns *sync.WaitGroup,
	sink *connSink,
	chat domain.SimpleChat,
	text string,
	model string,
) {
	defer turns.Done()

	err := h.simpleChatSvc.RunTurn(ctx, chat, text, model, sink)
	if err != nil {
		log.Warn().Err(err).Str("chat_id", chat.Uuid.String()).Msg("simple chat turn failed")
	}
}

// connSink is the per-connection chatprotocol.EventSink: it owns the websocket's write side and
// routes inbound permission decisions to whichever turn is waiting on them. One instance per
// connection, never shared.
type connSink struct {
	conn *websocket.Conn

	// writeMu serialises writes — gorilla permits only one concurrent writer, and a turn
	// goroutine writes while the read loop may also echo.
	writeMu sync.Mutex

	// seq stamps each outbound event with a per-connection monotonic sequence number, mirroring
	// the bridge hub's numbering so the client's existing seq-based de-duplication keeps working.
	// No cross-connection replay is needed here — history comes from Postgres, not a backlog.
	seq atomic.Uint64

	// waitersMu guards waiters, which maps a permission request id to the channel its parked
	// AwaitPermissionDecision call is listening on.
	waitersMu sync.Mutex
	waiters   map[string]chan chatprotocol.PermissionDecision
}

func newConnSink(conn *websocket.Conn) *connSink {
	return &connSink{
		conn:    conn,
		waiters: make(map[string]chan chatprotocol.PermissionDecision),
	}
}

func (c *connSink) Send(event chatprotocol.Event) error {
	event.Seq = c.seq.Add(1)

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	err := c.conn.WriteJSON(event)
	if err != nil {
		return rerrors.Wrap(err, "error writing simple chat event")
	}

	return nil
}

func (c *connSink) AwaitPermissionDecision(
	ctx context.Context, requestId string,
) (chatprotocol.PermissionDecision, error) {
	// Buffered so resolvePermission never blocks on a waiter that has already given up.
	decisions := make(chan chatprotocol.PermissionDecision, 1)

	c.waitersMu.Lock()
	c.waiters[requestId] = decisions
	c.waitersMu.Unlock()

	defer c.forgetWaiter(requestId)

	select {
	case decision := <-decisions:
		return decision, nil
	case <-ctx.Done():
		return "", rerrors.Wrap(ctx.Err(), "error awaiting permission decision")
	}
}

// resolvePermission hands decision to the turn parked on requestId, and echoes it back to the
// client so its pending permission card resolves — matching the bridge, which also rebroadcasts
// a decision it accepted. An unknown id is ignored: the turn it belonged to is already gone.
func (c *connSink) resolvePermission(requestId string, decision chatprotocol.PermissionDecision) {
	c.waitersMu.Lock()
	waiter, ok := c.waiters[requestId]
	c.waitersMu.Unlock()

	if !ok {
		return
	}

	waiter <- decision

	echo := chatprotocol.Event{
		Type:     chatprotocol.EventPermissionDecision,
		ID:       requestId,
		Decision: decision,
	}

	err := c.Send(echo)
	if err != nil {
		log.Debug().Err(err).Msg("error echoing simple chat permission decision")
	}
}

func (c *connSink) forgetWaiter(requestId string) {
	c.waitersMu.Lock()
	delete(c.waiters, requestId)
	c.waitersMu.Unlock()
}
