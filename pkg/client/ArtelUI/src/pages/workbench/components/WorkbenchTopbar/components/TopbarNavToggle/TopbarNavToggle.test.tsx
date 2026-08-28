import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import TopbarNavToggle from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarNavToggle/TopbarNavToggle.tsx"

vi.mock("morphicons/react", () => ({MorphIcon: () => null}))

// The container is `display: none` above the mobile breakpoint (jsdom has no
// viewport, so the CSS-modules `css: true` transform always applies the default) —
// query with `hidden: true` so the burger is still found in the a11y tree.
describe("TopbarNavToggle", () => {
    it("renders a button with the accessible name 'Open menu'", () => {
        render(<TopbarNavToggle onToggle={vi.fn()}/>)

        expect(screen.getByRole("button", {name: "Open menu", hidden: true})).toBeInTheDocument()
    })

    it("calls onToggle once when the button is clicked", () => {
        const onToggle = vi.fn()
        render(<TopbarNavToggle onToggle={onToggle}/>)

        fireEvent.click(screen.getByRole("button", {name: "Open menu", hidden: true}))

        expect(onToggle).toHaveBeenCalledTimes(1)
    })
})
