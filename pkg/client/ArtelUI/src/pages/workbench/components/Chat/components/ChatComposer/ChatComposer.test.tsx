import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ChatComposer from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.tsx"

function renderComposer(disabled = false) {
    const onChange = vi.fn()
    const onSend = vi.fn()
    const onNewChat = vi.fn()

    render(
        <ChatComposer
            value=""
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
        renderComposer(true)

        expect(screen.getByLabelText("New chat")).toBeDisabled()
    })
})
