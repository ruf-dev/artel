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

describe("useChatSession", () => {
    beforeEach(() => {
        MockWebSocket.instances = []
        vi.stubGlobal("WebSocket", MockWebSocket)
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it("connects to the workbench chat route and folds an incoming event into items", async () => {
        const {result, unmount} = renderHook(() => useChatSession("v1"))

        expect(MockWebSocket.instances).toHaveLength(1)
        expect(MockWebSocket.instances[0].url).toContain("/api/vaults/workbench/v1/terminal/")

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            MockWebSocket.instances[0].emit({type: "assistant_text_done", id: "a1", text: "hello"})
        })
        await waitFor(() => expect(result.current.items).toHaveLength(1))
        expect(result.current.items[0]).toMatchObject({kind: "assistant_message", text: "hello"})

        unmount()
    })

    it("optimistically applies and sends a user message even while the socket is still connecting", () => {
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            result.current.sendMessage("hi there")
        })

        expect(result.current.items).toEqual([{kind: "user_message", key: "user_message-0", text: "hi there"}])
        // Socket not open yet -> nothing should have been written to the wire.
        expect(MockWebSocket.instances[0].sent).toHaveLength(0)
    })

    it("sends over the wire once open", () => {
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            MockWebSocket.instances[0].open()
        })

        act(() => {
            result.current.sendPermissionDecision("p1", "allow_once")
        })

        expect(MockWebSocket.instances[0].sent).toEqual([
            JSON.stringify({type: "permission_decision", id: "p1", decision: "allow_once"}),
        ])
    })

    it("reconnects with a fresh socket after the connection drops", async () => {
        const {result} = renderHook(() => useChatSession("v1"))
        const first = MockWebSocket.instances[0]

        act(() => {
            first.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            first.onclose?.()
        })
        expect(result.current.status).toBe("reconnecting")

        await waitFor(() => expect(MockWebSocket.instances).toHaveLength(2), {timeout: 3000})
    }, 10000)
})
