import {JSX} from "react"

import RenameDialog from "@/pages/notes/components/RenameDialog/RenameDialog.tsx"

interface RenameHandlerArgs {
    selectedPath: string | null
    openDialog: (...children: JSX.Element[]) => void
    moveNote: (newPath: string) => Promise<void>
    bakeError: (title: string, err: unknown) => void
}

export function buildRenameHandler({selectedPath, openDialog, moveNote, bakeError}: RenameHandlerArgs) {
    return function handleRename() {
        if (!selectedPath) return
        openDialog(
            <RenameDialog
                currentPath={selectedPath}
                onConfirm={(newPath: string) =>
                    moveNote(newPath).catch(err => bakeError("Failed to move note", err))
                }
            />
        )
    }
}
