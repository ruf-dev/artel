import {describe, expect, it} from "vitest"

import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

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

    it("resolves the most recent open auth_code_needed on auth_code_submit", () => {
        let items: ChatItem[] = []
        items = applyEvent(items, {type: "auth_link", url: "https://example.com/authorize"})
        items = applyEvent(items, {type: "auth_code_needed"})
        expect(items[1]).toEqual({kind: "auth_code_needed", key: "auth_code_needed-1", resolved: false})

        items = applyEvent(items, {type: "auth_code_submit", code: "123456"})
        expect(items[1]).toEqual({kind: "auth_code_needed", key: "auth_code_needed-1", resolved: true})
    })

    it("appends an error item", () => {
        const items = applyEvent([], {type: "error", text: "boom"})
        expect(items).toEqual([{kind: "error", key: "error-0", text: "boom"}])
    })
})
