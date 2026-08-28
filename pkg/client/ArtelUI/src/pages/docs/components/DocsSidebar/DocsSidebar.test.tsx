import { describe, expect, it, vi } from "vitest"
import { fireEvent, render, screen } from "@testing-library/react"

import DocsSidebar from "@/pages/docs/components/DocsSidebar/DocsSidebar.tsx"
import rowCls from "@/components/FileTree/FileTreeRow.module.css"

const folders = ["Guides"]
const notes = [{ path: "Guides/intro.md" }, { path: "readme.md" }]

function renderSidebar(overrides: Partial<Parameters<typeof DocsSidebar>[0]> = {}) {
    const onSelectNote = vi.fn()
    render(
        <DocsSidebar
            vaultName="My Vault"
            folders={folders}
            notes={notes}
            selectedPath={null}
            onSelectNote={onSelectNote}
            {...overrides}
        />,
    )
    return { onSelectNote }
}

describe("DocsSidebar", () => {
    it("shows the vault name and keeps folders collapsed by default", () => {
        renderSidebar()

        expect(screen.getByText("My Vault")).toBeInTheDocument()
        expect(screen.getByText("Guides")).toBeInTheDocument()
        expect(screen.getByText("readme")).toBeInTheDocument()
        expect(screen.queryByText("intro")).not.toBeInTheDocument()
    })

    it("falls back to 'Docs' when no vault name is given", () => {
        renderSidebar({ vaultName: null })

        expect(screen.getByText("Docs")).toBeInTheDocument()
    })

    it("expands and collapses a folder when its row is clicked", () => {
        renderSidebar()

        fireEvent.click(screen.getByText("Guides"))
        expect(screen.getByText("intro")).toBeInTheDocument()

        fireEvent.click(screen.getByText("Guides"))
        expect(screen.queryByText("intro")).not.toBeInTheDocument()
    })

    it("calls onSelectNote with the full path when a note row is clicked", () => {
        const { onSelectNote } = renderSidebar()

        fireEvent.click(screen.getByText("Guides"))
        fireEvent.click(screen.getByText("intro"))

        expect(onSelectNote).toHaveBeenCalledWith("Guides/intro.md")
    })

    it("marks the selected note row active", () => {
        renderSidebar({ selectedPath: "readme.md" })

        expect(screen.getByText("readme").closest("[data-path]")).toHaveClass(rowCls.FileTreeRowActive)
    })
})
