import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

export type TurnState = "idle" | "working" | "streaming" | "awaiting_input"

// Derives what the model/turn is doing right now from the folded item list plus the
// session's pendingTurn flag (set on send, cleared on turn_done/error — see
// useChatSession.ts). "streaming" = assistant text is actively arriving; "awaiting_input"
// = a permission_request is on screen undecided (the card itself is the CTA, no spinner);
// "working" = a turn is in flight but nothing is rendering yet (queued, reasoning, between
// a tool result and the next token); "idle" = nothing in flight.
export function deriveTurnState(items: ChatItem[], pendingTurn: boolean): TurnState {
    const last = items[items.length - 1]
    if (last?.kind === "assistant_message" && !last.done) return "streaming"
    if (last?.kind === "permission_request" && !last.decision) return "awaiting_input"
    if (pendingTurn) return "working"
    return "idle"
}
