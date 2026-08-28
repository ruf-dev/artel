import {describe, expect, it, vi} from "vitest"
import {render, screen} from "@testing-library/react"
import {MemoryRouter} from "react-router-dom"

import WorkbenchHeader from "@/pages/workbench/components/WorkbenchHeader/WorkbenchHeader.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"

vi.mock("@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx", () => ({
    default: () => <div className="workbench-toolbar" data-testid="workbench-toolbar"/>,
}))

const defaultProps = {
    isLoading: false,
    effectiveMode: "docker" as WorkbenchMode | "picking",
    exists: true,
    vaultId: "vault-123",
    status: "running",
    onStart: vi.fn(),
    onStop: vi.fn(),
    stopping: false,
    starting: false,
    view: "chat" as const,
    onViewChange: vi.fn(),
    chatLocked: false,
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
    })

    it("renders nothing for api mode", () => {
        renderHeader({effectiveMode: "api"})

        expect(screen.queryByTestId("workbench-toolbar")).not.toBeInTheDocument()
    })

    it("renders nothing while picking a mode", () => {
        renderHeader({effectiveMode: "picking"})

        expect(screen.queryByTestId("workbench-toolbar")).not.toBeInTheDocument()
    })

    it("renders nothing while still loading, regardless of mode", () => {
        renderHeader({isLoading: true})

        expect(screen.queryByTestId("workbench-toolbar")).not.toBeInTheDocument()
    })

    it("renders nothing when there is no vaultId", () => {
        renderHeader({vaultId: undefined})

        expect(screen.queryByTestId("workbench-toolbar")).not.toBeInTheDocument()
    })
})
