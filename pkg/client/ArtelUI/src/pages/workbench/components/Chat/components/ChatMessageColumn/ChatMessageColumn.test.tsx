import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ChatMessageColumn from "@/pages/workbench/components/Chat/components/ChatMessageColumn/ChatMessageColumn.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

// Shared render helper - pulled out so each `it` below stays short enough to
// keep the describe callback under the max-lines-per-function budget.
function renderColumn(items: ChatItem[], pending?: {bucket: "normal" | "slow" | "stuck"}) {
    const onRetryMessage = vi.fn()
    const onResendMessage = vi.fn()

    render(
        <ChatMessageColumn
            items={items}
            vaultId="v1"
            onRetryMessage={onRetryMessage}
            onResendMessage={onResendMessage}
            onPermissionDecision={vi.fn()}
            pending={pending}
        />,
    )

    return {onRetryMessage, onResendMessage}
}

describe("ChatMessageColumn", () => {
    it("trailing retry on an error preceded by a user_message with an id"
        + " calls onResendMessage, not onRetryMessage", () => {
        const {onRetryMessage, onResendMessage} = renderColumn([
            {kind: "user_message", key: "u1", id: "m1", text: "hello"},
            {kind: "error", key: "e1", text: "turn failed"},
        ])

        fireEvent.click(screen.getByLabelText("Retry message"))

        expect(onResendMessage).toHaveBeenCalledWith("m1", "hello")
        expect(onRetryMessage).not.toHaveBeenCalled()
    })

    it("trailing retry on a dangling user_message with an id calls onResendMessage, not onRetryMessage", () => {
        const {onRetryMessage, onResendMessage} = renderColumn([
            {kind: "user_message", key: "u1", id: "m1", text: "hello"},
        ])

        fireEvent.click(screen.getByLabelText("Retry message"))

        expect(onResendMessage).toHaveBeenCalledWith("m1", "hello")
        expect(onRetryMessage).not.toHaveBeenCalled()
    })

    it("disables the trailing retry when the failed user_message has no id", () => {
        const {onRetryMessage, onResendMessage} = renderColumn([
            {kind: "user_message", key: "u1", text: "hello"},
        ])

        const retryBtn = screen.getByLabelText("Retry message")
        expect(retryBtn).toHaveAttribute("aria-disabled", "true")
        expect(retryBtn).toHaveAttribute("data-tooltip-content", "Nothing to retry")

        fireEvent.click(retryBtn)
        expect(onResendMessage).not.toHaveBeenCalled()
        expect(onRetryMessage).not.toHaveBeenCalled()
    })

    it("mid-transcript hover-retry on a completed assistant_message still calls onRetryMessage, unchanged", () => {
        const {onRetryMessage, onResendMessage} = renderColumn([
            {kind: "user_message", key: "u1", id: "m1", text: "hi"},
            {kind: "assistant_message", key: "a1", id: "a1", text: "reply", done: true},
        ])

        // The assistant message is the last item, so no trailing-failure retry is
        // rendered - only the mid-transcript hover-retry, backed by onRetryMessage.
        fireEvent.click(screen.getByLabelText("Retry message"))

        expect(onRetryMessage).toHaveBeenCalledWith("hi")
        expect(onResendMessage).not.toHaveBeenCalled()
    })

    it("suppresses the trailing-failure retry entirely while a turn is pending", () => {
        const {onRetryMessage, onResendMessage} = renderColumn([
            {kind: "user_message", key: "u1", id: "m1", text: "hello"},
        ], {bucket: "stuck"})

        // The dangling user_message itself shows no retry action while pending...
        expect(screen.queryByLabelText("Retry message")).not.toBeInTheDocument()

        // ...only PendingAssistantBubble's own stuck-turn retry, which is untouched
        // by this change and still goes through onRetryMessage.
        fireEvent.click(screen.getByText("Retry"))
        expect(onRetryMessage).toHaveBeenCalledWith("hello")
        expect(onResendMessage).not.toHaveBeenCalled()
    })
})
