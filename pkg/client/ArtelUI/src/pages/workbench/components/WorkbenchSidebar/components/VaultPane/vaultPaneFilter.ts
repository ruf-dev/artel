import {getAncestorFolderPaths} from "@/components/FileTree/fileTree.ts"
import {NoteItem} from "@/app/api/artel/notes.pb.ts"

// Pure derivation for the Vault pane's client-side search. With an empty query the
// vault's folders/notes pass through untouched; with a query we keep the notes
// whose path contains it (case-insensitive) and rebuild the folder list as the
// union of every matched note's ancestor folders, so the FileTree only shows the
// branches that still hold a match.
export interface FilteredVaultTree {
    folders: string[]
    notes: NoteItem[]
}

export function filterVaultTree(query: string, folders: string[], notes: NoteItem[]): FilteredVaultTree {
    const needle = query.trim().toLowerCase()
    if (!needle) return {folders, notes}

    const matched = notes.filter(note => (note.path ?? "").toLowerCase().includes(needle))

    const folderSet = new Set<string>()
    matched.forEach(note => {
        if (note.path) getAncestorFolderPaths(note.path).forEach(folder => folderSet.add(folder))
    })

    return {folders: [...folderSet], notes: matched}
}
