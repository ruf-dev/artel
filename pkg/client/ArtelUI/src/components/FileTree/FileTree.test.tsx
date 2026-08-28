import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import FileTree from "@/components/FileTree/FileTree.tsx"
import rowCls from "@/components/FileTree/FileTreeRow.module.css"

const folders = ["Projects"]
const items = [{path: "Projects/spec.md"}, {path: "root.md"}]

describe("FileTree", () => {
    it("renders root-level items and keeps folders collapsed by default", () => {
        render(
            <FileTree
                folders={folders}
                items={items}
                isActive={() => false}
                onSelectItem={vi.fn()}
            />,
        )

        expect(screen.getByText("root")).toBeInTheDocument()
        expect(screen.getByText("Projects")).toBeInTheDocument()
        expect(screen.queryByText("spec")).not.toBeInTheDocument()
    })

    it("expands a folder when its row is clicked", () => {
        render(<FileTree folders={folders} items={items} onSelectItem={vi.fn()}/>)

        fireEvent.click(screen.getByText("Projects"))

        expect(screen.getByText("spec")).toBeInTheDocument()
    })

    it("calls onSelectItem with the full path when a file row is clicked", () => {
        const onSelectItem = vi.fn()
        render(<FileTree folders={folders} items={items} onSelectItem={onSelectItem}/>)

        fireEvent.click(screen.getByText("Projects"))
        fireEvent.click(screen.getByText("spec"))

        expect(onSelectItem).toHaveBeenCalledWith("Projects/spec.md")
    })

    it("applies the active class to a row the isActive predicate matches", () => {
        render(
            <FileTree
                folders={folders}
                items={items}
                isActive={path => path === "root.md"}
                onSelectItem={vi.fn()}
            />,
        )

        expect(screen.getByText("root").closest("[data-path]"))
            .toHaveClass(rowCls.FileTreeRowActive)
    })

    it("renders renderFolderTrailing output next to a folder row", () => {
        render(
            <FileTree
                folders={folders}
                items={items}
                onSelectItem={vi.fn()}
                renderFolderTrailing={folderPath => <span data-testid="trail">{folderPath}</span>}
            />,
        )

        expect(screen.getByTestId("trail")).toHaveTextContent("Projects")
    })

    it("marks a row draggable when itemDragProps returns draggable", () => {
        render(
            <FileTree
                folders={folders}
                items={items}
                onSelectItem={vi.fn()}
                itemDragProps={() => ({draggable: true})}
            />,
        )

        expect(screen.getByText("root").closest("[data-path]"))
            .toHaveAttribute("draggable", "true")
    })
})
