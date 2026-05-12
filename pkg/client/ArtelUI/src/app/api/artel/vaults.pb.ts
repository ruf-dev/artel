/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type CreateVaultRequest = {
  name?: string;
};

export type CreateVaultResponse = {
  id?: string;
  name?: string;
  dbUrl?: string;
};

export type CreateVault = Record<string, never>;

export type GetVaultRequest = {
  id?: string;
};

export type GetVaultResponse = {
  id?: string;
  name?: string;
  dbUrl?: string;
};

export type GetVault = Record<string, never>;

export type ListVaultsRequest = Record<string, never>;

export type ListVaultsResponse = {
  vaults?: VaultItem[];
};

export type ListVaults = Record<string, never>;

export type DeleteVaultRequest = {
  id?: string;
};

export type DeleteVaultResponse = Record<string, never>;

export type DeleteVault = Record<string, never>;

export type AddMemberRequest = {
  vaultId?: string;
  userId?: string;
};

export type AddMemberResponse = Record<string, never>;

export type AddMember = Record<string, never>;

export type RemoveMemberRequest = {
  vaultId?: string;
  userId?: string;
};

export type RemoveMemberResponse = Record<string, never>;

export type RemoveMember = Record<string, never>;

export type VaultItem = {
  id?: string;
  name?: string;
  dbUrl?: string;
};

export class VaultsAPI {
  static CreateVault(this:void, req: CreateVaultRequest, initReq?: fm.InitReq): Promise<CreateVaultResponse> {
    return fm.fetchRequest<CreateVaultResponse>(`/api/vaults/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetVault(this:void, req: GetVaultRequest, initReq?: fm.InitReq): Promise<GetVaultResponse> {
    return fm.fetchRequest<GetVaultResponse>(`/api/vaults/${req.id}?${fm.renderURLSearchParams(req, ["id"])}`, {...initReq, method: "GET"});
  }
  static ListVaults(this:void, req: ListVaultsRequest, initReq?: fm.InitReq): Promise<ListVaultsResponse> {
    return fm.fetchRequest<ListVaultsResponse>(`/api/vaults/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteVault(this:void, req: DeleteVaultRequest, initReq?: fm.InitReq): Promise<DeleteVaultResponse> {
    return fm.fetchRequest<DeleteVaultResponse>(`/api/vaults/${req.id}/delete`, {...initReq, method: "POST"});
  }
  static AddMember(this:void, req: AddMemberRequest, initReq?: fm.InitReq): Promise<AddMemberResponse> {
    return fm.fetchRequest<AddMemberResponse>(`/api/vaults/${req.vaultId}/members`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RemoveMember(this:void, req: RemoveMemberRequest, initReq?: fm.InitReq): Promise<RemoveMemberResponse> {
    return fm.fetchRequest<RemoveMemberResponse>(`/api/vaults/${req.vaultId}/members/${req.userId}?${fm.renderURLSearchParams(req, ["vaultId", "userId"])}`, {...initReq, method: "DELETE"});
  }
}