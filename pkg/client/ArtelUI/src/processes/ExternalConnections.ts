import {
    AddAnthropicConnectionRequest, AddEmailConnectionRequest, AddGenericConnectionRequest,
    AddGitlabConnectionRequest, AddTelegramConnectionRequest, AddTrelloConnectionRequest,
    AddOpenAIConnectionRequest, AddS3ConnectionRequest, AddCouchDBConnectionRequest,
    AddPostgresConnectionRequest,
    CheckAnthropicConnectionRequest, CheckAnthropicConnectionResponse,
    CheckOpenAIConnectionRequest, CheckOpenAIConnectionResponse,
    CheckS3ConnectionRequest, CheckS3ConnectionResponse,
    CheckCouchDBConnectionRequest, CheckCouchDBConnectionResponse,
    CheckPostgresConnectionRequest, CheckPostgresConnectionResponse,
    ExternalConnectionInfo, ExternalProvider,
    ExternalConnectionsAPI, GetProviderStatisticsResponse, Spreadsheet,
} from "@/app/api/artel/external_connections.pb.ts"
import * as fm from "@/app/api/artel/fetch.pb.ts"
import useUser from "@/hooks/user/User.ts"

export interface IExternalConnectionsService {
    list: () => Promise<ExternalConnectionInfo[]>
    initiateGoogleOAuth: () => Promise<string>
    exchangeGoogleOAuth: (code: string, state: string) => Promise<void>
    disconnect: (provider: string) => Promise<void>
    disconnectConnection: (id: string) => Promise<void>
    getPickerToken: () => Promise<string>
    addSpreadsheet: (spreadsheetId: string, name: string) => Promise<Spreadsheet>
    listSpreadsheets: () => Promise<Spreadsheet[]>
    removeSpreadsheet: (spreadsheetId: string) => Promise<void>
    addEmailConnection: (req: AddEmailConnectionRequest) => Promise<ExternalConnectionInfo>
    addGitlabConnection: (req: AddGitlabConnectionRequest) => Promise<ExternalConnectionInfo>
    generateGitlabWebhookSecret: () => Promise<{connection: ExternalConnectionInfo; webhookSecret: string}>
    addTelegramConnection: (req: AddTelegramConnectionRequest) => Promise<ExternalConnectionInfo>
    ensureArtelTelegramConnection: () => Promise<string>
    addTrelloConnection: (req: AddTrelloConnectionRequest) => Promise<ExternalConnectionInfo>
    addAnthropicConnection: (req: AddAnthropicConnectionRequest) => Promise<ExternalConnectionInfo>
    checkAnthropicConnection: (req: CheckAnthropicConnectionRequest) => Promise<CheckAnthropicConnectionResponse>
    addOpenAIConnection: (req: AddOpenAIConnectionRequest) => Promise<ExternalConnectionInfo>
    checkOpenAIConnection: (req: CheckOpenAIConnectionRequest) => Promise<CheckOpenAIConnectionResponse>
    addS3Connection: (req: AddS3ConnectionRequest) => Promise<ExternalConnectionInfo>
    checkS3Connection: (req: CheckS3ConnectionRequest) => Promise<CheckS3ConnectionResponse>
    addCouchDBConnection: (req: AddCouchDBConnectionRequest) => Promise<ExternalConnectionInfo>
    checkCouchDBConnection: (req: CheckCouchDBConnectionRequest) => Promise<CheckCouchDBConnectionResponse>
    addPostgresConnection: (req: AddPostgresConnectionRequest) => Promise<ExternalConnectionInfo>
    checkPostgresConnection: (req: CheckPostgresConnectionRequest) => Promise<CheckPostgresConnectionResponse>
    addGenericConnection: (req: AddGenericConnectionRequest) => Promise<ExternalConnectionInfo>
    getProviderStatistics: (provider: ExternalProvider) => Promise<GetProviderStatisticsResponse>
}

export class ExternalConnectionsService implements IExternalConnectionsService {
    async list(): Promise<ExternalConnectionInfo[]> {
        const res = await ExternalConnectionsAPI.ListConnections({}, useUser.getState().auth.getInitReq())
        return res.connections ?? []
    }

    async initiateGoogleOAuth(): Promise<string> {
        const res = await ExternalConnectionsAPI.InitiateGoogleOAuth(
            {origin: window.location.origin}, useUser.getState().auth.getInitReq(),
        )
        return res.authUrl ?? ""
    }

    async exchangeGoogleOAuth(code: string, state: string): Promise<void> {
        const initReq = useUser.getState().auth.getInitReq()
        await fm.fetchRequest("/api/external-connections/google/exchange", {
            ...initReq,
            method: "POST",
            body: JSON.stringify({code, state}),
        })
    }

    async disconnect(provider: string): Promise<void> {
        await ExternalConnectionsAPI.DisconnectProvider({provider}, useUser.getState().auth.getInitReq())
    }

    async disconnectConnection(id: string): Promise<void> {
        await ExternalConnectionsAPI.DisconnectConnection({id}, useUser.getState().auth.getInitReq())
    }

    async getPickerToken(): Promise<string> {
        const res = await ExternalConnectionsAPI.GetGooglePickerToken({}, useUser.getState().auth.getInitReq())
        return res.accessToken ?? ""
    }

    async addSpreadsheet(spreadsheetId: string, name: string): Promise<Spreadsheet> {
        const res = await ExternalConnectionsAPI.AddSpreadsheet(
            {spreadsheetId, name}, useUser.getState().auth.getInitReq(),
        )
        return res.spreadsheet!
    }

    async listSpreadsheets(): Promise<Spreadsheet[]> {
        const res = await ExternalConnectionsAPI.ListSpreadsheets({}, useUser.getState().auth.getInitReq())
        return res.spreadsheets ?? []
    }

    async removeSpreadsheet(spreadsheetId: string): Promise<void> {
        await ExternalConnectionsAPI.RemoveSpreadsheet({spreadsheetId}, useUser.getState().auth.getInitReq())
    }

    async addEmailConnection(req: AddEmailConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddEmailConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async addGitlabConnection(req: AddGitlabConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddGitlabConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async generateGitlabWebhookSecret(): Promise<{connection: ExternalConnectionInfo; webhookSecret: string}> {
        const res = await ExternalConnectionsAPI.GenerateGitlabWebhookSecret({}, useUser.getState().auth.getInitReq())
        return {connection: res.connection!, webhookSecret: res.webhookSecret ?? ""}
    }

    async addTelegramConnection(req: AddTelegramConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddTelegramConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async ensureArtelTelegramConnection(): Promise<string> {
        const res = await ExternalConnectionsAPI.EnsureArtelTelegramConnection({}, useUser.getState().auth.getInitReq())
        return res.externalConnectionId ?? ""
    }

    async addTrelloConnection(req: AddTrelloConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddTrelloConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async addAnthropicConnection(req: AddAnthropicConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddAnthropicConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async checkAnthropicConnection(req: CheckAnthropicConnectionRequest): Promise<CheckAnthropicConnectionResponse> {
        return ExternalConnectionsAPI.CheckAnthropicConnection(req, useUser.getState().auth.getInitReq())
    }

    async addOpenAIConnection(req: AddOpenAIConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddOpenAIConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async checkOpenAIConnection(req: CheckOpenAIConnectionRequest): Promise<CheckOpenAIConnectionResponse> {
        return ExternalConnectionsAPI.CheckOpenAIConnection(req, useUser.getState().auth.getInitReq())
    }

    async addS3Connection(req: AddS3ConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddS3Connection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async checkS3Connection(req: CheckS3ConnectionRequest): Promise<CheckS3ConnectionResponse> {
        return ExternalConnectionsAPI.CheckS3Connection(req, useUser.getState().auth.getInitReq())
    }

    async addCouchDBConnection(req: AddCouchDBConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddCouchDBConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async checkCouchDBConnection(req: CheckCouchDBConnectionRequest): Promise<CheckCouchDBConnectionResponse> {
        return ExternalConnectionsAPI.CheckCouchDBConnection(req, useUser.getState().auth.getInitReq())
    }

    async addPostgresConnection(req: AddPostgresConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddPostgresConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async checkPostgresConnection(req: CheckPostgresConnectionRequest): Promise<CheckPostgresConnectionResponse> {
        return ExternalConnectionsAPI.CheckPostgresConnection(req, useUser.getState().auth.getInitReq())
    }

    async addGenericConnection(req: AddGenericConnectionRequest): Promise<ExternalConnectionInfo> {
        const res = await ExternalConnectionsAPI.AddGenericConnection(req, useUser.getState().auth.getInitReq())
        return res.connection!
    }

    async getProviderStatistics(provider: ExternalProvider): Promise<GetProviderStatisticsResponse> {
        return ExternalConnectionsAPI.GetProviderStatistics({provider}, useUser.getState().auth.getInitReq())
    }
}

export const externalConnectionsService = new ExternalConnectionsService()
