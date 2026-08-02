import { PublicDocsNoteItem } from "@/app/api/artel"
import DocsTreeItem from "@/pages/docs/components/DocsTreeItem/DocsTreeItem.tsx"
import {
    DocsFolderNode as DocsFolderNodeType, getDirectNotes, getNoteName,
} from "@/pages/docs/processes/docsTree.ts"

interface DocsFolderNodeProps {
    node: DocsFolderNodeType
    notes: PublicDocsNoteItem[]
    openFolders: Set<string>
    selectedPath: string | null
    depth: number
    onToggle: (path: string) => void
    onSelectNote: (path: string) => void
}

export default function DocsFolderNode(props: DocsFolderNodeProps) {
    const { node, notes, openFolders, selectedPath, depth } = props
    const { onToggle, onSelectNote } = props
    const isOpen = openFolders.has(node.path)
    const directNotes = getDirectNotes(node.path, notes)

    return (
        <>
            <DocsTreeItem
                name={node.name}
                path={node.path}
                isFolder
                isOpen={isOpen}
                depth={depth}
                onClick={() => onToggle(node.path)}
            />
            {isOpen && (
                <>
                    {node.children.map(child => (
                        <DocsFolderNode
                            key={child.path}
                            node={child}
                            notes={notes}
                            openFolders={openFolders}
                            selectedPath={selectedPath}
                            depth={depth + 1}
                            onToggle={onToggle}
                            onSelectNote={onSelectNote}
                        />
                    ))}
                    {directNotes.map(note => (
                        <DocsTreeItem
                            key={note.path}
                            name={getNoteName(note)}
                            path={note.path}
                            depth={depth + 1}
                            active={selectedPath === note.path}
                            onClick={() => note.path && onSelectNote(note.path)}
                        />
                    ))}
                </>
            )}
        </>
    )
}
