import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ChatMessageColumn from "@/pages/workbench/components/Chat/components/ChatMessageColumn/ChatMessageColumn.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

describe("ChatMessageColumn", () => {
    it("trailing retry on an error preceded by a user_message with an id"
        + " calls onResendMessage, not onRetryMessage", () => {
        const items: ChatItem[] = [
            {kind: "user_message", key: "u1", id: "m1", text: "hello"},
            {kind: "error", key: "e1", text: "turn failed"},
        ]
        const onRetryMessage = vi.fn()
        const onResendMessage = vi.fn()

        render(
            <ChatMessageColumn
                items={items}
                onRetryMessage={onRetryMessage}
                onResendMessage={onResendMessage}
                onPermissionDecision={vi.fn()}
            />,
        )

        fireEvent.click(screen.getByLabelText("Retry message"))

        expect(onResendMessage).toHaveBeenCalledWith("m1", "hello")
        expect(onRetryMessage).not.toHaveBeenCalled()
    })

    it("trailing retry on a dangling user_message with an id calls onResendMessage, not onRetryMessage", () => {
        const items: ChatItem[] = [
            {kind: "user_message", key: "u1", id: "m1", text: "hello"},
        ]
        const onRetryMessage = vi.fn()
        const onResendMessage = vi.fn()

        render(
            <ChatMessageColumn
                items={items}
                onRetryMessage={onRetryMessage}
                onResendMessage={onResendMessage}
                onPermissionDecision={vi.fn()}
            />,
        )

        fireEvent.click(screen.getByLabelText("Retry message"))

        expect(onResendMessage).toHaveBeenCalledWith("m1", "hello")
        expect(onRetryMessage).not.toHaveBeenCalled()
    })

    it("disables the trailing retry when the failed user_message has no id", () => {
        const items: ChatItem[] = [
            {kind: "user_message", key: "u1", text: "hello"},
        ]
        const onRetryMessage = vi.fn()
        const onResendMessage = vi.fn()

        render(
            <ChatMessageColumn
                items={items}
                onRetryMessage={onRetryMessage}
                onResendMessage={onResendMessage}
                onPermissionDecision={vi.fn()}
            />,
        )

        const retryBtn = screen.getByLabelText("Retry message")
        expect(retryBtn).toHaveAttribute("aria-disabled", "true")
        expect(retryBtn).toHaveAttribute("data-tooltip-content", "Nothing to retry")

        fireEvent.click(retryBtn)
        expect(onResendMessage).not.toHaveBeenCalled()
        expect(onRetryMessage).not.toHaveBeenCalled()
    })

    it("mid-transcript hover-retry on a completed assistant_message still calls onRetryMessage, unchanged", () => {
        const items: ChatItem[] = [
            {kind: "user_message", key: "u1", id: "m1", text: "hi"},
            {kind: "assistant_message", key: "a1", id: "a1", text: "reply", done: true},
        ]
        const onRetryMessage = vi.fn()
        const onResendMessage = vi.fn()

        render(
            <ChatMessageColumn
                items={items}
                onRetryMessage={onRetryMessage}
                onResendMessage={onResendMessage}
                onPermissionDecision={vi.fn()}
            />,
        )

        // The assistant message is the last item, so no trailing-failure retry is
        // rendered - only the mid-transcript hover-retry, backed by onRetryMessage.
        fireEvent.click(screen.getByLabelText("Retry message"))

        expect(onRetryMessage).toHaveBeenCalledWith("hi")
        expect(onResendMessage).not.toHaveBeenCalled()
    })

    it("suppresses the trailing-failure retry entirely while a turn is pending", () => {
        const items: ChatItem[] = [
            {kind: "user_message", key: "u1", id: "m1", text: "hello"},
        ]
        const onRetryMessage = vi.fn()
        const onResendMessage = vi.fn()

        render(
            <ChatMessageColumn
                items={items}
                onRetryMessage={onRetryMessage}
                onResendMessage={onResendMessage}
                onPermissionDecision={vi.fn()}
                pending={{bucket: "stuck"}}
            />,
        )

        // The dangling user_message itself shows no retry action while pending...
        expect(screen.queryByLabelText("Retry message")).not.toBeInTheDocument()

        // ...only PendingAssistantBubble's own stuck-turn retry, which is untouched
        // by this change and still goes through onRetryMessage.
        fireEvent.click(screen.getByText("Retry"))
        expect(onRetryMessage).toHaveBeenCalledWith("hello")
        expect(onResendMessage).not.toHaveBeenCalled()
    })
})
