import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ChatComposer from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.tsx"

function renderComposer(value = "", disabled = false) {
    const onChange = vi.fn()
    const onSend = vi.fn()
    const onNewChat = vi.fn()

    render(
        <ChatComposer
            value={value}
            onChange={onChange}
            onSend={onSend}
            onNewChat={onNewChat}
            disabled={disabled}
            placeholder="Message the workbench…"
        />,
    )

    return {onChange, onSend, onNewChat}
}

describe("ChatComposer", () => {
    it("calls onNewChat when the new chat button is clicked", () => {
        const {onNewChat} = renderComposer()

        fireEvent.click(screen.getByLabelText("New chat"))

        expect(onNewChat).toHaveBeenCalledTimes(1)
    })

    it("disables the new chat button when disabled is true", () => {
        renderComposer("", true)

        expect(screen.getByLabelText("New chat")).toBeDisabled()
    })

    it("renders a textarea with the placeholder", () => {
        renderComposer()

        expect(screen.getByRole("textbox")).toBe(screen.getByPlaceholderText("Message the workbench…"))
    })

    it("edits the draft through the textarea", () => {
        const {onChange} = renderComposer()

        fireEvent.change(screen.getByRole("textbox"), {target: {value: "hi"}})

        expect(onChange).toHaveBeenCalledWith("hi")
    })

    it("sends on Enter when the draft is non-empty", () => {
        const {onSend} = renderComposer("hello")

        fireEvent.keyDown(screen.getByRole("textbox"), {key: "Enter"})

        expect(onSend).toHaveBeenCalledTimes(1)
    })

    it("does not send on Shift+Enter", () => {
        const {onSend} = renderComposer("hello")

        fireEvent.keyDown(screen.getByRole("textbox"), {key: "Enter", shiftKey: true})

        expect(onSend).not.toHaveBeenCalled()
    })

    it("does not send on Enter when the draft is empty", () => {
        const {onSend} = renderComposer("")

        fireEvent.keyDown(screen.getByRole("textbox"), {key: "Enter"})

        expect(onSend).not.toHaveBeenCalled()
    })
})
