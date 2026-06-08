import {create} from 'zustand'

import {NoteItem, NotesAPI} from "@/app/api/artel/notes.pb.ts"
import useUser from "@/hooks/user/User.ts"

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
            const {auth} = useUser.getState()
            const [foldersRes, notesRes] = await Promise.all([
                NotesAPI.ListFolders({vaultId}, auth.getInitReq()),
                NotesAPI.ListNotes({vaultId}, auth.getInitReq()),
            ])
            set({
                vaultId,
                folders: foldersRes.folders ?? [],
                notes: notesRes.notes ?? [],
            })
        } catch (err) {
            set({error: String(err)})
        } finally {
            set({loading: false})
        }
    },

    selectNote: async (vaultId: string, path: string) => {
        set({loading: true, error: null})
        try {
            const {auth} = useUser.getState()
            const res = await NotesAPI.GetNote({vaultId, path}, auth.getInitReq())
            set({selectedPath: path, noteContent: res.content ?? null})
        } catch (err) {
            set({error: String(err)})
        } finally {
            set({loading: false})
        }
    },

    listTags: async (vaultId: string) => {
        const {auth} = useUser.getState()
        const res = await NotesAPI.ListTags({vaultId}, auth.getInitReq())
        return res.tags ?? []
    },

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
