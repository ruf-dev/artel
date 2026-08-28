import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import SidebarToggleButton from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarToggleButton/SidebarToggleButton.tsx"

vi.mock("morphicons/react", () => ({MorphIcon: () => null}))

describe("SidebarToggleButton", () => {
    it("reflects the open state on aria-expanded", () => {
        const {rerender} = render(<SidebarToggleButton open={false} onToggle={vi.fn()}/>)

        expect(screen.getByRole("button", {name: "Toggle conversations"}))
            .toHaveAttribute("aria-expanded", "false")

        rerender(<SidebarToggleButton open onToggle={vi.fn()}/>)

        expect(screen.getByRole("button", {name: "Toggle conversations"}))
            .toHaveAttribute("aria-expanded", "true")
    })

    it("calls onToggle when clicked", () => {
        const onToggle = vi.fn()
        render(<SidebarToggleButton open={false} onToggle={onToggle}/>)

        fireEvent.click(screen.getByRole("button", {name: "Toggle conversations"}))

        expect(onToggle).toHaveBeenCalledTimes(1)
    })
})
