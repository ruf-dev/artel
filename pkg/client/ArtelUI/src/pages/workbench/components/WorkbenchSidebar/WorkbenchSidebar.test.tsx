import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import WorkbenchSidebar from "@/pages/workbench/components/WorkbenchSidebar/WorkbenchSidebar.tsx"
import {WorkbenchHistory} from "@/pages/workbench/processes/useWorkbenchHistory.ts"

let brandProps: Record<string, unknown> | undefined

vi.mock("@/pages/workbench/components/WorkbenchSidebar/components/SidebarBrand/SidebarBrand.tsx", () => ({
    default: (props: Record<string, unknown>) => {
        brandProps = props
        return <div className="brand" data-testid="sidebar-brand"/>
    },
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

describe("WorkbenchSidebar", () => {
    it("forwards showCloseButton to the brand as showClose", () => {
        render(<WorkbenchSidebar history={makeHistory()} showCloseButton/>)

        expect(brandProps?.showClose).toBe(true)
    })

    it("leaves showClose undefined on the brand when not provided", () => {
        render(<WorkbenchSidebar history={makeHistory()}/>)

        expect(brandProps?.showClose).toBeUndefined()
    })

    it("calls history.onNewChat when the New chat button is clicked", () => {
        const onNewChat = vi.fn()
        render(<WorkbenchSidebar history={makeHistory({onNewChat})}/>)

        fireEvent.click(screen.getByRole("button", {name: /new chat/i}))

        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("renders the History pane and keeps the footer and brand", () => {
        render(<WorkbenchSidebar history={makeHistory()}/>)

        expect(screen.getByTestId("history-pane")).toBeInTheDocument()
        expect(screen.getByTestId("sidebar-brand")).toBeInTheDocument()
        expect(screen.getByTestId("sidebar-footer")).toBeInTheDocument()
    })

    it("keeps the History pane when the disabled Tools segment is clicked", () => {
        render(<WorkbenchSidebar history={makeHistory()}/>)

        fireEvent.click(screen.getByRole("button", {name: /tools/i}))

        expect(screen.getByTestId("history-pane")).toBeInTheDocument()
    })

    it("keeps the History pane when the disabled Vault segment is clicked", () => {
        render(<WorkbenchSidebar history={makeHistory()}/>)

        fireEvent.click(screen.getByRole("button", {name: /vault/i}))

        expect(screen.getByTestId("history-pane")).toBeInTheDocument()
    })
})
