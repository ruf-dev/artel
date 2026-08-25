import {SimpleChatMessage} from "@/processes/SimpleChat.ts"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

// Maps a Simple Chat thread's persisted message rows (SimpleChatMessage — one flat
// row per turn, from GetSimpleChat) into the same render-ready ChatItem shape the
// live WebSocket session folds ChatEvents into (chatReducer.ts's applyEvent), so
// ChatMessageList can render a resumed thread's history exactly like a live one.
//
// Assumption (backend `role` values not finalized as of writing — the backend
// track is building the persistence layer in parallel): "user" / "assistant" map
// 1:1 to the live event kinds; a tool invocation is persisted as a single flat row
// (no separate started/result rows to correlate, unlike the live tool_call_started
// / tool_call_result pair), so it maps straight to a `done: true` tool_call item
// with `content` as the output and `toolName` as the tool name; "error" maps to an
// error item. Any other/unknown role is dropped rather than rendered as junk.
export function simpleChatMessagesToItems(messages: SimpleChatMessage[]): ChatItem[] {
    return messages.reduce<ChatItem[]>((items, m, index) => {
        const item = messageToItem(m, index)
        return item ? [...items, item] : items
    }, [])
}

function messageToItem(m: SimpleChatMessage, index: number): ChatItem | null {
    switch (m.role) {
        case "user":
            return {kind: "user_message", key: `user_message-${m.id || index}`, text: m.content, id: m.id}
        case "assistant":
            return {
                kind: "assistant_message",
                key: `assistant_message-${m.id || index}`,
                id: m.id || `assistant-${index}`,
                text: m.content,
                done: true,
            }
        case "tool":
        case "tool_call":
            return {
                kind: "tool_call",
                key: `tool_call-${m.id || index}`,
                id: m.id || `tool-${index}`,
                toolName: m.toolName ?? "",
                output: m.content,
                isError: m.isError,
                done: true,
            }
        case "error":
            return {kind: "error", key: `error-${m.id || index}`, text: m.content}
        default:
            return null
    }
}
