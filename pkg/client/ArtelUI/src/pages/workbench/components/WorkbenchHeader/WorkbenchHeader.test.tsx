import {describe, expect, it, vi} from "vitest"
import {render, screen} from "@testing-library/react"
import {MemoryRouter} from "react-router-dom"

import WorkbenchHeader from "@/pages/workbench/components/WorkbenchHeader/WorkbenchHeader.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"

vi.mock("@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx", () => ({
    default: () => <div className="workbench-toolbar" data-testid="workbench-toolbar"/>,
}))

vi.mock("@/pages/workbench/components/SimpleChatTopBar/SimpleChatTopBar.tsx", () => ({
    default: () => <div className="simple-chat-top-bar" data-testid="simple-chat-top-bar"/>,
}))

const simpleChatTopBar = {
    models: ["m1"],
    currentModel: "m1",
    onChangeModel: vi.fn(),
    onNewChat: vi.fn(),
    onToggleHistory: vi.fn(),
}

const defaultProps = {
    isLoading: false,
    effectiveMode: "docker" as WorkbenchMode | "picking",
    exists: true,
    vaultId: "vault-123",
    vaultName: "Test Vault",
    status: "running",
    onStart: vi.fn(),
    onStop: vi.fn(),
    stopping: false,
    starting: false,
    view: "chat" as const,
    onViewChange: vi.fn(),
    chatLocked: false,
    onToggleHistory: vi.fn(),
    simpleChatTopBar,
}

function renderHeader(overrides: Partial<typeof defaultProps> = {}) {
    return render(
        <MemoryRouter>
            <WorkbenchHeader {...defaultProps} {...overrides}/>
        </MemoryRouter>,
    )
}

describe("WorkbenchHeader", () => {
    it("renders WorkbenchToolbar for docker mode", () => {
        renderHeader()

        expect(screen.getByTestId("workbench-toolbar")).toBeInTheDocument()
        expect(screen.queryByTestId("simple-chat-top-bar")).not.toBeInTheDocument()
    })

    it("renders SimpleChatTopBar for simple-chat mode", () => {
        renderHeader({effectiveMode: "simple-chat"})

        expect(screen.getByTestId("simple-chat-top-bar")).toBeInTheDocument()
        expect(screen.queryByTestId("workbench-toolbar")).not.toBeInTheDocument()
    })

    it("renders the back-to-vaults fallback link while picking a mode", () => {
        renderHeader({effectiveMode: "picking"})

        expect(screen.getByText("← Back to vaults")).toBeInTheDocument()
        expect(screen.queryByTestId("workbench-toolbar")).not.toBeInTheDocument()
        expect(screen.queryByTestId("simple-chat-top-bar")).not.toBeInTheDocument()
    })

    it("renders the fallback link while still loading, regardless of mode", () => {
        renderHeader({isLoading: true})

        expect(screen.getByText("← Back to vaults")).toBeInTheDocument()
    })

    it("renders the fallback link when there is no vaultId", () => {
        renderHeader({vaultId: undefined})

        expect(screen.getByText("← Back to vaults")).toBeInTheDocument()
    })
})
