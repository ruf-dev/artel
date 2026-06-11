import {useEffect, useMemo} from "react"

import cls from "./NotesPage.module.css"
import {useNotes} from "@/app/hooks/Notes.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import NotesSidebar from "@/pages/notes/components/NotesSidebar/NotesSidebar.tsx"
import NoteViewer from "@/pages/notes/components/NoteViewer/NoteViewer.tsx"

export default function NotesPage() {
    const {vaultId, noteContent, selectVault} = useNotes()
    const {vaults} = useVaults()

    const vaultOptions = useMemo(
        () => vaults
            .filter(v => v.id && v.name)
            .map(v => ({id: v.id!, name: v.name!})) || [],
        [vaults],
    )

    useEffect(() => {
        const singleVault = vaultOptions.length === 1 ? vaultOptions[0] : null
        if (singleVault && !vaultId) {
            void selectVault(singleVault.id)
        }
    }, [vaultOptions, vaultId, selectVault])

    return (
        <div className={cls.NotesPageContainer}>
            <div className={cls.SidebarWrapper}>
                <NotesSidebar vaults={vaultOptions}/>
            </div>
            <div className={cls.ViewerWrapper}>
                <NoteViewer content={noteContent}/>
            </div>
        </div>
    )
}
