import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import WorkbenchToolbar from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx"

vi.mock(
    "@/pages/workbench/components/WorkbenchToolbar/components/WorkbenchSettingsMenu/WorkbenchSettingsMenu.tsx",
    () => ({
        default: () => (
            <div className="workbench-settings-menu" data-testid="workbench-settings-menu">Settings Menu</div>
        ),
    }),
)

// WorkbenchToolbar composes WorkbenchTopBarShell, which calls useNavigate() for its
// back button — mock it out rather than wrapping every render/rerender call below in
// a <MemoryRouter>, same approach WorkbenchTopBarShell.test.tsx uses directly.
vi.mock("react-router-dom", async () => {
    const actual = await vi.importActual("react-router-dom")
    return {...actual, useNavigate: () => vi.fn()}
})

describe("WorkbenchToolbar", () => {
    const mockOnStart = vi.fn()
    const mockOnStop = vi.fn()
    const mockOnViewChange = vi.fn()
    const mockOnToggleHistory = vi.fn()

    const defaultProps = {
        vaultName: "Test Vault",
        status: "running" as const,
        exists: true,
        vaultId: "vault-123",
        onStart: mockOnStart,
        onStop: mockOnStop,
        stopping: false,
        starting: false,
        view: "chat" as const,
        onViewChange: mockOnViewChange,
        chatLocked: false,
        onToggleHistory: mockOnToggleHistory,
    }

    it("renders both Chat and Terminal toggle buttons when status is running and exists is true", () => {
        render(<WorkbenchToolbar {...defaultProps} />)

        const chatButton = screen.getByRole("button", {name: "Chat"})
        const terminalButton = screen.getByRole("button", {name: "Terminal"})

        expect(chatButton).toBeInTheDocument()
        expect(terminalButton).toBeInTheDocument()
    })

    it("does not render toggle buttons when status is not running", () => {
        render(<WorkbenchToolbar {...defaultProps} status="stopped"/>)

        const chatButton = screen.queryByRole("button", {name: "Chat"})
        const terminalButton = screen.queryByRole("button", {name: "Terminal"})

        expect(chatButton).not.toBeInTheDocument()
        expect(terminalButton).not.toBeInTheDocument()
    })

    it("calls onViewChange with terminal when Terminal button is clicked", () => {
        render(<WorkbenchToolbar {...defaultProps} />)

        const terminalButton = screen.getByRole("button", {name: "Terminal"})
        fireEvent.click(terminalButton)

        expect(mockOnViewChange).toHaveBeenCalledWith("terminal")
    })

    it("sets aria-pressed=true on Terminal button when view is terminal", () => {
        render(<WorkbenchToolbar {...defaultProps} view="terminal"/>)

        const terminalButton = screen.getByRole("button", {name: "Terminal"})
        const chatButton = screen.getByRole("button", {name: "Chat"})

        expect(terminalButton).toHaveAttribute("aria-pressed", "true")
        expect(chatButton).toHaveAttribute("aria-pressed", "false")
    })

    it("sets aria-pressed=true on Chat button when view is chat", () => {
        render(<WorkbenchToolbar {...defaultProps} view="chat"/>)

        const chatButton = screen.getByRole("button", {name: "Chat"})
        const terminalButton = screen.getByRole("button", {name: "Terminal"})

        expect(chatButton).toHaveAttribute("aria-pressed", "true")
        expect(terminalButton).toHaveAttribute("aria-pressed", "false")
    })

    it("disables the Chat button when chatLocked is true", () => {
        render(<WorkbenchToolbar {...defaultProps} chatLocked/>)

        const chatButton = screen.getByRole("button", {name: "Chat"})

        expect(chatButton).toBeDisabled()
    })

    it("enables the Chat button when chatLocked is false", () => {
        render(<WorkbenchToolbar {...defaultProps} chatLocked={false}/>)

        const chatButton = screen.getByRole("button", {name: "Chat"})

        expect(chatButton).not.toBeDisabled()
    })

    it("keeps the Terminal button always enabled regardless of chatLocked", () => {
        const {rerender} = render(<WorkbenchToolbar {...defaultProps} chatLocked/>)

        expect(screen.getByRole("button", {name: "Terminal"})).not.toBeDisabled()

        rerender(<WorkbenchToolbar {...defaultProps} chatLocked={false}/>)

        expect(screen.getByRole("button", {name: "Terminal"})).not.toBeDisabled()
    })

    it("shows the history button when the chat view is active", () => {
        render(<WorkbenchToolbar {...defaultProps} view="chat"/>)

        expect(screen.getByLabelText("View chat history")).toBeInTheDocument()
    })

    it("hides the history button when the terminal view is active", () => {
        render(<WorkbenchToolbar {...defaultProps} view="terminal"/>)

        expect(screen.queryByLabelText("View chat history")).not.toBeInTheDocument()
    })

    it("calls onToggleHistory when the history button is clicked", () => {
        render(<WorkbenchToolbar {...defaultProps} view="chat"/>)

        fireEvent.click(screen.getByLabelText("View chat history"))

        expect(mockOnToggleHistory).toHaveBeenCalledTimes(1)
    })
})
