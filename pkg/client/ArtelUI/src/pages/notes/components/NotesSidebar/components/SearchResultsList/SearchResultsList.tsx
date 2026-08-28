import { NoteItem } from "@/app/hooks/Notes.ts"
import FileTreeRow from "@/components/FileTree/FileTreeRow.tsx"
import { getItemName } from "@/components/FileTree/fileTree.ts"

interface SearchResultsListProps {
    notes: NoteItem[]
    searchQuery: string
    selectedPath: string | null
    highlightedPath: string | null
    onSelectNote: (path: string) => void
}

// Flat search results (list mode): one FileTreeRow per match, with the containing folder shown
// as the row subtitle. No folders, no tree recursion.
export default function SearchResultsList(props: SearchResultsListProps) {
    const { notes, searchQuery, selectedPath, highlightedPath, onSelectNote } = props
    const query = searchQuery.trim().toLowerCase()
    const results = notes.filter(n => (n.path ?? "").toLowerCase().includes(query))

    return (
        <>
            {results.map(note => {
                const dir = note.path ? note.path.split("/").slice(0, -1).join("/") : undefined
                return (
                    <FileTreeRow
                        key={note.path}
                        name={getItemName(note)}
                        path={note.path}
                        subtitle={dir || undefined}
                        active={selectedPath === note.path}
                        highlighted={!!note.path && highlightedPath === note.path}
                        onClick={() => note.path && onSelectNote(note.path)}
                    />
                )
            })}
        </>
    )
}
