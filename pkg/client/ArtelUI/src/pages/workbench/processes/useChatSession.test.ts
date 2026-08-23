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

    it("does nothing when sending a user message while the socket is still connecting", () => {
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            result.current.sendMessage("hi there")
        })

        // Socket not open yet -> nothing should have been written to the wire, and no
        // optimistic item should have been applied locally either (a send that never
        // went out must not render as if it had).
        expect(MockWebSocket.instances[0].sent).toHaveLength(0)
        expect(result.current.items).toEqual([])
    })

    it("silently absorbs a stale auth_code_needed replayed from the bridge backlog", () => {
        // The bridge no longer emits this event type (the setup-token pty flow that
        // produced it is gone), but a still-running bridge's in-memory backlog can
        // still replay old instances of it on reconnect — it must fold to nothing
        // rather than surfacing a permanent "authorize" item.
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            MockWebSocket.instances[0].open()
        })
        act(() => {
            MockWebSocket.instances[0].emit({type: "auth_code_needed"})
        })

        expect(result.current.items).toEqual([])
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

describe("useChatSession authComplete", () => {
    beforeEach(() => {
        MockWebSocket.instances = []
        vi.stubGlobal("WebSocket", MockWebSocket)
    })

    afterEach(() => {
        vi.unstubAllGlobals()
    })

    it("marks authComplete once an auth_complete event arrives", async () => {
        const {result} = renderHook(() => useChatSession("v1"))

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        expect(result.current.authComplete).toBe(false)

        act(() => {
            MockWebSocket.instances[0].emit({type: "auth_complete"})
        })
        await waitFor(() => expect(result.current.authComplete).toBe(true))
    })

    it("keeps authComplete true across a reconnect, same as items", async () => {
        // A dropped/reconnecting socket doesn't wipe `items` either (only a fresh
        // effect run — i.e. a vaultId change — does that): the hub replays backlog
        // to newly-attached consumers, so authComplete should survive a reconnect
        // the same way the already-applied items do.
        const {result} = renderHook(() => useChatSession("v1"))
        const first = MockWebSocket.instances[0]

        act(() => {
            first.open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            first.emit({type: "auth_complete"})
        })
        await waitFor(() => expect(result.current.authComplete).toBe(true))

        act(() => {
            first.onclose?.()
        })
        expect(result.current.status).toBe("reconnecting")
        expect(result.current.authComplete).toBe(true)

        await waitFor(() => expect(MockWebSocket.instances).toHaveLength(2), {timeout: 3000})
    }, 10000)

    it("resets authComplete when the connect-effect re-runs for a new vaultId", async () => {
        const {result, rerender} = renderHook(({vaultId}) => useChatSession(vaultId), {
            initialProps: {vaultId: "v1"},
        })

        act(() => {
            MockWebSocket.instances[0].open()
        })
        await waitFor(() => expect(result.current.status).toBe("open"))

        act(() => {
            MockWebSocket.instances[0].emit({type: "auth_complete"})
        })
        await waitFor(() => expect(result.current.authComplete).toBe(true))

        rerender({vaultId: "v2"})

        expect(result.current.authComplete).toBe(false)
    })
})
