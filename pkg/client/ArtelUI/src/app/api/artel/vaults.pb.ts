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
  role?: string;
};

export type AddMemberResponse = Record<string, never>;

export type AddMember = Record<string, never>;

export type RemoveMemberRequest = {
  vaultId?: string;
  userId?: string;
};

export type RemoveMemberResponse = Record<string, never>;

export type RemoveMember = Record<string, never>;

export type ListMembersRequest = {
  vaultId?: string;
};

export type ListMembersResponse = {
  members?: VaultMemberInfo[];
};

export type ListMembers = Record<string, never>;

export type CreateInviteLinkRequest = {
  vaultId?: string;
  role?: string;
};

export type CreateInviteLinkResponse = {
  invite?: VaultInviteItem;
};

export type CreateInviteLink = Record<string, never>;

export type ListInviteLinksRequest = {
  vaultId?: string;
};

export type ListInviteLinksResponse = {
  invites?: VaultInviteItem[];
};

export type ListInviteLinks = Record<string, never>;

export type RevokeInviteLinkRequest = {
  vaultId?: string;
  inviteId?: string;
};

export type RevokeInviteLinkResponse = Record<string, never>;

export type RevokeInviteLink = Record<string, never>;

export type AcceptInviteRequest = {
  token?: string;
};

export type AcceptInviteResponse = {
  vaultId?: string;
};

export type AcceptInvite = Record<string, never>;

export type VaultItem = {
  id?: string;
  name?: string;
  dbUrl?: string;
  livesyncPassphrase?: string;
};

export type VaultMemberInfo = {
  id?: string;
  userId?: string;
  email?: string;
  username?: string;
  role?: string;
};

export type VaultInviteItem = {
  id?: string;
  vaultId?: string;
  role?: string;
  token?: string;
  revoked?: boolean;
  createdAt?: string;
};

export class VaultsAPI {
  static CreateVault(this:void, req: CreateVaultRequest, initReq?: fm.InitReq): Promise<CreateVaultResponse> {
    return fm.fetchRequest<CreateVaultResponse>(`/api/vaults/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetVault(this:void, req: GetVaultRequest, initReq?: fm.InitReq): Promise<GetVaultResponse> {
    return fm.fetchRequest<GetVaultResponse>(`/api/vaults/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListVaults(this:void, req: ListVaultsRequest, initReq?: fm.InitReq): Promise<ListVaultsResponse> {
    return fm.fetchRequest<ListVaultsResponse>(`/api/vaults/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteVault(this:void, req: DeleteVaultRequest, initReq?: fm.InitReq): Promise<DeleteVaultResponse> {
    return fm.fetchRequest<DeleteVaultResponse>(`/api/vaults/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddMember(this:void, req: AddMemberRequest, initReq?: fm.InitReq): Promise<AddMemberResponse> {
    return fm.fetchRequest<AddMemberResponse>(`/api/vaults/members/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RemoveMember(this:void, req: RemoveMemberRequest, initReq?: fm.InitReq): Promise<RemoveMemberResponse> {
    return fm.fetchRequest<RemoveMemberResponse>(`/api/vaults/members/remove`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListMembers(this:void, req: ListMembersRequest, initReq?: fm.InitReq): Promise<ListMembersResponse> {
    return fm.fetchRequest<ListMembersResponse>(`/api/vaults/members/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CreateInviteLink(this:void, req: CreateInviteLinkRequest, initReq?: fm.InitReq): Promise<CreateInviteLinkResponse> {
    return fm.fetchRequest<CreateInviteLinkResponse>(`/api/vaults/invites/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListInviteLinks(this:void, req: ListInviteLinksRequest, initReq?: fm.InitReq): Promise<ListInviteLinksResponse> {
    return fm.fetchRequest<ListInviteLinksResponse>(`/api/vaults/invites/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RevokeInviteLink(this:void, req: RevokeInviteLinkRequest, initReq?: fm.InitReq): Promise<RevokeInviteLinkResponse> {
    return fm.fetchRequest<RevokeInviteLinkResponse>(`/api/vaults/invites/revoke`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AcceptInvite(this:void, req: AcceptInviteRequest, initReq?: fm.InitReq): Promise<AcceptInviteResponse> {
    return fm.fetchRequest<AcceptInviteResponse>(`/api/vaults/join`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}