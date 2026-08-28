import { describe, expect, it, vi } from "vitest"
import { fireEvent, render, screen } from "@testing-library/react"

import FolderSection from "@/pages/notes/components/NotesSidebar/components/FolderSection/FolderSection.tsx"
import rowCls from "@/components/FileTree/FileTreeRow.module.css"

const folders = ["Projects", "Projects/Deep"]
const notes = [
    { path: "Projects/spec.md" },
    { path: "Projects/Deep/plan.md" },
    { path: "root.md" },
]

function renderSection(overrides: Partial<Parameters<typeof FolderSection>[0]> = {}) {
    const props = {
        folders,
        notes,
        selectedPath: null,
        highlightedPath: null,
        revealPath: null,
        vaultId: "v1",
        onSelectNote: vi.fn(),
        onCreateNote: vi.fn(),
        onDownloadFolder: vi.fn(),
        onDeleteFolder: vi.fn(),
        onUpload: vi.fn(),
        onItemDragStart: vi.fn(() => vi.fn()),
        onFolderDrop: vi.fn(() => vi.fn()),
        ...overrides,
    }
    render(<FolderSection {...props}/>)
    return props
}

describe("FolderSection", () => {
    it("keeps folders collapsed and shows root notes by default", () => {
        renderSection()

        expect(screen.getByText("Projects")).toBeInTheDocument()
        expect(screen.getByText("root")).toBeInTheDocument()
        expect(screen.queryByText("spec")).not.toBeInTheDocument()
    })

    it("toggles a folder open and closed on click", () => {
        renderSection()

        fireEvent.click(screen.getByText("Projects"))
        expect(screen.getByText("spec")).toBeInTheDocument()

        fireEvent.click(screen.getByText("Projects"))
        expect(screen.queryByText("spec")).not.toBeInTheDocument()
    })

    it("auto-expands every ancestor folder of revealPath", () => {
        renderSection({ revealPath: "Projects/Deep/plan.md" })

        expect(screen.getByText("Deep")).toBeInTheDocument()
        expect(screen.getByText("plan")).toBeInTheDocument()
    })

    it("selects a note with the vault id and path", () => {
        const props = renderSection()

        fireEvent.click(screen.getByText("root"))

        expect(props.onSelectNote).toHaveBeenCalledWith("v1", "root.md")
    })

    it("marks the selected note active", () => {
        renderSection({ selectedPath: "root.md" })

        expect(screen.getByText("root").closest("[data-path]")).toHaveClass(rowCls.FileTreeRowActive)
    })

    it("wires drag-start through onItemDragStart for notes and folders", () => {
        const props = renderSection()

        const rootRow = screen.getByText("root").closest("[data-path]") as HTMLElement
        expect(rootRow).toHaveAttribute("draggable", "true")
        fireEvent.dragStart(rootRow)
        expect(props.onItemDragStart).toHaveBeenCalledWith("root.md", false)

        fireEvent.dragStart(screen.getByText("Projects").closest("[data-path]") as HTMLElement)
        expect(props.onItemDragStart).toHaveBeenCalledWith("Projects", true)
    })

    it("routes a drop on a folder row to onFolderDrop for that folder", () => {
        const drop = vi.fn()
        const onFolderDrop = vi.fn(() => drop)
        renderSection({ onFolderDrop })

        const folderRow = screen.getByText("Projects").closest("[data-path]") as HTMLElement
        fireEvent.dragOver(folderRow)
        fireEvent.drop(folderRow)

        expect(onFolderDrop).toHaveBeenCalledWith("Projects")
        expect(drop).toHaveBeenCalled()
    })

    it("renders the folder action cluster (add / copy / kebab) on a folder row", () => {
        renderSection()

        expect(screen.getByTitle("New note here")).toBeInTheDocument()
        expect(screen.getByTitle("Copy path")).toBeInTheDocument()
        expect(screen.getByTitle("Folder actions")).toBeInTheDocument()
    })
})
