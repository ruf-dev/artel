import { beforeEach, describe, expect, it, vi } from "vitest"
import { fireEvent, render, screen } from "@testing-library/react"

import NotesFolderActions from "@/pages/notes/components/NotesSidebar/components/NotesFolderActions/NotesFolderActions.tsx" // eslint-disable-line max-len

const writeText = vi.fn(() => Promise.resolve())

beforeEach(() => {
    writeText.mockClear()
    Object.assign(navigator, { clipboard: { writeText } })
})

describe("NotesFolderActions", () => {
    it("copies the folder path and flips the button title to 'Copied!'", () => {
        render(<NotesFolderActions folderPath="Projects/2024" onCreateNoteInFolder={vi.fn()}/>)

        fireEvent.click(screen.getByTitle("Copy path"))

        expect(writeText).toHaveBeenCalledWith("Projects/2024")
        expect(screen.getByTitle("Copied!")).toBeInTheDocument()
    })

    it("calls onCreateNoteInFolder with the folder path from the add button", () => {
        const onCreateNoteInFolder = vi.fn()
        render(<NotesFolderActions folderPath="Projects" onCreateNoteInFolder={onCreateNoteInFolder}/>)

        fireEvent.click(screen.getByTitle("New note here"))

        expect(onCreateNoteInFolder).toHaveBeenCalledWith("Projects")
    })

    it("omits the kebab menu when neither download nor delete is supplied", () => {
        render(<NotesFolderActions folderPath="Projects" onCreateNoteInFolder={vi.fn()}/>)

        expect(screen.queryByTitle("Folder actions")).not.toBeInTheDocument()
    })

    it("exposes download and delete actions through the kebab menu", () => {
        const onDownloadFolder = vi.fn()
        const onDeleteFolder = vi.fn()
        render(
            <NotesFolderActions
                folderPath="Projects"
                onCreateNoteInFolder={vi.fn()}
                onDownloadFolder={onDownloadFolder}
                onDeleteFolder={onDeleteFolder}
            />,
        )

        fireEvent.click(screen.getByTitle("Folder actions"))
        fireEvent.click(screen.getByText("Download as .zip"))
        expect(onDownloadFolder).toHaveBeenCalledWith("Projects")

        fireEvent.click(screen.getByTitle("Folder actions"))
        fireEvent.click(screen.getByText("Delete folder"))
        expect(onDeleteFolder).toHaveBeenCalledWith("Projects")
    })
})
