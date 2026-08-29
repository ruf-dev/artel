import {renderHook, waitFor} from "@testing-library/react"
import {beforeEach, describe, expect, it, vi} from "vitest"

import {
    toSimpleChatSessionBundle,
    useSimpleChatController,
} from "@/pages/workbench/processes/useSimpleChatController.ts"

const mockUseSimpleChat = vi.fn()
const mockUseOpenRouterModels = vi.fn()
const mockCreate = vi.fn()
const mockCreating = vi.fn(() => false)
const mockUseSimpleChatSession = vi.fn()
const mockUseLastUsedModel = vi.fn()
const mockSetLastUsedModel = vi.fn()

vi.mock("@/app/hooks/SimpleChat.ts", () => ({
    useSimpleChat: (chatId?: string) => mockUseSimpleChat(chatId),
    useSimpleChatMutations: (vaultId?: string) => ({create: mockCreate, creating: mockCreating(), vaultId}),
}))

vi.mock("@/app/hooks/useLastUsedModel.ts", () => ({
    useLastUsedModel: () => mockUseLastUsedModel(),
}))

vi.mock("@/pages/workbench/processes/useOpenRouterModels.ts", () => ({
    useOpenRouterModels: (enabled: boolean) => mockUseOpenRouterModels(enabled),
}))

vi.mock("@/pages/workbench/processes/useSimpleChatSession.ts", () => ({
    useSimpleChatSession: (chatId: string | undefined, initialModel: string, initialItems: unknown[]) =>
        mockUseSimpleChatSession(chatId, initialModel, initialItems),
}))

beforeEach(() => {
    vi.clearAllMocks()
    mockUseSimpleChat.mockReturnValue({chat: undefined, messages: []})
    mockUseOpenRouterModels.mockReturnValue({models: ["m1"], recommendedDefaultModel: "m1", isLoading: false})
    mockUseSimpleChatSession.mockReturnValue({
        items: [],
        status: "closed",
        sendMessage: vi.fn(),
        resendMessage: vi.fn(),
        sendPermissionDecision: vi.fn(),
        currentModel: "m1",
        setModel: vi.fn(),
        pendingTurn: false,
    })
    mockUseLastUsedModel.mockReturnValue({lastUsedModel: undefined, setLastUsedModel: mockSetLastUsedModel})
    mockCreate.mockResolvedValue({id: "new-chat-id"})
    mockCreating.mockReturnValue(false)
})

describe("useSimpleChatController", () => {
    it("does not pass chatId through to the underlying hooks when inactive", () => {
        renderHook(() => useSimpleChatController({
            chatId: "c1",
            vaultId: "v1",
            active: false,
            onChatIdChange: vi.fn(),
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
            onChatIdChange: vi.fn(),
            bakeError: vi.fn(),
        }))

        expect(mockUseSimpleChat).toHaveBeenCalledWith("c1")
        expect(mockUseOpenRouterModels).toHaveBeenCalledWith(true)
        expect(mockUseSimpleChatSession).toHaveBeenCalledWith("c1", expect.anything(), expect.anything())
    })

    it("creates a chat with the current model and reports the new id on success", async () => {
        const onChatIdChange = vi.fn()
        const {result} = renderHook(() => useSimpleChatController({
            chatId: undefined,
            vaultId: "v1",
            active: true,
            onChatIdChange,
            bakeError: vi.fn(),
        }))

        result.current.handleNewChat()

        await waitFor(() => expect(onChatIdChange).toHaveBeenCalledWith("new-chat-id"))
        expect(mockCreate).toHaveBeenCalledWith({model: "m1", vaultAccess: true})
    })

    it("clears chatId before the create call resolves, so no send can race onto the old chat", () => {
        const onChatIdChange = vi.fn()
        let resolveCreate!: (chat: {id: string}) => void
        mockCreate.mockReturnValue(new Promise(resolve => {
            resolveCreate = resolve
        }))

        const {result} = renderHook(() => useSimpleChatController({
            chatId: "old-chat-id",
            vaultId: "v1",
            active: true,
            onChatIdChange,
            bakeError: vi.fn(),
        }))

        result.current.handleNewChat()

        expect(onChatIdChange).toHaveBeenCalledWith(undefined)
        expect(onChatIdChange).not.toHaveBeenCalledWith("new-chat-id")

        resolveCreate({id: "new-chat-id"})
    })

    it("exposes creatingChat from the create mutation's pending state", () => {
        mockCreating.mockReturnValue(true)

        const {result} = renderHook(() => useSimpleChatController({
            chatId: undefined,
            vaultId: "v1",
            active: true,
            onChatIdChange: vi.fn(),
            bakeError: vi.fn(),
        }))

        expect(result.current.creatingChat).toBe(true)
    })

    it("bakes an error toast when chat creation fails", async () => {
        mockCreate.mockRejectedValue(new Error("boom"))
        const bakeError = vi.fn()
        const {result} = renderHook(() => useSimpleChatController({
            chatId: undefined,
            vaultId: "v1",
            active: true,
            onChatIdChange: vi.fn(),
            bakeError,
        }))

        result.current.handleNewChat()

        await waitFor(() => expect(bakeError).toHaveBeenCalledWith("Failed to create chat", expect.any(Error)))
    })
})

describe("useSimpleChatController - last used model caching", () => {
    it("prefers the cached last-used model over recommendedDefaultModel for a chat with no model of its own", () => {
        mockUseLastUsedModel.mockReturnValue({lastUsedModel: "cached-model", setLastUsedModel: mockSetLastUsedModel})

        renderHook(() => useSimpleChatController({
            chatId: undefined,
            vaultId: "v1",
            active: true,
            onChatIdChange: vi.fn(),
            bakeError: vi.fn(),
        }))

        expect(mockUseSimpleChatSession).toHaveBeenCalledWith(undefined, "cached-model", expect.anything())
    })

    it("a chat's own persisted model still wins over the cached last-used model", () => {
        mockUseSimpleChat.mockReturnValue({chat: {model: "chat-model"}, messages: []})
        mockUseLastUsedModel.mockReturnValue({lastUsedModel: "cached-model", setLastUsedModel: mockSetLastUsedModel})

        renderHook(() => useSimpleChatController({
            chatId: "c1",
            vaultId: "v1",
            active: true,
            onChatIdChange: vi.fn(),
            bakeError: vi.fn(),
        }))

        expect(mockUseSimpleChatSession).toHaveBeenCalledWith("c1", "chat-model", expect.anything())
    })

    it("setModel caches the newly selected model as last-used", () => {
        const sessionSetModel = vi.fn()
        mockUseSimpleChatSession.mockReturnValue({
            items: [],
            status: "closed",
            sendMessage: vi.fn(),
            resendMessage: vi.fn(),
            sendPermissionDecision: vi.fn(),
            currentModel: "m1",
            setModel: sessionSetModel,
            pendingTurn: false,
        })

        const {result} = renderHook(() => useSimpleChatController({
            chatId: "c1",
            vaultId: "v1",
            active: true,
            onChatIdChange: vi.fn(),
            bakeError: vi.fn(),
        }))

        result.current.setModel("new-model")

        expect(sessionSetModel).toHaveBeenCalledWith("new-model")
        expect(mockSetLastUsedModel).toHaveBeenCalledWith("new-model")
    })

    it("handleNewChat caches the model it creates the chat with", async () => {
        const {result} = renderHook(() => useSimpleChatController({
            chatId: undefined,
            vaultId: "v1",
            active: true,
            onChatIdChange: vi.fn(),
            bakeError: vi.fn(),
        }))

        result.current.handleNewChat()

        await waitFor(() => expect(mockCreate).toHaveBeenCalled())
        expect(mockSetLastUsedModel).toHaveBeenCalledWith("m1")
    })
})

describe("toSimpleChatSessionBundle", () => {
    it("maps the controller's fields into the session bundle shape, renaming handleNewChat to onNewChat", () => {
        const sendMessage = vi.fn()
        const resendMessage = vi.fn()
        const sendPermissionDecision = vi.fn()
        const handleNewChat = vi.fn()
        const setModel = vi.fn()

        const bundle = toSimpleChatSessionBundle({
            items: [],
            status: "open",
            sendMessage,
            resendMessage,
            sendPermissionDecision,
            currentModel: "m1",
            setModel,
            models: ["m1"],
            modelsLoading: false,
            handleNewChat,
            creatingChat: false,
            pendingTurn: false,
        })

        expect(bundle).toEqual({
            items: [],
            status: "open",
            sendMessage,
            resendMessage,
            sendPermissionDecision,
            onNewChat: handleNewChat,
            models: ["m1"],
            currentModel: "m1",
            modelsLoading: false,
            onChangeModel: setModel,
            pendingTurn: false,
        })
    })
})
