import {
    AddEmailConnectionRequest, AddGitlabConnectionRequest, ExternalConnectionInfo, ExternalConnectionsAPI, Spreadsheet,
} from "@/app/api/artel/external_connections.pb.ts"
import * as fm from "@/app/api/artel/fetch.pb.ts"
import useUser from "@/hooks/user/User.ts"

export interface IExternalConnectionsService {
    list: () => Promise<ExternalConnectionInfo[]>
    initiateGoogleOAuth: () => Promise<string>
    exchangeGoogleOAuth: (code: string, state: string) => Promise<void>
    disconnect: (provider: string) => Promise<void>
    getPickerToken: () => Promise<string>
    addSpreadsheet: (spreadsheetId: string, name: string) => Promise<Spreadsheet>
    listSpreadsheets: () => Promise<Spreadsheet[]>
    removeSpreadsheet: (spreadsheetId: string) => Promise<void>
    addEmailConnection: (req: AddEmailConnectionRequest) => Promise<ExternalConnectionInfo>
    addGitlabConnection: (req: AddGitlabConnectionRequest) => Promise<ExternalConnectionInfo>
    generateGitlabWebhookSecret: () => Promise<{connection: ExternalConnectionInfo; webhookSecret: string}>
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
}

export const externalConnectionsService = new ExternalConnectionsService()
