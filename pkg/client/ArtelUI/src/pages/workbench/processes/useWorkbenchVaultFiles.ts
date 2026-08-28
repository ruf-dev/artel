import {useQuery} from "@tanstack/react-query"

import {notesService} from "@/processes/Notes.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import {NoteItem} from "@/app/api/artel/notes.pb.ts"
import useUser from "@/hooks/user/User.ts"

// Page-local read-only hook for the workbench sidebar's Vault pane — lists the
// current vault's folders + notes so a file can be attached to the chat as a
// context chip. Mirrors pages/docs/hooks/usePublicDocs.ts (two react-query reads,
// derived flat return), but authenticated: notesService injects the init req.

export interface WorkbenchVaultFiles {
    folders: string[]
    notes: NoteItem[]
    isLoading: boolean
    error: unknown
}

export function useWorkbenchVaultFiles(vaultId: string | undefined): WorkbenchVaultFiles {
    const {auth} = useUser()

    const foldersQuery = useQuery({
        queryKey: ['workbench-vault-files', 'folders', vaultId],
        queryFn: () => notesService.listFolders(vaultId!),
        enabled: !!vaultId && auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    const notesQuery = useQuery({
        queryKey: ['workbench-vault-files', 'notes', vaultId],
        queryFn: () => notesService.listNotes(vaultId!),
        enabled: !!vaultId && auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    return {
        folders: foldersQuery.data ?? [],
        notes: notesQuery.data ?? [],
        isLoading: foldersQuery.isLoading || notesQuery.isLoading,
        error: foldersQuery.error || notesQuery.error,
    }
}
