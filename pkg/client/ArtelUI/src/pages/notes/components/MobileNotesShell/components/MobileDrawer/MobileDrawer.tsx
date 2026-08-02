import {Button} from "@vervstack/chures"

import { useNotes } from "@/app/hooks/Notes.ts"
import { useDialog } from "@/app/hooks/Dialog.ts"
import { useBakeError } from "@/app/hooks/useErrorToast.ts"
import NotesSidebar from "@/pages/notes/components/NotesSidebar/NotesSidebar.tsx"
import CreateNoteDialog from "@/pages/notes/components/CreateNoteDialog/CreateNoteDialog.tsx"
import ArtelLogoIcon from "@/pages/notes/components/icons/ArtelLogoIcon.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import DrawerCloseButton from "@/pages/notes/components/MobileNotesShell/components/MobileDrawer/components/DrawerCloseButton/DrawerCloseButton.tsx"
import cls from "@/pages/notes/components/MobileNotesShell/components/MobileDrawer/MobileDrawer.module.css"

interface VaultOption {
    id: string
    name: string
    isPublic?: boolean
}

interface MobileDrawerProps {
    open: boolean
    onClose: () => void
    vaultOptions: VaultOption[]
}

export default function MobileDrawer({ open, onClose, vaultOptions }: MobileDrawerProps) {
    const { vaultId, folders, createNote } = useNotes()
    const { OpenDialog } = useDialog()
    const bakeError = useBakeError()

    function handleNewNote() {
        OpenDialog(
            <CreateNoteDialog
                initialPath=""
                folders={folders}
                onConfirm={(path: string) =>
                    createNote(path).catch(err => bakeError("Failed to create note", err))
                }
            />
        )
        onClose()
    }

    const backdropClass = `${cls.Backdrop}${open ? ` ${cls.BackdropOpen}` : ""}`
    const panelClass = `${cls.DrawerPanel}${open ? ` ${cls.DrawerPanelOpen}` : ""}`

    return (
        <>
            <div className={backdropClass} onClick={onClose} />
            <div className={panelClass}>
                <div className={cls.DrawerStatusSpacer} />
                <div className={cls.DrawerHeader}>
                    <ArtelLogoIcon/>
                    <span className={cls.DrawerWordmark}>artel</span>
                    <span className={cls.DrawerNotesLabel}>notes</span>
                    <span className={cls.DrawerSpacer} />
                    <DrawerCloseButton onClose={onClose}/>
                </div>
                <div className={cls.DrawerBody}>
                    <NotesSidebar vaults={vaultOptions} showCreateButton={false} />
                </div>
                {vaultId && (
                    <div className={cls.DrawerFooter}>
                        <Button variant="primary" className={cls.NewNoteBtn} onClick={handleNewNote}>
                            <svg
                                viewBox="0 0 12 12"
                                width={11}
                                height={11}
                                fill="none"
                                stroke="currentColor"
                                strokeWidth="1.8"
                                strokeLinecap="round"
                            >
                                <path d="M6 1v10M1 6h10" />
                            </svg>
                            New Note
                        </Button>
                    </div>
                )}
            </div>
        </>
    )
}
