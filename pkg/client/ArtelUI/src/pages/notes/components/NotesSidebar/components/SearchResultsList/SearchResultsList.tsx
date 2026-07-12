import { NoteItem } from "@/app/hooks/Notes.ts"
import TreeItem from "@/pages/notes/components/NotesSidebar/components/TreeItem/TreeItem.tsx"
import { getNoteName } from "@/pages/notes/components/NotesSidebar/processes/notesTreeHelpers.ts"

interface SearchResultsListProps {
    notes: NoteItem[]
    searchQuery: string
    selectedPath: string | null
    highlightedPath: string | null
    onSelectNote: (path: string) => void
}

export default function SearchResultsList(props: SearchResultsListProps) {
    const { notes, searchQuery, selectedPath, highlightedPath, onSelectNote } = props
    const query = searchQuery.trim().toLowerCase()
    const results = notes.filter(n => (n.path ?? "").toLowerCase().includes(query))

    return (
        <>
            {results.map(note => {
                const dir = note.path ? note.path.split("/").slice(0, -1).join("/") : undefined
                return (
                    <TreeItem
                        key={note.path}
                        name={getNoteName(note)}
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
