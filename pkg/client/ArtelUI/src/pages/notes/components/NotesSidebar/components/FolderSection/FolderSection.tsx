import { useState } from "react"
import {Button} from "@vervstack/chures"

import { NoteItem } from "@/app/hooks/Notes.ts"
import PlusIcon from "@/pages/notes/components/icons/PlusIcon.tsx"
import TreeItem from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.tsx"
import FolderNodeItem from "@/pages/notes/components/NotesSidebar/components/FolderNodeItem/FolderNodeItem.tsx"
import { buildFolderTree, getNoteName } from "@/pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts"
import cls from "@/pages/notes/components/NotesSidebar/components/FolderSection/FolderSection.module.css"

interface FolderSectionProps {
    folders: string[]
    notes: NoteItem[]
    selectedPath: string | null
    vaultId: string
    onSelectNote: (vaultId: string, path: string) => void
    onCreateNote: (folderPath?: string) => void
    showCreateButton?: boolean
}

export default function FolderSection(props: FolderSectionProps) {
    const { folders, notes, selectedPath, vaultId } = props
    const { onSelectNote, onCreateNote, showCreateButton = true } = props
    const [openFolders, setOpenFolders] = useState<Set<string>>(new Set())

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

    const tree = buildFolderTree(folders)
    const rootNotes = notes.filter(n => n.path && !n.path.includes("/"))

    return (
        <>
            <div className={cls.SectionHeaderContainer}>
                <span className={cls.SectionLabel}>All Notes</span>
                {showCreateButton && (
                    <Button
                        variant="ghost"
                        className={cls.CreateNoteBtn}
                        onClick={() => onCreateNote()}
                        title="New note"
                    >
                        <PlusIcon/>
                    </Button>
                )}
            </div>
            {tree.map(node => (
                <FolderNodeItem
                    key={node.path}
                    node={node}
                    notes={notes}
                    openFolders={openFolders}
                    selectedPath={selectedPath}
                    depth={0}
                    onToggle={toggleFolder}
                    onSelectNote={path => onSelectNote(vaultId, path)}
                    onCreateNoteInFolder={onCreateNote}
                />
            ))}
            {rootNotes.map(note => (
                <TreeItem
                    key={note.path}
                    name={getNoteName(note)}
                    active={selectedPath === note.path}
                    onClick={() => note.path && onSelectNote(vaultId, note.path)}
                />
            ))}
        </>
    )
}
