import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import UserMessageBubble from "@/pages/workbench/components/Chat/components/UserMessageBubble/UserMessageBubble.tsx"

const OpenDialog = vi.fn()

vi.mock("@/app/hooks/Dialog", () => ({
    useDialog: () => ({OpenDialog}),
}))

vi.mock("@/pages/workbench/components/AttachedNoteDialog/AttachedNoteDialog.tsx", () => ({
    default: () => <div className="test-mock" data-testid="attached-note-dialog"/>,
}))

describe("UserMessageBubble", () => {
    it("renders text without attachments", () => {
        render(<UserMessageBubble text="Hello world"/>)

        expect(screen.getByText("Hello world")).toBeInTheDocument()
    })

    it("renders stripped caption when attachments are present", () => {
        const text = "[Attached vault files: notes/file.md]\n\nHello world"
        render(
            <UserMessageBubble
                text={text}
                attachments={["notes/file.md"]}
            />,
        )

        expect(screen.queryByText("Hello world")).toBeInTheDocument()
        expect(screen.queryByText(/Attached vault files/)).not.toBeInTheDocument()
    })

    it("renders one chip per attachment", () => {
        render(
            <UserMessageBubble
                text="[Attached vault files: notes/a.md, notes/b.md]\n\nContent"
                attachments={["notes/a.md", "notes/b.md"]}
                vaultId="vault-1"
            />,
        )

        expect(screen.getByText("a")).toBeInTheDocument()
        expect(screen.getByText("b")).toBeInTheDocument()
    })

    it("opens AttachedNoteDialog when a chip is clicked", () => {
        OpenDialog.mockReset()

        render(
            <UserMessageBubble
                text="[Attached vault files: notes/file.md]\n\nContent"
                attachments={["notes/file.md"]}
                vaultId="vault-1"
            />,
        )

        fireEvent.click(screen.getByRole("button", {name: /View attached file/}))

        expect(OpenDialog).toHaveBeenCalled()
    })

    it("does not open dialog when vaultId is missing", () => {
        OpenDialog.mockReset()

        render(
            <UserMessageBubble
                text="[Attached vault files: notes/file.md]\n\nContent"
                attachments={["notes/file.md"]}
                vaultId={undefined}
            />,
        )

        fireEvent.click(screen.getByRole("button", {name: /View attached file/}))

        expect(OpenDialog).not.toHaveBeenCalled()
    })
})
