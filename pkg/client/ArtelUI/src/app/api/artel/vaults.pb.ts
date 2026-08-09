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
  s3InstanceId?: string;
  s3BucketName?: string;
  workbenchExists?: boolean;
  workbenchStatus?: string;
  postgresEnabled?: boolean;
  postgresStatus?: string;
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

export type LinkS3BucketRequest = {
  vaultId?: string;
  s3InstanceId?: string;
  bucketName?: string;
};

export type LinkS3BucketResponse = Record<string, never>;

export type LinkS3Bucket = Record<string, never>;

export type UnlinkS3BucketRequest = {
  vaultId?: string;
};

export type UnlinkS3BucketResponse = Record<string, never>;

export type UnlinkS3Bucket = Record<string, never>;

export type SetVaultBinaryStorageRequest = {
  vaultId?: string;
  useCouchdb?: boolean;
};

export type SetVaultBinaryStorageResponse = Record<string, never>;

export type SetVaultBinaryStorage = Record<string, never>;

export type PublishVaultRequest = {
  vaultId?: string;
  slug?: string;
};

export type PublishVaultResponse = {
  vault?: VaultItem;
};

export type PublishVault = Record<string, never>;

export type UnpublishVaultRequest = {
  vaultId?: string;
};

export type UnpublishVaultResponse = {
  vault?: VaultItem;
};

export type UnpublishVault = Record<string, never>;

export type CreateWorkbenchRequest = {
  vaultId?: string;
};

export type CreateWorkbenchResponse = {
  status?: string;
};

export type CreateWorkbench = Record<string, never>;

export type StartWorkbenchRequest = {
  vaultId?: string;
  authMode?: string;
};

export type StartWorkbenchResponse = {
  status?: string;
  authMode?: string;
};

export type StartWorkbench = Record<string, never>;

export type StopWorkbenchRequest = {
  vaultId?: string;
};

export type StopWorkbenchResponse = {
  status?: string;
};

export type StopWorkbench = Record<string, never>;

export type WatchWorkbenchLoginRequest = {
  vaultId?: string;
};

export type WatchWorkbenchLoginResponse = {
  state?: string;
  url?: string;
  errorMessage?: string;
};

export type WatchWorkbenchLogin = Record<string, never>;

export type SubmitWorkbenchLoginCodeRequest = {
  vaultId?: string;
  code?: string;
};

export type SubmitWorkbenchLoginCodeResponse = Record<string, never>;

export type SubmitWorkbenchLoginCode = Record<string, never>;

export type EnablePostgresDatabaseRequest = {
  vaultId?: string;
};

export type EnablePostgresDatabaseResponse = {
  status?: string;
  errorMessage?: string;
};

export type EnablePostgresDatabase = Record<string, never>;

export type GetPostgresDatabaseRequest = {
  vaultId?: string;
};

export type GetPostgresDatabaseResponse = {
  status?: string;
  errorMessage?: string;
};

export type GetPostgresDatabase = Record<string, never>;

export type DisablePostgresDatabaseRequest = {
  vaultId?: string;
};

export type DisablePostgresDatabaseResponse = {
  status?: string;
};

export type DisablePostgresDatabase = Record<string, never>;

export type VaultItem = {
  id?: string;
  name?: string;
  dbUrl?: string;
  livesyncPassphrase?: string;
  s3InstanceId?: string;
  s3BucketName?: string;
  useCouchdbForBinaries?: boolean;
  isPublic?: boolean;
  slug?: string;
  role?: string;
  postgresEnabled?: boolean;
  postgresStatus?: string;
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
  static LinkS3Bucket(this:void, req: LinkS3BucketRequest, initReq?: fm.InitReq): Promise<LinkS3BucketResponse> {
    return fm.fetchRequest<LinkS3BucketResponse>(`/api/vaults/s3/link`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UnlinkS3Bucket(this:void, req: UnlinkS3BucketRequest, initReq?: fm.InitReq): Promise<UnlinkS3BucketResponse> {
    return fm.fetchRequest<UnlinkS3BucketResponse>(`/api/vaults/s3/unlink`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SetVaultBinaryStorage(this:void, req: SetVaultBinaryStorageRequest, initReq?: fm.InitReq): Promise<SetVaultBinaryStorageResponse> {
    return fm.fetchRequest<SetVaultBinaryStorageResponse>(`/api/vaults/binary-storage`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static PublishVault(this:void, req: PublishVaultRequest, initReq?: fm.InitReq): Promise<PublishVaultResponse> {
    return fm.fetchRequest<PublishVaultResponse>(`/api/vaults/${req.vaultId}/publish`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static UnpublishVault(this:void, req: UnpublishVaultRequest, initReq?: fm.InitReq): Promise<UnpublishVaultResponse> {
    return fm.fetchRequest<UnpublishVaultResponse>(`/api/vaults/${req.vaultId}/unpublish`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CreateWorkbench(this:void, req: CreateWorkbenchRequest, initReq?: fm.InitReq): Promise<CreateWorkbenchResponse> {
    return fm.fetchRequest<CreateWorkbenchResponse>(`/api/vaults/workbench/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static StartWorkbench(this:void, req: StartWorkbenchRequest, initReq?: fm.InitReq): Promise<StartWorkbenchResponse> {
    return fm.fetchRequest<StartWorkbenchResponse>(`/api/vaults/workbench/start`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static StopWorkbench(this:void, req: StopWorkbenchRequest, initReq?: fm.InitReq): Promise<StopWorkbenchResponse> {
    return fm.fetchRequest<StopWorkbenchResponse>(`/api/vaults/workbench/stop`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static WatchWorkbenchLogin(this:void, req: WatchWorkbenchLoginRequest, entityNotifier?: fm.NotifyStreamEntityArrival<WatchWorkbenchLoginResponse>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<WatchWorkbenchLoginResponse>(`/api/vaults/workbench/login/watch`, entityNotifier, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SubmitWorkbenchLoginCode(this:void, req: SubmitWorkbenchLoginCodeRequest, initReq?: fm.InitReq): Promise<SubmitWorkbenchLoginCodeResponse> {
    return fm.fetchRequest<SubmitWorkbenchLoginCodeResponse>(`/api/vaults/workbench/login/submit`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static EnablePostgresDatabase(this:void, req: EnablePostgresDatabaseRequest, initReq?: fm.InitReq): Promise<EnablePostgresDatabaseResponse> {
    return fm.fetchRequest<EnablePostgresDatabaseResponse>(`/api/vaults/postgres/enable`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetPostgresDatabase(this:void, req: GetPostgresDatabaseRequest, initReq?: fm.InitReq): Promise<GetPostgresDatabaseResponse> {
    return fm.fetchRequest<GetPostgresDatabaseResponse>(`/api/vaults/postgres/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DisablePostgresDatabase(this:void, req: DisablePostgresDatabaseRequest, initReq?: fm.InitReq): Promise<DisablePostgresDatabaseResponse> {
    return fm.fetchRequest<DisablePostgresDatabaseResponse>(`/api/vaults/postgres/disable`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}