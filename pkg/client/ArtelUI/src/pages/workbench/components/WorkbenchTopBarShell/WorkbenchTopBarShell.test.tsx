import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"
import {MemoryRouter} from "react-router-dom"

import WorkbenchTopBarShell from "@/pages/workbench/components/WorkbenchTopBarShell/WorkbenchTopBarShell.tsx"

const mockNavigate = vi.fn()

vi.mock("react-router-dom", async () => {
    const actual = await vi.importActual("react-router-dom")
    return {...actual, useNavigate: () => mockNavigate}
})

function renderShell(statusBadge?: React.ReactNode, actions?: React.ReactNode) {
    return render(
        <MemoryRouter>
            <WorkbenchTopBarShell vaultName="Test Vault" statusBadge={statusBadge} actions={actions}/>
        </MemoryRouter>,
    )
}

describe("WorkbenchTopBarShell", () => {
    it("renders the vault name", () => {
        renderShell()

        expect(screen.getByText("Test Vault")).toBeInTheDocument()
    })

    it("navigates back to vaults when the back button is clicked", () => {
        renderShell()

        fireEvent.click(screen.getByLabelText("Back to vaults"))

        expect(mockNavigate).toHaveBeenCalledWith("/")
    })

    it("renders the statusBadge slot when passed", () => {
        renderShell(<div className="status-badge" data-testid="status-badge">Running</div>)

        expect(screen.getByTestId("status-badge")).toBeInTheDocument()
    })

    it("renders the actions slot when passed", () => {
        renderShell(undefined, <div className="actions" data-testid="actions">Actions</div>)

        expect(screen.getByTestId("actions")).toBeInTheDocument()
    })
})
