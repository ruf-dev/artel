import {act, renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

const h = vi.hoisted(() => ({
    openDialog: vi.fn(),
    closeDialog: vi.fn(),
    deleteChat: vi.fn(() => Promise.resolve()),
    state: {chats: [] as unknown[]},
}))

vi.mock("@/app/hooks/SimpleChat.ts", () => ({
    useSimpleChats: () => ({chats: h.state.chats, isLoading: false, refetch: vi.fn()}),
    useSimpleChatMutations: () => ({
        create: vi.fn(),
        delete: h.deleteChat,
        creating: false,
        deleting: false,
    }),
}))

vi.mock("@/app/hooks/Dialog", () => ({
    useDialog: () => ({OpenDialog: h.openDialog, CloseDialog: h.closeDialog}),
}))

vi.mock("@/app/hooks/useErrorToast.ts", () => ({useBakeError: () => vi.fn()}))

vi.mock("@/pages/workbench/components/WorkbenchHistoryTranscriptDialog/WorkbenchHistoryTranscriptDialog.tsx", () => ({
    default: () => null,
}))

vi.mock("@/pages/workbench/processes/chatHistoryApi.ts", () => ({
    listChatSessions: vi.fn(() => Promise.resolve([])),
    getChatSession: vi.fn(() => Promise.resolve([])),
}))

import * as chatHistoryApi from "@/pages/workbench/processes/chatHistoryApi.ts"
import {useWorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"

const listChatSessions = vi.mocked(chatHistoryApi.listChatSessions)
const getChatSession = vi.mocked(chatHistoryApi.getChatSession)

function apiChat(over: Record<string, unknown> = {}) {
    return {
        id: "c1",
        vaultId: "v1",
        title: "Refactor auth",
        model: "anthropic/claude",
        vaultAccess: true,
        createdAt: "",
        updatedAt: "",
        lastActivityAt: "2026-08-20T10:00:00Z",
        ...over,
    }
}

beforeEach(() => {
    vi.clearAllMocks()
    h.state.chats = []
    listChatSessions.mockResolvedValue([])
    getChatSession.mockResolvedValue([])
})

const noop = vi.fn()

describe("useWorkbenchHistory — api mode", () => {
    it("maps saved chats to rows with model · date subtitle and api source", () => {
        h.state.chats = [apiChat()]

        const {result} = renderHook(() => useWorkbenchHistory({
            mode: "api", vaultId: "v1", activeApiChatId: "c1", onSelectApiChat: noop, onNewChat: noop,
        }))

        expect(result.current.rows).toHaveLength(1)
        expect(result.current.rows[0]).toMatchObject({id: "c1", title: "Refactor auth", source: "api"})
        expect(result.current.rows[0].subtitle).toContain("anthropic/claude")
        expect(result.current.activeId).toBe("c1")
        expect(typeof result.current.onDelete).toBe("function")
    })

    it("falls back to 'Untitled chat' when a chat has no title", () => {
        h.state.chats = [apiChat({title: ""})]

        const {result} = renderHook(() => useWorkbenchHistory({
            mode: "api", vaultId: "v1", onSelectApiChat: noop, onNewChat: noop,
        }))

        expect(result.current.rows[0].title).toBe("Untitled chat")
    })

    it("opens a ConfirmDialog on delete and deletes + clears the active chat on confirm", async () => {
        const onSelectApiChat = vi.fn()
        h.state.chats = [apiChat()]

        const {result} = renderHook(() => useWorkbenchHistory({
            mode: "api", vaultId: "v1", activeApiChatId: "c1", onSelectApiChat, onNewChat: noop,
        }))

        act(() => result.current.onDelete?.("c1"))

        expect(h.openDialog).toHaveBeenCalledTimes(1)
        const el = h.openDialog.mock.calls[0][0] as {props: Record<string, unknown>}
        expect(el.props.title).toBe("Delete Chat")

        await (el.props.onConfirm as () => Promise<void>)()

        expect(h.deleteChat).toHaveBeenCalledWith("c1")
        expect(onSelectApiChat).toHaveBeenCalledWith(undefined)
    })
})

describe("useWorkbenchHistory — docker mode", () => {
    it("fetches past sessions, maps them to docker rows, and exposes no delete", async () => {
        listChatSessions.mockResolvedValue([
            {id: "s1", firstUserMessage: "Debug websocket", lastActivityAt: "2026-08-20T09:00:00Z"},
        ])

        const {result} = renderHook(() => useWorkbenchHistory({
            mode: "docker", vaultId: "v1", onSelectApiChat: noop, onNewChat: noop,
        }))

        await waitFor(() => expect(result.current.rows).toHaveLength(1))
        expect(listChatSessions).toHaveBeenCalledWith("v1")
        expect(result.current.rows[0]).toMatchObject({id: "s1", title: "Debug websocket", source: "docker"})
        expect(result.current.onDelete).toBeUndefined()
    })

    it("loads the transcript and opens the transcript dialog on select", async () => {
        listChatSessions.mockResolvedValue([
            {id: "s1", firstUserMessage: "Question", lastActivityAt: "2026-08-20T09:00:00Z"},
        ])
        getChatSession.mockResolvedValue([{type: "user_message", text: "Question", seq: 1}])

        const {result} = renderHook(() => useWorkbenchHistory({
            mode: "docker", vaultId: "v1", onSelectApiChat: noop, onNewChat: noop,
        }))

        await waitFor(() => expect(result.current.rows).toHaveLength(1))

        result.current.onSelect("s1")

        await waitFor(() => expect(h.openDialog).toHaveBeenCalled())
        expect(getChatSession).toHaveBeenCalledWith("v1", "s1")
        const el = h.openDialog.mock.calls[0][0] as {props: Record<string, unknown>}
        expect(el.props.title).toBe("Question")
        expect(Array.isArray(el.props.items)).toBe(true)
    })
})
