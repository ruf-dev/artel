import {renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

import {
    toSimpleChatSessionBundle,
    toSimpleChatTopBarProps,
    useSimpleChatController,
} from "@/pages/workbench/processes/useSimpleChatController.ts"

const mockUseSimpleChat = vi.fn()
const mockUseOpenRouterModels = vi.fn()
const mockCreate = vi.fn()
const mockUseSimpleChatSession = vi.fn()

vi.mock("@/app/hooks/SimpleChat.ts", () => ({
    useSimpleChat: (chatId?: string) => mockUseSimpleChat(chatId),
    useSimpleChatMutations: (vaultId?: string) => ({create: mockCreate, vaultId}),
}))

vi.mock("@/pages/workbench/processes/useOpenRouterModels.ts", () => ({
    useOpenRouterModels: (enabled: boolean) => mockUseOpenRouterModels(enabled),
}))

vi.mock("@/pages/workbench/processes/useSimpleChatSession.ts", () => ({
    useSimpleChatSession: (chatId: string | undefined, initialModel: string, initialItems: unknown[]) =>
        mockUseSimpleChatSession(chatId, initialModel, initialItems),
}))

describe("useSimpleChatController", () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mockUseSimpleChat.mockReturnValue({chat: undefined, messages: []})
        mockUseOpenRouterModels.mockReturnValue({models: ["m1"], recommendedDefaultModel: "m1", isLoading: false})
        mockUseSimpleChatSession.mockReturnValue({
            items: [],
            status: "closed",
            sendMessage: vi.fn(),
            sendPermissionDecision: vi.fn(),
            currentModel: "m1",
            setModel: vi.fn(),
        })
        mockCreate.mockResolvedValue({id: "new-chat-id"})
    })

    it("does not pass chatId through to the underlying hooks when inactive", () => {
        renderHook(() => useSimpleChatController({
            chatId: "c1",
            vaultId: "v1",
            active: false,
            onChatCreated: vi.fn(),
            bakeError: vi.fn(),
        }))

        expect(mockUseSimpleChat).toHaveBeenCalledWith(undefined)
        expect(mockUseOpenRouterModels).toHaveBeenCalledWith(false)
        expect(mockUseSimpleChatSession).toHaveBeenCalledWith(undefined, expect.anything(), expect.anything())
    })

    it("passes chatId through to the underlying hooks when active", () => {
        renderHook(() => useSimpleChatController({
            chatId: "c1",
            vaultId: "v1",
            active: true,
            onChatCreated: vi.fn(),
            bakeError: vi.fn(),
        }))

        expect(mockUseSimpleChat).toHaveBeenCalledWith("c1")
        expect(mockUseOpenRouterModels).toHaveBeenCalledWith(true)
        expect(mockUseSimpleChatSession).toHaveBeenCalledWith("c1", expect.anything(), expect.anything())
    })

    it("creates a chat with the current model and reports the new id on success", async () => {
        const onChatCreated = vi.fn()
        const {result} = renderHook(() => useSimpleChatController({
            chatId: undefined,
            vaultId: "v1",
            active: true,
            onChatCreated,
            bakeError: vi.fn(),
        }))

        result.current.handleNewChat()

        await waitFor(() => expect(onChatCreated).toHaveBeenCalledWith("new-chat-id"))
        expect(mockCreate).toHaveBeenCalledWith({model: "m1", vaultAccess: true})
    })

    it("bakes an error toast when chat creation fails", async () => {
        mockCreate.mockRejectedValue(new Error("boom"))
        const bakeError = vi.fn()
        const {result} = renderHook(() => useSimpleChatController({
            chatId: undefined,
            vaultId: "v1",
            active: true,
            onChatCreated: vi.fn(),
            bakeError,
        }))

        result.current.handleNewChat()

        await waitFor(() => expect(bakeError).toHaveBeenCalledWith("Failed to create chat", expect.any(Error)))
    })
})

describe("toSimpleChatSessionBundle", () => {
    it("maps the controller's fields into the session bundle shape, renaming handleNewChat to onNewChat", () => {
        const sendMessage = vi.fn()
        const sendPermissionDecision = vi.fn()
        const handleNewChat = vi.fn()

        const bundle = toSimpleChatSessionBundle({
            items: [],
            status: "open",
            sendMessage,
            sendPermissionDecision,
            currentModel: "m1",
            setModel: vi.fn(),
            models: ["m1"],
            modelsLoading: false,
            handleNewChat,
        })

        expect(bundle).toEqual({
            items: [],
            status: "open",
            sendMessage,
            sendPermissionDecision,
            onNewChat: handleNewChat,
        })
    })
})

describe("toSimpleChatTopBarProps", () => {
    it("maps the controller's fields into the top-bar props shape, renaming setModel/handleNewChat", () => {
        const setModel = vi.fn()
        const handleNewChat = vi.fn()
        const onToggleHistory = vi.fn()

        const props = toSimpleChatTopBarProps({
            items: [],
            status: "open",
            sendMessage: vi.fn(),
            sendPermissionDecision: vi.fn(),
            currentModel: "m1",
            setModel,
            models: ["m1", "m2"],
            modelsLoading: true,
            handleNewChat,
        }, onToggleHistory)

        expect(props).toEqual({
            models: ["m1", "m2"],
            currentModel: "m1",
            modelsLoading: true,
            onChangeModel: setModel,
            onNewChat: handleNewChat,
            onToggleHistory,
        })
    })
})
