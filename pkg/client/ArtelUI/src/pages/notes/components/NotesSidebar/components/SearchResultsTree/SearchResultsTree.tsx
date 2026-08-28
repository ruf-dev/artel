import { useEffect, useState } from "react"

import { NoteItem } from "@/app/hooks/Notes.ts"
import FileTreeLeafRow from "@/components/FileTree/FileTreeLeafRow.tsx"
import FileTreeNode from "@/components/FileTree/FileTreeNode.tsx"
import { buildFolderTree, getAncestorFolderPaths, sortItemsByName } from "@/components/FileTree/fileTree.ts"
import NotesFolderActions from "@/pages/notes/components/NotesSidebar/components/NotesFolderActions/NotesFolderActions.tsx" // eslint-disable-line max-len

interface SearchResultsTreeProps {
    folders: string[]
    notes: NoteItem[]
    searchQuery: string
    selectedPath: string | null
    highlightedPath: string | null
    onSelectNote: (path: string) => void
    onCreateNote: (folderPath?: string) => void
}

// Filtered variant of FolderSection's tree: same shared FileTree building blocks, but folders
// carry only the add-note action (no download/delete, no drag-and-drop — moving into a
// filtered subset is ambiguous).
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
    const rootNotes = sortItemsByName(matchedNotes.filter(n => n.path && !n.path.includes("/")))

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

    function renderFolderTrailing(folderPath: string) {
        return <NotesFolderActions folderPath={folderPath} onCreateNoteInFolder={onCreateNote}/>
    }

    return (
        <>
            {tree.map(node => (
                <FileTreeNode
                    key={node.path}
                    node={node}
                    items={matchedNotes}
                    openFolders={openFolders}
                    depth={0}
                    isActive={p => p === selectedPath}
                    isHighlighted={p => p === highlightedPath}
                    onToggleFolder={toggleFolder}
                    onSelectItem={onSelectNote}
                    renderFolderTrailing={renderFolderTrailing}
                />
            ))}
            {rootNotes.map(note => (
                <FileTreeLeafRow
                    key={note.path}
                    item={note}
                    depth={0}
                    isActive={p => p === selectedPath}
                    isHighlighted={p => p === highlightedPath}
                    onSelectItem={onSelectNote}
                />
            ))}
        </>
    )
}
