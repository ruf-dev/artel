import {create} from 'zustand'

import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import {externalConnectionsService} from "@/processes/ExternalConnections.ts"

interface ExternalConnectionsState {
    connections: ExternalConnectionInfo[]
    loading: boolean
    fetch: () => Promise<void>
    disconnect: (provider: string) => Promise<void>
    initiateGoogleOAuth: () => Promise<string>
}

export const useExternalConnections = create<ExternalConnectionsState>((set, get) => ({
    connections: [],
    loading: false,

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
        await get().fetch()
    },

    initiateGoogleOAuth: async () => {
        return externalConnectionsService.initiateGoogleOAuth()
    },
}))
