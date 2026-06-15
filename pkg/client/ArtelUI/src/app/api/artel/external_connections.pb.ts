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

export type InitiateGoogleOAuthRequest = {
  origin?: string;
};

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

export type GooglePickerTokenRequest = Record<string, never>;

export type GooglePickerTokenResponse = {
  accessToken?: string;
};

export type GooglePickerToken = Record<string, never>;

export type Spreadsheet = {
  id?: string;
  name?: string;
  createdAt?: string;
};

export type AddSpreadsheetRequest = {
  spreadsheetId?: string;
  name?: string;
};

export type AddSpreadsheetResponse = {
  spreadsheet?: Spreadsheet;
};

export type AddSpreadsheet = Record<string, never>;

export type ListSpreadsheetsRequest = Record<string, never>;

export type ListSpreadsheetsResponse = {
  spreadsheets?: Spreadsheet[];
};

export type ListSpreadsheets = Record<string, never>;

export type RemoveSpreadsheetRequest = {
  spreadsheetId?: string;
};

export type RemoveSpreadsheetResponse = Record<string, never>;

export type RemoveSpreadsheet = Record<string, never>;

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
  static GetGooglePickerToken(this:void, req: GooglePickerTokenRequest, initReq?: fm.InitReq): Promise<GooglePickerTokenResponse> {
    return fm.fetchRequest<GooglePickerTokenResponse>(`/api/external-connections/google/picker-token`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddSpreadsheet(this:void, req: AddSpreadsheetRequest, initReq?: fm.InitReq): Promise<AddSpreadsheetResponse> {
    return fm.fetchRequest<AddSpreadsheetResponse>(`/api/external-connections/google/spreadsheets/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListSpreadsheets(this:void, req: ListSpreadsheetsRequest, initReq?: fm.InitReq): Promise<ListSpreadsheetsResponse> {
    return fm.fetchRequest<ListSpreadsheetsResponse>(`/api/external-connections/google/spreadsheets/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RemoveSpreadsheet(this:void, req: RemoveSpreadsheetRequest, initReq?: fm.InitReq): Promise<RemoveSpreadsheetResponse> {
    return fm.fetchRequest<RemoveSpreadsheetResponse>(`/api/external-connections/google/spreadsheets/remove`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}