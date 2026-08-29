import {act, renderHook, waitFor} from "@testing-library/react"
import {afterEach, beforeEach, describe, expect, it, vi} from "vitest"

import {useSimpleChatSession} from "@/pages/workbench/processes/useSimpleChatSession.ts"

// Mirrors useChatSession.test.ts's MockWebSocket harness - no shared test helper
// exists for it yet, so it's duplicated here rather than invented as a new shared
// utility for a single extra caller.
class MockWebSocket {
    static OPEN = 1
    static instances: MockWebSocket[] = []

    readyState = 0
    sent: string[] = []
    onopen: (() => void) | null = null
    onmessage: ((event: {data: string}) => void) | null = null
    onclose: (() => void) | null = null
    onerror: (() => void) | null = null

    constructor(public url: string) {
        MockWebSocket.instances.push(this)
    }

    send(data: string) {
        this.sent.push(data)
    }

    close() {
        this.readyState = 3
        this.onclose?.()
    }

    open() {
        this.readyState = MockWebSocket.OPEN
        this.onopen?.()
    }

    emit(event: unknown) {
        this.onmessage?.({data: JSON.stringify(event)})
    }
}

describe("useSimpleChatSession", () => {
    beforeEach(() => {
        MockWebSocket.instances = []
        vi.stubGlobal("WebSocket", MockWebSocket)
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it("connects to the simple-chat route for the given chatId", async () => {
        const {result} = renderHook(() => useSimpleChatSession("c1", "anthropic/claude-sonnet-4"))

        expect(MockWebSocket.instances).toHaveLength(1)
        expect(MockWebSocket.instances[0].url).toContain("/api/simple-chats/c1/ws")

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))
    })

    it("sendMessage stamps the current model onto the outbound user_message", async () => {
        const {result} = renderHook(() => useSimpleChatSession("c1", "anthropic/claude-sonnet-4"))
        const ws = MockWebSocket.instances[0]

        act(() => {
            ws.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            result.current.sendMessage("hello")
        })

        await waitFor(() => expect(result.current.items).toHaveLength(1))
        expect(result.current.items[0]).toMatchObject({kind: "user_message", text: "hello"})

        expect(ws.sent).toHaveLength(1)
        const sentMessage = JSON.parse(ws.sent[0])
        expect(sentMessage).toMatchObject({
            type: "user_message",
            text: "hello",
            model: "anthropic/claude-sonnet-4",
        })
        expect(sentMessage.id).toBeDefined()
    })

    it("resendMessage reuses the id, stamps the model, and truncates trailing items", async () => {
        const {result} = renderHook(() => useSimpleChatSession("c1", "anthropic/claude-sonnet-4"))
        const ws = MockWebSocket.instances[0]

        act(() => {
            ws.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            result.current.sendMessage("hello")
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))
        const messageId = JSON.parse(ws.sent[0]).id

        // The turn fails: a trailing error item is left dangling after the message.
        act(() => {
            ws.emit({type: "error", text: "turn failed"})
        })
        await waitFor(() => expect(result.current.items).toHaveLength(2))
        expect(result.current.items[1]).toMatchObject({kind: "error"})

        act(() => {
            result.current.resendMessage(messageId, "hello")
        })

        // The trailing error is gone, and the resend frame reuses the id and stamps
        // the model, same as sendMessage.
        expect(result.current.items).toHaveLength(1)
        expect(result.current.items[0]).toMatchObject({kind: "user_message", text: "hello", id: messageId})
        expect(ws.sent).toHaveLength(2)
        expect(JSON.parse(ws.sent[1])).toEqual({
            type: "user_message",
            text: "hello",
            id: messageId,
            model: "anthropic/claude-sonnet-4",
        })

        // The backend echoes the resent message back with the same id: dedup
        // absorbs it rather than appending a duplicate bubble.
        act(() => {
            ws.emit({type: "user_message", text: "hello", id: messageId, model: "anthropic/claude-sonnet-4", seq: 1})
        })
        expect(result.current.items).toHaveLength(1)
    })
})

describe("useSimpleChatSession - attachments", () => {
    beforeEach(() => {
        MockWebSocket.instances = []
        vi.stubGlobal("WebSocket", MockWebSocket)
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it("includes attachments on the dispatched WS payload when sendMessage is passed them", async () => {
        const {result} = renderHook(() => useSimpleChatSession("c1", "anthropic/claude-sonnet-4"))
        const ws = MockWebSocket.instances[0]

        act(() => {
            ws.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            result.current.sendMessage("check these out", ["notes/a.md", "b.md"])
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))
        expect(result.current.items[0]).toMatchObject({
            kind: "user_message", text: "check these out", attachments: ["notes/a.md", "b.md"],
        })
        expect(JSON.parse(ws.sent[0])).toMatchObject({
            attachments: ["notes/a.md", "b.md"], model: "anthropic/claude-sonnet-4",
        })
    })

    it("includes attachments on the dispatched WS payload when resendMessage is passed them", async () => {
        const {result} = renderHook(() => useSimpleChatSession("c1", "anthropic/claude-sonnet-4"))
        const ws = MockWebSocket.instances[0]

        act(() => {
            ws.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            result.current.sendMessage("hello")
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))
        const messageId = JSON.parse(ws.sent[0]).id

        act(() => {
            result.current.resendMessage(messageId, "hello", ["notes/a.md"])
        })
        // Dispatched wire frame carries the new attachments (the local item stays
        // as truncateAfterUserMessage/dedup left it - resend always reuses the id).
        expect(JSON.parse(ws.sent[1])).toMatchObject({id: messageId, attachments: ["notes/a.md"]})
    })
})

describe("useSimpleChatSession - seeding items once persisted history resolves", () => {
    beforeEach(() => {
        MockWebSocket.instances = []
        vi.stubGlobal("WebSocket", MockWebSocket)
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it("picks up initialItems arriving after connect, without a chatId change", async () => {
        // Simulates a reload/direct-nav to an already-known chatId: the connect
        // effect fires immediately with an empty initialItems (GetSimpleChat still
        // loading), then a rerender later supplies the resolved history.
        const {result, rerender} = renderHook(
            ({initialItems}) => useSimpleChatSession("c1", "anthropic/claude-sonnet-4", initialItems),
            {initialProps: {initialItems: [] as Parameters<typeof useSimpleChatSession>[2]}},
        )

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))
        expect(result.current.items).toEqual([])

        const persisted: Parameters<typeof useSimpleChatSession>[2] = [
            {kind: "user_message", key: "user_message-0", id: "m1", text: "earlier message"},
        ]
        rerender({initialItems: persisted})

        await waitFor(() => expect(result.current.items).toEqual(persisted))
    })

    it("does not let a later-resolving initialItems clobber a live item that arrived first", async () => {
        const {result, rerender} = renderHook(
            ({initialItems}) => useSimpleChatSession("c1", "anthropic/claude-sonnet-4", initialItems),
            {initialProps: {initialItems: [] as Parameters<typeof useSimpleChatSession>[2]}},
        )
        const ws = MockWebSocket.instances[0]

        act(() => {
            ws.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        // A live event arrives before the persisted-history fetch resolves (e.g. the
        // user reloaded mid-turn).
        act(() => {
            ws.emit({type: "user_message", text: "live message", id: "live-1", seq: 1})
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))

        const persisted: Parameters<typeof useSimpleChatSession>[2] = [
            {kind: "user_message", key: "user_message-0", id: "m1", text: "earlier message"},
        ]
        rerender({initialItems: persisted})

        // The prev.length === 0 guard means the late initialItems is dropped rather
        // than overwriting the already-populated (live) items.
        expect(result.current.items).toHaveLength(1)
        expect(result.current.items[0]).toMatchObject({kind: "user_message", id: "live-1", text: "live message"})
    })
})
