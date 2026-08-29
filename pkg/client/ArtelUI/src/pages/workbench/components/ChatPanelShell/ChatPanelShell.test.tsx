import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ChatPanelShell from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

function renderShell(
    items: ChatItem[] = [],
    bannerStatus: ChatConnectionStatus = "open",
    attachedPaths: string[] = [],
) {
    const sendMessage = vi.fn()
    const sendPermissionDecision = vi.fn()
    const onNewChat = vi.fn()
    const onRemoveAttachment = vi.fn()
    const onClearAttachments = vi.fn()

    const ctx = {tweaksOpen: false, tweaksSection: undefined, openTweaks: vi.fn(), closeTweaks: vi.fn()}

    const result = render(
        <ChatPanelShell
            items={items}
            vaultId="v1"
            bannerStatus={bannerStatus}
            sendMessage={sendMessage}
            sendPermissionDecision={sendPermissionDecision}
            onNewChat={onNewChat}
            pendingTurn={false}
            composerDisabled={bannerStatus !== "open"}
            composerPlaceholder="Message the workbench…"
            ctx={ctx}
            attachedPaths={attachedPaths}
            onRemoveAttachment={onRemoveAttachment}
            onClearAttachments={onClearAttachments}
        />,
    )

    return {
        sendMessage, sendPermissionDecision, onNewChat,
        onRemoveAttachment, onClearAttachments, container: result.container,
    }
}

describe("ChatPanelShell", () => {
    it("sends a user message via the composer and clears the input", () => {
        const {sendMessage} = renderShell()

        const input = screen.getByPlaceholderText("Message the workbench…")
        fireEvent.change(input, {target: {value: "hello there"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendMessage).toHaveBeenCalledWith("hello there", [])
    })

    it("sends the raw text when no vault files are attached", () => {
        const {sendMessage, onClearAttachments} = renderShell([], "open", [])

        fireEvent.change(screen.getByPlaceholderText("Message the workbench…"), {target: {value: "hello"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendMessage).toHaveBeenCalledWith("hello", [])
        expect(onClearAttachments).toHaveBeenCalledTimes(1)
    })

    it("prepends an attached-vault-files preamble, passes attachments separately, and clears on send", () => {
        const {sendMessage, onClearAttachments} = renderShell([], "open", ["a/b.md"])

        fireEvent.change(screen.getByPlaceholderText("Message the workbench…"), {target: {value: "hello"}})
        fireEvent.click(screen.getByLabelText("Send message"))

        expect(sendMessage).toHaveBeenCalledWith("[Attached vault files: a/b.md]\n\nhello", ["a/b.md"])
        expect(onClearAttachments).toHaveBeenCalledTimes(1)
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
