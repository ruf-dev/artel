import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

function renderChat(items: ChatItem[] = [], status: ChatConnectionStatus = "open") {
    const sendMessage = vi.fn()
    const sendPermissionDecision = vi.fn()
    const sendAuthCode = vi.fn()

    render(
        <Chat
            items={items}
            status={status}
            sendMessage={sendMessage}
            sendPermissionDecision={sendPermissionDecision}
            sendAuthCode={sendAuthCode}
        />,
    )

    return {sendMessage, sendPermissionDecision, sendAuthCode}
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

    it("switches the composer to auth-code mode while a code prompt is pending", () => {
        const {sendAuthCode, sendMessage} = renderChat(
            [{kind: "auth_code_needed", key: "auth_code_needed-0", resolved: false}],
        )

        const input = screen.getByPlaceholderText("Enter the authentication code…")
        fireEvent.change(input, {target: {value: "123456"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendAuthCode).toHaveBeenCalledWith("123456")
        expect(sendMessage).not.toHaveBeenCalled()
    })
})
