import {describe, expect, it} from "vitest"
import {render, screen} from "@testing-library/react"

import WorkbenchTopBarShell from "@/pages/workbench/components/WorkbenchTopBarShell/WorkbenchTopBarShell.tsx"

function renderShell(statusBadge?: React.ReactNode, actions?: React.ReactNode) {
    return render(
        <WorkbenchTopBarShell statusBadge={statusBadge} actions={actions}/>
    )
}

describe("WorkbenchTopBarShell", () => {
    it("renders the statusBadge slot when passed", () => {
        renderShell(<div className="status-badge" data-testid="status-badge">Running</div>)

        expect(screen.getByTestId("status-badge")).toBeInTheDocument()
    })

    it("renders the actions slot when passed", () => {
        renderShell(undefined, <div className="actions" data-testid="actions">Actions</div>)

        expect(screen.getByTestId("actions")).toBeInTheDocument()
    })
})
