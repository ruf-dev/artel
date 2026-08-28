import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ChatComposer from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.tsx"

function renderComposer(value = "", disabled = false, attachedPaths: string[] = []) {
    const onChange = vi.fn()
    const onSend = vi.fn()
    const onNewChat = vi.fn()
    const onOpenTweaks = vi.fn()
    const onRemoveAttachment = vi.fn()

    render(
        <ChatComposer
            value={value}
            onChange={onChange}
            onSend={onSend}
            onNewChat={onNewChat}
            disabled={disabled}
            placeholder="Message the workbench…"
            onOpenTweaks={onOpenTweaks}
            attachedPaths={attachedPaths}
            onRemoveAttachment={onRemoveAttachment}
        />,
    )

    return {onChange, onSend, onNewChat, onOpenTweaks, onRemoveAttachment}
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

    it("a composer chip opens the tweaks panel at its section", () => {
        const {onOpenTweaks} = renderComposer()

        fireEvent.click(screen.getByRole("button", {name: "Connections"}))

        expect(onOpenTweaks).toHaveBeenCalledWith("connections")
    })

    it("renders an attached-file chip per attachedPaths entry", () => {
        renderComposer("", false, ["notes/a.md", "b.md"])

        expect(screen.getByText("a")).toBeInTheDocument()
        expect(screen.getByText("b")).toBeInTheDocument()
    })

    it("renders no attached-file chips when attachedPaths is empty", () => {
        renderComposer("", false, [])

        expect(screen.queryByRole("button", {name: /^Remove /})).not.toBeInTheDocument()
    })

    it("calls onRemoveAttachment with the path when a chip's × is clicked", () => {
        const {onRemoveAttachment} = renderComposer("", false, ["notes/a.md"])

        fireEvent.click(screen.getByRole("button", {name: "Remove notes/a.md"}))

        expect(onRemoveAttachment).toHaveBeenCalledWith("notes/a.md")
    })
})
