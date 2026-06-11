import {create} from 'zustand'

import {NoteItem} from "@/app/api/artel/notes.pb.ts"
import {notesService} from "@/processes/Notes.ts"

export type {NoteItem}

interface NotesState {
    vaultId: string | null
    folders: string[]
    notes: NoteItem[]
    selectedPath: string | null
    noteContent: string | null
    loading: boolean
    error: string | null
    selectVault: (vaultId: string) => Promise<void>
    selectNote: (vaultId: string, path: string) => Promise<void>
    listTags: (vaultId: string) => Promise<string[]>
    reset: () => void
}

export const useNotes = create<NotesState>((set) => ({
    vaultId: null,
    folders: [],
    notes: [],
    selectedPath: null,
    noteContent: null,
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
            set({selectedPath: path, noteContent})
        } catch (err) {
            set({error: String(err)})
        } finally {
            set({loading: false})
        }
    },

    listTags: (vaultId: string) => notesService.listTags(vaultId),

    reset: () => {
        set({
            vaultId: null,
            folders: [],
            notes: [],
            selectedPath: null,
            noteContent: null,
            loading: false,
            error: null,
        })
    },
}))
