import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import PendingAssistantBubble
    from "@/pages/workbench/components/Chat/components/PendingAssistantBubble/PendingAssistantBubble.tsx"

describe("PendingAssistantBubble", () => {
    it("renders label when provided", () => {
        render(
            <PendingAssistantBubble
                label="Claude Code"
                bucket="normal"
                onRetry={vi.fn()}
            />,
        )

        expect(screen.getByText("Claude Code")).toBeInTheDocument()
    })

    it("shows 'Thinking…' for normal bucket", () => {
        render(
            <PendingAssistantBubble
                bucket="normal"
                onRetry={vi.fn()}
            />,
        )

        expect(screen.getByText("Thinking…")).toBeInTheDocument()
    })

    it("shows 'Still working…' for slow bucket", () => {
        render(
            <PendingAssistantBubble
                bucket="slow"
                onRetry={vi.fn()}
            />,
        )

        expect(screen.getByText("Still working…")).toBeInTheDocument()
    })

    it("shows stuck message and Retry button for stuck bucket", () => {
        render(
            <PendingAssistantBubble
                bucket="stuck"
                onRetry={vi.fn()}
            />,
        )

        expect(screen.getByText("This is taking longer than usual — the model may be stuck.")).toBeInTheDocument()
        expect(screen.getByText("Retry")).toBeInTheDocument()
    })

    it("calls onRetry when stuck Retry button is clicked", () => {
        const onRetry = vi.fn()
        render(
            <PendingAssistantBubble
                bucket="stuck"
                onRetry={onRetry}
            />,
        )

        fireEvent.click(screen.getByText("Retry"))
        expect(onRetry).toHaveBeenCalled()
    })

    it("does not show Retry button for normal bucket", () => {
        render(
            <PendingAssistantBubble
                bucket="normal"
                onRetry={vi.fn()}
            />,
        )

        expect(screen.queryByText("Retry")).not.toBeInTheDocument()
    })

    it("does not show Retry button for slow bucket", () => {
        render(
            <PendingAssistantBubble
                bucket="slow"
                onRetry={vi.fn()}
            />,
        )

        expect(screen.queryByText("Retry")).not.toBeInTheDocument()
    })

    it("disables Retry button when retryDisabled is true", () => {
        const onRetry = vi.fn()
        render(
            <PendingAssistantBubble
                bucket="stuck"
                onRetry={onRetry}
                retryDisabled={true}
            />,
        )

        const retryButton = screen.getByText("Retry").closest("button")
        expect(retryButton).toBeDisabled()

        fireEvent.click(retryButton!)
        expect(onRetry).not.toHaveBeenCalled()
    })
})
