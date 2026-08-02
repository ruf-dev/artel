/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type GetVaultBySlugRequest = {
  slug?: string;
};

export type GetVaultBySlugResponse = {
  id?: string;
  name?: string;
};

export type GetVaultBySlug = Record<string, never>;

export type PublicDocsNoteItem = {
  path?: string;
  mtime?: string;
};

export type PublicDocsListFoldersRequest = {
  slug?: string;
};

export type PublicDocsListFoldersResponse = {
  folders?: string[];
};

export type PublicDocsListFolders = Record<string, never>;

export type PublicDocsListNotesRequest = {
  slug?: string;
};

export type PublicDocsListNotesResponse = {
  notes?: PublicDocsNoteItem[];
};

export type PublicDocsListNotes = Record<string, never>;

export type PublicDocsGetNoteRequest = {
  slug?: string;
  path?: string;
};

export type PublicDocsGetNoteResponse = {
  content?: string;
};

export type PublicDocsGetNote = Record<string, never>;

export type PublicDocsListTagsRequest = {
  slug?: string;
};

export type PublicDocsListTagsResponse = {
  tags?: string[];
};

export type PublicDocsListTags = Record<string, never>;

export class PublicDocsAPI {
  static GetVaultBySlug(this:void, req: GetVaultBySlugRequest, initReq?: fm.InitReq): Promise<GetVaultBySlugResponse> {
    return fm.fetchRequest<GetVaultBySlugResponse>(`/api/public-docs/vault`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListFolders(this:void, req: PublicDocsListFoldersRequest, initReq?: fm.InitReq): Promise<PublicDocsListFoldersResponse> {
    return fm.fetchRequest<PublicDocsListFoldersResponse>(`/api/public-docs/folders`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListNotes(this:void, req: PublicDocsListNotesRequest, initReq?: fm.InitReq): Promise<PublicDocsListNotesResponse> {
    return fm.fetchRequest<PublicDocsListNotesResponse>(`/api/public-docs/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetNote(this:void, req: PublicDocsGetNoteRequest, initReq?: fm.InitReq): Promise<PublicDocsGetNoteResponse> {
    return fm.fetchRequest<PublicDocsGetNoteResponse>(`/api/public-docs/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListTags(this:void, req: PublicDocsListTagsRequest, initReq?: fm.InitReq): Promise<PublicDocsListTagsResponse> {
    return fm.fetchRequest<PublicDocsListTagsResponse>(`/api/public-docs/tags`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}