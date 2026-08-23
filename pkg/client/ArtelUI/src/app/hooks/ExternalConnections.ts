import {create} from 'zustand'

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
    GetProviderStatisticsResponse,
    Spreadsheet,
} from "@/app/api/artel/external_connections.pb.ts"
import {externalConnectionsService} from "@/processes/ExternalConnections.ts"

interface ExternalConnectionsState {
    connections: ExternalConnectionInfo[]
    loading: boolean
    spreadsheets: Spreadsheet[]
    spreadsheetsLoading: boolean
    fetch: () => Promise<void>
    disconnect: (provider: string) => Promise<void>
    disconnectConnection: (id: string) => Promise<void>
    initiateGoogleOAuth: () => Promise<string>
    getPickerToken: () => Promise<string>
    fetchSpreadsheets: () => Promise<void>
    addSpreadsheet: (spreadsheetId: string, name: string) => Promise<void>
    removeSpreadsheet: (spreadsheetId: string) => Promise<void>
    addEmailConnection: (req: AddEmailConnectionRequest) => Promise<void>
    addGitlabConnection: (req: AddGitlabConnectionRequest) => Promise<void>
    generateGitlabWebhookSecret: () => Promise<string>
    addTelegramConnection: (req: AddTelegramConnectionRequest) => Promise<void>
    addTrelloConnection: (req: AddTrelloConnectionRequest) => Promise<void>
    addAnthropicConnection: (req: AddAnthropicConnectionRequest) => Promise<void>
    checkAnthropicConnection: (req: CheckAnthropicConnectionRequest) => Promise<CheckAnthropicConnectionResponse>
    addOpenAIConnection: (req: AddOpenAIConnectionRequest) => Promise<void>
    checkOpenAIConnection: (req: CheckOpenAIConnectionRequest) => Promise<CheckOpenAIConnectionResponse>
    addS3Connection: (req: AddS3ConnectionRequest) => Promise<void>
    checkS3Connection: (req: CheckS3ConnectionRequest) => Promise<CheckS3ConnectionResponse>
    addCouchDBConnection: (req: AddCouchDBConnectionRequest) => Promise<void>
    checkCouchDBConnection: (req: CheckCouchDBConnectionRequest) => Promise<CheckCouchDBConnectionResponse>
    addPostgresConnection: (req: AddPostgresConnectionRequest) => Promise<void>
    checkPostgresConnection: (req: CheckPostgresConnectionRequest) => Promise<CheckPostgresConnectionResponse>
    addGenericConnection: (req: AddGenericConnectionRequest) => Promise<void>
    getProviderStatistics: (provider: ExternalProvider) => Promise<GetProviderStatisticsResponse>
}

// eslint-disable-next-line max-lines-per-function -- Zustand store with many methods
export const useExternalConnections = create<ExternalConnectionsState>((set, get) => ({
    connections: [],
    loading: false,
    spreadsheets: [],
    spreadsheetsLoading: false,

    fetch: async () => {
        set({loading: true})
        try {
            const connections = await externalConnectionsService.list()
            set({connections})
        } finally {
            set({loading: false})
        }
    },

    disconnect: async (provider: string) => {
        await externalConnectionsService.disconnect(provider)
        set({spreadsheets: []})
        await get().fetch()
    },

    disconnectConnection: async (id: string) => {
        await externalConnectionsService.disconnectConnection(id)
        await get().fetch()
    },

    initiateGoogleOAuth: async () => {
        return externalConnectionsService.initiateGoogleOAuth()
    },

    getPickerToken: () => externalConnectionsService.getPickerToken(),

    fetchSpreadsheets: async () => {
        set({spreadsheetsLoading: true})
        try {
            const spreadsheets = await externalConnectionsService.listSpreadsheets()
            set({spreadsheets})
        } finally {
            set({spreadsheetsLoading: false})
        }
    },

    addSpreadsheet: async (spreadsheetId: string, name: string) => {
        await externalConnectionsService.addSpreadsheet(spreadsheetId, name)
        await get().fetchSpreadsheets()
    },

    removeSpreadsheet: async (spreadsheetId: string) => {
        await externalConnectionsService.removeSpreadsheet(spreadsheetId)
        await get().fetchSpreadsheets()
    },

    addEmailConnection: async (req: AddEmailConnectionRequest) => {
        await externalConnectionsService.addEmailConnection(req)
        await get().fetch()
    },

    addGitlabConnection: async (req: AddGitlabConnectionRequest) => {
        await externalConnectionsService.addGitlabConnection(req)
        await get().fetch()
    },

    generateGitlabWebhookSecret: async () => {
        const {webhookSecret} = await externalConnectionsService.generateGitlabWebhookSecret()
        await get().fetch()
        return webhookSecret
    },

    addTelegramConnection: async (req: AddTelegramConnectionRequest) => {
        await externalConnectionsService.addTelegramConnection(req)
        await get().fetch()
    },

    addTrelloConnection: async (req: AddTrelloConnectionRequest) => {
        await externalConnectionsService.addTrelloConnection(req)
        await get().fetch()
    },

    addAnthropicConnection: async (req: AddAnthropicConnectionRequest) => {
        await externalConnectionsService.addAnthropicConnection(req)
        await get().fetch()
    },

    checkAnthropicConnection: async (req: CheckAnthropicConnectionRequest) => {
        return externalConnectionsService.checkAnthropicConnection(req)
    },

    addOpenAIConnection: async (req: AddOpenAIConnectionRequest) => {
        await externalConnectionsService.addOpenAIConnection(req)
        await get().fetch()
    },

    checkOpenAIConnection: async (req: CheckOpenAIConnectionRequest) => {
        return externalConnectionsService.checkOpenAIConnection(req)
    },

    addS3Connection: async (req: AddS3ConnectionRequest) => {
        await externalConnectionsService.addS3Connection(req)
        await get().fetch()
    },

    checkS3Connection: async (req: CheckS3ConnectionRequest) => {
        return externalConnectionsService.checkS3Connection(req)
    },

    addCouchDBConnection: async (req: AddCouchDBConnectionRequest) => {
        await externalConnectionsService.addCouchDBConnection(req)
        await get().fetch()
    },

    checkCouchDBConnection: async (req: CheckCouchDBConnectionRequest) => {
        return externalConnectionsService.checkCouchDBConnection(req)
    },

    addPostgresConnection: async (req: AddPostgresConnectionRequest) => {
        await externalConnectionsService.addPostgresConnection(req)
        await get().fetch()
    },

    checkPostgresConnection: async (req: CheckPostgresConnectionRequest) => {
        return externalConnectionsService.checkPostgresConnection(req)
    },

    addGenericConnection: async (req: AddGenericConnectionRequest) => {
        await externalConnectionsService.addGenericConnection(req)
        await get().fetch()
    },

    getProviderStatistics: (provider: ExternalProvider) => {
        return externalConnectionsService.getProviderStatistics(provider)
    },
}))
