import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import SegmentedControl from "@/components/atoms/SegmentedControl/SegmentedControl.tsx"

const OPTIONS = [
    {key: "history", label: "History"},
    {key: "tools", label: "Tools"},
    {key: "vault", label: "Vault"},
]

describe("SegmentedControl", () => {
    it("renders every option label", () => {
        render(<SegmentedControl options={OPTIONS} active="history" onChange={vi.fn()}/>)

        OPTIONS.forEach(({label}) => {
            expect(screen.getByRole("button", {name: label})).toBeInTheDocument()
        })
    })

    it("marks only the active option with the active class", () => {
        render(<SegmentedControl options={OPTIONS} active="tools" onChange={vi.fn()}/>)

        expect(screen.getByRole("button", {name: "Tools"}).className).toMatch(/SegmentActive/)
        expect(screen.getByRole("button", {name: "History"}).className).not.toMatch(/SegmentActive/)
        expect(screen.getByRole("button", {name: "Vault"}).className).not.toMatch(/SegmentActive/)
    })

    it("calls onChange with the option key when a non-active option is clicked", () => {
        const onChange = vi.fn()
        render(<SegmentedControl options={OPTIONS} active="history" onChange={onChange}/>)

        fireEvent.click(screen.getByRole("button", {name: "Vault"}))

        expect(onChange).toHaveBeenCalledTimes(1)
        expect(onChange).toHaveBeenCalledWith("vault")
    })

    it("does not call onChange when a disabled option is clicked", () => {
        const onChange = vi.fn()
        const options = [
            {key: "history", label: "History"},
            {key: "tools", label: "Tools", disabled: true},
        ]
        render(<SegmentedControl options={options} active="history" onChange={onChange}/>)

        fireEvent.click(screen.getByRole("button", {name: "Tools"}))

        expect(onChange).not.toHaveBeenCalled()
    })

    it("renders aria-disabled and tooltip data attributes on a disabled option with a tooltip", () => {
        const options = [
            {key: "history", label: "History"},
            {key: "tools", label: "Tools", disabled: true, tooltip: "Coming soon"},
        ]
        render(<SegmentedControl options={options} active="history" onChange={vi.fn()}/>)

        const tools = screen.getByRole("button", {name: "Tools"})
        expect(tools).toHaveAttribute("aria-disabled", "true")
        expect(tools).toHaveAttribute("data-tooltip-content", "Coming soon")
    })
})
