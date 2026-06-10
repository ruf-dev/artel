/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type CouchUserEntry = {
  name?: string;
  roles?: string[];
};

export type ListCouchUsersRequest = {
  instanceId?: string;
};

export type ListCouchUsersResponse = {
  users?: CouchUserEntry[];
};

export type ListCouchUsers = Record<string, never>;

export type DeleteCouchUserRequest = {
  instanceId?: string;
  username?: string;
};

export type DeleteCouchUserResponse = Record<string, never>;

export type DeleteCouchUser = Record<string, never>;

export type ChangeCouchUserPasswordRequest = {
  instanceId?: string;
  username?: string;
  newPassword?: string;
};

export type ChangeCouchUserPasswordResponse = Record<string, never>;

export type ChangeCouchUserPassword = Record<string, never>;

export type ListCouchDatabasesRequest = {
  instanceId?: string;
};

export type ListCouchDatabasesResponse = {
  databases?: string[];
};

export type ListCouchDatabases = Record<string, never>;

export type GrantDatabaseAccessRequest = {
  instanceId?: string;
  dbName?: string;
  username?: string;
};

export type GrantDatabaseAccessResponse = Record<string, never>;

export type GrantDatabaseAccess = Record<string, never>;

export type RevokeDatabaseAccessRequest = {
  instanceId?: string;
  dbName?: string;
  username?: string;
};

export type RevokeDatabaseAccessResponse = Record<string, never>;

export type RevokeDatabaseAccess = Record<string, never>;

export type GetUserDatabaseAccessRequest = {
  instanceId?: string;
  username?: string;
};

export type GetUserDatabaseAccessResponse = {
  databases?: string[];
};

export type GetUserDatabaseAccess = Record<string, never>;

export class AdminCouchAPI {
  static ListCouchUsers(this:void, req: ListCouchUsersRequest, initReq?: fm.InitReq): Promise<ListCouchUsersResponse> {
    return fm.fetchRequest<ListCouchUsersResponse>(`/api/admin_couch/users/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DeleteCouchUser(this:void, req: DeleteCouchUserRequest, initReq?: fm.InitReq): Promise<DeleteCouchUserResponse> {
    return fm.fetchRequest<DeleteCouchUserResponse>(`/api/admin_couch/users/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ChangeCouchUserPassword(this:void, req: ChangeCouchUserPasswordRequest, initReq?: fm.InitReq): Promise<ChangeCouchUserPasswordResponse> {
    return fm.fetchRequest<ChangeCouchUserPasswordResponse>(`/api/admin_couch/users/change-password`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListCouchDatabases(this:void, req: ListCouchDatabasesRequest, initReq?: fm.InitReq): Promise<ListCouchDatabasesResponse> {
    return fm.fetchRequest<ListCouchDatabasesResponse>(`/api/admin_couch/db/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GrantDatabaseAccess(this:void, req: GrantDatabaseAccessRequest, initReq?: fm.InitReq): Promise<GrantDatabaseAccessResponse> {
    return fm.fetchRequest<GrantDatabaseAccessResponse>(`/api/admin_couch/db/grant`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RevokeDatabaseAccess(this:void, req: RevokeDatabaseAccessRequest, initReq?: fm.InitReq): Promise<RevokeDatabaseAccessResponse> {
    return fm.fetchRequest<RevokeDatabaseAccessResponse>(`/api/admin_couch/db/revoke`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetUserDatabaseAccess(this:void, req: GetUserDatabaseAccessRequest, initReq?: fm.InitReq): Promise<GetUserDatabaseAccessResponse> {
    return fm.fetchRequest<GetUserDatabaseAccessResponse>(`/api/admin_couch/db/user-access`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}