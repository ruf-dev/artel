import {describe, expect, it, vi} from "vitest"
import {render, screen} from "@testing-library/react"

import SimpleChat from "@/pages/workbench/components/SimpleChat/SimpleChat.tsx"

let lastChatPanelShellProps: Record<string, unknown> | undefined

vi.mock("@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx", () => ({
    default: (props: Record<string, unknown>) => {
        lastChatPanelShellProps = props
        return <div className="chat-panel-shell" data-testid="chat-panel-shell"/>
    },
}))

vi.mock("@/pages/workbench/components/SimpleChatHistorySidebar/SimpleChatHistorySidebar.tsx", () => ({
    default: () => <div className="simple-chat-history-sidebar" data-testid="simple-chat-history-sidebar"/>,
}))

describe("SimpleChat", () => {
    const baseSession = {
        items: [],
        status: "open" as const,
        sendMessage: vi.fn(),
        sendPermissionDecision: vi.fn(),
        onNewChat: vi.fn(),
        models: ["m1"],
        currentModel: "m1",
        modelsLoading: false,
        onChangeModel: vi.fn(),
    }

    const defaultProps = {
        chatId: "chat-1",
        vaultId: "vault-1",
        onSelectChat: vi.fn(),
        onClose: vi.fn(),
        session: baseSession,
    }

    it("always hides ChatPanelShell's own new-chat button, since the top bar owns it", () => {
        render(<SimpleChat {...defaultProps}/>)

        expect(lastChatPanelShellProps?.hideNewChatButton).toBe(true)
    })

    it("uses the session status as the banner status when a chat is selected", () => {
        render(<SimpleChat {...defaultProps} session={{...baseSession, status: "reconnecting"}}/>)

        expect(lastChatPanelShellProps?.bannerStatus).toBe("reconnecting")
    })

    it("forces the banner status to closed when no chat is selected", () => {
        render(<SimpleChat {...defaultProps} chatId={undefined} session={{...baseSession, status: "open"}}/>)

        expect(lastChatPanelShellProps?.bannerStatus).toBe("closed")
    })

    it("disables the composer when there is no active chat", () => {
        render(<SimpleChat {...defaultProps} chatId={undefined}/>)

        expect(lastChatPanelShellProps?.composerDisabled).toBe(true)
    })

    it("disables the composer when the session isn't open", () => {
        render(<SimpleChat {...defaultProps} session={{...baseSession, status: "connecting"}}/>)

        expect(lastChatPanelShellProps?.composerDisabled).toBe(true)
    })

    it("enables the composer when a chat is selected and the session is open", () => {
        render(<SimpleChat {...defaultProps}/>)

        expect(lastChatPanelShellProps?.composerDisabled).toBe(false)
    })

    it("renders SimpleChatHistorySidebar alongside the chat panel", () => {
        render(<SimpleChat {...defaultProps}/>)

        expect(screen.getByTestId("simple-chat-history-sidebar")).toBeInTheDocument()
        expect(screen.getByTestId("chat-panel-shell")).toBeInTheDocument()
    })
})
