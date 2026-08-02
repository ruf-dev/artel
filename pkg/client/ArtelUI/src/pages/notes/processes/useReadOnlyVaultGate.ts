import {useEffect} from "react"

import {NoteMode} from "@/app/hooks/Notes.ts"
import {VaultItem} from "@/app/hooks/Vaults.ts"

// A vault is read-only here only when it's public AND the caller has no membership role
// in it at all (role is empty/undefined) — i.e. a random signed-in user browsing someone
// else's published vault. A caller with any role (owner/reader/maintainer) keeps the
// pre-existing unrestricted editability, same as private vaults always had.
//
// Also forces the store back to 'preview' if the vault selection changes out from under
// an open editor (e.g. switching from a private vault mid-edit into a public one).
export function useReadOnlyVaultGate(
    vaults: Pick<VaultItem, 'id' | 'isPublic' | 'role'>[], vaultId: string | null,
    mode: NoteMode, setMode: (mode: NoteMode) => void,
): boolean {
    const vault = vaults.find(v => v.id === vaultId)
    const isReadOnlyVault = !!vault?.isPublic && !vault?.role

    useEffect(() => {
        if (isReadOnlyVault && mode === 'edit') {
            setMode('preview')
        }
    }, [isReadOnlyVault, mode, setMode])

    return isReadOnlyVault
}
