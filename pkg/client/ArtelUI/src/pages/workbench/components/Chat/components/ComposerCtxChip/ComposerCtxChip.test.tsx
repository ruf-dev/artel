import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import ComposerCtxChip from "@/pages/workbench/components/Chat/components/ComposerCtxChip/ComposerCtxChip.tsx"

describe("ComposerCtxChip", () => {
    it("renders the file basename as the label", () => {
        render(<ComposerCtxChip path="notes/sub/report.md" onRemove={vi.fn()}/>)

        expect(screen.getByText("report")).toBeInTheDocument()
    })

    it("exposes the full vault-relative path as the title", () => {
        const {container} = render(<ComposerCtxChip path="notes/sub/report.md" onRemove={vi.fn()}/>)

        expect(container.querySelector("[title='notes/sub/report.md']")).toBeInTheDocument()
    })

    it("calls onRemove with the path when the × button is clicked", () => {
        const onRemove = vi.fn()
        render(<ComposerCtxChip path="notes/sub/report.md" onRemove={onRemove}/>)

        fireEvent.click(screen.getByRole("button", {name: "Remove notes/sub/report.md"}))

        expect(onRemove).toHaveBeenCalledWith("notes/sub/report.md")
    })
})
