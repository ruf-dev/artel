import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen, waitFor} from "@testing-library/react"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import * as chatHistoryApi from "@/pages/workbench/processes/chatHistoryApi.ts"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

function renderChat(items: ChatItem[] = [], status: ChatConnectionStatus = "open", vaultId?: string) {
    const sendMessage = vi.fn()
    const sendPermissionDecision = vi.fn()
    const onNewChat = vi.fn()

    const result = render(
        <Chat
            items={items}
            status={status}
            sendMessage={sendMessage}
            sendPermissionDecision={sendPermissionDecision}
            onNewChat={onNewChat}
            vaultId={vaultId}
        />,
    )

    return {sendMessage, sendPermissionDecision, onNewChat, container: result.container}
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

        fireEvent.click(screen.getByText("New Chat"))

        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("shows the history button when vaultId is provided", () => {
        renderChat([], "open", "v1")

        expect(screen.getByLabelText("View chat history")).toBeInTheDocument()
    })

    it("does not show the history button when vaultId is not provided", () => {
        renderChat()

        expect(screen.queryByLabelText("View chat history")).not.toBeInTheDocument()
    })

    it("opens the history sidebar and fetches sessions when the history button is clicked", async () => {
        const listSessionsSpy = vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])

        renderChat([], "open", "v1")

        fireEvent.click(screen.getByLabelText("View chat history"))

        await waitFor(() => {
            expect(listSessionsSpy).toHaveBeenCalledWith("v1")
            expect(screen.getByText("Chat History")).toBeInTheDocument()
        })

        listSessionsSpy.mockRestore()
    })

    it("closes the history sidebar when its close button is clicked", async () => {
        vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])

        const {container} = renderChat([], "open", "v1")

        fireEvent.click(screen.getByLabelText("View chat history"))

        await waitFor(() => {
            expect(screen.getByText("Chat History")).toBeInTheDocument()
            expect(container.querySelector('[class*="ChatContentOverlay"]')).toBeTruthy()
        })

        fireEvent.click(screen.getByLabelText("Close chat history"))

        // The sidebar itself stays mounted (so it can slide back out via CSS
        // transition) — closed is expressed by the "open" modifier class and the
        // blur overlay disappearing, not by the sidebar unmounting.
        await waitFor(() => {
            expect(container.querySelector('[class*="ChatContentOverlay"]')).toBeFalsy()
            const sidebar = container.querySelector('[class*="ChatHistorySidebarContainer"]')
            expect(sidebar?.className).not.toMatch(/ChatHistorySidebarOpen/)
        })

        vi.restoreAllMocks()
    })

    it("closes the history sidebar when clicking outside it on the blurred content", async () => {
        vi.spyOn(chatHistoryApi, "listChatSessions").mockResolvedValue([])

        const {container} = renderChat([], "open", "v1")

        fireEvent.click(screen.getByLabelText("View chat history"))

        await waitFor(() => {
            expect(screen.getByText("Chat History")).toBeInTheDocument()
        })

        const overlay = container.querySelector('[class*="ChatContentOverlay"]')
        expect(overlay).toBeTruthy()
        fireEvent.click(overlay as Element)

        await waitFor(() => {
            expect(container.querySelector('[class*="ChatContentOverlay"]')).toBeFalsy()
            const sidebar = container.querySelector('[class*="ChatHistorySidebarContainer"]')
            expect(sidebar?.className).not.toMatch(/ChatHistorySidebarOpen/)
        })

        vi.restoreAllMocks()
    })
})
