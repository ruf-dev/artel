// Package permissions implements the human-in-the-loop tool approval path.
//
// The `claude` CLI has no `--permission-prompt-tool` and no MCP-based approval flag. The one
// extension point that can block a tool call on an out-of-band decision is the **PreToolUse hook
// of type "http"**: Claude Code POSTs the tool call to a URL and blocks synchronously on the
// JSON response. This package is both ends of that: the HTTP endpoint Claude Code calls, and the
// broker that parks the request until a WebSocket consumer answers.
//
// Registration lives in ../claudesettings, and both must agree on the URL and on the relationship
// between the two timeouts: the hook's own timeout must be comfortably longer than DecisionWait,
// so a request that nobody answers is resolved *here*, as an explicit deny, rather than by Claude
// Code giving up on the hook. That matters because a hook that errors or times out is treated as
// non-blocking by Claude Code and falls through to its normal permission flow — which is not a
// decision this bridge ever wants made on its behalf.
//
// Request and response shapes (both confirmed against CLI 2.1.234):
//
//	POST -> {"session_id","prompt_id","transcript_path","cwd","permission_mode",
//	         "hook_event_name":"PreToolUse","tool_name","tool_input","tool_use_id"}
//	 <-     {"hookSpecificOutput":{"hookEventName":"PreToolUse",
//	         "permissionDecision":"allow"|"deny","permissionDecisionReason":"..."}}
package permissions

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"workbenchbridge/internal/chatprotocol"
)

const (
	// ListenAddr is the loopback-only address the hook endpoint binds to. Loopback because the
	// only legitimate caller is a `claude` process the bridge itself spawned inside this same
	// container — nothing outside it has any business approving tool calls, and the port is
	// deliberately not published by internal/clients/workbenchdocker's CreateContainer.
	ListenAddr = "127.0.0.1:8787"

	// HookPath is the path component of the hook URL, shared with the settings writer.
	HookPath = "/hooks/pre-tool-use"

	// HookUrl is what gets written into the `claude` settings file.
	HookUrl = "http://" + ListenAddr + HookPath

	// DecisionWait is how long a parked hook request waits for a consumer's decision before the
	// broker itself denies it. HookTimeoutSeconds is what the settings file tells Claude Code to
	// allow; the gap between them is the safety margin described in the package doc.
	DecisionWait       = 290 * time.Second
	HookTimeoutSeconds = 300
)

// hookRequest is the subset of Claude Code's PreToolUse payload the broker uses.
type hookRequest struct {
	ToolName  string          `json:"tool_name"`
	ToolInput json.RawMessage `json:"tool_input"`
	ToolUseId string          `json:"tool_use_id"`
}

// hookResponse is Claude Code's PreToolUse decision envelope.
type hookResponse struct {
	HookSpecificOutput hookSpecificOutput `json:"hookSpecificOutput"`
}

type hookSpecificOutput struct {
	HookEventName            string `json:"hookEventName"`
	PermissionDecision       string `json:"permissionDecision"`
	PermissionDecisionReason string `json:"permissionDecisionReason"`
}

// Broker parks PreToolUse hook requests until a consumer decides them, and remembers
// "allow always" answers per tool name for the lifetime of the process.
type Broker struct {
	// broadcast publishes a permission_request to every attached consumer. Injected rather than
	// depending on the hub directly, so the broker is testable with no WebSocket in sight.
	broadcast func(event chatprotocol.Event)

	// wait is how long a parked request blocks; a field rather than the DecisionWait constant so
	// tests can exercise the timeout path without waiting minutes.
	wait time.Duration

	mu      sync.Mutex
	pending map[string]chan chatprotocol.PermissionDecision
	// alwaysAllowed holds every tool name a consumer answered allow_always for. Per tool name
	// (not per tool name + input) and process-lifetime-scoped, matching what the decision means
	// to the person clicking it: "stop asking me about this tool in this session".
	alwaysAllowed map[string]struct{}
}

// NewBroker returns a Broker publishing permission requests through broadcast.
func NewBroker(broadcast func(event chatprotocol.Event)) *Broker {
	return &Broker{
		broadcast:     broadcast,
		wait:          DecisionWait,
		pending:       map[string]chan chatprotocol.PermissionDecision{},
		alwaysAllowed: map[string]struct{}{},
	}
}

// SetWait overrides how long a parked request waits before being denied. For tests.
func (b *Broker) SetWait(wait time.Duration) {
	b.wait = wait
}

// Decide records a consumer's answer to the permission request with the given id, unblocking the
// parked hook request. An id nobody is waiting on (a duplicate click, a decision arriving after
// the timeout already denied it) is ignored — deliberately: the alternative, letting a late
// "allow" resurrect an already-denied call, would approve a tool call whose result can no longer
// reach the model.
//
// allow_always is recorded before unblocking, so a tool call the model makes immediately
// afterwards is answered from memory rather than racing the consumer.
//
// Decide returns true if id had a live waiter that this call unblocked, false if it was a no-op
// (id wasn't in b.pending). Decide itself never broadcasts — it only sees an id and a decision,
// not the original chatprotocol.Event, so it stays decoupled from that shape. The caller is
// responsible for broadcasting the decision when (and only when) this returns true, so a decision
// that actually applied is the only kind that gets persisted/replayed (see main.go's
// EventPermissionDecision case).
func (b *Broker) Decide(id string, decision chatprotocol.PermissionDecision) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	waiter, ok := b.pending[id]
	if !ok {
		return false
	}

	delete(b.pending, id)

	waiter <- decision
	close(waiter)

	return true
}

// ServeHTTP is the PreToolUse hook endpoint. It always answers 2xx with an explicit allow or
// deny — never a non-2xx or a timeout, both of which Claude Code treats as a non-blocking hook
// failure and resolves through its own permission flow instead.
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var request hookRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("permissions: error decoding hook request: %v", err)
		writeDecision(w, "deny", "the workbench bridge could not read this tool call")

		return
	}

	decision, reason := b.resolve(r, request)

	writeDecision(w, decision, reason)
}

// resolve answers one hook request: from the allow-always memory if possible, otherwise by
// publishing a permission_request and parking until a consumer decides, the wait elapses, or
// Claude Code hangs up (r.Context() cancelling, e.g. because the user interrupted the turn).
func (b *Broker) resolve(r *http.Request, request hookRequest) (string, string) {
	b.mu.Lock()
	_, always := b.alwaysAllowed[request.ToolName]
	b.mu.Unlock()

	if always {
		return "allow", "previously approved for the rest of this session"
	}

	id := request.ToolUseId
	if id == "" {
		id = randomId()
	}

	waiter := b.park(id)

	event := chatprotocol.Event{
		Type:         chatprotocol.EventPermissionRequest,
		ID:           id,
		ToolName:     request.ToolName,
		Input:        request.ToolInput,
		InputSummary: summarise(request.ToolInput),
		Options: []chatprotocol.PermissionDecision{
			chatprotocol.DecisionAllowOnce,
			chatprotocol.DecisionAllowAlways,
			chatprotocol.DecisionDeny,
		},
	}

	b.broadcast(event)

	timer := time.NewTimer(b.wait)
	defer timer.Stop()

	select {
	case decision := <-waiter:
		return b.apply(request.ToolName, decision)
	case <-timer.C:
		b.abandon(id)
		b.broadcast(chatprotocol.Event{
			Type:     chatprotocol.EventPermissionDecision,
			ID:       id,
			Decision: chatprotocol.DecisionDeny,
		})

		return "deny", "no one approved this tool call in time"
	case <-r.Context().Done():
		b.abandon(id)
		b.broadcast(chatprotocol.Event{
			Type:     chatprotocol.EventPermissionDecision,
			ID:       id,
			Decision: chatprotocol.DecisionDeny,
		})

		return "deny", "the request was cancelled before it was approved"
	}
}

// apply turns a consumer's decision into a hook answer, recording allow_always on the way
// through.
func (b *Broker) apply(toolName string, decision chatprotocol.PermissionDecision) (string, string) {
	switch decision {
	case chatprotocol.DecisionAllowAlways:
		b.mu.Lock()
		b.alwaysAllowed[toolName] = struct{}{}
		b.mu.Unlock()

		return "allow", "approved for the rest of this session"
	case chatprotocol.DecisionAllowOnce:
		return "allow", "approved"
	default:
		return "deny", "denied by the user"
	}
}

// park registers a waiter for id. The channel is buffered so Decide never blocks on a waiter that
// has already given up.
func (b *Broker) park(id string) chan chatprotocol.PermissionDecision {
	waiter := make(chan chatprotocol.PermissionDecision, 1)

	b.mu.Lock()
	b.pending[id] = waiter
	b.mu.Unlock()

	return waiter
}

// abandon drops a waiter that timed out or was cancelled, so a later Decide for it is a no-op
// rather than a send into a channel nobody reads.
func (b *Broker) abandon(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.pending, id)
}

// Pending reports how many hook requests are currently parked. For tests.
func (b *Broker) Pending() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	return len(b.pending)
}

func writeDecision(w http.ResponseWriter, decision, reason string) {
	response := hookResponse{
		HookSpecificOutput: hookSpecificOutput{
			HookEventName:            "PreToolUse",
			PermissionDecision:       decision,
			PermissionDecisionReason: reason,
		},
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("permissions: error writing hook decision: %v", err)
	}
}

// summarise renders a tool input compactly for display alongside the tool name in the approval
// prompt. The full input travels untruncated in Event.Input.
func summarise(input json.RawMessage) string {
	const limit = 400

	if len(input) == 0 {
		return ""
	}

	runes := []rune(string(input))
	if len(runes) <= limit {
		return string(runes)
	}

	return string(runes[:limit]) + "… (truncated)"
}

// randomId produces a correlation id for the rare hook request that arrives without a
// tool_use_id, so two such requests can never collide in the pending map.
func randomId() string {
	buf := make([]byte, 16)

	_, err := rand.Read(buf)
	if err != nil {
		return "perm-" + time.Now().Format(time.RFC3339Nano)
	}

	return "perm-" + hex.EncodeToString(buf)
}
