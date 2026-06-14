/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);

export enum ExternalProvider {
  EXTERNAL_PROVIDER_UNSPECIFIED = "EXTERNAL_PROVIDER_UNSPECIFIED",
  EXTERNAL_PROVIDER_GOOGLE_SHEETS = "EXTERNAL_PROVIDER_GOOGLE_SHEETS",
  EXTERNAL_PROVIDER_TRELLO = "EXTERNAL_PROVIDER_TRELLO",
  EXTERNAL_PROVIDER_MIRO = "EXTERNAL_PROVIDER_MIRO",
}

export type GoogleConnectionInfo = {
  email?: string;
  scopes?: string;
};

export type GenericConnection = {
  fields?: Record<string, string>;
};

type BaseExternalConnectionInfo = {
  id?: string;
  provider?: ExternalProvider;
  createdAt?: string;
  updatedAt?: string;
};

export type ExternalConnectionInfo = BaseExternalConnectionInfo &
  OneOf<{
    google: GoogleConnectionInfo;
    generic: GenericConnection;
  }>;

export type InitiateGoogleOAuthRequest = Record<string, never>;

export type InitiateGoogleOAuthResponse = {
  authUrl?: string;
};

export type InitiateGoogleOAuth = Record<string, never>;

export type ListConnectionsRequest = Record<string, never>;

export type ListConnectionsResponse = {
  connections?: ExternalConnectionInfo[];
};

export type ListConnections = Record<string, never>;

export type DisconnectProviderRequest = {
  provider?: string;
};

export type DisconnectProviderResponse = Record<string, never>;

export type DisconnectProvider = Record<string, never>;

export class ExternalConnectionsAPI {
  static InitiateGoogleOAuth(this:void, req: InitiateGoogleOAuthRequest, initReq?: fm.InitReq): Promise<InitiateGoogleOAuthResponse> {
    return fm.fetchRequest<InitiateGoogleOAuthResponse>(`/api/external-connections/google/initiate`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListConnections(this:void, req: ListConnectionsRequest, initReq?: fm.InitReq): Promise<ListConnectionsResponse> {
    return fm.fetchRequest<ListConnectionsResponse>(`/api/external-connections/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static DisconnectProvider(this:void, req: DisconnectProviderRequest, initReq?: fm.InitReq): Promise<DisconnectProviderResponse> {
    return fm.fetchRequest<DisconnectProviderResponse>(`/api/external-connections/disconnect`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}