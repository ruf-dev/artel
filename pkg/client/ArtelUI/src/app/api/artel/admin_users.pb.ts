/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";
import * as GoogleProtobufTimestamp from "./google/protobuf/timestamp.pb";


export type ArtelUserEntry = {
  userId?: string;
  username?: string;
  email?: string;
  createdAt?: GoogleProtobufTimestamp.Timestamp;
};

export type ListArtelUsersRequest = {
  search?: string;
  limit?: number;
  offset?: number;
};

export type ListArtelUsersResponse = {
  users?: ArtelUserEntry[];
  total?: string;
};

export type ListArtelUsers = Record<string, never>;

export type ArtelUserDetails = {
  userId?: string;
  username?: string;
  email?: string;
  isAdministrator?: boolean;
  hasEmails?: boolean;
  hasTaskTrackers?: boolean;
  hasNotes?: boolean;
  subscriptionActive?: boolean;
  hasSpreadsheets?: boolean;
};

export type GetArtelUserRequest = {
  userId?: string;
};

export type GetArtelUserResponse = {
  user?: ArtelUserDetails;
};

export type GetArtelUser = Record<string, never>;

export type UserSession = {
  sessionId?: string;
  createdAt?: GoogleProtobufTimestamp.Timestamp;
  expiresAt?: GoogleProtobufTimestamp.Timestamp;
};

export type GetUserSessionsRequest = {
  userId?: string;
};

export type GetUserSessionsResponse = {
  sessions?: UserSession[];
};

export type GetUserSessions = Record<string, never>;

export class AdminUsersAPI {
  static ListArtelUsers(this:void, req: ListArtelUsersRequest, initReq?: fm.InitReq): Promise<ListArtelUsersResponse> {
    return fm.fetchRequest<ListArtelUsersResponse>(`/api/admin_users/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetArtelUser(this:void, req: GetArtelUserRequest, initReq?: fm.InitReq): Promise<GetArtelUserResponse> {
    return fm.fetchRequest<GetArtelUserResponse>(`/api/admin_users/get`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetUserSessions(this:void, req: GetUserSessionsRequest, initReq?: fm.InitReq): Promise<GetUserSessionsResponse> {
    return fm.fetchRequest<GetUserSessionsResponse>(`/api/admin_users/sessions`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}