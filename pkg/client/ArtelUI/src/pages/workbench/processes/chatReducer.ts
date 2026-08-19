import {ChatEvent, PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

// Derived, render-ready items built by folding the raw ChatEvent stream (both directions
// — server pushes and our own outgoing sends run through the same reducer, see
// useChatSession.ts) into a flat list keyed for React. One item per visible unit in the
// conversation; system_init/turn_done carry no visible unit and are absorbed silently.
export type ChatItem =
    | {kind: "user_message"; key: string; text: string}
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
    | {kind: "auth_link"; key: string; url: string}
    | {kind: "auth_code_needed"; key: string; resolved: boolean}
    | {kind: "error"; key: string; text: string}

function applyUserMessage(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    const item: ChatItem = {kind: "user_message", key: `user_message-${prev.length}`, text: event.text ?? ""}
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

function applyAuthCodeSubmit(prev: ChatItem[]): ChatItem[] {
    // No correlation id on auth events — resolve the most recent still-open prompt.
    const idx = prev.map(it => it.kind === "auth_code_needed" && !it.resolved).lastIndexOf(true)
    if (idx === -1) return prev
    return prev.map((it, i) => i === idx && it.kind === "auth_code_needed" ? {...it, resolved: true} : it)
}

export function applyEvent(prev: ChatItem[], event: ChatEvent): ChatItem[] {
    switch (event.type) {
        case "system_init":
        case "turn_done":
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
        case "auth_link":
            return [...prev, {kind: "auth_link", key: `auth_link-${prev.length}`, url: event.url ?? ""}]
        case "auth_code_needed":
            return [...prev, {kind: "auth_code_needed", key: `auth_code_needed-${prev.length}`, resolved: false}]
        case "auth_code_submit":
            return applyAuthCodeSubmit(prev)
        case "error":
            return [...prev, {kind: "error", key: `error-${prev.length}`, text: event.text ?? ""}]
        default:
            return prev
    }
}
