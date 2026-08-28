import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import MessageActions from "@/pages/workbench/components/Chat/components/MessageActions/MessageActions.tsx"

describe("MessageActions", () => {
    it("renders a copy and a retry button", () => {
        render(<MessageActions onCopy={vi.fn()} onRetry={vi.fn()}/>)

        expect(screen.getByLabelText("Copy message")).toBeInTheDocument()
        expect(screen.getByLabelText("Retry message")).toBeInTheDocument()
    })

    it("calls onCopy when Copy is clicked", () => {
        const onCopy = vi.fn()
        render(<MessageActions onCopy={onCopy} onRetry={vi.fn()}/>)

        fireEvent.click(screen.getByLabelText("Copy message"))

        expect(onCopy).toHaveBeenCalledTimes(1)
    })

    it("calls onRetry when Retry is clicked", () => {
        const onRetry = vi.fn()
        render(<MessageActions onCopy={vi.fn()} onRetry={onRetry}/>)

        fireEvent.click(screen.getByLabelText("Retry message"))

        expect(onRetry).toHaveBeenCalledTimes(1)
    })

    it("marks Retry aria-disabled and does not call onRetry when retryDisabled", () => {
        const onRetry = vi.fn()
        render(<MessageActions onCopy={vi.fn()} onRetry={onRetry} retryDisabled/>)

        const retry = screen.getByLabelText("Retry message")
        expect(retry).toHaveAttribute("aria-disabled", "true")

        fireEvent.click(retry)

        expect(onRetry).not.toHaveBeenCalled()
    })
})
