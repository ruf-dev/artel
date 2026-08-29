import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import AttachmentChip from "@/pages/workbench/components/Chat/components/AttachmentChip/AttachmentChip.tsx"

describe("AttachmentChip", () => {
    it("renders the file basename as the label", () => {
        render(<AttachmentChip path="notes/sub/report.md" onClick={vi.fn()}/>)

        expect(screen.getByText("report")).toBeInTheDocument()
    })

    it("exposes the full vault-relative path as the title", () => {
        const {container} = render(<AttachmentChip path="notes/sub/report.md" onClick={vi.fn()}/>)

        expect(container.querySelector("[title='notes/sub/report.md']")).toBeInTheDocument()
    })

    it("calls onClick with the path when clicked", () => {
        const onClick = vi.fn()
        render(<AttachmentChip path="notes/sub/report.md" onClick={onClick}/>)

        fireEvent.click(screen.getByRole("button"))

        expect(onClick).toHaveBeenCalledWith("notes/sub/report.md")
    })
})
