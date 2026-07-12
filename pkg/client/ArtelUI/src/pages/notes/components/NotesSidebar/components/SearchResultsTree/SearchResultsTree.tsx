import { useEffect, useState } from "react"

import { NoteItem } from "@/app/hooks/Notes.ts"
import TreeItem from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.tsx"
import FolderNodeItem from "@/pages/notes/components/NotesSidebar/components/FolderNodeItem/FolderNodeItem.tsx"
import {
    buildFolderTree, getAncestorFolderPaths, getNoteName,
} from "@/pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts"

interface SearchResultsTreeProps {
    folders: string[]
    notes: NoteItem[]
    searchQuery: string
    selectedPath: string | null
    highlightedPath: string | null
    onSelectNote: (path: string) => void
    onCreateNote: (folderPath?: string) => void
}

export default function SearchResultsTree(props: SearchResultsTreeProps) {
    const { folders, notes, searchQuery, selectedPath, highlightedPath } = props
    const { onSelectNote, onCreateNote } = props
    const query = searchQuery.trim().toLowerCase()

    const matchedNotes = notes.filter(n => (n.path ?? "").toLowerCase().includes(query))
    const matchedFolders = folders.filter(f => f.toLowerCase().includes(query))

    const relevantFolders = new Set<string>(matchedFolders)
    matchedNotes.forEach(n => {
        if (!n.path) return
        getAncestorFolderPaths(n.path).forEach(a => relevantFolders.add(a))
    })

    const tree = buildFolderTree(Array.from(relevantFolders))
    const rootNotes = matchedNotes.filter(n => n.path && !n.path.includes("/"))

    const [openFolders, setOpenFolders] = useState<Set<string>>(relevantFolders)

    useEffect(() => {
        setOpenFolders(relevantFolders)
    }, [searchQuery])

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

    return (
        <>
            {tree.map(node => (
                <FolderNodeItem
                    key={node.path}
                    node={node}
                    notes={matchedNotes}
                    openFolders={openFolders}
                    selectedPath={selectedPath}
                    highlightedPath={highlightedPath}
                    depth={0}
                    onToggle={toggleFolder}
                    onSelectNote={onSelectNote}
                    onCreateNoteInFolder={onCreateNote}
                />
            ))}
            {rootNotes.map(note => (
                <TreeItem
                    key={note.path}
                    name={getNoteName(note)}
                    path={note.path}
                    active={selectedPath === note.path}
                    highlighted={!!note.path && highlightedPath === note.path}
                    onClick={() => note.path && onSelectNote(note.path)}
                />
            ))}
        </>
    )
}
