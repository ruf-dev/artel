import { DragEvent } from "react"
import { Button } from "@vervstack/chures"

import { cn } from "@/app/utils/cn.ts"
import { NoteItem } from "@/app/hooks/Notes.ts"
import FileTreeLeafRow from "@/components/FileTree/FileTreeLeafRow.tsx"
import FileTreeNode from "@/components/FileTree/FileTreeNode.tsx"
import { buildFolderTree, sortItemsByName } from "@/components/FileTree/fileTree.ts"
import PlusIcon from "@/pages/notes/components/icons/PlusIcon.tsx"
import UploadIcon from "@/pages/notes/components/icons/UploadIcon.tsx"
import NotesFolderActions from "@/pages/notes/components/NotesSidebar/components/NotesFolderActions/NotesFolderActions.tsx" // eslint-disable-line max-len
import { useFolderOpenState } from "@/pages/notes/components/NotesSidebar/processes/useFolderOpenState.ts"
import { useRootDropTarget } from "@/pages/notes/components/NotesSidebar/processes/useRootDropTarget.ts"
import cls from "@/pages/notes/components/NotesSidebar/components/FolderSection/FolderSection.module.css"

interface FolderSectionProps {
    folders: string[]
    notes: NoteItem[]
    selectedPath: string | null
    highlightedPath: string | null
    revealPath: string | null
    vaultId: string
    onSelectNote: (vaultId: string, path: string) => void
    onCreateNote: (folderPath?: string) => void
    onDownloadFolder?: (path: string) => void
    onDeleteFolder?: (path: string) => void
    onUpload?: () => void
    showCreateButton?: boolean
    onItemDragStart: (path: string, isFolder: boolean) => (e: DragEvent<HTMLDivElement>) => void
    onFolderDrop: (targetFolderPath: string) => (e: DragEvent<HTMLDivElement>) => void
}

// Owns the /notes tree's expand/collapse + reveal-in-tree state (via useFolderOpenState);
// per-row rendering is delegated to the shared FileTree building blocks (FileTreeNode for
// folders, FileTreeLeafRow for notes). Drag-and-drop and the folder action cluster are wired
// in via `itemDragProps` / `renderFolderTrailing`.
export default function FolderSection(props: FolderSectionProps) {
    const { folders, notes, selectedPath, highlightedPath, revealPath, vaultId } = props
    const { onSelectNote, onCreateNote, onDownloadFolder, onDeleteFolder, onUpload } = props
    const { showCreateButton = true, onItemDragStart, onFolderDrop } = props
    const { openFolders, toggleFolder } = useFolderOpenState(revealPath)
    const rootDrop = useRootDropTarget(onFolderDrop)

    function itemDragProps(path: string, isFolder: boolean) {
        return {
            draggable: true,
            onDragStart: onItemDragStart(path, isFolder),
            onDropTarget: isFolder ? onFolderDrop(path) : undefined,
        }
    }

    function renderFolderTrailing(folderPath: string) {
        return (
            <NotesFolderActions
                folderPath={folderPath}
                onCreateNoteInFolder={onCreateNote}
                onDownloadFolder={onDownloadFolder}
                onDeleteFolder={onDeleteFolder}
            />
        )
    }

    const tree = buildFolderTree(folders)
    const rootNotes = sortItemsByName(notes.filter(n => n.path && !n.path.includes("/")))

    return (
        <>
            <div
                className={cn(cls.SectionHeaderContainer, rootDrop.isDragOver && cls.SectionHeaderDropTarget)}
                onDragOver={rootDrop.onDragOver}
                onDragLeave={rootDrop.onDragLeave}
                onDrop={rootDrop.onDrop}
            >
                <span className={cls.SectionLabel}>All Notes</span>
                <div className={cls.SectionActions}>
                    {onUpload && (
                        <Button
                            variant="ghost"
                            className={cls.CreateNoteBtn}
                            onClick={onUpload}
                            data-tooltip-id="root-tooltip"
                            data-tooltip-content="Upload folder from .zip"
                        >
                            <UploadIcon/>
                        </Button>
                    )}
                    {showCreateButton && (
                        <Button variant="ghost" className={cls.CreateNoteBtn} onClick={() => onCreateNote()}>
                            <PlusIcon/>
                        </Button>
                    )}
                </div>
            </div>
            {tree.map(node => (
                <FileTreeNode
                    key={node.path}
                    node={node}
                    items={notes}
                    openFolders={openFolders}
                    depth={0}
                    isActive={p => p === selectedPath}
                    isHighlighted={p => p === highlightedPath}
                    onToggleFolder={toggleFolder}
                    onSelectItem={path => onSelectNote(vaultId, path)}
                    renderFolderTrailing={renderFolderTrailing}
                    itemDragProps={itemDragProps}
                />
            ))}
            {rootNotes.map(note => (
                <FileTreeLeafRow
                    key={note.path}
                    item={note}
                    depth={0}
                    isActive={p => p === selectedPath}
                    isHighlighted={p => p === highlightedPath}
                    onSelectItem={path => onSelectNote(vaultId, path)}
                    itemDragProps={itemDragProps}
                />
            ))}
        </>
    )
}
