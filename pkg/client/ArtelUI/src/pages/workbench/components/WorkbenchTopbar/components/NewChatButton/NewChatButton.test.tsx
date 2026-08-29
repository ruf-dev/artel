import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import NewChatButton
    from "@/pages/workbench/components/WorkbenchTopbar/components/NewChatButton/NewChatButton.tsx"

vi.mock("morphicons/react", () => ({
    // MessageSquare has 1 path segment, SquarePen has 2 — length tells the two
    // lucide icon data arrays apart without needing a hoisted reference to either.
    MorphIcon: ({icon}: {icon: unknown[]}) => <span data-testid="icon" data-segments={icon.length}/>,
}))

describe("NewChatButton", () => {
    it("shows the chat-bubble glyph by default", () => {
        render(<NewChatButton onClick={vi.fn()}/>)

        expect(screen.getByTestId("icon")).toHaveAttribute("data-segments", "1")
    })

    it("morphs to the compose glyph on hover and back on hover-out", () => {
        render(<NewChatButton onClick={vi.fn()}/>)

        const btn = screen.getByRole("button", {name: "New chat"})

        fireEvent.mouseEnter(btn)
        expect(screen.getByTestId("icon")).toHaveAttribute("data-segments", "2")

        fireEvent.mouseLeave(btn)
        expect(screen.getByTestId("icon")).toHaveAttribute("data-segments", "1")
    })

    it("morphs on focus and back on blur", () => {
        render(<NewChatButton onClick={vi.fn()}/>)

        const btn = screen.getByRole("button", {name: "New chat"})

        fireEvent.focus(btn)
        expect(screen.getByTestId("icon")).toHaveAttribute("data-segments", "2")

        fireEvent.blur(btn)
        expect(screen.getByTestId("icon")).toHaveAttribute("data-segments", "1")
    })

    it("calls onClick when clicked", () => {
        const onClick = vi.fn()
        render(<NewChatButton onClick={onClick}/>)

        fireEvent.click(screen.getByRole("button", {name: "New chat"}))

        expect(onClick).toHaveBeenCalledTimes(1)
    })

    it("is disabled when disabled is true", () => {
        render(<NewChatButton onClick={vi.fn()} disabled/>)

        expect(screen.getByRole("button", {name: "New chat"})).toBeDisabled()
    })
})
