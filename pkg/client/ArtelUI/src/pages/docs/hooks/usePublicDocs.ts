import {useState} from "react"
import {useQuery} from "@tanstack/react-query"

import {PublicDocsAPI, PublicDocsNoteItem} from "@/app/api/artel"

// Page-local hook for the public /docs/:slug reader — not a generic reusable hook, so it
// lives here instead of app/hooks/. All PublicDocsAPI calls are anonymous (no initReq /
// auth.getInitReq()): these are public endpoints.

export function usePublicDocs(slug: string) {
    const [selectedPath, setSelectedPath] = useState<string | null>(null)

    const vaultQuery = useQuery({
        queryKey: ['public-docs', 'vault', slug],
        queryFn: () => PublicDocsAPI.GetVaultBySlug({slug}),
        enabled: !!slug,
    })

    const foldersQuery = useQuery({
        queryKey: ['public-docs', 'folders', slug],
        queryFn: () => PublicDocsAPI.ListFolders({slug}),
        enabled: !!slug,
    })

    const notesQuery = useQuery({
        queryKey: ['public-docs', 'notes', slug],
        queryFn: () => PublicDocsAPI.ListNotes({slug}),
        enabled: !!slug,
    })

    const noteQuery = useQuery({
        queryKey: ['public-docs', 'note', slug, selectedPath],
        queryFn: () => PublicDocsAPI.GetNote({slug, path: selectedPath!}),
        enabled: !!slug && !!selectedPath,
    })

    const notes: PublicDocsNoteItem[] = notesQuery.data?.notes ?? []

    return {
        vaultName: vaultQuery.data?.name ?? null,
        folders: foldersQuery.data?.folders ?? [],
        notes,
        selectedPath,
        selectNote: (path: string) => setSelectedPath(path),
        noteContent: noteQuery.data?.content ?? null,
        isLoading: vaultQuery.isLoading || foldersQuery.isLoading || notesQuery.isLoading,
        error: vaultQuery.error || foldersQuery.error || notesQuery.error || noteQuery.error,
    }
}
