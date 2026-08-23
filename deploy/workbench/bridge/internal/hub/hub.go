// Package hub is the bridge's WebSocket fan-out: it accepts any number of simultaneous consumers
// (the web chat UI, a Telegram relay, ...), replays this bridge process's event backlog to each
// new one, broadcasts every outbound event to all of them, and funnels every inbound event from
// any of them into a single channel the rest of the bridge consumes.
//
// Consumers are deliberately equals — there is no primary/secondary distinction. A permission
// decision or a user message is honoured whichever socket it arrives on, which is what makes
// "start a conversation in the browser, approve the tool call from Telegram" work without any
// session-ownership bookkeeping.
package hub

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"

	"workbenchbridge/internal/chatprotocol"
)

// backlogLimit caps how many events are retained for replay to a late-joining consumer. A long
// conversation with verbose tool output would otherwise grow this without bound inside a
// memory-limited container; once full, the oldest events are dropped, so a very late joiner sees
// a truncated — never a wrong — history.
const backlogLimit = 4000

// sendQueueSize is the per-connection outbound buffer. A consumer that stops reading (a wedged
// browser tab, a stalled proxy) fills its queue and is then dropped rather than being allowed to
// block the whole broadcast path, which would stall the turn producing the events.
const sendQueueSize = 4096

// Hub fans events out to every attached consumer and collects their inbound events.
type Hub struct {
	upgrader websocket.Upgrader

	// mu guards conns, backlog, nextSeq, and the active session file, together with every
	// connection's send channel and closed flag. One lock for all of it, rather than a lock per
	// connection, is what makes "queue the backlog and register the connection" atomic against a
	// concurrent Broadcast: a consumer attaching mid-turn sees each event exactly once, either as
	// replayed backlog or as a live broadcast.
	mu      sync.Mutex
	conns   map[*connection]struct{}
	backlog []chatprotocol.Event

	// nextSeq is the Seq stamped onto the next published event (see publish). It starts at 1 for
	// a brand-new Hub with no prior history, or continues past the highest Seq found in an
	// existing session file when one is rehydrated — see loadHistory. It is deliberately never
	// reset by a session rotation (Reset/rotateSession): Seq is monotonic for this Hub's entire
	// lifetime, spanning any number of chat sessions, not just the current one.
	nextSeq uint64

	// historyDir is where each chat session's events are persisted as one append-only JSONL file
	// — see history.go. It should live on the workbench's persistent named volume so history
	// survives a bridge process restart (container restart/redeploy); see New's caller.
	historyDir string

	// sessionFile is the active session's append-only JSONL file that every published event is
	// appended to (see publish/appendToSession). Rotated — closed and replaced — by
	// rotateSession, which Reset calls before writing its event.
	sessionFile *os.File

	// sessionID is the active session file's name without its .jsonl extension, kept only for
	// log messages (see rotateSession).
	sessionID string

	inbound chan chatprotocol.Event
}

type connection struct {
	ws   *websocket.Conn
	send chan []byte

	// closed is set (under Hub.mu) exactly once, immediately before send is closed, so no
	// broadcast can ever send on a closed channel.
	closed bool
}

// New returns a Hub that persists every published event as JSONL under historyDir and rehydrates
// its backlog and sequence numbering from the most recent session file already there, if any —
// see loadHistory. historyDir is created if missing. Pass a path on the workbench's persistent
// named volume so history survives a bridge process restart (container restart/redeploy); the
// caller (main.go) hardcodes the one that satisfies this.
func New(historyDir string) *Hub {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		// Every request reaching the bridge has already passed the backend's authenticated
		// reverse proxy (internal/transport/vaults_api/workbench_terminal.go), and the bridge is
		// only reachable through it, so there is no additional origin to check here — and the
		// proxied Origin header is the end user's browser's, not something this process can
		// meaningfully validate.
		CheckOrigin: alwaysAllowOrigin,
	}

	h := &Hub{
		upgrader:   upgrader,
		conns:      map[*connection]struct{}{},
		inbound:    make(chan chatprotocol.Event, 64),
		historyDir: historyDir,
	}

	h.loadHistory()

	return h
}

func alwaysAllowOrigin(r *http.Request) bool {
	return true
}

// Inbound is the stream of events consumers sent to the bridge (user_message,
// permission_decision, auth_code_submit). Never closed — the bridge runs for the container's
// lifetime.
func (h *Hub) Inbound() <-chan chatprotocol.Event {
	return h.inbound
}

// ConsumerCount reports how many consumers are currently attached.
func (h *Hub) ConsumerCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	return len(h.conns)
}

// Broadcast records event in the backlog and delivers it to every currently attached consumer.
// Safe to call from any goroutine, including concurrently with a consumer attaching.
func (h *Hub) Broadcast(event chatprotocol.Event) {
	h.publish(event, func(stamped chatprotocol.Event) {
		h.backlog = append(h.backlog, stamped)
		if len(h.backlog) > backlogLimit {
			h.backlog = h.backlog[len(h.backlog)-backlogLimit:]
		}
	})
}

// Reset discards the entire backlog in favor of event alone, then delivers event to every
// currently attached consumer. Use it for starting a new chat: existing history should not be
// replayed to a consumer connecting after this point, only event (typically the new_chat event
// itself) and whatever is broadcast from here on. It also rotates to a brand-new session file
// (see rotateSession) before event is written to disk, so the discarded chat's history stays
// intact and browsable under its own file rather than being appended to.
func (h *Hub) Reset(event chatprotocol.Event) {
	h.publish(event, func(stamped chatprotocol.Event) {
		h.rotateSession()
		h.backlog = []chatprotocol.Event{stamped}
	})
}

// publish stamps event with the next sequence number and — under h.mu — runs mutateBacklog to
// apply Broadcast's or Reset's differing backlog-mutation step (Reset's also rotates the session
// file), appends event to the active session file, marshals it, and delivers it to every
// currently attached consumer. Safe to call from any goroutine, including concurrently with a
// consumer attaching.
//
// mutateBacklog receives the stamped event rather than closing over the caller's original copy:
// Broadcast/Reset's event parameter and this function's are distinct copies (Go passes by value),
// so the Seq assigned here would otherwise never reach what gets appended to h.backlog or written
// to disk.
func (h *Hub) publish(event chatprotocol.Event, mutateBacklog func(stamped chatprotocol.Event)) {
	h.mu.Lock()
	defer h.mu.Unlock()

	event.Seq = h.nextSeq
	h.nextSeq++

	mutateBacklog(event)

	h.appendToSession(event)

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("hub: error marshalling %s event: %v", event.Type, err)

		return
	}

	for conn := range h.conns {
		h.deliverLocked(conn, payload)
	}
}

// deliverLocked queues payload on conn, dropping conn entirely if its queue is full — see
// sendQueueSize. Caller must hold h.mu.
func (h *Hub) deliverLocked(conn *connection, payload []byte) {
	if conn.closed {
		return
	}

	select {
	case conn.send <- payload:
	default:
		log.Print("hub: dropping a consumer that stopped reading")
		h.removeLocked(conn)
	}
}

// ServeHTTP upgrades an incoming request to a WebSocket and serves it for as long as the consumer
// stays attached. It is registered as a catch-all: which path a consumer connects on is
// deliberately not significant, since the backend proxy forwards everything under a workbench's
// route prefix here verbatim and each consumer track picks its own path.
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		h.serveNonUpgrade(w, r)

		return
	}

	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Upgrade has already written its own error response by this point.
		log.Printf("hub: error upgrading connection: %v", err)

		return
	}

	conn := &connection{
		ws:   ws,
		send: make(chan []byte, sendQueueSize),
	}

	h.add(conn)

	go h.writeLoop(conn)

	h.readLoop(conn)
}

// serveNonUpgrade routes every plain (non-upgrade) request: the /history session listing and
// /history/<id> session detail endpoints (see history.go), or — for anything else — the status
// response (see serveStatus), exactly as before those two were added. Routed by inspecting
// r.Method/r.URL.Path directly rather than an http.ServeMux, since that's the whole surface: two
// GET routes plus a catch-all fallback don't warrant a second router nested inside Hub.
func (h *Hub) serveNonUpgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/history" {
		h.serveHistoryList(w)

		return
	}

	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/history/") {
		id := strings.TrimPrefix(r.URL.Path, "/history/")
		h.serveHistoryDetail(w, id)

		return
	}

	h.serveStatus(w)
}

// serveStatus answers a plain (non-upgrade) GET that isn't a /history route, so a health check or
// a curious operator hitting the bridge through the backend proxy gets something intelligible
// instead of a protocol error.
func (h *Hub) serveStatus(w http.ResponseWriter) {
	h.mu.Lock()
	status := map[string]any{
		"service":   "workbench-bridge",
		"consumers": len(h.conns),
		"backlog":   len(h.backlog),
	}
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		log.Printf("hub: error writing status response: %v", err)
	}
}

// add registers conn and seeds its send queue with the backlog, under the same lock — see Hub.mu.
func (h *Hub) add(conn *connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, event := range h.backlog {
		payload, err := json.Marshal(event)
		if err != nil {
			continue
		}

		select {
		case conn.send <- payload:
		default:
			// A backlog longer than the send queue can't be replayed in full; the consumer still
			// gets everything live from the moment it attached.
		}
	}

	h.conns[conn] = struct{}{}
}

func (h *Hub) remove(conn *connection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.removeLocked(conn)
}

// removeLocked detaches conn and closes its send channel, ending its writeLoop. Caller must hold
// h.mu.
func (h *Hub) removeLocked(conn *connection) {
	if conn.closed {
		return
	}

	conn.closed = true
	delete(h.conns, conn)
	close(conn.send)
}

// readLoop decodes inbound frames until the consumer disconnects. A frame that isn't a valid
// Event is reported back to that one consumer as an error event rather than dropping the whole
// connection — a consumer sending malformed JSON is a bug worth surfacing, not worth
// disconnecting over.
func (h *Hub) readLoop(conn *connection) {
	defer h.remove(conn)

	for {
		_, payload, err := conn.ws.ReadMessage()
		if err != nil {
			return
		}

		var event chatprotocol.Event

		err = json.Unmarshal(payload, &event)
		if err != nil {
			h.replyMalformed(conn, err)

			continue
		}

		h.inbound <- event
	}
}

// replyMalformed tells the one consumer that sent unparseable JSON about it, without disturbing
// any other consumer's stream.
func (h *Hub) replyMalformed(conn *connection, cause error) {
	event := chatprotocol.Event{
		Type: chatprotocol.EventError,
		Text: "malformed event: " + cause.Error(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.deliverLocked(conn, payload)
}

// writeLoop drains conn's send queue. It owns the websocket's write side exclusively — gorilla
// permits only one concurrent writer — which is why Broadcast queues rather than writing inline.
func (h *Hub) writeLoop(conn *connection) {
	defer func() {
		err := conn.ws.Close()
		if err != nil {
			log.Printf("hub: error closing consumer connection: %v", err)
		}
	}()

	for payload := range conn.send {
		err := conn.ws.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			h.remove(conn)

			return
		}
	}
}
