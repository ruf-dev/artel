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

export class NotesAPI {
  static ListFolders(this:void, req: ListFoldersRequest, initReq?: fm.InitReq): Promise<ListFoldersResponse> {
    return fm.fetchRequest<ListFoldersResponse>(`/api/notes/${req.vaultId}/folders?${fm.renderURLSearchParams(req, ["vaultId"])}`, {...initReq, method: "GET"});
  }
  static ListNotes(this:void, req: ListNotesRequest, initReq?: fm.InitReq): Promise<ListNotesResponse> {
    return fm.fetchRequest<ListNotesResponse>(`/api/notes/${req.vaultId}?${fm.renderURLSearchParams(req, ["vaultId"])}`, {...initReq, method: "GET"});
  }
  static GetNote(this:void, req: GetNoteRequest, initReq?: fm.InitReq): Promise<GetNoteResponse> {
    return fm.fetchRequest<GetNoteResponse>(`/api/notes/${req.vaultId}/note`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTags(this:void, req: ListTagsRequest, initReq?: fm.InitReq): Promise<ListTagsResponse> {
    return fm.fetchRequest<ListTagsResponse>(`/api/notes/${req.vaultId}/tags?${fm.renderURLSearchParams(req, ["vaultId"])}`, {...initReq, method: "GET"});
  }
}