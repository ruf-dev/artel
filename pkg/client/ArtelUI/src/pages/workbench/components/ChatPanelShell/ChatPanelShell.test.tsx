import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ChatPanelShell from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

function renderShell(
    items: ChatItem[] = [],
    bannerStatus: ChatConnectionStatus = "open",
) {
    const sendMessage = vi.fn()
    const sendPermissionDecision = vi.fn()
    const onNewChat = vi.fn()

    const result = render(
        <ChatPanelShell
            items={items}
            bannerStatus={bannerStatus}
            sendMessage={sendMessage}
            sendPermissionDecision={sendPermissionDecision}
            onNewChat={onNewChat}
            composerDisabled={bannerStatus !== "open"}
            composerPlaceholder="Message the workbench…"
        />,
    )

    return {sendMessage, sendPermissionDecision, onNewChat, container: result.container}
}

describe("ChatPanelShell", () => {
    it("sends a user message via the composer and clears the input", () => {
        const {sendMessage} = renderShell()

        const input = screen.getByPlaceholderText("Message the workbench…")
        fireEvent.change(input, {target: {value: "hello there"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendMessage).toHaveBeenCalledWith("hello there")
    })

    it("shows the reconnecting banner when the socket drops", () => {
        renderShell([], "reconnecting")

        expect(screen.getByText("Reconnecting…")).toBeInTheDocument()
    })

    it("calls onNewChat when the New Chat button is clicked", () => {
        const {onNewChat} = renderShell()

        fireEvent.click(screen.getByLabelText("New chat"))

        expect(onNewChat).toHaveBeenCalledTimes(1)
    })
})
