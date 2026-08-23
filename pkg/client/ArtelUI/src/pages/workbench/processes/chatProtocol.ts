// Hand-written TypeScript mirror of the backend's ChatEvent wire protocol
// (internal/chatprotocol/events.go). There is no proto/codegen for this WebSocket —
// keep this file's field names and JSON shape byte-identical to the Go struct's `json`
// tags on any backend change.

export type EventType =
    | "system_init"
    | "user_message"
    | "assistant_text_delta"
    | "assistant_text_done"
    | "tool_call_started"
    | "tool_call_result"
    | "permission_request"
    | "permission_decision"
    | "turn_done"
    | "error"
    | "auth_complete"
    | "new_chat"

export type PermissionDecision = "allow_once" | "allow_always" | "deny"

// Single envelope type for every message on the chat WebSocket, in both directions.
// Only the fields relevant to `type` are populated; the rest are undefined.
export interface ChatEvent {
    type: EventType

    // Correlation id: ties tool_call_started <-> tool_call_result, and
    // permission_request <-> permission_decision.
    id?: string

    // assistant_text_delta / assistant_text_done / user_message / error
    text?: string

    // tool_call_started / tool_call_result / permission_request
    tool_name?: string
    input?: unknown
    input_summary?: string
    output?: string
    is_error?: boolean

    // permission_request / permission_decision
    options?: PermissionDecision[]
    decision?: PermissionDecision

    // turn_done
    session_id?: string
    cost_usd?: number
}
