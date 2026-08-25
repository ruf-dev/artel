import {describe, expect, it} from "vitest"
import {render, screen} from "@testing-library/react"

import HistorySidebarShell from "@/pages/workbench/components/HistorySidebarShell/HistorySidebarShell.tsx"

describe("HistorySidebarShell", () => {
    it("renders its children", () => {
        render(
            <HistorySidebarShell open={false}>
                <div className="sidebar-content">Sidebar content</div>
            </HistorySidebarShell>,
        )

        expect(screen.getByText("Sidebar content")).toBeInTheDocument()
    })

    it("applies the open class when open is true", () => {
        const {container} = render(
            <HistorySidebarShell open={true}>
                <div className="sidebar-content">Sidebar content</div>
            </HistorySidebarShell>,
        )

        expect(container.querySelector('[class*="HistorySidebarShellOpen"]')).toBeTruthy()
    })

    it("does not apply the open class when open is false", () => {
        const {container} = render(
            <HistorySidebarShell open={false}>
                <div className="sidebar-content">Sidebar content</div>
            </HistorySidebarShell>,
        )

        expect(container.querySelector('[class*="HistorySidebarShellOpen"]')).toBeFalsy()
    })
})
