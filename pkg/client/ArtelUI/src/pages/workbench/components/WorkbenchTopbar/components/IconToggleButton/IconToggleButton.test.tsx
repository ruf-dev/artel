import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import IconToggleButton
    from "@/pages/workbench/components/WorkbenchTopbar/components/IconToggleButton/IconToggleButton.tsx"

describe("IconToggleButton", () => {
    it("renders the icon and exposes the label as aria-label", () => {
        render(<IconToggleButton icon={<span data-testid="icon">i</span>} label="Toggle conversations"/>)

        expect(screen.getByTestId("icon")).toBeInTheDocument()
        expect(screen.getByRole("button", {name: "Toggle conversations"})).toBeInTheDocument()
    })

    it("fires onClick when enabled", () => {
        const onClick = vi.fn()
        render(<IconToggleButton icon={null} label="Tweaks" onClick={onClick}/>)

        fireEvent.click(screen.getByRole("button", {name: "Tweaks"}))

        expect(onClick).toHaveBeenCalledTimes(1)
    })

    it("blocks onClick and sets aria-disabled + tooltip content when disabled", () => {
        const onClick = vi.fn()
        render(
            <IconToggleButton icon={null} label="Tweaks" onClick={onClick} disabled tooltip="Coming soon"/>,
        )

        const btn = screen.getByRole("button", {name: "Tweaks"})
        fireEvent.click(btn)

        expect(onClick).not.toHaveBeenCalled()
        expect(btn).toHaveAttribute("aria-disabled", "true")
        expect(btn).toHaveAttribute("data-tooltip-id", "root-tooltip")
        expect(btn).toHaveAttribute("data-tooltip-content", "Coming soon")
    })

    it("adds the active class when active", () => {
        render(<IconToggleButton icon={null} label="Toggle conversations" active/>)

        expect(screen.getByRole("button", {name: "Toggle conversations"}).className).toMatch(/BtnActive/)
    })
})
