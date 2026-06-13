/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


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
}