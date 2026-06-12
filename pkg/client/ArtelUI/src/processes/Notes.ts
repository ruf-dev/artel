import useUser from "@/hooks/user/User.ts"
import {NoteItem, NotesAPI} from "@/app/api/artel/notes.pb.ts"

export interface INotesService {
    listFolders: (vaultId: string) => Promise<string[]>
    listNotes: (vaultId: string) => Promise<NoteItem[]>
    getNote: (vaultId: string, path: string) => Promise<string | null>
    listTags: (vaultId: string) => Promise<string[]>
    saveNote: (vaultId: string, path: string, content: string) => Promise<void>
}

export class NotesService implements INotesService {
    async listFolders(vaultId: string): Promise<string[]> {
        const res = await NotesAPI.ListFolders({vaultId}, useUser.getState().auth.getInitReq())
        return res.folders ?? []
    }

    async listNotes(vaultId: string): Promise<NoteItem[]> {
        const res = await NotesAPI.ListNotes({vaultId}, useUser.getState().auth.getInitReq())
        return res.notes ?? []
    }

    async getNote(vaultId: string, path: string): Promise<string | null> {
        const res = await NotesAPI.GetNote({vaultId, path}, useUser.getState().auth.getInitReq())
        return res.content ?? null
    }

    async listTags(vaultId: string): Promise<string[]> {
        const res = await NotesAPI.ListTags({vaultId}, useUser.getState().auth.getInitReq())
        return res.tags ?? []
    }

    async saveNote(vaultId: string, path: string, content: string): Promise<void> {
        await NotesAPI.SaveNote({vaultId, path, content}, useUser.getState().auth.getInitReq())
    }
}

export const notesService = new NotesService()
