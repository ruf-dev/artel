import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import WorkbenchSidebar from "@/pages/workbench/components/WorkbenchSidebar/WorkbenchSidebar.tsx"
import {WorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"

vi.mock("morphicons/react", () => ({MorphIcon: () => null}))

vi.mock("@/pages/workbench/components/WorkbenchSidebar/components/SidebarBrand/SidebarBrand.tsx", () => ({
    default: () => <div className="brand" data-testid="sidebar-brand"/>,
}))

vi.mock("@/pages/workbench/components/WorkbenchSidebar/components/SidebarFooter/SidebarFooter.tsx", () => ({
    default: () => <div className="footer" data-testid="sidebar-footer"/>,
}))

vi.mock("@/pages/workbench/components/WorkbenchSidebar/components/HistoryPane/HistoryPane.tsx", () => ({
    default: () => <div className="history-pane" data-testid="history-pane"/>,
}))

function makeHistory(over: Partial<WorkbenchHistory> = {}): WorkbenchHistory {
    return {
        rows: [],
        loading: false,
        activeId: undefined,
        onSelect: vi.fn(),
        onDelete: undefined,
        onNewChat: vi.fn(),
        ...over,
    }
}

function makeProps(over: Partial<Parameters<typeof WorkbenchSidebar>[0]> = {}) {
    return {history: makeHistory(), navOpen: true, onToggleNav: vi.fn(), ...over}
}

describe("WorkbenchSidebar", () => {
    it("renders the History pane and keeps the footer and brand", () => {
        render(<WorkbenchSidebar {...makeProps()}/>)

        expect(screen.getByTestId("history-pane")).toBeInTheDocument()
        expect(screen.getByTestId("sidebar-brand")).toBeInTheDocument()
        expect(screen.getByTestId("sidebar-footer")).toBeInTheDocument()
    })

    it("renders the collapse toggle and calls onToggleNav when it is clicked", () => {
        const onToggleNav = vi.fn()
        render(<WorkbenchSidebar {...makeProps({onToggleNav})}/>)

        fireEvent.click(screen.getByRole("button", {name: "Toggle conversations"}))

        expect(onToggleNav).toHaveBeenCalledTimes(1)
    })

    it("keeps the History pane when the disabled Tools segment is clicked", () => {
        render(<WorkbenchSidebar {...makeProps()}/>)

        fireEvent.click(screen.getByRole("button", {name: /tools/i}))

        expect(screen.getByTestId("history-pane")).toBeInTheDocument()
    })

    it("keeps the History pane when the disabled Vault segment is clicked", () => {
        render(<WorkbenchSidebar {...makeProps()}/>)

        fireEvent.click(screen.getByRole("button", {name: /vault/i}))

        expect(screen.getByTestId("history-pane")).toBeInTheDocument()
    })

    it("hides the History pane and brand when collapsed, keeping the toggle", () => {
        render(<WorkbenchSidebar {...makeProps({navOpen: false})}/>)

        expect(screen.queryByTestId("history-pane")).not.toBeInTheDocument()
        expect(screen.queryByTestId("sidebar-brand")).not.toBeInTheDocument()
        expect(screen.getByRole("button", {name: "Toggle conversations"})).toBeInTheDocument()
    })

    it("shows the History pane and brand when expanded", () => {
        render(<WorkbenchSidebar {...makeProps({navOpen: true})}/>)

        expect(screen.getByTestId("history-pane")).toBeInTheDocument()
        expect(screen.getByTestId("sidebar-brand")).toBeInTheDocument()
        expect(screen.getByRole("button", {name: "Toggle conversations"})).toBeInTheDocument()
    })

    it("always renders the backdrop, whether the drawer is open or closed", () => {
        const {rerender} = render(<WorkbenchSidebar {...makeProps({navOpen: true})}/>)
        expect(screen.getByTestId("sidebar-backdrop")).toBeInTheDocument()

        rerender(<WorkbenchSidebar {...makeProps({navOpen: false})}/>)
        expect(screen.getByTestId("sidebar-backdrop")).toBeInTheDocument()
    })

    it("clicking the backdrop calls onToggleNav", () => {
        const onToggleNav = vi.fn()
        render(<WorkbenchSidebar {...makeProps({navOpen: true, onToggleNav})}/>)

        fireEvent.click(screen.getByTestId("sidebar-backdrop"))

        expect(onToggleNav).toHaveBeenCalledTimes(1)
    })
})
