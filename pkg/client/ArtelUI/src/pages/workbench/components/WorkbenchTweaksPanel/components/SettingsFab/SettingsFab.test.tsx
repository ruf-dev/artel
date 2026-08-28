import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import SettingsFab from "@/pages/workbench/components/WorkbenchTweaksPanel/components/SettingsFab/SettingsFab.tsx"

vi.mock("morphicons/react", () => ({MorphIcon: () => null}))

describe("SettingsFab", () => {
    it("reflects the open state on aria-expanded and swaps the label", () => {
        const {rerender} = render(<SettingsFab open={false} onToggle={vi.fn()}/>)

        const closed = screen.getByRole("button", {name: "Settings"})
        expect(closed).toHaveAttribute("aria-expanded", "false")

        rerender(<SettingsFab open onToggle={vi.fn()}/>)

        const opened = screen.getByRole("button", {name: "Close settings"})
        expect(opened).toHaveAttribute("aria-expanded", "true")
    })

    it("calls onToggle when clicked", () => {
        const onToggle = vi.fn()
        render(<SettingsFab open={false} onToggle={onToggle}/>)

        fireEvent.click(screen.getByRole("button", {name: "Settings"}))

        expect(onToggle).toHaveBeenCalledTimes(1)
    })
})
