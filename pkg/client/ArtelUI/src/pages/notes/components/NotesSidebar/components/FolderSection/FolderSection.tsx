import { DragEvent, useEffect, useState } from "react"
import {Button} from "@vervstack/chures"

import { cn } from "@/app/utils/cn.ts"
import { NoteItem } from "@/app/hooks/Notes.ts"
import PlusIcon from "@/pages/notes/components/icons/PlusIcon.tsx"
import UploadIcon from "@/pages/notes/components/icons/UploadIcon.tsx"
import TreeItem from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.tsx"
import FolderNodeItem from "@/pages/notes/components/NotesSidebar/components/FolderNodeItem/FolderNodeItem.tsx"
import {
    buildFolderTree, getAncestorFolderPaths, getNoteName, sortNotesByName,
} from "@/pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts"
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

export default function FolderSection(props: FolderSectionProps) {
    const { folders, notes, selectedPath, highlightedPath, revealPath, vaultId } = props
    const { onSelectNote, onCreateNote, onDownloadFolder, onDeleteFolder, onUpload } = props
    const { showCreateButton = true, onItemDragStart, onFolderDrop } = props
    const [openFolders, setOpenFolders] = useState<Set<string>>(new Set())
    const rootDrop = useRootDropTarget(onFolderDrop)

    function toggleFolder(path: string) {
        setOpenFolders(prev => {
            const next = new Set(prev)
            if (next.has(path)) {
                next.delete(path)
            } else {
                next.add(path)
            }
            return next
        })
    }

    useEffect(() => {
        if (!revealPath) return
        const ancestors = getAncestorFolderPaths(revealPath)
        if (ancestors.length === 0) return
        setOpenFolders(prev => {
            const next = new Set(prev)
            ancestors.forEach(a => next.add(a))
            return next
        })
    }, [revealPath])

    const tree = buildFolderTree(folders)
    const rootNotes = sortNotesByName(notes.filter(n => n.path && !n.path.includes("/")))

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
                        <Button
                            variant="ghost"
                            className={cls.CreateNoteBtn}
                            onClick={() => onCreateNote()}
                        >
                            <PlusIcon/>
                        </Button>
                    )}
                </div>
            </div>
            {tree.map(node => (
                <FolderNodeItem
                    key={node.path}
                    node={node}
                    notes={notes}
                    openFolders={openFolders}
                    selectedPath={selectedPath}
                    highlightedPath={highlightedPath}
                    depth={0}
                    onToggle={toggleFolder}
                    onSelectNote={path => onSelectNote(vaultId, path)}
                    onCreateNoteInFolder={onCreateNote}
                    onDownloadFolder={onDownloadFolder}
                    onDeleteFolder={onDeleteFolder}
                    onItemDragStart={onItemDragStart}
                    onFolderDrop={onFolderDrop}
                />
            ))}
            {rootNotes.map(note => (
                <TreeItem
                    key={note.path}
                    name={getNoteName(note)}
                    path={note.path}
                    active={selectedPath === note.path}
                    highlighted={!!note.path && highlightedPath === note.path}
                    onClick={() => note.path && onSelectNote(vaultId, note.path)}
                    onDragStart={note.path ? onItemDragStart(note.path, false) : undefined}
                />
            ))}
        </>
    )
}
