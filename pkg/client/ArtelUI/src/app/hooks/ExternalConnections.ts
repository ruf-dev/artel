import {create} from 'zustand'

import {
    AddEmailConnectionRequest, AddGitlabConnectionRequest, ExternalConnectionInfo, Spreadsheet,
} from "@/app/api/artel/external_connections.pb.ts"
import {externalConnectionsService} from "@/processes/ExternalConnections.ts"

interface ExternalConnectionsState {
    connections: ExternalConnectionInfo[]
    loading: boolean
    spreadsheets: Spreadsheet[]
    spreadsheetsLoading: boolean
    fetch: () => Promise<void>
    disconnect: (provider: string) => Promise<void>
    initiateGoogleOAuth: () => Promise<string>
    getPickerToken: () => Promise<string>
    fetchSpreadsheets: () => Promise<void>
    addSpreadsheet: (spreadsheetId: string, name: string) => Promise<void>
    removeSpreadsheet: (spreadsheetId: string) => Promise<void>
    addEmailConnection: (req: AddEmailConnectionRequest) => Promise<void>
    addGitlabConnection: (req: AddGitlabConnectionRequest) => Promise<void>
    generateGitlabWebhookSecret: () => Promise<string>
}

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
}))
