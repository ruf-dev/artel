import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import Textarea from "@/components/atoms/Textarea/Textarea.tsx"

describe("Textarea", () => {
    it("renders the current value", () => {
        render(<Textarea value="hello" setValue={vi.fn()}/>)

        expect(screen.getByRole("textbox")).toHaveValue("hello")
    })

    it("calls setValue on change", () => {
        const setValue = vi.fn()
        render(<Textarea value="" setValue={setValue}/>)

        fireEvent.change(screen.getByRole("textbox"), {target: {value: "typed"}})

        expect(setValue).toHaveBeenCalledWith("typed")
    })

    it("forwards onKeyDown", () => {
        const onKeyDown = vi.fn()
        render(<Textarea value="" setValue={vi.fn()} onKeyDown={onKeyDown}/>)

        fireEvent.keyDown(screen.getByRole("textbox"), {key: "Enter"})

        expect(onKeyDown).toHaveBeenCalledTimes(1)
    })

    it("renders a disabled textarea", () => {
        render(<Textarea value="" setValue={vi.fn()} disabled/>)

        expect(screen.getByRole("textbox")).toBeDisabled()
    })
})
