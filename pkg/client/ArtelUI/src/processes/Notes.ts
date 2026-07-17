import useUser from "@/hooks/user/User.ts"
import {b64Decode} from "@/app/api/artel/fetch.pb.ts"
import {ImportConflictAction, ImportResolution, NoteItem, NotesAPI} from "@/app/api/artel/notes.pb.ts"

export {ImportConflictAction}

export interface ImportResolutionInput {
    path: string
    action: ImportConflictAction
    renameTo?: string
}

export interface ImportResult {
    imported: number
    skipped: number
}

export interface DeleteFolderResult {
    deleted: number
    failed: string[]
}

export interface INotesService {
    listFolders: (vaultId: string) => Promise<string[]>
    listNotes: (vaultId: string) => Promise<NoteItem[]>
    getNote: (vaultId: string, path: string) => Promise<string | null>
    listTags: (vaultId: string) => Promise<string[]>
    saveNote: (vaultId: string, path: string, content: string) => Promise<void>
    moveNote: (vaultId: string, oldPath: string, newPath: string) => Promise<void>
    moveFolder: (vaultId: string, oldPath: string, newPath: string) => Promise<void>
    exportFolder: (vaultId: string, path: string) => Promise<Uint8Array>
    checkImportConflicts: (vaultId: string, destPath: string, zipData: Uint8Array) => Promise<string[]>
    commitImport: (
        vaultId: string, destPath: string, zipData: Uint8Array, resolutions: ImportResolutionInput[],
    ) => Promise<ImportResult>
    deleteFolder: (vaultId: string, path: string) => Promise<DeleteFolderResult>
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

    async moveNote(vaultId: string, oldPath: string, newPath: string): Promise<void> {
        await NotesAPI.MoveNote({vaultId, oldPath, newPath}, useUser.getState().auth.getInitReq())
    }

    async moveFolder(vaultId: string, oldPath: string, newPath: string): Promise<void> {
        await NotesAPI.MoveFolder({vaultId, oldPath, newPath}, useUser.getState().auth.getInitReq())
    }

    async exportFolder(vaultId: string, path: string): Promise<Uint8Array> {
        const res = await NotesAPI.ExportFolder({vaultId, path}, useUser.getState().auth.getInitReq())
        // grpc-gateway sends `bytes` fields as base64 strings over the wire; the generated
        // Uint8Array type only holds once b64Decode has actually run.
        return b64Decode((res.zipData ?? "") as unknown as string)
    }

    async checkImportConflicts(vaultId: string, destPath: string, zipData: Uint8Array): Promise<string[]> {
        const req = {vaultId, destPath, zipData}
        const res = await NotesAPI.CheckImportConflicts(req, useUser.getState().auth.getInitReq())
        return res.conflictingPaths ?? []
    }

    async commitImport(
        vaultId: string, destPath: string, zipData: Uint8Array, resolutions: ImportResolutionInput[],
    ): Promise<ImportResult> {
        const pbResolutions: ImportResolution[] = resolutions.map(r => ({
            path: r.path,
            action: r.action,
            renameTo: r.renameTo ?? "",
        }))

        const res = await NotesAPI.CommitImport(
            {vaultId, destPath, zipData, resolutions: pbResolutions},
            useUser.getState().auth.getInitReq(),
        )

        return {imported: res.importedCount ?? 0, skipped: res.skippedCount ?? 0}
    }

    async deleteFolder(vaultId: string, path: string): Promise<DeleteFolderResult> {
        const res = await NotesAPI.DeleteFolder({vaultId, path}, useUser.getState().auth.getInitReq())
        return {deleted: res.deletedCount ?? 0, failed: res.failedPaths ?? []}
    }
}

export const notesService = new NotesService()
