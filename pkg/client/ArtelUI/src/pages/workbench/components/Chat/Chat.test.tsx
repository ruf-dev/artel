import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

function renderChat(items: ChatItem[] = [], status: ChatConnectionStatus = "open", vaultId?: string) {
    const sendMessage = vi.fn()
    const sendPermissionDecision = vi.fn()
    const onNewChat = vi.fn()

    render(
        <Chat
            items={items}
            status={status}
            sendMessage={sendMessage}
            sendPermissionDecision={sendPermissionDecision}
            onNewChat={onNewChat}
            vaultId={vaultId}
        />,
    )

    return {sendMessage, sendPermissionDecision, onNewChat}
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
})
