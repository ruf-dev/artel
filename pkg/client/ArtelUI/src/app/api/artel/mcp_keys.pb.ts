/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as ArtelApiExternalConnections from "./external_connections.pb";
import * as fm from "./fetch.pb";

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);

export enum SmtpOperation {
  SMTP_OP_UNSPECIFIED = "SMTP_OP_UNSPECIFIED",
  SMTP_OP_SEND = "SMTP_OP_SEND",
}

export enum ImapOperation {
  IMAP_OP_UNSPECIFIED = "IMAP_OP_UNSPECIFIED",
  IMAP_OP_LIST_FOLDERS = "IMAP_OP_LIST_FOLDERS",
  IMAP_OP_LIST_MESSAGES = "IMAP_OP_LIST_MESSAGES",
  IMAP_OP_FETCH_MESSAGE = "IMAP_OP_FETCH_MESSAGE",
}

export type McpKeyInfo = {
  id?: string;
  vaultId?: string;
  name?: string;
  keyPreview?: string;
  createdAt?: string;
  lastAccessedAt?: string;
  emailAccountId?: string;
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

export type ListUserMcpKeysRequest = Record<string, never>;

export type ListUserMcpKeysResponse = {
  keys?: McpKeyInfo[];
};

export type ListUserMcpKeys = Record<string, never>;

export type SetMcpKeyAccessRequest = {
  keyId?: string;
  vaultId?: string;
  emailAccountId?: string;
};

export type SetMcpKeyAccessResponse = Record<string, never>;

export type SetMcpKeyAccess = Record<string, never>;

export type McpConnectorInfo = {
  mcpKeyId?: string;
  mcpName?: string;
  externalConnectionId?: string;
  createdAt?: string;
};

export type SmtpToolAction = {
  operation?: SmtpOperation;
};

export type ImapToolAction = {
  operation?: ImapOperation;
};

export type StringParam = Record<string, never>;

export type IntegerParam = Record<string, never>;

export type EnumParam = {
  values?: string[];
};

type BaseToolParamDef = {
  description?: string;
};

export type ToolParamDef = BaseToolParamDef &
  OneOf<{
    stringParam: StringParam;
    integerParam: IntegerParam;
    enumParam: EnumParam;
  }>;

type BaseMcpToolInfo = {
  name?: string;
  description?: string;
  params?: Record<string, ToolParamDef>;
};

export type McpToolInfo = BaseMcpToolInfo &
  OneOf<{
    smtp: SmtpToolAction;
    imap: ImapToolAction;
  }>;

export type MomCandidate = {
  name?: string;
  author?: string;
  description?: string;
  connections?: ArtelApiExternalConnections.ExternalConnectionInfo[];
  tools?: McpToolInfo[];
};

export type ListMcpConnectorsRequest = {
  keyId?: string;
};

export type ListMcpConnectorsResponse = {
  connectors?: McpConnectorInfo[];
};

export type ListMcpConnectors = Record<string, never>;

export type AddMcpConnectorRequest = {
  keyId?: string;
  mcpName?: string;
  externalConnectionId?: string;
};

export type AddMcpConnectorResponse = {
  connector?: McpConnectorInfo;
};

export type AddMcpConnector = Record<string, never>;

export type RemoveMcpConnectorRequest = {
  keyId?: string;
  mcpName?: string;
};

export type RemoveMcpConnectorResponse = Record<string, never>;

export type RemoveMcpConnector = Record<string, never>;

export type ListMomCandidatesRequest = Record<string, never>;

export type ListMomCandidatesResponse = {
  candidates?: MomCandidate[];
};

export type ListMomCandidates = Record<string, never>;

export type ExecuteMomToolRequest = {
  mcpName?: string;
  toolName?: string;
  externalConnectionId?: string;
  paramsJson?: string;
};

export type ExecuteMomToolResponse = {
  result?: string;
};

export type ExecuteMomTool = Record<string, never>;

export class McpKeysAPI {
  static CreateMcpKey(this:void, req: CreateMcpKeyRequest, initReq?: fm.InitReq): Promise<CreateMcpKeyResponse> {
    return fm.fetchRequest<CreateMcpKeyResponse>(`/api/mcp/keys/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListMcpKeys(this:void, req: ListMcpKeysRequest, initReq?: fm.InitReq): Promise<ListMcpKeysResponse> {
    return fm.fetchRequest<ListMcpKeysResponse>(`/api/mcp/keys/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RevokeMcpKey(this:void, req: RevokeMcpKeyRequest, initReq?: fm.InitReq): Promise<RevokeMcpKeyResponse> {
    return fm.fetchRequest<RevokeMcpKeyResponse>(`/api/mcp/keys/revoke`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListUserMcpKeys(this:void, req: ListUserMcpKeysRequest, initReq?: fm.InitReq): Promise<ListUserMcpKeysResponse> {
    return fm.fetchRequest<ListUserMcpKeysResponse>(`/api/mcp/keys/list-user`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SetMcpKeyAccess(this:void, req: SetMcpKeyAccessRequest, initReq?: fm.InitReq): Promise<SetMcpKeyAccessResponse> {
    return fm.fetchRequest<SetMcpKeyAccessResponse>(`/api/mcp/keys/access`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListMcpConnectors(this:void, req: ListMcpConnectorsRequest, initReq?: fm.InitReq): Promise<ListMcpConnectorsResponse> {
    return fm.fetchRequest<ListMcpConnectorsResponse>(`/api/mcp/keys/connectors/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static AddMcpConnector(this:void, req: AddMcpConnectorRequest, initReq?: fm.InitReq): Promise<AddMcpConnectorResponse> {
    return fm.fetchRequest<AddMcpConnectorResponse>(`/api/mcp/keys/connectors/add`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RemoveMcpConnector(this:void, req: RemoveMcpConnectorRequest, initReq?: fm.InitReq): Promise<RemoveMcpConnectorResponse> {
    return fm.fetchRequest<RemoveMcpConnectorResponse>(`/api/mcp/keys/connectors/remove`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ListMomCandidates(this:void, req: ListMomCandidatesRequest, initReq?: fm.InitReq): Promise<ListMomCandidatesResponse> {
    return fm.fetchRequest<ListMomCandidatesResponse>(`/api/mcp/moms/candidates`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static ExecuteMomTool(this:void, req: ExecuteMomToolRequest, initReq?: fm.InitReq): Promise<ExecuteMomToolResponse> {
    return fm.fetchRequest<ExecuteMomToolResponse>(`/api/mcp/moms/execute`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}