import {ChatEvent, PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

// Derived, render-ready items built by folding the raw ChatEvent stream (both directions
// — server pushes and our own outgoing sends run through the same reducer, see
// useChatSession.ts) into a flat list keyed for React. One item per visible unit in the
// conversation; system_init/turn_done carry no visible unit and are absorbed silently.
export type ChatItem =
    | {kind: "user_message"; key: string; text: string; id?: string}
    | {kind: "assistant_message"; key: string; id: string; text: string; done: boolean}
    | {
        kind: "tool_call"
        key: string
        id: string
        toolName: string
        input?: unknown
        inputSummary?: string
        output?: string
        isError?: boolean
        done: boolean
    }
    | {
        kind: "permission_request"
        key: string
        id: string
        toolName: string
        input?: unknown
        inputSummary?: string
        options: PermissionDecision[]
        decision?: PermissionDecision
    }
    | {kind: "error"; key: string; text: string}

function applyUserMessage(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    // Deduplicate by id if present: if an existing item with the same id is found,
    // skip appending (it's the optimistic copy already rendered, or a duplicate replay).
    // Events with no id at all must still append normally (backward compatibility).
    if (event.id && prev.some(i => i.kind === "user_message" && i.id === event.id)) {
        return prev
    }
    const item: ChatItem = {
        kind: "user_message",
        key: `user_message-${prev.length}`,
        text: event.text ?? "",
        id: event.id,
    }
    return [...prev, item]
}

function applyAssistantDelta(prev: ChatItem[], event: ChatEvent, done: boolean): ChatItem[] {
    const id = event.id ?? `assistant-${prev.length}`
    const idx = prev.findIndex(i => i.kind === "assistant_message" && i.id === id)
    if (idx === -1) {
        const item: ChatItem = {
            kind: "assistant_message", key: `assistant_message-${id}`, id, text: event.text ?? "", done,
        }
        return [...prev, item]
    }
    return prev.map((it, i) => {
        if (i !== idx || it.kind !== "assistant_message") return it
        return done
            ? {...it, text: event.text ?? it.text, done: true}
            : {...it, text: it.text + (event.text ?? "")}
    })
}

function applyToolCallStarted(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    const id = event.id ?? `tool-${prev.length}`
    const item: ChatItem = {
        kind: "tool_call",
        key: `tool_call-${id}`,
        id,
        toolName: event.tool_name ?? "",
        input: event.input,
        inputSummary: event.input_summary,
        done: false,
    }
    return [...prev, item]
}

function applyToolCallResult(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    const id = event.id ?? ""
    const idx = prev.findIndex(i => i.kind === "tool_call" && i.id === id)
    if (idx === -1) {
        const item: ChatItem = {
            kind: "tool_call",
            key: `tool_call-${id || prev.length}`,
            id,
            toolName: event.tool_name ?? "",
            output: event.output,
            isError: event.is_error,
            done: true,
        }
        return [...prev, item]
    }
    return prev.map((it, i) => i === idx && it.kind === "tool_call"
        ? {...it, output: event.output, isError: event.is_error, done: true}
        : it)
}

function applyPermissionRequest(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    const id = event.id ?? `permission-${prev.length}`
    // Deduplicate by id: if an item with the same id already exists, skip appending.
    const idx = prev.findIndex(i => i.kind === "permission_request" && i.id === id)
    if (idx !== -1) {
        return prev
    }
    const item: ChatItem = {
        kind: "permission_request",
        key: `permission_request-${id}`,
        id,
        toolName: event.tool_name ?? "",
        input: event.input,
        inputSummary: event.input_summary,
        options: event.options ?? ["allow_once", "allow_always", "deny"],
    }
    return [...prev, item]
}

function applyPermissionDecision(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    const id = event.id ?? ""
    return prev.map(it => it.kind === "permission_request" && it.id === id
        ? {...it, decision: event.decision}
        : it)
}

// Truncates the item list to end just after the user_message with the given id -
// used by resendMessage (useChatSession.ts/useSimpleChatSession.ts) to drop any
// trailing dangling/error state before re-sending in place. A no-op (returns the
// same array reference) when no such id is found.
export function truncateAfterUserMessage(items: ChatItem[], id: string): ChatItem[] {
    const idx = items.findIndex(i => i.kind === "user_message" && i.id === id)
    return idx === -1 ? items : items.slice(0, idx + 1)
}

export function applyEvent(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    switch (event.type) {
        case "system_init":
        case "turn_done":
        case "auth_complete":
            // auth_link/auth_code_needed/auth_code_submit are handled by the default
            // case below: the bridge no longer emits the first two (the setup-token
            // pty flow that produced them is gone), but a still-running bridge's
            // in-memory event backlog can still replay old instances of them to a
            // newly-attached consumer — absorb silently rather than resurrecting a
            // permanent "authorize" card for an already-completed login.
            return prev
        case "user_message":
            return applyUserMessage(prev, event)
        case "assistant_text_delta":
            return applyAssistantDelta(prev, event, false)
        case "assistant_text_done":
            return applyAssistantDelta(prev, event, true)
        case "tool_call_started":
            return applyToolCallStarted(prev, event)
        case "tool_call_result":
            return applyToolCallResult(prev, event)
        case "permission_request":
            return applyPermissionRequest(prev, event)
        case "permission_decision":
            return applyPermissionDecision(prev, event)
        case "error":
            return [...prev, {kind: "error", key: `error-${prev.length}`, text: event.text ?? ""}]
        case "new_chat":
            // Starts a fresh conversation, mirroring the backend's Hub.Reset — that call
            // makes this event the sole backlog entry from that point on, so wiping local
            // items here keeps replay-on-reconnect consistent whether it arrives live or
            // via backlog replay.
            return []
        default:
            return prev
    }
}
