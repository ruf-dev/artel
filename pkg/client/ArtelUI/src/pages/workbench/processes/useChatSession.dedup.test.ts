import {act, renderHook, waitFor} from "@testing-library/react"
import {afterEach, beforeEach, describe, expect, it, vi} from "vitest"

import {useChatSession} from "@/pages/workbench/processes/useChatSession.ts"

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

describe("useChatSession seq-based deduplication", () => {
    beforeEach(() => {
        MockWebSocket.instances = []
        vi.stubGlobal("WebSocket", MockWebSocket)
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it("deduplicates events with the same seq value (replayed from backlog)", async () => {
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            MockWebSocket.instances[0].emit({type: "assistant_text_done", id: "a1", text: "hello", seq: 1})
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))
        expect(result.current.items[0]).toMatchObject({kind: "assistant_message", text: "hello"})

        act(() => {
            MockWebSocket.instances[0].emit({type: "assistant_text_done", id: "a1", text: "hello", seq: 1})
        })
        expect(result.current.items).toHaveLength(1)
    })

    it("handles reconnect that replays already-seen seqs and adds new ones", async () => {
        const {result} = renderHook(() => useChatSession("v1"))
        const first = MockWebSocket.instances[0]

        act(() => {
            first.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            first.emit({type: "assistant_text_done", id: "a1", text: "msg1", seq: 1})
            first.emit({type: "assistant_text_done", id: "a2", text: "msg2", seq: 2})
            first.emit({type: "assistant_text_done", id: "a3", text: "msg3", seq: 3})
        })
        await waitFor(() => expect(result.current.items).toHaveLength(3))

        act(() => {
            first.onclose?.()
        })
        expect(result.current.status).toBe("reconnecting")

        await waitFor(() => expect(MockWebSocket.instances).toHaveLength(2), {timeout: 3000})

        const second = MockWebSocket.instances[1]
        act(() => {
            second.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            second.emit({type: "assistant_text_done", id: "a1", text: "msg1", seq: 1})
            second.emit({type: "assistant_text_done", id: "a2", text: "msg2", seq: 2})
            second.emit({type: "assistant_text_done", id: "a3", text: "msg3", seq: 3})
            second.emit({type: "assistant_text_done", id: "a4", text: "msg4", seq: 4})
        })

        await waitFor(() => expect(result.current.items).toHaveLength(4))
        expect(result.current.items.map(i => i.kind === "assistant_message" ? i.text : "")).toEqual([
            "msg1", "msg2", "msg3", "msg4"
        ])
    }, 10000)

    it("applies events with increasing seq values in order", async () => {
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            MockWebSocket.instances[0].emit({type: "assistant_text_done", id: "a1", text: "first", seq: 1})
            MockWebSocket.instances[0].emit({type: "assistant_text_done", id: "a2", text: "second", seq: 2})
            MockWebSocket.instances[0].emit({type: "assistant_text_done", id: "a3", text: "third", seq: 3})
        })

        await waitFor(() => expect(result.current.items).toHaveLength(3))
        expect(result.current.items.map(i => i.kind === "assistant_message" ? i.text : "")).toEqual([
            "first", "second", "third"
        ])
    })

    it("applies events with no seq (backward compatibility)", async () => {
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            MockWebSocket.instances[0].emit({type: "assistant_text_done", id: "a1", text: "no seq"})
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))
        expect(result.current.items[0]).toMatchObject({kind: "assistant_message", text: "no seq"})
    })
})

describe("useChatSession - user_message deduplication", () => {
    beforeEach(() => {
        MockWebSocket.instances = []
        vi.stubGlobal("WebSocket", MockWebSocket)
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it("deduplicates user_message by id when echoed back from the bridge", async () => {
        const {result} = renderHook(() => useChatSession("v1"))
        const ws = MockWebSocket.instances[0]

        act(() => {
            ws.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        // Send a user message via sendMessage (generates an id, applies optimistically).
        act(() => {
            result.current.sendMessage("hello")
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))
        expect(result.current.items[0]).toMatchObject({kind: "user_message", text: "hello"})

        // Extract the id from what was sent to the wire.
        const sentMessage = JSON.parse(ws.sent[0])
        expect(sentMessage.id).toBeDefined()
        const messageId = sentMessage.id

        // The bridge echoes the message back with the same id but a different seq.
        act(() => {
            ws.emit({type: "user_message", text: "hello", id: messageId, seq: 1})
        })

        // Should still be exactly one item, not two (the optimistic copy + the echo).
        expect(result.current.items).toHaveLength(1)
        expect(result.current.items[0]).toMatchObject({kind: "user_message", text: "hello", id: messageId})
    })
})
