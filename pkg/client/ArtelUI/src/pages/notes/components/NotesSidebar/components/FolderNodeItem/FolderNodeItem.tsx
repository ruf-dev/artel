import { NoteItem } from "@/app/hooks/Notes.ts"
import TreeItem from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.tsx"
import {
    FolderNode, getDirectNotes, getNoteName,
} from "@/pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts"

interface FolderNodeItemProps {
    node: FolderNode
    notes: NoteItem[]
    openFolders: Set<string>
    selectedPath: string | null
    highlightedPath: string | null
    depth: number
    onToggle: (path: string) => void
    onSelectNote: (path: string) => void
    onCreateNoteInFolder: (path: string) => void
}

export default function FolderNodeItem(props: FolderNodeItemProps) {
    const { node, notes, openFolders, selectedPath, highlightedPath } = props
    const { depth, onToggle, onSelectNote, onCreateNoteInFolder } = props
    const isOpen = openFolders.has(node.path)
    const directNotes = getDirectNotes(node.path, notes)

    return (
        <>
            <TreeItem
                name={node.name}
                path={node.path}
                isFolder
                isOpen={isOpen}
                depth={depth}
                onClick={() => onToggle(node.path)}
                onAddInFolder={() => onCreateNoteInFolder(node.path)}
                onCopyPath={() => navigator.clipboard.writeText(node.path)}
            />
            {isOpen && (
                <>
                    {node.children.map(child => (
                        <FolderNodeItem
                            key={child.path}
                            node={child}
                            notes={notes}
                            openFolders={openFolders}
                            selectedPath={selectedPath}
                            highlightedPath={highlightedPath}
                            depth={depth + 1}
                            onToggle={onToggle}
                            onSelectNote={onSelectNote}
                            onCreateNoteInFolder={onCreateNoteInFolder}
                        />
                    ))}
                    {directNotes.map(note => (
                        <TreeItem
                            key={note.path}
                            name={getNoteName(note)}
                            path={note.path}
                            depth={depth + 1}
                            active={selectedPath === note.path}
                            highlighted={!!note.path && highlightedPath === note.path}
                            onClick={() => note.path && onSelectNote(note.path)}
                        />
                    ))}
                </>
            )}
        </>
    )
}
