import {useEffect, useMemo, useRef, useState} from "react"
import {useLocation, useNavigate, useParams} from "react-router-dom"

import {useNotes} from "@/app/hooks/Notes.ts"
import {useVaults} from "@/app/hooks/Vaults.ts"
import {usePortrait} from "@/app/hooks/usePortrait.ts"
import MobileNotesShell from "@/pages/notes/components/MobileNotesShell/MobileNotesShell.tsx"
import DesktopNotesShell from "@/pages/notes/components/DesktopNotesShell/DesktopNotesShell.tsx"
import {useAutosave} from "@/pages/notes/hooks/useAutosave.ts"
import {NoteMode} from "@/app/hooks/Notes.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {buildNotesUrl, decodeNotePath} from "@/pages/notes/processes/notesUrl.ts"
import {useReadOnlyVaultGate} from "@/pages/notes/processes/useReadOnlyVaultGate.ts"
import {buildRenameHandler} from "@/pages/notes/processes/renameHandler.tsx"


export default function NotesPage() {
    const notesStore = useNotes()
    const { OpenDialog } = useDialog()
    const bakeError = useBakeError()
    const {vaults} = useVaults()
    const scrollTopRef = useRef<number>(0)
    const [fontScale, setFontScale] = useState(1)
    const [previewEditing, setPreviewEditing] = useState(false)
    const isPortrait = usePortrait()
    const navigate = useNavigate()
    const location = useLocation()
    const routeParams = useParams()
    const vaultIdParam = routeParams.vaultId ?? null
    const notePathParam = routeParams["*"] ? decodeNotePath(routeParams["*"]) : null

    function zoomIn()  { setFontScale(s => Math.min(2.5, +(s + 0.15).toFixed(2))) }
    function zoomOut() { setFontScale(s => Math.max(0.5, +(s - 0.15).toFixed(2))) }

    const vaultOptions = useMemo(
        () => vaults
            .filter(v => v.id && v.name)
            .map(v => ({id: v.id!, name: v.name!, isPublic: v.isPublic})) || [],
        [vaults],
    )

    const isReadOnlyVault = useReadOnlyVaultGate(vaults, notesStore.vaultId, notesStore.mode, notesStore.setMode)

    useEffect(() => {
        if (vaultIdParam) return
        const singleVault = vaultOptions.length === 1 ? vaultOptions[0] : null
        if (singleVault && !notesStore.vaultId) {
            void notesStore.selectVault(singleVault.id)
        }
    }, [vaultOptions, notesStore.vaultId, notesStore.selectVault, vaultIdParam])

    // URL -> store: loads the vault/note named in the URL (initial load, pasted link, back/forward).
    useEffect(() => {
        if (vaultIdParam && vaultIdParam !== notesStore.vaultId) {
            void notesStore.selectVault(vaultIdParam)
        }
    }, [vaultIdParam])

    useEffect(() => {
        if (vaultIdParam && notePathParam && notePathParam !== notesStore.selectedPath) {
            void notesStore.selectNote(vaultIdParam, notePathParam)
        }
    }, [vaultIdParam, notePathParam])

    // store -> URL: keeps the address bar in sync whenever the selected vault/note changes.
    useEffect(() => {
        const targetPath = buildNotesUrl(notesStore.vaultId, notesStore.selectedPath)
        if (targetPath !== location.pathname) {
            navigate({pathname: targetPath, search: location.search}, {replace: true})
        }
    }, [notesStore.vaultId, notesStore.selectedPath])

    useEffect(() => { setPreviewEditing(false) }, [notesStore.selectedPath])

    useEffect(() => {
        if (isReadOnlyVault) return
        if (notesStore.selectedPath && notesStore.noteContent === '' && notesStore.mode === 'preview') {
            setPreviewEditing(true)
        }
    }, [notesStore.selectedPath, notesStore.noteContent, isReadOnlyVault])

    const { saveStatus, saveError, forceSave } = useAutosave({
        noteId: notesStore.selectedPath,
        vaultId: notesStore.vaultId,
        content: notesStore.noteContent ?? '',
        savedContent: notesStore.savedContent ?? '',
        onSave: notesStore.saveNote,
        onSavedContentChange: notesStore.setSavedContent,
    })

    function handleModeChange(newMode: NoteMode) {
        if (newMode === 'edit' && isReadOnlyVault) return
        if (notesStore.mode === 'edit' || (notesStore.mode === 'preview' && previewEditing)) {
            forceSave()
        }
        setPreviewEditing(false)
        notesStore.setMode(newMode)
    }

    const handleRename = buildRenameHandler({
        selectedPath: notesStore.selectedPath, openDialog: OpenDialog, moveNote: notesStore.moveNote, bakeError,
    })

    const mode = notesStore.mode
    const selectedPath = notesStore.selectedPath
    const showEditor = (mode === 'edit' || (mode === 'preview' && previewEditing)) && selectedPath !== null

    const shellProps = {
        vaultOptions, mode, selectedPath, showEditor, scrollTopRef, fontScale,
        noteContent: notesStore.noteContent,
        saveStatus: showEditor ? saveStatus : 'idle' as const,
        saveError, onModeChange: handleModeChange, onChange: notesStore.setContent,
        onContentClick: mode === 'preview' && selectedPath && !isReadOnlyVault
            ? () => setPreviewEditing(true) : undefined,
        onEscape: mode === 'preview' ? () => setPreviewEditing(false) : undefined,
        onRename: handleRename,
        hideEdit: isReadOnlyVault,
    }

    if (isPortrait) {
        return <MobileNotesShell {...shellProps}/>
    }

    return <DesktopNotesShell {...shellProps} onZoomIn={zoomIn} onZoomOut={zoomOut}/>
}
