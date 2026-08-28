import {describe, expect, it} from "vitest"

import {deriveTurnState} from "@/pages/workbench/processes/turnState.ts"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

describe("deriveTurnState", () => {
    it("returns 'idle' when list is empty and pendingTurn is false", () => {
        expect(deriveTurnState([], false)).toBe("idle")
    })

    it("returns 'working' when list is empty but pendingTurn is true", () => {
        expect(deriveTurnState([], true)).toBe("working")
    })

    it("returns 'working' when tail is user_message and pendingTurn is true", () => {
        const items: ChatItem[] = [
            {kind: "user_message", key: "u1", text: "hello", id: "msg-1"},
        ]
        expect(deriveTurnState(items, true)).toBe("working")
    })

    it("returns 'streaming' when tail is assistant_message with done:false, regardless of pendingTurn", () => {
        const items: ChatItem[] = [
            {kind: "assistant_message", key: "a1", id: "a1", text: "Hel", done: false},
        ]
        expect(deriveTurnState(items, false)).toBe("streaming")
        expect(deriveTurnState(items, true)).toBe("streaming")
    })

    it("returns 'idle' when tail is assistant_message with done:true and pendingTurn is false", () => {
        const items: ChatItem[] = [
            {kind: "assistant_message", key: "a1", id: "a1", text: "Hello", done: true},
        ]
        expect(deriveTurnState(items, false)).toBe("idle")
    })

    it("returns 'working' when tail is assistant_message with done:true but pendingTurn is true", () => {
        const items: ChatItem[] = [
            {kind: "assistant_message", key: "a1", id: "a1", text: "Hello", done: true},
        ]
        expect(deriveTurnState(items, true)).toBe("working")
    })

    it("returns 'awaiting_input' when tail is permission_request with no decision", () => {
        const items: ChatItem[] = [
            {
                kind: "permission_request",
                key: "p1",
                id: "p1",
                toolName: "Bash",
                inputSummary: "rm -rf /tmp/x",
                options: ["allow_once", "allow_always", "deny"],
            },
        ]
        expect(deriveTurnState(items, false)).toBe("awaiting_input")
        expect(deriveTurnState(items, true)).toBe("awaiting_input")
    })

    it("returns 'working' when tail is permission_request WITH decision and pendingTurn is true", () => {
        const items: ChatItem[] = [
            {
                kind: "permission_request",
                key: "p1",
                id: "p1",
                toolName: "Bash",
                inputSummary: "rm -rf /tmp/x",
                options: ["allow_once", "allow_always", "deny"],
                decision: "allow_once",
            },
        ]
        expect(deriveTurnState(items, true)).toBe("working")
    })

    it("returns 'idle' when tail is permission_request WITH decision and pendingTurn is false", () => {
        const items: ChatItem[] = [
            {
                kind: "permission_request",
                key: "p1",
                id: "p1",
                toolName: "Bash",
                inputSummary: "rm -rf /tmp/x",
                options: ["allow_once", "allow_always", "deny"],
                decision: "allow_once",
            },
        ]
        expect(deriveTurnState(items, false)).toBe("idle")
    })

    it("returns 'working' when tail is tool_call with done:true and pendingTurn is true", () => {
        const items: ChatItem[] = [
            {
                kind: "tool_call",
                key: "t1",
                id: "t1",
                toolName: "Read",
                output: "content",
                done: true,
            },
        ]
        expect(deriveTurnState(items, true)).toBe("working")
    })

    it("returns 'idle' when tail is error and pendingTurn is false", () => {
        const items: ChatItem[] = [
            {kind: "error", key: "e1", text: "something went wrong"},
        ]
        expect(deriveTurnState(items, false)).toBe("idle")
    })

    it("returns 'working' when tail is error and pendingTurn is true", () => {
        const items: ChatItem[] = [
            {kind: "error", key: "e1", text: "something went wrong"},
        ]
        expect(deriveTurnState(items, true)).toBe("working")
    })
})
