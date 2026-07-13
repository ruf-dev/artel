import {create} from 'zustand'

import {NoteItem} from "@/app/api/artel/notes.pb.ts"
import {
    DeleteFolderResult, ImportConflictAction, ImportResolutionInput, ImportResult, notesService,
} from "@/processes/Notes.ts"

export {ImportConflictAction}
export type {NoteItem, DeleteFolderResult, ImportResolutionInput, ImportResult}
export type NoteMode = 'preview' | 'edit' | 'read'

const HIGHLIGHT_DURATION_MS = 2700

function requireVaultId(vaultId: string | null): string {
    if (!vaultId) throw new Error("no vault selected")
    return vaultId
}

async function commitImportAndRefresh(
    get: () => NotesState, destPath: string, zipData: Uint8Array, resolutions: ImportResolutionInput[],
): Promise<ImportResult> {
    const vaultId = requireVaultId(get().vaultId)
    const result = await notesService.commitImport(vaultId, destPath, zipData, resolutions)
    await get().selectVault(vaultId)
    return result
}

async function deleteFolderAndRefresh(get: () => NotesState, path: string): Promise<DeleteFolderResult> {
    const vaultId = requireVaultId(get().vaultId)
    const result = await notesService.deleteFolder(vaultId, path)
    await get().selectVault(vaultId)
    return result
}

interface NotesState {
    vaultId: string | null
    folders: string[]
    notes: NoteItem[]
    selectedPath: string | null
    noteContent: string | null
    savedContent: string | null
    mode: NoteMode
    loading: boolean
    error: string | null
    highlightedPath: string | null
    selectVault: (vaultId: string) => Promise<void>
    selectNote: (vaultId: string, path: string) => Promise<void>
    setMode: (mode: NoteMode) => void
    setContent: (content: string) => void
    setSavedContent: (content: string) => void
    saveNote: (content: string) => Promise<void>
    moveNote: (newPath: string) => Promise<void>
    createNote: (path: string) => Promise<void>
    listTags: (vaultId: string) => Promise<string[]>
    downloadFolder: (path: string) => Promise<Uint8Array>
    checkImportConflicts: (destPath: string, zipData: Uint8Array) => Promise<string[]>
    commitImport: (destPath: string, zipData: Uint8Array, resolutions: ImportResolutionInput[]) => Promise<ImportResult>
    deleteFolder: (path: string) => Promise<DeleteFolderResult>
    highlightNote: (path: string) => void
    reset: () => void
}

export const useNotes = create<NotesState>((set, get) => ({
    vaultId: null,
    folders: [],
    notes: [],
    selectedPath: null,
    noteContent: null,
    savedContent: null,
    mode: 'preview',
    loading: false,
    error: null,
    highlightedPath: null,

    selectVault: async (vaultId: string) => {
        set({loading: true, error: null, selectedPath: null, noteContent: null, savedContent: null, mode: 'preview'})
        try {
            const [folders, notes] = await Promise.all([
                notesService.listFolders(vaultId),
                notesService.listNotes(vaultId),
            ])
            set({vaultId, folders, notes})
        } catch (err) {
            set({error: String(err)})
        } finally {
            set({loading: false})
        }
    },

    selectNote: async (vaultId: string, path: string) => {
        set({loading: true, error: null})
        try {
            const noteContent = await notesService.getNote(vaultId, path)
            const savedMode = (localStorage.getItem(`artel.notes.mode.${path}`) as NoteMode) ?? 'preview'
            set({selectedPath: path, noteContent, savedContent: noteContent, mode: savedMode})
        } catch (err) {
            set({error: String(err)})
        } finally {
            set({loading: false})
        }
    },

    setMode: (mode: NoteMode) => {
        const { selectedPath } = get()
        if (selectedPath) {
            localStorage.setItem(`artel.notes.mode.${selectedPath}`, mode)
        }
        set({mode})
    },

    setContent: (content: string) => set({noteContent: content}),

    setSavedContent: (content: string) => set({savedContent: content}),

    saveNote: async (content: string) => {
        const { vaultId, selectedPath } = get()
        if (!vaultId || !selectedPath) return
        await notesService.saveNote(vaultId, selectedPath, content)
    },

    moveNote: async (newPath: string) => {
        const { vaultId, selectedPath } = get()
        if (!vaultId || !selectedPath) return
        await notesService.moveNote(vaultId, selectedPath, newPath)
        await get().selectVault(vaultId)
        await get().selectNote(vaultId, newPath)
    },

    createNote: async (path: string) => {
        const { vaultId } = get()
        if (!vaultId) return
        await notesService.saveNote(vaultId, path, '')
        await get().selectVault(vaultId)
        await get().selectNote(vaultId, path)
        get().setMode('edit')
    },

    listTags: (vaultId: string) => notesService.listTags(vaultId),

    downloadFolder: async (path: string) => notesService.exportFolder(requireVaultId(get().vaultId), path),

    checkImportConflicts: (destPath: string, zipData: Uint8Array) =>
        notesService.checkImportConflicts(requireVaultId(get().vaultId), destPath, zipData),

    commitImport: (destPath: string, zipData: Uint8Array, resolutions: ImportResolutionInput[]) =>
        commitImportAndRefresh(get, destPath, zipData, resolutions),

    deleteFolder: (path: string) => deleteFolderAndRefresh(get, path),

    highlightNote: (path: string) => {
        set({highlightedPath: null})
        requestAnimationFrame(() => {
            set({highlightedPath: path})
            setTimeout(() => {
                if (get().highlightedPath === path) set({highlightedPath: null})
            }, HIGHLIGHT_DURATION_MS)
        })
    },

    reset: () => {
        set({
            vaultId: null,
            folders: [],
            notes: [],
            selectedPath: null,
            noteContent: null,
            savedContent: null,
            mode: 'preview',
            loading: false,
            error: null,
            highlightedPath: null,
        })
    },
}))
