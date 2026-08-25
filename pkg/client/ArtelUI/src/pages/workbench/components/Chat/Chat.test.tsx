import {describe, expect, it, vi} from "vitest"
import {render, screen} from "@testing-library/react"
import {ReactNode} from "react"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

interface MockChatPanelShellProps {
    bannerStatus: ChatConnectionStatus
    composerPlaceholder: string
    composerDisabled: boolean
    historySidebar?: ReactNode
}

// Chat.tsx is now a thin prop-passthrough wrapper around ChatPanelShell — the
// behavioral scenarios (composer send/clear, overlay click, history fetch, etc.)
// live in ChatPanelShell.test.tsx, so this only asserts the right props reach it.
vi.mock("@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx", () => ({
    default: (props: MockChatPanelShellProps) => (
        <div className="chat-panel-shell" data-testid="chat-panel-shell">
            <span data-testid="banner-status">{props.bannerStatus}</span>
            <span data-testid="composer-placeholder">{props.composerPlaceholder}</span>
            <span data-testid="composer-disabled">{String(props.composerDisabled)}</span>
            <span data-testid="history-sidebar-present">{String(Boolean(props.historySidebar))}</span>
        </div>
    ),
}))

function renderChat(items: ChatItem[] = [], status: ChatConnectionStatus = "open", vaultId?: string) {
    render(
        <Chat
            items={items}
            status={status}
            sendMessage={vi.fn()}
            sendPermissionDecision={vi.fn()}
            onNewChat={vi.fn()}
            vaultId={vaultId}
            historyOpen={false}
            onCloseHistory={vi.fn()}
        />,
    )
}

describe("Chat", () => {
    it("passes status through to ChatPanelShell as bannerStatus", () => {
        renderChat([], "reconnecting")

        expect(screen.getByTestId("banner-status")).toHaveTextContent("reconnecting")
    })

    it("passes the workbench composer placeholder", () => {
        renderChat()

        expect(screen.getByTestId("composer-placeholder")).toHaveTextContent("Message the workbench…")
    })

    it("disables the composer when status is not open", () => {
        renderChat([], "connecting")

        expect(screen.getByTestId("composer-disabled")).toHaveTextContent("true")
    })

    it("enables the composer when status is open", () => {
        renderChat([], "open")

        expect(screen.getByTestId("composer-disabled")).toHaveTextContent("false")
    })

    it("passes a history sidebar when vaultId is set", () => {
        renderChat([], "open", "v1")

        expect(screen.getByTestId("history-sidebar-present")).toHaveTextContent("true")
    })

    it("passes no history sidebar when vaultId is undefined", () => {
        renderChat([], "open", undefined)

        expect(screen.getByTestId("history-sidebar-present")).toHaveTextContent("false")
    })
})
