import {create} from 'zustand'

import {NoteItem} from "@/app/api/artel/notes.pb.ts"
import {notesService} from "@/processes/Notes.ts"

export type {NoteItem}
export type NoteMode = 'preview' | 'edit' | 'read'

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
    selectVault: (vaultId: string) => Promise<void>
    selectNote: (vaultId: string, path: string) => Promise<void>
    setMode: (mode: NoteMode) => void
    setContent: (content: string) => void
    setSavedContent: (content: string) => void
    saveNote: (content: string) => Promise<void>
    moveNote: (newPath: string) => Promise<void>
    listTags: (vaultId: string) => Promise<string[]>
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

    selectVault: async (vaultId: string) => {
        set({loading: true, error: null})
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

    setContent: (content: string) => {
        set({noteContent: content})
    },

    setSavedContent: (content: string) => {
        set({savedContent: content})
    },

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

    listTags: (vaultId: string) => notesService.listTags(vaultId),

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
        })
    },
}))
