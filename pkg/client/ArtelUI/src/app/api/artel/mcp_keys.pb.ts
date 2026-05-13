/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type McpKeyInfo = {
  id?: string;
  vaultId?: string;
  name?: string;
  keyPreview?: string;
  createdAt?: string;
};

export type CreateMcpKeyRequest = {
  vaultId?: string;
  name?: string;
};

export type CreateMcpKeyResponse = {
  key?: McpKeyInfo;
  rawToken?: string;
};

export type CreateMcpKey = Record<string, never>;

export type ListMcpKeysRequest = {
  vaultId?: string;
};

export type ListMcpKeysResponse = {
  keys?: McpKeyInfo[];
};

export type ListMcpKeys = Record<string, never>;

export type RevokeMcpKeyRequest = {
  vaultId?: string;
  keyId?: string;
};

export type RevokeMcpKeyResponse = Record<string, never>;

export type RevokeMcpKey = Record<string, never>;

export class McpKeysAPI {
  static CreateMcpKey(this:void, req: CreateMcpKeyRequest, initReq?: fm.InitReq): Promise<CreateMcpKeyResponse> {
    return fm.fetchRequest<CreateMcpKeyResponse>(`/api/vaults/${req.vaultId}/mcp-keys`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListMcpKeys(this:void, req: ListMcpKeysRequest, initReq?: fm.InitReq): Promise<ListMcpKeysResponse> {
    return fm.fetchRequest<ListMcpKeysResponse>(`/api/vaults/${req.vaultId}/mcp-keys?${fm.renderURLSearchParams(req, ["vaultId"])}`, {...initReq, method: "GET"});
  }
  static RevokeMcpKey(this:void, req: RevokeMcpKeyRequest, initReq?: fm.InitReq): Promise<RevokeMcpKeyResponse> {
    return fm.fetchRequest<RevokeMcpKeyResponse>(`/api/vaults/${req.vaultId}/mcp-keys/${req.keyId}?${fm.renderURLSearchParams(req, ["vaultId", "keyId"])}`, {...initReq, method: "DELETE"});
  }
}