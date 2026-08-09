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
  EXTERNAL_PROVIDER_EMAIL = "EXTERNAL_PROVIDER_EMAIL",
  EXTERNAL_PROVIDER_GITLAB = "EXTERNAL_PROVIDER_GITLAB",
  EXTERNAL_PROVIDER_ANTHROPIC = "EXTERNAL_PROVIDER_ANTHROPIC",
  EXTERNAL_PROVIDER_TELEGRAM = "EXTERNAL_PROVIDER_TELEGRAM",
  EXTERNAL_PROVIDER_OPENAI = "EXTERNAL_PROVIDER_OPENAI",
  EXTERNAL_PROVIDER_S3 = "EXTERNAL_PROVIDER_S3",
  EXTERNAL_PROVIDER_COUCHDB = "EXTERNAL_PROVIDER_COUCHDB",
  EXTERNAL_PROVIDER_POSTGRES = "EXTERNAL_PROVIDER_POSTGRES",
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
  providerName?: string;
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

export type DisconnectConnectionRequest = {
  id?: string;
};

export type DisconnectConnectionResponse = Record<string, never>;

export type DisconnectConnection = Record<string, never>;

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

export type AddEmailConnectionRequest = {
  email?: string;
  imapHost?: string;
  imapPort?: number;
  smtpHost?: string;
  smtpPort?: number;
  password?: string;
};

export type AddEmailConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddEmailConnection = Record<string, never>;

export type CheckEmailConnectionRequest = {
  email?: string;
  imapHost?: string;
  imapPort?: number;
  smtpHost?: string;
  smtpPort?: number;
  password?: string;
};

export type CheckEmailConnectionResponse = Record<string, never>;

export type CheckEmailConnection = Record<string, never>;

export type MailServerSuggestion = {
  domain?: string;
  smtp?: string;
  smtpPort?: number;
  imap?: string;
  imapPort?: number;
  appPasswordUrl?: string;
};

export type ListMailServerSuggestionsRequest = {
  domain?: string;
};

export type ListMailServerSuggestionsResponse = {
  suggestions?: MailServerSuggestion[];
};

export type ListMailServerSuggestions = Record<string, never>;

export type AddGitlabConnectionRequest = {
  personalAccessToken?: string;
  instanceUrl?: string;
};

export type AddGitlabConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddGitlabConnection = Record<string, never>;

export type CheckGitlabConnectionRequest = {
  personalAccessToken?: string;
  instanceUrl?: string;
};

export type CheckGitlabConnectionResponse = {
  username?: string;
};

export type CheckGitlabConnection = Record<string, never>;

export type GenerateGitlabWebhookSecretRequest = Record<string, never>;

export type GenerateGitlabWebhookSecretResponse = {
  connection?: ExternalConnectionInfo;
  webhookSecret?: string;
};

export type GenerateGitlabWebhookSecret = Record<string, never>;

export type AddTrelloConnectionRequest = {
  apiKey?: string;
  apiToken?: string;
};

export type AddTrelloConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddTrelloConnection = Record<string, never>;

export type CheckTrelloConnectionRequest = {
  apiKey?: string;
  apiToken?: string;
};

export type CheckTrelloConnectionResponse = {
  fullName?: string;
};

export type CheckTrelloConnection = Record<string, never>;

export type AddTelegramConnectionRequest = {
  botToken?: string;
};

export type AddTelegramConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddTelegramConnection = Record<string, never>;

export type CheckTelegramConnectionRequest = {
  botToken?: string;
};

export type CheckTelegramConnectionResponse = {
  botUsername?: string;
};

export type CheckTelegramConnection = Record<string, never>;

export type AddAnthropicConnectionRequest = {
  apiKey?: string;
  baseUrl?: string;
  defaultModel?: string;
};

export type AddAnthropicConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddAnthropicConnection = Record<string, never>;

export type CheckAnthropicConnectionRequest = {
  apiKey?: string;
  baseUrl?: string;
  defaultModel?: string;
};

export type CheckAnthropicConnectionResponse = {
  availableModels?: string[];
  recommendedDefaultModel?: string;
};

export type CheckAnthropicConnection = Record<string, never>;

export type AddOpenAIConnectionRequest = {
  apiKey?: string;
  baseUrl?: string;
  defaultModel?: string;
};

export type AddOpenAIConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddOpenAIConnection = Record<string, never>;

export type CheckOpenAIConnectionRequest = {
  apiKey?: string;
  baseUrl?: string;
  defaultModel?: string;
};

export type CheckOpenAIConnectionResponse = {
  availableModels?: string[];
  recommendedDefaultModel?: string;
};

export type CheckOpenAIConnection = Record<string, never>;

export type AddS3ConnectionRequest = {
  endpoint?: string;
  region?: string;
  accessKey?: string;
  secretKey?: string;
  useSsl?: boolean;
  pathStyle?: boolean;
};

export type AddS3ConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddS3Connection = Record<string, never>;

export type CheckS3ConnectionRequest = {
  endpoint?: string;
  region?: string;
  accessKey?: string;
  secretKey?: string;
  useSsl?: boolean;
  pathStyle?: boolean;
};

export type CheckS3ConnectionResponse = Record<string, never>;

export type CheckS3Connection = Record<string, never>;

export type AddCouchDBConnectionRequest = {
  url?: string;
  username?: string;
  password?: string;
};

export type AddCouchDBConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddCouchDBConnection = Record<string, never>;

export type CheckCouchDBConnectionRequest = {
  url?: string;
  username?: string;
  password?: string;
};

export type CheckCouchDBConnectionResponse = Record<string, never>;

export type CheckCouchDBConnection = Record<string, never>;

export type AddPostgresConnectionRequest = {
  host?: string;
  port?: number;
  database?: string;
  username?: string;
  password?: string;
  sslMode?: string;
};

export type AddPostgresConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddPostgresConnection = Record<string, never>;

export type CheckPostgresConnectionRequest = {
  host?: string;
  port?: number;
  database?: string;
  username?: string;
  password?: string;
  sslMode?: string;
};

export type CheckPostgresConnectionResponse = Record<string, never>;

export type CheckPostgresConnection = Record<string, never>;

export type AddGenericConnectionRequest = {
  provider?: string;
  credentials?: Record<string, string>;
};

export type AddGenericConnectionResponse = {
  connection?: ExternalConnectionInfo;
};

export type AddGenericConnection = Record<string, never>;

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
  static DisconnectConnection(this:void, req: DisconnectConnectionRequest, initReq?: fm.InitReq): Promise<DisconnectConnectionResponse> {
    return fm.fetchRequest<DisconnectConnectionResponse>(`/api/external-connections/disconnect-connection`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
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
  static AddEmailConnection(this:void, req: AddEmailConnectionRequest, initReq?: fm.InitReq): Promise<AddEmailConnectionResponse> {
    return fm.fetchRequest<AddEmailConnectionResponse>(`/api/external-connections/email/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckEmailConnection(this:void, req: CheckEmailConnectionRequest, initReq?: fm.InitReq): Promise<CheckEmailConnectionResponse> {
    return fm.fetchRequest<CheckEmailConnectionResponse>(`/api/external-connections/email/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListMailServerSuggestions(this:void, req: ListMailServerSuggestionsRequest, initReq?: fm.InitReq): Promise<ListMailServerSuggestionsResponse> {
    return fm.fetchRequest<ListMailServerSuggestionsResponse>(`/api/external-connections/email/suggestions`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddGitlabConnection(this:void, req: AddGitlabConnectionRequest, initReq?: fm.InitReq): Promise<AddGitlabConnectionResponse> {
    return fm.fetchRequest<AddGitlabConnectionResponse>(`/api/external-connections/gitlab/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckGitlabConnection(this:void, req: CheckGitlabConnectionRequest, initReq?: fm.InitReq): Promise<CheckGitlabConnectionResponse> {
    return fm.fetchRequest<CheckGitlabConnectionResponse>(`/api/external-connections/gitlab/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GenerateGitlabWebhookSecret(this:void, req: GenerateGitlabWebhookSecretRequest, initReq?: fm.InitReq): Promise<GenerateGitlabWebhookSecretResponse> {
    return fm.fetchRequest<GenerateGitlabWebhookSecretResponse>(`/api/external-connections/gitlab/webhook-secret`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddTrelloConnection(this:void, req: AddTrelloConnectionRequest, initReq?: fm.InitReq): Promise<AddTrelloConnectionResponse> {
    return fm.fetchRequest<AddTrelloConnectionResponse>(`/api/external-connections/trello/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckTrelloConnection(this:void, req: CheckTrelloConnectionRequest, initReq?: fm.InitReq): Promise<CheckTrelloConnectionResponse> {
    return fm.fetchRequest<CheckTrelloConnectionResponse>(`/api/external-connections/trello/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddTelegramConnection(this:void, req: AddTelegramConnectionRequest, initReq?: fm.InitReq): Promise<AddTelegramConnectionResponse> {
    return fm.fetchRequest<AddTelegramConnectionResponse>(`/api/external-connections/telegram/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckTelegramConnection(this:void, req: CheckTelegramConnectionRequest, initReq?: fm.InitReq): Promise<CheckTelegramConnectionResponse> {
    return fm.fetchRequest<CheckTelegramConnectionResponse>(`/api/external-connections/telegram/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddAnthropicConnection(this:void, req: AddAnthropicConnectionRequest, initReq?: fm.InitReq): Promise<AddAnthropicConnectionResponse> {
    return fm.fetchRequest<AddAnthropicConnectionResponse>(`/api/external-connections/anthropic/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckAnthropicConnection(this:void, req: CheckAnthropicConnectionRequest, initReq?: fm.InitReq): Promise<CheckAnthropicConnectionResponse> {
    return fm.fetchRequest<CheckAnthropicConnectionResponse>(`/api/external-connections/anthropic/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddOpenAIConnection(this:void, req: AddOpenAIConnectionRequest, initReq?: fm.InitReq): Promise<AddOpenAIConnectionResponse> {
    return fm.fetchRequest<AddOpenAIConnectionResponse>(`/api/external-connections/openai/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckOpenAIConnection(this:void, req: CheckOpenAIConnectionRequest, initReq?: fm.InitReq): Promise<CheckOpenAIConnectionResponse> {
    return fm.fetchRequest<CheckOpenAIConnectionResponse>(`/api/external-connections/openai/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddGenericConnection(this:void, req: AddGenericConnectionRequest, initReq?: fm.InitReq): Promise<AddGenericConnectionResponse> {
    return fm.fetchRequest<AddGenericConnectionResponse>(`/api/external-connections/generic`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddS3Connection(this:void, req: AddS3ConnectionRequest, initReq?: fm.InitReq): Promise<AddS3ConnectionResponse> {
    return fm.fetchRequest<AddS3ConnectionResponse>(`/api/external-connections/s3/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckS3Connection(this:void, req: CheckS3ConnectionRequest, initReq?: fm.InitReq): Promise<CheckS3ConnectionResponse> {
    return fm.fetchRequest<CheckS3ConnectionResponse>(`/api/external-connections/s3/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddCouchDBConnection(this:void, req: AddCouchDBConnectionRequest, initReq?: fm.InitReq): Promise<AddCouchDBConnectionResponse> {
    return fm.fetchRequest<AddCouchDBConnectionResponse>(`/api/external-connections/couchdb/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckCouchDBConnection(this:void, req: CheckCouchDBConnectionRequest, initReq?: fm.InitReq): Promise<CheckCouchDBConnectionResponse> {
    return fm.fetchRequest<CheckCouchDBConnectionResponse>(`/api/external-connections/couchdb/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddPostgresConnection(this:void, req: AddPostgresConnectionRequest, initReq?: fm.InitReq): Promise<AddPostgresConnectionResponse> {
    return fm.fetchRequest<AddPostgresConnectionResponse>(`/api/external-connections/postgres/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CheckPostgresConnection(this:void, req: CheckPostgresConnectionRequest, initReq?: fm.InitReq): Promise<CheckPostgresConnectionResponse> {
    return fm.fetchRequest<CheckPostgresConnectionResponse>(`/api/external-connections/postgres/check`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}