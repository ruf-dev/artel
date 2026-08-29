// Package chatprotocol defines the ChatEvent wire protocol exchanged over the workbench
// chat WebSocket between the in-container bridge process and any consumer (web frontend,
// Telegram bot relay). Canonical source — deploy/workbench/bridge/internal/chatprotocol/events.go
// is a hand-mirrored copy for the bridge's standalone Go module (it can't import this module
// directly; keep the two files byte-identical on any change).
package chatprotocol

import "encoding/json"

type EventType string

const (
	// EventSystemInit: bridge → consumer, carries no payload fields.
	EventSystemInit EventType = "system_init"
	// EventUserMessage: consumer → bridge, carries Text.
	EventUserMessage EventType = "user_message"
	// EventNewChat: consumer → bridge, carries no payload fields. Broadcast back to every attached
	// consumer via Hub.Reset (see hub.go) as the sole backlog entry going forward — a consumer
	// reconnecting afterward replays only the new conversation, not the one this discarded.
	EventNewChat EventType = "new_chat"
	// EventAssistantTextDelta: bridge → consumer, carries Text and ID.
	EventAssistantTextDelta EventType = "assistant_text_delta"
	// EventAssistantTextDone: bridge → consumer, carries Text and ID.
	EventAssistantTextDone EventType = "assistant_text_done"
	// EventToolCallStarted: bridge → consumer, carries ID/ToolName/InputSummary.
	EventToolCallStarted EventType = "tool_call_started"
	// EventToolCallResult: bridge → consumer, carries ID/ToolName/Output/IsError.
	EventToolCallResult EventType = "tool_call_result"
	// EventPermissionRequest: bridge → consumer, carries ID/ToolName/InputSummary/Options.
	EventPermissionRequest EventType = "permission_request"
	// EventPermissionDecision: consumer → bridge, carries ID/Decision.
	EventPermissionDecision EventType = "permission_decision"
	// EventTurnDone: bridge → consumer, carries SessionID/CostUSD.
	EventTurnDone EventType = "turn_done"
	// EventError: bridge → consumer, carries Text.
	EventError EventType = "error"
	// EventAuthLink: bridge → consumer, carries URL.
	EventAuthLink EventType = "auth_link"
	// EventAuthCodeNeeded: bridge → consumer, carries no payload fields.
	EventAuthCodeNeeded EventType = "auth_code_needed"
	// EventAuthCodeSubmit: consumer → bridge, carries Code.
	EventAuthCodeSubmit EventType = "auth_code_submit"
	// EventAuthComplete: bridge → consumer, carries no payload fields. Broadcast once, immediately
	// after the bridge's startup auth check finishes — whether a subscription login actually ran,
	// was skipped (api_key mode / token already present), or failed — since the bridge falls
	// through to normal operation in every case. This is the one unambiguous signal a consumer can
	// use to switch off an "auth pending" UI, instead of inferring completion from the absence of
	// further auth_link/auth_code_needed events.
	EventAuthComplete EventType = "auth_complete"
)

type PermissionDecision string

const (
	DecisionAllowOnce   PermissionDecision = "allow_once"
	DecisionAllowAlways PermissionDecision = "allow_always"
	DecisionDeny        PermissionDecision = "deny"
)

// Event is the single envelope type for every message on the chat WebSocket, in both
// directions. Only the fields relevant to Type are populated; the rest are zero/omitted.
type Event struct {
	Type EventType `json:"type"`

	// Correlation id: ties tool_call_started -> tool_call_result, and
	// permission_request -> permission_decision.
	ID string `json:"id,omitempty"`

	// assistant_text_delta / assistant_text_done / user_message / error
	Text string `json:"text,omitempty"`

	// Attachments carries the vault-relative paths a user_message attached, display-only —
	// separate from Text so the attached-files preamble is no longer folded into the message
	// text itself. Populated only on user_message; every other event type leaves it empty.
	Attachments []string `json:"attachments,omitempty"`

	// tool_call_started / tool_call_result / permission_request
	ToolName     string          `json:"tool_name,omitempty"`
	Input        json.RawMessage `json:"input,omitempty"`
	InputSummary string          `json:"input_summary,omitempty"`
	Output       string          `json:"output,omitempty"`
	IsError      bool            `json:"is_error,omitempty"`

	// permission_request / permission_decision
	Options  []PermissionDecision `json:"options,omitempty"`
	Decision PermissionDecision   `json:"decision,omitempty"`

	// turn_done
	SessionID string  `json:"session_id,omitempty"`
	CostUSD   float64 `json:"cost_usd,omitempty"`

	// auth_link / auth_code_submit
	URL  string `json:"url,omitempty"`
	Code string `json:"code,omitempty"`

	// Seq is a monotonically increasing sequence number assigned by Hub.publish, used by a
	// consumer to skip events it has already applied when the backlog is replayed to it
	// (on first attach or after a reconnect), instead of re-rendering/duplicating them.
	Seq uint64 `json:"seq,omitempty"`
}
