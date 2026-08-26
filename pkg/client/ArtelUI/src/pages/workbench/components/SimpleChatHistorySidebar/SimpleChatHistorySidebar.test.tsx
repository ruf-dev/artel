import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"
import {Button} from "@vervstack/chures"

import SimpleChatHistorySidebar from
    "@/pages/workbench/components/SimpleChatHistorySidebar/SimpleChatHistorySidebar.tsx"

let collapsedRailProps: Record<string, unknown> | undefined

vi.mock(
    // eslint-disable-next-line max-len -- mock path too long to wrap under 120 chars
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistoryListScreen/SimpleChatHistoryListScreen.tsx",
    () => ({
        default: (props: Record<string, unknown>) => (
            <div className="list-screen" data-testid="list-screen">
                <Button
                    onClick={() => {
                        (props as {onCollapse: () => void}).onCollapse()
                    }}
                    data-testid="collapse-button"
                >
                    Collapse
                </Button>
            </div>
        ),
    }),
)

vi.mock(
    // eslint-disable-next-line max-len -- mock path too long to wrap under 120 chars
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistorySidebarCollapsedRail/SimpleChatHistorySidebarCollapsedRail.tsx",
    () => ({
        default: (props: Record<string, unknown>) => {
            collapsedRailProps = props
            return (
                <div className="collapsed-rail" data-testid="collapsed-rail">
                    <Button
                        onClick={() => {
                            (props as {onExpand: () => void}).onExpand()
                        }}
                        data-testid="expand-button"
                    >
                        Expand
                    </Button>
                    <Button
                        onClick={() => {
                            (props as {onNewChat: () => void}).onNewChat()
                        }}
                        data-testid="rail-new-chat-button"
                    >
                        New Chat
                    </Button>
                </div>
            )
        },
    }),
)

vi.mock("@/app/hooks/SimpleChat.ts", () => ({
    useSimpleChats: () => ({
        chats: [],
        isLoading: false,
    }),
    useSimpleChatMutations: () => ({
        delete: vi.fn(),
    }),
}))

vi.mock("@/app/hooks/Dialog", () => ({
    useDialog: () => ({
        OpenDialog: vi.fn(),
        CloseDialog: vi.fn(),
    }),
}))

vi.mock("@/app/hooks/useErrorToast.ts", () => ({
    useBakeError: () => vi.fn(),
}))

vi.mock("@/hooks/user/User.ts")

describe("SimpleChatHistorySidebar", () => {
    const defaultProps = {
        vaultId: "vault-123",
        activeChatId: undefined,
        onSelectChat: vi.fn(),
        onNewChat: vi.fn(),
        models: ["m1"],
        currentModel: "m1",
        modelsLoading: false,
        onChangeModel: vi.fn(),
    }

    it("renders the list screen initially", () => {
        render(<SimpleChatHistorySidebar {...defaultProps}/>)

        expect(screen.getByTestId("list-screen")).toBeInTheDocument()
        expect(screen.queryByTestId("collapsed-rail")).not.toBeInTheDocument()
    })

    it("switches to the collapsed rail when the collapse button is clicked", () => {
        render(<SimpleChatHistorySidebar {...defaultProps}/>)

        const collapseButton = screen.getByTestId("collapse-button")
        fireEvent.click(collapseButton)

        expect(screen.getByTestId("collapsed-rail")).toBeInTheDocument()
        expect(screen.queryByTestId("list-screen")).not.toBeInTheDocument()
    })

    it("switches back to the list screen when the expand button is clicked", () => {
        render(<SimpleChatHistorySidebar {...defaultProps}/>)

        const collapseButton = screen.getByTestId("collapse-button")
        fireEvent.click(collapseButton)

        expect(screen.getByTestId("collapsed-rail")).toBeInTheDocument()

        const expandButton = screen.getByTestId("expand-button")
        fireEvent.click(expandButton)

        expect(screen.getByTestId("list-screen")).toBeInTheDocument()
        expect(screen.queryByTestId("collapsed-rail")).not.toBeInTheDocument()
    })

    it("calls onNewChat when the new chat button in the collapsed rail is clicked", () => {
        const mockOnNewChat = vi.fn()
        render(<SimpleChatHistorySidebar {...defaultProps} onNewChat={mockOnNewChat}/>)

        const collapseButton = screen.getByTestId("collapse-button")
        fireEvent.click(collapseButton)

        const railNewChatButton = screen.getByTestId("rail-new-chat-button")
        fireEvent.click(railNewChatButton)

        expect(mockOnNewChat).toHaveBeenCalled()
    })

    it("passes onNewChat to the collapsed rail component", () => {
        const mockOnNewChat = vi.fn()
        render(<SimpleChatHistorySidebar {...defaultProps} onNewChat={mockOnNewChat}/>)

        const collapseButton = screen.getByTestId("collapse-button")
        fireEvent.click(collapseButton)

        expect(collapsedRailProps?.onNewChat).toBe(mockOnNewChat)
    })
})
