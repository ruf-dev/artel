/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export enum ImportConflictAction {
  SKIP = "SKIP",
  OVERWRITE = "OVERWRITE",
  RENAME = "RENAME",
}

export type NoteItem = {
  path?: string;
  mtime?: string;
};

export type ListFoldersRequest = {
  vaultId?: string;
};

export type ListFoldersResponse = {
  folders?: string[];
};

export type ListFolders = Record<string, never>;

export type ListNotesRequest = {
  vaultId?: string;
};

export type ListNotesResponse = {
  notes?: NoteItem[];
};

export type ListNotes = Record<string, never>;

export type GetNoteRequest = {
  vaultId?: string;
  path?: string;
};

export type GetNoteResponse = {
  content?: string;
};

export type GetNote = Record<string, never>;

export type ListTagsRequest = {
  vaultId?: string;
};

export type ListTagsResponse = {
  tags?: string[];
};

export type ListTags = Record<string, never>;

export type SaveNoteRequest = {
  vaultId?: string;
  path?: string;
  content?: string;
};

export type SaveNoteResponse = Record<string, never>;

export type SaveNote = Record<string, never>;

export type MoveNoteRequest = {
  vaultId?: string;
  oldPath?: string;
  newPath?: string;
};

export type MoveNoteResponse = Record<string, never>;

export type MoveNote = Record<string, never>;

export type ExportFolderRequest = {
  vaultId?: string;
  path?: string;
};

export type ExportFolderResponse = {
  zipData?: Uint8Array;
};

export type ExportFolder = Record<string, never>;

export type ImportResolution = {
  path?: string;
  action?: ImportConflictAction;
  renameTo?: string;
};

export type CheckImportConflictsRequest = {
  vaultId?: string;
  destPath?: string;
  zipData?: Uint8Array;
};

export type CheckImportConflictsResponse = {
  conflictingPaths?: string[];
};

export type CheckImportConflicts = Record<string, never>;

export type CommitImportRequest = {
  vaultId?: string;
  destPath?: string;
  zipData?: Uint8Array;
  resolutions?: ImportResolution[];
};

export type CommitImportResponse = {
  importedCount?: number;
  skippedCount?: number;
};

export type CommitImport = Record<string, never>;

export type DeleteFolderRequest = {
  vaultId?: string;
  path?: string;
};

export type DeleteFolderResponse = {
  deletedCount?: number;
  failedPaths?: string[];
};

export type DeleteFolder = Record<string, never>;

export class NotesAPI {
  static ListFolders(this:void, req: ListFoldersRequest, initReq?: fm.InitReq): Promise<ListFoldersResponse> {
    return fm.fetchRequest<ListFoldersResponse>(`/api/notes/folders`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListNotes(this:void, req: ListNotesRequest, initReq?: fm.InitReq): Promise<ListNotesResponse> {
    return fm.fetchRequest<ListNotesResponse>(`/api/notes/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetNote(this:void, req: GetNoteRequest, initReq?: fm.InitReq): Promise<GetNoteResponse> {
    return fm.fetchRequest<GetNoteResponse>(`/api/notes/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTags(this:void, req: ListTagsRequest, initReq?: fm.InitReq): Promise<ListTagsResponse> {
    return fm.fetchRequest<ListTagsResponse>(`/api/notes/tags`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SaveNote(this:void, req: SaveNoteRequest, initReq?: fm.InitReq): Promise<SaveNoteResponse> {
    return fm.fetchRequest<SaveNoteResponse>(`/api/notes/save`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static MoveNote(this:void, req: MoveNoteRequest, initReq?: fm.InitReq): Promise<MoveNoteResponse> {
    return fm.fetchRequest<MoveNoteResponse>(`/api/notes/move`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ExportFolder(this:void, req: ExportFolderRequest, initReq?: fm.InitReq): Promise<ExportFolderResponse> {
    return fm.fetchRequest<ExportFolderResponse>(`/api/notes/export`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckImportConflicts(this:void, req: CheckImportConflictsRequest, initReq?: fm.InitReq): Promise<CheckImportConflictsResponse> {
    return fm.fetchRequest<CheckImportConflictsResponse>(`/api/notes/import/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CommitImport(this:void, req: CommitImportRequest, initReq?: fm.InitReq): Promise<CommitImportResponse> {
    return fm.fetchRequest<CommitImportResponse>(`/api/notes/import/commit`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteFolder(this:void, req: DeleteFolderRequest, initReq?: fm.InitReq): Promise<DeleteFolderResponse> {
    return fm.fetchRequest<DeleteFolderResponse>(`/api/notes/folder/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}