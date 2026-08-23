package permissions

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"workbenchbridge/internal/chatprotocol"
)

// recordingBroadcaster is a minimal broadcast func stand-in that records every event it's given,
// safe for concurrent use since Decide/resolve run on whatever goroutine calls them and the test
// asserting on the recording runs on another.
type recordingBroadcaster struct {
	mu     sync.Mutex
	events []chatprotocol.Event
}

func (r *recordingBroadcaster) broadcast(event chatprotocol.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

func (r *recordingBroadcaster) snapshot() []chatprotocol.Event {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]chatprotocol.Event, len(r.events))
	copy(out, r.events)

	return out
}

// find returns the first recorded event of the given type with the given id, if any.
func (r *recordingBroadcaster) find(eventType chatprotocol.EventType, id string) (chatprotocol.Event, bool) {
	for _, event := range r.snapshot() {
		if event.Type == eventType && event.ID == id {
			return event, true
		}
	}

	return chatprotocol.Event{}, false
}

func hookRequestBody(toolName, toolUseId string) string {
	return `{"tool_name":"` + toolName + `","tool_use_id":"` + toolUseId + `","tool_input":{}}`
}

func decodeHookResponse(t *testing.T, body string) hookResponse {
	t.Helper()

	var response hookResponse

	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		t.Fatalf("error decoding hook response %q: %v", body, err)
	}

	return response
}

// TestDecide_FoundWaiter proves Decide returns true and unblocks a parked waiter when a decision
// arrives for a request that's actually pending.
func TestDecide_FoundWaiter(t *testing.T) {
	rec := &recordingBroadcaster{}
	broker := NewBroker(rec.broadcast)

	waiter := broker.park("req-1")

	ok := broker.Decide("req-1", chatprotocol.DecisionAllowOnce)
	if !ok {
		t.Fatalf("Decide returned false for a request with a live waiter")
	}

	select {
	case decision := <-waiter:
		if decision != chatprotocol.DecisionAllowOnce {
			t.Fatalf("waiter received decision %q, want %q", decision, chatprotocol.DecisionAllowOnce)
		}
	case <-time.After(time.Second):
		t.Fatalf("waiter never received the decision")
	}
}

// TestDecide_UnknownId proves Decide is a no-op — and reports it via its return value — for an id
// nobody is waiting on (duplicate click, or a decision arriving after the timeout already denied
// it).
func TestDecide_UnknownId(t *testing.T) {
	rec := &recordingBroadcaster{}
	broker := NewBroker(rec.broadcast)

	ok := broker.Decide("unknown-id", chatprotocol.DecisionAllowOnce)
	if ok {
		t.Fatalf("Decide returned true for an id with no pending waiter")
	}
}

// TestResolve_DecidedWhileParked drives the full ServeHTTP hook path and simulates a consumer
// deciding while the request is still parked, proving the "allow" decision reaches the HTTP
// response. This also covers the "normal decide-while-parked path" the package previously had no
// test for.
func TestResolve_DecidedWhileParked(t *testing.T) {
	rec := &recordingBroadcaster{}
	broker := NewBroker(rec.broadcast)
	broker.SetWait(5 * time.Second)

	req := httptest.NewRequest("POST", HookPath, strings.NewReader(hookRequestBody("Bash", "req-parked")))
	recorder := httptest.NewRecorder()

	done := make(chan struct{})

	go func() {
		broker.ServeHTTP(recorder, req)
		close(done)
	}()

	// Wait for the permission_request to be broadcast (proves the request actually parked) before
	// deciding it.
	deadline := time.After(2 * time.Second)

	for {
		if _, ok := rec.find(chatprotocol.EventPermissionRequest, "req-parked"); ok {
			break
		}

		select {
		case <-deadline:
			t.Fatalf("permission_request for req-parked was never broadcast")
		case <-time.After(10 * time.Millisecond):
		}
	}

	ok := broker.Decide("req-parked", chatprotocol.DecisionAllowOnce)
	if !ok {
		t.Fatalf("Decide returned false for the parked request")
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("ServeHTTP never returned after Decide")
	}

	response := decodeHookResponse(t, recorder.Body.String())
	if response.HookSpecificOutput.PermissionDecision != "allow" {
		t.Fatalf("hook response decision = %q, want %q", response.HookSpecificOutput.PermissionDecision, "allow")
	}
}

// TestResolve_Timeout proves that a parked request nobody decides eventually resolves as a deny —
// both in the HTTP hook response and, per the bug this package is being fixed for, via a
// broadcast permission_decision event so the outcome is persisted/replayed.
func TestResolve_Timeout(t *testing.T) {
	rec := &recordingBroadcaster{}
	broker := NewBroker(rec.broadcast)
	broker.SetWait(20 * time.Millisecond)

	req := httptest.NewRequest("POST", HookPath, strings.NewReader(hookRequestBody("Bash", "req-timeout")))
	recorder := httptest.NewRecorder()

	broker.ServeHTTP(recorder, req)

	response := decodeHookResponse(t, recorder.Body.String())
	if response.HookSpecificOutput.PermissionDecision != "deny" {
		t.Fatalf("hook response decision = %q, want %q", response.HookSpecificOutput.PermissionDecision, "deny")
	}

	event, ok := rec.find(chatprotocol.EventPermissionDecision, "req-timeout")
	if !ok {
		t.Fatalf("no permission_decision event was broadcast for the timed-out request; got %+v", rec.snapshot())
	}

	if event.Decision != chatprotocol.DecisionDeny {
		t.Fatalf("broadcast permission_decision.Decision = %q, want %q", event.Decision, chatprotocol.DecisionDeny)
	}
}

// TestResolve_ContextCancelled proves the same deny-and-broadcast behavior when Claude Code hangs
// up on the request (r.Context() cancelled) instead of the wait timer elapsing.
func TestResolve_ContextCancelled(t *testing.T) {
	rec := &recordingBroadcaster{}
	broker := NewBroker(rec.broadcast)
	broker.SetWait(10 * time.Second) // long enough that only cancellation can resolve this

	ctx, cancel := context.WithCancel(context.Background())

	req := httptest.NewRequest("POST", HookPath, strings.NewReader(hookRequestBody("Bash", "req-cancelled")))
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()

	done := make(chan struct{})

	go func() {
		broker.ServeHTTP(recorder, req)
		close(done)
	}()

	// Give the request a moment to actually park before cancelling it.
	deadline := time.After(2 * time.Second)

	for {
		if _, ok := rec.find(chatprotocol.EventPermissionRequest, "req-cancelled"); ok {
			break
		}

		select {
		case <-deadline:
			t.Fatalf("permission_request for req-cancelled was never broadcast")
		case <-time.After(10 * time.Millisecond):
		}
	}

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("ServeHTTP never returned after context cancellation")
	}

	response := decodeHookResponse(t, recorder.Body.String())
	if response.HookSpecificOutput.PermissionDecision != "deny" {
		t.Fatalf("hook response decision = %q, want %q", response.HookSpecificOutput.PermissionDecision, "deny")
	}

	event, ok := rec.find(chatprotocol.EventPermissionDecision, "req-cancelled")
	if !ok {
		t.Fatalf("no permission_decision event was broadcast for the cancelled request; got %+v", rec.snapshot())
	}

	if event.Decision != chatprotocol.DecisionDeny {
		t.Fatalf("broadcast permission_decision.Decision = %q, want %q", event.Decision, chatprotocol.DecisionDeny)
	}
}

// TestAllowAlwaysMemory proves that once a decision of allow_always is recorded for a tool name,
// a later hook request for that same tool name is answered immediately from memory — no
// permission_request broadcast, no parking — matching resolve's short-circuit at the top.
func TestAllowAlwaysMemory(t *testing.T) {
	rec := &recordingBroadcaster{}
	broker := NewBroker(rec.broadcast)
	broker.SetWait(5 * time.Second)

	// First request: parks, then gets decided allow_always.
	req1 := httptest.NewRequest("POST", HookPath, strings.NewReader(hookRequestBody("Bash", "req-always-1")))
	recorder1 := httptest.NewRecorder()

	done1 := make(chan struct{})

	go func() {
		broker.ServeHTTP(recorder1, req1)
		close(done1)
	}()

	deadline := time.After(2 * time.Second)

	for {
		if _, ok := rec.find(chatprotocol.EventPermissionRequest, "req-always-1"); ok {
			break
		}

		select {
		case <-deadline:
			t.Fatalf("permission_request for req-always-1 was never broadcast")
		case <-time.After(10 * time.Millisecond):
		}
	}

	ok := broker.Decide("req-always-1", chatprotocol.DecisionAllowAlways)
	if !ok {
		t.Fatalf("Decide returned false for req-always-1")
	}

	<-done1

	response1 := decodeHookResponse(t, recorder1.Body.String())
	if response1.HookSpecificOutput.PermissionDecision != "allow" {
		t.Fatalf("first hook response decision = %q, want %q", response1.HookSpecificOutput.PermissionDecision, "allow")
	}

	// Second request for the same tool name: should be answered from memory, with no parking and
	// no new permission_request broadcast.
	req2 := httptest.NewRequest("POST", HookPath, strings.NewReader(hookRequestBody("Bash", "req-always-2")))
	recorder2 := httptest.NewRecorder()

	broker.ServeHTTP(recorder2, req2)

	response2 := decodeHookResponse(t, recorder2.Body.String())
	if response2.HookSpecificOutput.PermissionDecision != "allow" {
		t.Fatalf("second hook response decision = %q, want %q", response2.HookSpecificOutput.PermissionDecision, "allow")
	}

	if _, ok := rec.find(chatprotocol.EventPermissionRequest, "req-always-2"); ok {
		t.Fatalf("a permission_request was broadcast for req-always-2, but it should have been answered from allow_always memory")
	}

	if broker.Pending() != 0 {
		t.Fatalf("broker.Pending() = %d, want 0 after both requests resolved", broker.Pending())
	}
}
