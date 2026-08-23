package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"workbenchbridge/internal/chatprotocol"
	"workbenchbridge/internal/claudecli"
	"workbenchbridge/internal/envdrop"
	"workbenchbridge/internal/hub"
	"workbenchbridge/internal/permissions"
)

// wsReadTimeout bounds every websocket read in this test so a bug that reintroduces the missing
// broadcast (the user_message never showing up) fails fast with a clear error instead of hanging
// until `go test`'s own timeout.
const wsReadTimeout = 10 * time.Second

// writeStubClaude writes a fake `claude` binary to dir and returns its path. It ignores every
// argument it's given and prints a minimal, but valid, --output-format stream-json transcript: an
// init line (so translate.go's translateSystem fires) followed by a result line (so
// translateResult produces a turn_done). That's the least a real `claude -p ... --output-format
// stream-json --verbose --include-partial-messages` invocation could emit and still let
// claudecli.Runner.RunTurn complete a turn without emitting an error event of its own — see
// translate.go's package doc for the confirmed line shapes.
func writeStubClaude(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "claude")

	script := "#!/bin/sh\n" +
		`echo '{"type":"system","subtype":"init"}'` + "\n" +
		`echo '{"type":"result","session_id":"test-session","total_cost_usd":0}'` + "\n"

	err := os.WriteFile(path, []byte(script), 0o755)
	if err != nil {
		t.Fatalf("error writing stub claude script: %v", err)
	}

	return path
}

// dialChat opens a websocket consumer connection to the hub server at addr and returns it. addr
// is an http:// URL, translated to ws:// here so callers don't have to.
func dialChat(t *testing.T, addr string) *websocket.Conn {
	t.Helper()

	url := "ws" + strings.TrimPrefix(addr, "http") + "/ws"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("error dialing chat websocket at %s: %v", url, err)
	}

	return conn
}

// readEvents reads events off conn until predicate returns true for one of them (that event is
// included in the result) or wsReadTimeout elapses, in which case the test fails.
func readEvents(t *testing.T, conn *websocket.Conn, predicate func(chatprotocol.Event) bool) []chatprotocol.Event {
	t.Helper()

	err := conn.SetReadDeadline(time.Now().Add(wsReadTimeout))
	if err != nil {
		t.Fatalf("error setting read deadline: %v", err)
	}

	var events []chatprotocol.Event

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("error reading from chat websocket (got %d events first): %v", len(events), err)
		}

		var event chatprotocol.Event

		err = json.Unmarshal(payload, &event)
		if err != nil {
			t.Fatalf("error unmarshalling event %q: %v", payload, err)
		}

		events = append(events, event)

		if predicate(event) {
			return events
		}
	}
}

func isTurnDone(event chatprotocol.Event) bool {
	return event.Type == chatprotocol.EventTurnDone
}

// TestRunMainLoop_UserMessageSurvivesReconnect is the regression test for the bug where a
// reloading/reconnecting chat consumer would see the assistant's reply but not its own prior
// message. It reproduces the real flow end to end: an actual hub, an actual runMainLoop dispatch,
// and an actual claudecli.Runner driving a stub `claude` process — the only stand-in is the
// binary runner.RunTurn execs.
//
// Before the fix (chatHub.Broadcast(event) added ahead of runner.RunTurn in runMainLoop's
// EventUserMessage case), the hub's backlog only ever received events translated from claude's
// own stdout (system_init, turn_done, ...); the inbound user_message read off chatHub.Inbound()
// was used to start the turn but never itself broadcast, so it never entered the backlog and a
// consumer #2 attaching after the fact — modeling a page reload — would replay everything except
// the user's own message. Removing the added Broadcast call reproduces that: client #2's replay
// would then contain "system_init" and "turn_done" but no "user_message", and the assertion below
// on event types would fail.
func TestRunMainLoop_UserMessageSurvivesReconnect(t *testing.T) {
	chatHub := hub.New(t.TempDir())
	broker := permissions.NewBroker(chatHub.Broadcast)

	server := httptest.NewServer(chatHub)
	defer server.Close()

	stubBinary := writeStubClaude(t)
	envStore := envdrop.New(t.TempDir())
	runner := claudecli.NewRunner(stubBinary, t.TempDir(), "", envStore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runMainLoop(ctx, chatHub, runner, broker)

	// Consumer #1: opens the chat, sends a message, and stays around long enough to see the
	// resulting turn play out (mirroring the live tab that sent it).
	consumer1 := dialChat(t, server.URL)
	defer consumer1.Close()

	outbound := chatprotocol.Event{
		Type: chatprotocol.EventUserMessage,
		Text: "hello from consumer 1",
	}

	err := consumer1.WriteJSON(outbound)
	if err != nil {
		t.Fatalf("error sending user_message: %v", err)
	}

	readEvents(t, consumer1, isTurnDone)

	// Give the hub a moment to finish appending the turn's tail events to the backlog before a
	// second consumer attaches — readEvents already observed turn_done over the same connection
	// the backlog append happens under the same lock ahead of, so this is a formality, not a race
	// this test relies on.

	// Consumer #2: attaches fresh, modeling a page reload. It should replay the whole
	// conversation so far, including consumer #1's own message.
	consumer2 := dialChat(t, server.URL)
	defer consumer2.Close()

	replayed := readEvents(t, consumer2, isTurnDone)

	var sawUserMessage bool

	for _, event := range replayed {
		if event.Type == chatprotocol.EventUserMessage {
			sawUserMessage = true

			if event.Text != outbound.Text {
				t.Fatalf("replayed user_message text = %q, want %q", event.Text, outbound.Text)
			}
		}
	}

	if !sawUserMessage {
		t.Fatalf("consumer #2's replayed backlog %+v did not include the user_message consumer #1 sent — "+
			"this is the bug: a reconnecting consumer loses its own prior messages", replayed)
	}
}

func isNewChat(event chatprotocol.Event) bool {
	return event.Type == chatprotocol.EventNewChat
}

// TestRunMainLoop_NewChatResetsBacklogAndSession proves the two effects a new_chat event must have:
// the hub's backlog is reset (a consumer attaching afterward replays only the new_chat event, none
// of the prior conversation), and the runner's session id is cleared (the next turn starts a fresh
// `claude` conversation rather than --resume-ing the discarded one).
func TestRunMainLoop_NewChatResetsBacklogAndSession(t *testing.T) {
	chatHub := hub.New(t.TempDir())
	broker := permissions.NewBroker(chatHub.Broadcast)

	server := httptest.NewServer(chatHub)
	defer server.Close()

	stubBinary := writeStubClaude(t)
	envStore := envdrop.New(t.TempDir())
	runner := claudecli.NewRunner(stubBinary, t.TempDir(), "", envStore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runMainLoop(ctx, chatHub, runner, broker)

	// Establish a conversation, giving the runner a session id (the stub claude script always
	// reports "test-session" — see writeStubClaude).
	consumer1 := dialChat(t, server.URL)
	defer consumer1.Close()

	err := consumer1.WriteJSON(chatprotocol.Event{
		Type: chatprotocol.EventUserMessage,
		Text: "hello from consumer 1",
	})
	if err != nil {
		t.Fatalf("error sending user_message: %v", err)
	}

	readEvents(t, consumer1, isTurnDone)

	if runner.SessionId() == "" {
		t.Fatalf("runner has no session id after the first turn; stub claude should have set one")
	}

	// Start a new chat.
	err = consumer1.WriteJSON(chatprotocol.Event{Type: chatprotocol.EventNewChat})
	if err != nil {
		t.Fatalf("error sending new_chat: %v", err)
	}

	// Wait for runMainLoop to have actually processed the new_chat event (rather than sleeping):
	// consumer1 seeing it delivered live means the synchronous EventNewChat case — SetSessionId("")
	// followed by chatHub.Reset(event) — has already run to completion, since both happen in that
	// order on the same goroutine ahead of the broadcast reaching any consumer.
	readEvents(t, consumer1, isNewChat)

	if runner.SessionId() != "" {
		t.Fatalf("runner.SessionId() = %q after new_chat, want empty — session should be dropped so "+
			"the next turn does not --resume the discarded conversation", runner.SessionId())
	}

	// A freshly attached consumer, modeling a page reload after the new chat started, should replay
	// only the new_chat event — none of the earlier user_message/system_init/turn_done.
	consumer2 := dialChat(t, server.URL)
	defer consumer2.Close()

	replayed := readEvents(t, consumer2, isNewChat)

	if len(replayed) != 1 || replayed[0].Type != chatprotocol.EventNewChat {
		t.Fatalf("consumer #2's replayed backlog = %+v, want exactly one new_chat event — "+
			"the backlog should have been reset, not appended to", replayed)
	}
}
