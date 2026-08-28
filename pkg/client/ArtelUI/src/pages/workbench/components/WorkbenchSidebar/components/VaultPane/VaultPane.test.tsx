import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import VaultPane from "@/pages/workbench/components/WorkbenchSidebar/components/VaultPane/VaultPane.tsx"
import {VaultPaneBundle} from "@/pages/workbench/processes/useWorkbenchSidebar.ts"
import rowCls from "@/components/FileTree/FileTreeRow.module.css"

vi.mock("morphicons/react", () => ({MorphIcon: () => null}))

// The real FileTree is used (not a stub) — it's stable and gives the row
// active-state / onSelect wiring real coverage. Root-level notes (no "/" in the
// path) render immediately without expanding a folder.
function makeProps(over: Partial<VaultPaneBundle> = {}): VaultPaneBundle {
    return {
        vaultName: "My Vault",
        folders: [],
        notes: [{path: "alpha.md"}, {path: "beta.md"}],
        isLoading: false,
        attachedPaths: [],
        onToggleAttach: vi.fn(),
        ...over,
    }
}

describe("VaultPane", () => {
    it("renders the vault name and the note count in the header", () => {
        render(<VaultPane {...makeProps()}/>)

        expect(screen.getByText("My Vault")).toBeInTheDocument()
        expect(screen.getByText("2")).toBeInTheDocument()
    })

    it("filters the visible rows as you type in the search box", () => {
        render(<VaultPane {...makeProps()}/>)

        expect(screen.getByText("alpha")).toBeInTheDocument()
        expect(screen.getByText("beta")).toBeInTheDocument()

        fireEvent.change(screen.getByPlaceholderText("Search files…"), {target: {value: "alpha"}})

        expect(screen.getByText("alpha")).toBeInTheDocument()
        expect(screen.queryByText("beta")).not.toBeInTheDocument()
    })

    it("calls onToggleAttach with the note's full path when its row is clicked", () => {
        const onToggleAttach = vi.fn()
        render(<VaultPane {...makeProps({onToggleAttach})}/>)

        fireEvent.click(screen.getByText("alpha"))

        expect(onToggleAttach).toHaveBeenCalledWith("alpha.md")
    })

    it("marks a row active when its path is in attachedPaths", () => {
        render(<VaultPane {...makeProps({attachedPaths: ["alpha.md"]})}/>)

        expect(screen.getByText("alpha").closest("[data-path]")).toHaveClass(rowCls.FileTreeRowActive)
        expect(screen.getByText("beta").closest("[data-path]")).not.toHaveClass(rowCls.FileTreeRowActive)
    })

    it("shows the loading state and no tree while isLoading", () => {
        render(<VaultPane {...makeProps({isLoading: true})}/>)

        expect(screen.getByText("Loading vault…")).toBeInTheDocument()
        expect(screen.queryByText("alpha")).not.toBeInTheDocument()
    })

    it("shows the empty state when the vault has no notes", () => {
        render(<VaultPane {...makeProps({notes: []})}/>)

        expect(screen.getByText("This vault has no notes.")).toBeInTheDocument()
    })

    it("reports how many files are attached in the footer", () => {
        const {rerender} = render(<VaultPane {...makeProps()}/>)
        expect(screen.getByText("Nothing attached")).toBeInTheDocument()

        rerender(<VaultPane {...makeProps({attachedPaths: ["alpha.md"]})}/>)
        expect(screen.getByText("1 attached")).toBeInTheDocument()
    })
})
