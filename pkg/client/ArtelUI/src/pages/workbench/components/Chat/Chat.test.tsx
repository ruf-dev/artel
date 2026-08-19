import {beforeEach, describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

const sendMessage = vi.fn()
const sendPermissionDecision = vi.fn()
const sendAuthCode = vi.fn()

let mockItems: ChatItem[] = []
let mockStatus: "connecting" | "open" | "reconnecting" | "closed" = "open"

vi.mock("@/pages/workbench/processes/useChatSession.ts", () => ({
    useChatSession: () => ({
        items: mockItems,
        status: mockStatus,
        sendMessage,
        sendPermissionDecision,
        sendAuthCode,
    }),
}))

describe("Chat", () => {
    beforeEach(() => {
        mockItems = []
        mockStatus = "open"
        sendMessage.mockReset()
        sendPermissionDecision.mockReset()
        sendAuthCode.mockReset()
    })

    it("sends a user message via the composer and clears the input", () => {
        render(<Chat vaultId="v1"/>)

        const input = screen.getByPlaceholderText("Message the workbench…")
        fireEvent.change(input, {target: {value: "hello there"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendMessage).toHaveBeenCalledWith("hello there")
    })

    it("shows the reconnecting banner when the socket drops", () => {
        mockStatus = "reconnecting"
        render(<Chat vaultId="v1"/>)

        expect(screen.getByText("Reconnecting…")).toBeInTheDocument()
    })

    it("switches the composer to auth-code mode while a code prompt is pending", () => {
        mockItems = [{kind: "auth_code_needed", key: "auth_code_needed-0", resolved: false}]
        render(<Chat vaultId="v1"/>)

        const input = screen.getByPlaceholderText("Enter the authentication code…")
        fireEvent.change(input, {target: {value: "123456"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendAuthCode).toHaveBeenCalledWith("123456")
        expect(sendMessage).not.toHaveBeenCalled()
    })
})
