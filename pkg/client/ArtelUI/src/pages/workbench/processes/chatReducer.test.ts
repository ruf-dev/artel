import {describe, expect, it} from "vitest"

import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatEvent} from "@/pages/workbench/processes/chatProtocol.ts"

describe("applyEvent", () => {
    it("ignores system_init and turn_done", () => {
        const items: ChatItem[] = []
        expect(applyEvent(items, {type: "system_init"})).toBe(items)
        expect(applyEvent(items, {type: "turn_done", session_id: "s1", cost_usd: 0.1})).toBe(items)
    })

    it("appends a user_message item", () => {
        const items = applyEvent([], {type: "user_message", text: "hi"})
        expect(items).toEqual([{kind: "user_message", key: "user_message-0", text: "hi"}])
    })

    it("accumulates assistant_text_delta by id and finalizes on assistant_text_done", () => {
        let items: ChatItem[] = []
        items = applyEvent(items, {type: "assistant_text_delta", id: "a1", text: "Hel"})
        items = applyEvent(items, {type: "assistant_text_delta", id: "a1", text: "lo"})
        expect(items).toEqual([
            {kind: "assistant_message", key: "assistant_message-a1", id: "a1", text: "Hello", done: false},
        ])

        items = applyEvent(items, {type: "assistant_text_done", id: "a1", text: "Hello world"})
        expect(items).toEqual([
            {kind: "assistant_message", key: "assistant_message-a1", id: "a1", text: "Hello world", done: true},
        ])
    })

    it("correlates tool_call_started and tool_call_result by id", () => {
        let items: ChatItem[] = []
        items = applyEvent(items, {
            type: "tool_call_started", id: "t1", tool_name: "Read", input: {path: "a.md"}, input_summary: "a.md",
        })
        expect(items).toEqual([{
            kind: "tool_call", key: "tool_call-t1", id: "t1", toolName: "Read",
            input: {path: "a.md"}, inputSummary: "a.md", done: false,
        }])

        items = applyEvent(items, {type: "tool_call_result", id: "t1", output: "contents", is_error: false})
        expect(items).toEqual([{
            kind: "tool_call", key: "tool_call-t1", id: "t1", toolName: "Read",
            input: {path: "a.md"}, inputSummary: "a.md", output: "contents", isError: false, done: true,
        }])
    })

    it("creates a synthetic tool_call item when a result arrives with no matching started", () => {
        const items = applyEvent([], {type: "tool_call_result", id: "t9", tool_name: "Write", output: "ok"})
        expect(items).toEqual([{
            kind: "tool_call", key: "tool_call-t9", id: "t9", toolName: "Write",
            output: "ok", isError: undefined, done: true,
        }])
    })

    it("resolves a permission_request via a matching permission_decision", () => {
        let items: ChatItem[] = []
        items = applyEvent(items, {
            type: "permission_request", id: "p1", tool_name: "Bash", input_summary: "rm -rf /tmp/x",
            options: ["allow_once", "allow_always", "deny"],
        })
        items = applyEvent(items, {type: "permission_decision", id: "p1", decision: "allow_once"})
        expect(items).toEqual([{
            kind: "permission_request", key: "permission_request-p1", id: "p1", toolName: "Bash",
            input: undefined, inputSummary: "rm -rf /tmp/x", options: ["allow_once", "allow_always", "deny"],
            decision: "allow_once",
        }])
    })

    it("silently absorbs stale auth_link/auth_code_needed/auth_code_submit events via the default case", () => {
        // These event types no longer exist in the wire protocol (the setup-token pty
        // flow that produced them is gone — see chatProtocol.ts), but a still-running
        // bridge's in-memory backlog can still replay old instances of them on
        // reconnect. Cast past the type union, the same way a genuinely unknown/future
        // event type arriving over the wire would: it must render as nothing, not
        // resurrect a permanent "authorize" card.
        const items: ChatItem[] = []
        expect(applyEvent(items, {type: "auth_link"} as unknown as ChatEvent)).toBe(items)
        expect(applyEvent(items, {type: "auth_code_needed"} as unknown as ChatEvent)).toBe(items)
        expect(applyEvent(items, {type: "auth_code_submit"} as unknown as ChatEvent)).toBe(items)
    })

    it("appends an error item", () => {
        const items = applyEvent([], {type: "error", text: "boom"})
        expect(items).toEqual([{kind: "error", key: "error-0", text: "boom"}])
    })

    it("clears all items on new_chat regardless of prior contents", () => {
        let items: ChatItem[] = []
        items = applyEvent(items, {type: "user_message", text: "hi"})
        items = applyEvent(items, {type: "assistant_text_done", id: "a1", text: "hello"})
        items = applyEvent(items, {type: "error", text: "boom"})
        expect(items).toHaveLength(3)

        expect(applyEvent(items, {type: "new_chat"})).toEqual([])
    })
})

describe("applyEvent - deduplication", () => {
    it("deduplicates user_message by id when the same id is seen twice", () => {
        let items = applyEvent([], {type: "user_message", text: "hi", id: "msg-1"})
        expect(items).toEqual([{kind: "user_message", key: "user_message-0", text: "hi", id: "msg-1"}])

        // The same id arrives again (e.g., echoed from the bridge) — should skip appending.
        items = applyEvent(items, {type: "user_message", text: "hi", id: "msg-1"})
        expect(items).toHaveLength(1)
        expect(items[0]).toMatchObject({kind: "user_message", text: "hi", id: "msg-1"})
    })

    it("appends different user_message ids separately", () => {
        let items = applyEvent([], {type: "user_message", text: "first", id: "msg-1"})
        items = applyEvent(items, {type: "user_message", text: "second", id: "msg-2"})
        expect(items).toHaveLength(2)
        expect(items[0]).toMatchObject({kind: "user_message", text: "first", id: "msg-1"})
        expect(items[1]).toMatchObject({kind: "user_message", text: "second", id: "msg-2"})
    })

    it("appends user_message without id normally (backward compatibility)", () => {
        let items = applyEvent([], {type: "user_message", text: "legacy"})
        items = applyEvent(items, {type: "user_message", text: "legacy"})
        expect(items).toHaveLength(2)
    })

    it("deduplicates permission_request by id when the same id is seen twice", () => {
        let items: ChatItem[] = []
        items = applyEvent(items, {
            type: "permission_request", id: "p1", tool_name: "Bash", input_summary: "rm -rf /tmp/x",
            options: ["allow_once", "allow_always", "deny"],
        })
        expect(items).toHaveLength(1)

        // The same id arrives again (e.g., re-broadcast from the bridge) — should skip appending.
        items = applyEvent(items, {
            type: "permission_request", id: "p1", tool_name: "Bash", input_summary: "rm -rf /tmp/x",
            options: ["allow_once", "allow_always", "deny"],
        })
        expect(items).toHaveLength(1)
        expect(items[0]).toMatchObject({kind: "permission_request", id: "p1"})
    })

    it("appends different permission_request ids separately", () => {
        let items: ChatItem[] = []
        items = applyEvent(items, {
            type: "permission_request", id: "p1", tool_name: "Bash", input_summary: "cmd1",
            options: ["allow_once", "allow_always", "deny"],
        })
        items = applyEvent(items, {
            type: "permission_request", id: "p2", tool_name: "Write", input_summary: "file.txt",
            options: ["allow_once", "allow_always", "deny"],
        })
        expect(items).toHaveLength(2)
        expect(items[0]).toMatchObject({kind: "permission_request", id: "p1"})
        expect(items[1]).toMatchObject({kind: "permission_request", id: "p2"})
    })
})
