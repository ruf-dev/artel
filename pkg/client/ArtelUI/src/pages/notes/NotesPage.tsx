import {useEffect, useMemo, useRef} from "react"

import cls from "./NotesPage.module.css"
import {useNotes} from "@/app/hooks/Notes.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import NotesSidebar from "@/pages/notes/components/NotesSidebar/NotesSidebar.tsx"
import NoteViewer from "@/pages/notes/components/NoteViewer/NoteViewer.tsx"
import NoteEditor from "@/pages/notes/components/NoteEditor/NoteEditor.tsx"
import BreadcrumbBar from "@/pages/notes/components/BreadcrumbBar/BreadcrumbBar.tsx"
import RenameDialog from "@/pages/notes/components/RenameDialog/RenameDialog.tsx"
import {useAutosave} from "@/pages/notes/hooks/useAutosave.ts"
import {NoteMode} from "@/app/hooks/Notes.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

export default function NotesPage() {
    const {
        vaultId,
        noteContent,
        savedContent,
        mode,
        selectedPath,
        selectVault,
        setMode,
        setContent,
        setSavedContent,
        saveNote,
        moveNote,
    } = useNotes()
    const { OpenDialog } = useDialog()
    const bakeError = useBakeError()
    const {vaults} = useVaults()
    const scrollTopRef = useRef<number>(0)

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

    const { saveStatus, saveError, forceSave } = useAutosave({
        noteId: selectedPath,
        vaultId,
        content: noteContent ?? '',
        savedContent: savedContent ?? '',
        onSave: saveNote,
        onSavedContentChange: setSavedContent,
    })

    function handleModeChange(newMode: NoteMode) {
        if (mode === 'edit') {
            forceSave()
        }
        setMode(newMode)
    }

    function handleRename() {
        if (!selectedPath) return
        OpenDialog(
            <RenameDialog
                currentPath={selectedPath}
                onConfirm={async (newPath: string) => {
                    try {
                        await moveNote(newPath)
                    } catch (err) {
                        bakeError("Failed to move note", err)
                    }
                }}
            />
        )
    }

    const showEditor = mode === 'edit' && selectedPath !== null

    return (
        <div className={cls.NotesPageContainer}>
            <div className={cls.SidebarWrapper}>
                <NotesSidebar vaults={vaultOptions}/>
            </div>
            <div className={cls.ViewerWrapper}>
                {selectedPath && (
                    <BreadcrumbBar
                        path={selectedPath}
                        mode={mode}
                        onModeChange={handleModeChange}
                        saveStatus={showEditor ? saveStatus : 'idle'}
                        saveError={saveError}
                        onRename={handleRename}
                    />
                )}
                {showEditor ? (
                    <NoteEditor
                        content={noteContent ?? ''}
                        onChange={setContent}
                        scrollTopRef={scrollTopRef}
                    />
                ) : (
                    <NoteViewer content={noteContent}/>
                )}
            </div>
        </div>
    )
}
