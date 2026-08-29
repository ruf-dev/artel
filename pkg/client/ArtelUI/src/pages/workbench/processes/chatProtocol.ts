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

    // Monotonically increasing sequence number stamped by the bridge once per event
    // for the lifetime of the bridge process; used to deduplicate replayed events
    // on reconnect.
    seq?: number

    // assistant_text_delta / assistant_text_done / user_message / error
    text?: string

    // user_message only. Vault-relative paths of files attached to this message —
    // display metadata only, the referenced content is already folded into `text`
    // via the attachments preamble (see workbenchAttachments.ts's
    // withAttachmentsPreamble), so this field changes nothing about what's sent to
    // the model or persisted as message content.
    attachments?: string[]

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

    // Attributes which model produced an assistant_text_delta/assistant_text_done/turn_done
    // event. Populated only by Simple Chat's in-process engine — the Docker workbench bridge
    // never sets this field.
    model?: string
}
