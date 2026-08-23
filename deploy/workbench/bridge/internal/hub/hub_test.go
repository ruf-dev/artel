package hub

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"workbenchbridge/internal/chatprotocol"
)

// wsReadTimeout bounds every websocket read in this test so a bug that reintroduces unbounded
// backlog replay fails fast with a clear error instead of hanging until `go test`'s own timeout.
const wsReadTimeout = 10 * time.Second

// dialChat opens a websocket consumer connection to the hub server at addr and returns it. addr is
// an http:// URL, translated to ws:// here so callers don't have to. Mirrors main_test.go's helper
// of the same name in the module's other package.
func dialChat(t *testing.T, addr string) *websocket.Conn {
	t.Helper()

	url := "ws" + strings.TrimPrefix(addr, "http") + "/ws"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("error dialing chat websocket at %s: %v", url, err)
	}

	return conn
}

// readOneEvent reads a single event off conn, failing the test if none arrives within
// wsReadTimeout.
func readOneEvent(t *testing.T, conn *websocket.Conn) chatprotocol.Event {
	t.Helper()

	err := conn.SetReadDeadline(time.Now().Add(wsReadTimeout))
	if err != nil {
		t.Fatalf("error setting read deadline: %v", err)
	}

	_, payload, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("error reading from chat websocket: %v", err)
	}

	var event chatprotocol.Event

	err = json.Unmarshal(payload, &event)
	if err != nil {
		t.Fatalf("error unmarshalling event %q: %v", payload, err)
	}

	return event
}

// TestHub_ResetDiscardsBacklog proves Reset's contract: after Broadcasting a few events and then
// calling Reset with a new one, a consumer attaching afterward replays only the Reset event — none
// of the earlier Broadcast ones survive in the backlog.
func TestHub_ResetDiscardsBacklog(t *testing.T) {
	h := New()

	server := httptest.NewServer(h)
	defer server.Close()

	// An attached consumer to observe the live delivery side of both Broadcast and Reset, not just
	// the backlog replay side.
	live := dialChat(t, server.URL)
	defer live.Close()

	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventSystemInit})
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventUserMessage, Text: "hello"})
	h.Broadcast(chatprotocol.Event{Type: chatprotocol.EventTurnDone, SessionID: "abc"})

	for i := 0; i < 3; i++ {
		readOneEvent(t, live)
	}

	resetEvent := chatprotocol.Event{Type: chatprotocol.EventNewChat}
	h.Reset(resetEvent)

	// The already-attached consumer should see the Reset event delivered live too.
	delivered := readOneEvent(t, live)
	if delivered.Type != chatprotocol.EventNewChat {
		t.Fatalf("live consumer's post-Reset event type = %q, want %q", delivered.Type, chatprotocol.EventNewChat)
	}

	// A freshly attached consumer, modeling a page reload after Reset, should replay only the
	// Reset event — never the three Broadcast events from before it.
	fresh := dialChat(t, server.URL)
	defer fresh.Close()

	replayed := readOneEvent(t, fresh)
	if replayed.Type != chatprotocol.EventNewChat {
		t.Fatalf("fresh consumer's replayed event type = %q, want %q (backlog was not reset)",
			replayed.Type, chatprotocol.EventNewChat)
	}

	// Confirm there is nothing else queued behind it — the backlog held exactly one event.
	err := fresh.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if err != nil {
		t.Fatalf("error setting read deadline: %v", err)
	}

	_, _, err = fresh.ReadMessage()
	if err == nil {
		t.Fatalf("fresh consumer received a second replayed event; backlog should contain only the Reset event")
	}
}
