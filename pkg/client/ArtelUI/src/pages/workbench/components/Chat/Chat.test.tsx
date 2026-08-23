import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen, waitFor} from "@testing-library/react"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import * as chatHistoryApi from "@/pages/workbench/processes/chatHistoryApi.ts"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

function renderChat(
    items: ChatItem[] = [],
    status: ChatConnectionStatus = "open",
    vaultId?: string,
    historyOpen = false,
) {
    const sendMessage = vi.fn()
    const sendPermissionDecision = vi.fn()
    const onNewChat = vi.fn()
    const onCloseHistory = vi.fn()

    const result = render(
        <Chat
            items={items}
            status={status}
            sendMessage={sendMessage}
            sendPermissionDecision={sendPermissionDecision}
            onNewChat={onNewChat}
            vaultId={vaultId}
            historyOpen={historyOpen}
            onCloseHistory={onCloseHistory}
        />,
    )

    return {sendMessage, sendPermissionDecision, onNewChat, onCloseHistory, container: result.container}
}

describe("Chat", () => {
    it("sends a user message via the composer and clears the input", () => {
        const {sendMessage} = renderChat()

        const input = screen.getByPlaceholderText("Message the workbench…")
        fireEvent.change(input, {target: {value: "hello there"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendMessage).toHaveBeenCalledWith("hello there")
    })

    it("shows the reconnecting banner when the socket drops", () => {
        renderChat([], "reconnecting")

        expect(screen.getByText("Reconnecting…")).toBeInTheDocument()
    })

    it("calls onNewChat when the New Chat button is clicked", () => {
        const {onNewChat} = renderChat()

        fireEvent.click(screen.getByLabelText("New chat"))

        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("does not render the sidebar overlay when historyOpen is false", () => {
        const {container} = renderChat([], "open", "v1", false)

        expect(container.querySelector('[class*="ChatContentOverlay"]')).toBeFalsy()
    })

    it("renders the history sidebar and fetches sessions when historyOpen is true", async () => {
        const listSessionsSpy = vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])

        const {container} = renderChat([], "open", "v1", true)

        await waitFor(() => {
            expect(listSessionsSpy).toHaveBeenCalledWith("v1")
            expect(screen.getByText("Chat History")).toBeInTheDocument()
            expect(container.querySelector('[class*="ChatContentOverlay"]')).toBeTruthy()
        })

        listSessionsSpy.mockRestore()
    })

    it("calls onCloseHistory when the overlay is clicked", async () => {
        vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])

        const {container, onCloseHistory} = renderChat([], "open", "v1", true)

        await waitFor(() => {
            expect(container.querySelector('[class*="ChatContentOverlay"]')).toBeTruthy()
        })

        fireEvent.click(container.querySelector('[class*="ChatContentOverlay"]') as Element)

        expect(onCloseHistory).toHaveBeenCalledTimes(1)

        vi.restoreAllMocks()
    })

    it("calls onCloseHistory when the sidebar's own close button is clicked", async () => {
        vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])

        const {onCloseHistory} = renderChat([], "open", "v1", true)

        await waitFor(() => {
            expect(screen.getByText("Chat History")).toBeInTheDocument()
        })

        fireEvent.click(screen.getByLabelText("Close chat history"))

        expect(onCloseHistory).toHaveBeenCalledTimes(1)

        vi.restoreAllMocks()
    })
})
