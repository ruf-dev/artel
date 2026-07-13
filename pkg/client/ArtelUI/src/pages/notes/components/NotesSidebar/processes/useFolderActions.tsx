import {ConfirmDialog, useToaster} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {ImportResolutionInput, useNotes} from "@/app/hooks/Notes.ts"
import ImportZipDialog from "@/pages/notes/components/ImportZipDialog/ImportZipDialog.tsx"
import ImportConflictsDialog from "@/pages/notes/components/ImportConflictsDialog/ImportConflictsDialog.tsx"

function triggerZipDownload(zipData: Uint8Array, filename: string) {
    const blob = new Blob([zipData], {type: "application/zip"})
    const url = URL.createObjectURL(blob)
    const link = document.createElement("a")
    link.href = url
    link.download = filename
    link.click()
    URL.revokeObjectURL(url)
}

// Folder download/delete/import handlers for the sidebar's "⋮" folder menu and Upload
// button — split out of NotesSidebar.tsx to keep that component under the max-lines limit.
export function useFolderActions() {
    const notesStore = useNotes()
    const {OpenDialog, CloseDialog} = useDialog()
    const {bake} = useToaster()
    const bakeError = useBakeError()

    function handleDownloadFolder(path: string) {
        notesStore.downloadFolder(path)
            .then(zipData => triggerZipDownload(zipData, `${path.split("/").pop() || "vault"}.zip`))
            .catch(err => bakeError("Failed to export folder", err))
    }

    function handleDeleteFolder(path: string) {
        OpenDialog(
            <ConfirmDialog
                title="Delete folder"
                message={`Delete "${path}" and everything inside? This cannot be undone.`}
                confirmLabel="Delete"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    notesStore.deleteFolder(path)
                        .then(result => {
                            if (result.failed.length > 0) {
                                bake({
                                    title: "Folder partially deleted",
                                    description: `Deleted ${result.deleted} item(s); ${result.failed.length} failed`,
                                    level: "Warn",
                                })
                                return
                            }
                            bake({
                                title: "Folder deleted",
                                description: `Deleted ${result.deleted} item(s)`,
                                level: "Info",
                            })
                        })
                        .catch(err => bakeError("Failed to delete folder", err))
                }
            />
        )
    }

    function handleImportZip(destPath: string, zipData: Uint8Array) {
        return notesStore.checkImportConflicts(destPath, zipData)
            .then(conflicts => {
                if (conflicts.length === 0) {
                    return notesStore.commitImport(destPath, zipData, []).then(result => {
                        CloseDialog()
                        bake({
                            title: "Import complete",
                            description: `Imported ${result.imported} item(s)`,
                            level: "Info",
                        })
                    })
                }

                OpenDialog(
                    <ImportConflictsDialog
                        conflicts={conflicts}
                        onConfirm={(resolutions: ImportResolutionInput[]) =>
                            notesStore.commitImport(destPath, zipData, resolutions)
                                .then(result => {
                                    CloseDialog()
                                    bake({
                                        title: "Import complete",
                                        description: `Imported ${result.imported}, skipped ${result.skipped}`,
                                        level: "Info",
                                    })
                                })
                                .catch(err => bakeError("Failed to import zip", err))
                        }
                    />
                )
            })
            .catch(err => bakeError("Failed to check import conflicts", err))
    }

    function handleUpload() {
        OpenDialog(<ImportZipDialog onConfirm={handleImportZip} />)
    }

    return {handleDownloadFolder, handleDeleteFolder, handleUpload}
}
