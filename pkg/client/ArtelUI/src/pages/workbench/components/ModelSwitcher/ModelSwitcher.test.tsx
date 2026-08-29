import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ModelSwitcher from "@/pages/workbench/components/ModelSwitcher/ModelSwitcher.tsx"

interface MockDropdownProps {
    placeholder?: string
    onChange: (value: string[]) => void
}

vi.mock("@vervstack/chures", async (importActual) => {
    const actual = await importActual<typeof import("@vervstack/chures")>()
    return {
        ...actual,
        Dropdown: ({placeholder, onChange}: MockDropdownProps) => (
            <span
                role="button"
                data-testid="dropdown"
                data-placeholder={placeholder}
                onClick={() => onChange(["anthropic/model-b"])}
            >
                {placeholder}
            </span>
        ),
    }
})

describe("ModelSwitcher", () => {
    const models = ["anthropic/model-a", "anthropic/model-b"]

    it("defaults the dropdown placeholder to \"Model…\"", () => {
        render(<ModelSwitcher models={models} value="" onChange={vi.fn()}/>)

        expect(screen.getByTestId("dropdown")).toHaveAttribute("data-placeholder", "Model…")
    })

    it("forwards a custom placeholder to the dropdown", () => {
        render(<ModelSwitcher models={models} value="" onChange={vi.fn()} placeholder="Claude Code"/>)

        expect(screen.getByTestId("dropdown")).toHaveAttribute("data-placeholder", "Claude Code")
    })

    it("calls onChange with the picked model id", () => {
        const onChange = vi.fn()
        render(<ModelSwitcher models={models} value="" onChange={onChange}/>)

        fireEvent.click(screen.getByTestId("dropdown"))

        expect(onChange).toHaveBeenCalledWith("anthropic/model-b")
    })

    it("never calls onChange when disabled", () => {
        const onChange = vi.fn()
        render(<ModelSwitcher models={models} value="" onChange={onChange} disabled/>)

        fireEvent.click(screen.getByTestId("dropdown"))

        expect(onChange).not.toHaveBeenCalled()
    })

    it("marks the wrapper aria-disabled and adds the disabled class when disabled", () => {
        const {container} = render(<ModelSwitcher models={models} value="" onChange={vi.fn()} disabled/>)
        const wrapper = container.firstChild as HTMLElement

        expect(wrapper).toHaveAttribute("aria-disabled", "true")
        expect(wrapper.className).toMatch(/ModelSwitcherDisabled/)
    })

    it("adds the needs-attention class when needsAttention is true", () => {
        const {container} = render(<ModelSwitcher models={models} value="" onChange={vi.fn()} needsAttention/>)
        const wrapper = container.firstChild as HTMLElement

        expect(wrapper.className).toMatch(/ModelSwitcherNeedsAttention/)
    })

    it("omits the needs-attention class by default", () => {
        const {container} = render(<ModelSwitcher models={models} value="" onChange={vi.fn()}/>)
        const wrapper = container.firstChild as HTMLElement

        expect(wrapper.className).not.toMatch(/ModelSwitcherNeedsAttention/)
    })
})
