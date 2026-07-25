import {create} from 'zustand'

import {
    McpKeyInfo, CreateMcpKeyResponse, McpConnectorInfo, MomCandidate, CommunityConnectorInfo,
} from "@/app/api/artel/mcp_keys.pb.ts"
import {mcpKeysService} from "@/processes/McpKeys.ts"

interface McpKeysState {
    keys: McpKeyInfo[]
    loading: boolean
    connectorsByKey: Record<string, McpConnectorInfo[]>
    connectorsLoading: boolean
    momCandidates: MomCandidate[]
    momCandidatesLoading: boolean
    communityConnectors: CommunityConnectorInfo[]
    communityConnectorsLoading: boolean
    fetch: () => Promise<void>
    create: (name: string, vaultId: string) => Promise<CreateMcpKeyResponse>
    revoke: (keyId: string, vaultId: string) => Promise<void>
    setAccess: (keyId: string, vaultId: string) => Promise<void>
    fetchConnectors: (keyId: string) => Promise<void>
    addConnector: (keyId: string, mcpName: string, externalConnectionId: string) => Promise<void>
    removeConnector: (keyId: string, mcpName: string) => Promise<void>
    fetchMomCandidates: () => Promise<void>
    executeMomTool: (
        mcpName: string, toolName: string, externalConnectionId: string, params: Record<string, unknown>,
    ) => Promise<string>
    fetchCommunityConnectors: () => Promise<void>
    deleteCommunityConnector: (name: string) => Promise<void>
}

export const useMcpKeys = create<McpKeysState>((set, get) => ({
    keys: [],
    loading: false,
    connectorsByKey: {},
    connectorsLoading: false,
    momCandidates: [],
    momCandidatesLoading: false,
    communityConnectors: [],
    communityConnectorsLoading: false,

    fetch: async () => {
        set({loading: true})
        try {
            const keys = await mcpKeysService.list()
            set({keys})
        } finally {
            set({loading: false})
        }
    },

    create: async (name: string, vaultId: string) => {
        const resp = await mcpKeysService.create(name, vaultId)
        await get().fetch()
        return resp
    },

    revoke: async (keyId: string, vaultId: string) => {
        await mcpKeysService.revoke(keyId, vaultId)
        await get().fetch()
    },

    setAccess: async (keyId: string, vaultId: string) => {
        await mcpKeysService.setAccess(keyId, vaultId)
        await get().fetch()
    },

    fetchConnectors: async (keyId: string) => {
        set({connectorsLoading: true})
        try {
            const connectors = await mcpKeysService.listConnectors(keyId)
            set(state => ({connectorsByKey: {...state.connectorsByKey, [keyId]: connectors}}))
        } finally {
            set({connectorsLoading: false})
        }
    },

    addConnector: async (keyId: string, mcpName: string, externalConnectionId: string) => {
        await mcpKeysService.addConnector(keyId, mcpName, externalConnectionId)
        await get().fetchConnectors(keyId)
    },

    removeConnector: async (keyId: string, mcpName: string) => {
        await mcpKeysService.removeConnector(keyId, mcpName)
        await get().fetchConnectors(keyId)
    },

    fetchMomCandidates: async () => {
        set({momCandidatesLoading: true})
        try {
            const momCandidates = await mcpKeysService.listMomCandidates()
            set({momCandidates})
        } finally {
            set({momCandidatesLoading: false})
        }
    },

    executeMomTool: (
        mcpName: string, toolName: string, externalConnectionId: string, params: Record<string, unknown>,
    ) => {
        return mcpKeysService.executeMomTool(mcpName, toolName, externalConnectionId, params)
    },

    fetchCommunityConnectors: async () => {
        set({communityConnectorsLoading: true})
        try {
            const communityConnectors = await mcpKeysService.listCommunityConnectors()
            set({communityConnectors})
        } finally {
            set({communityConnectorsLoading: false})
        }
    },

    deleteCommunityConnector: async (name: string) => {
        await mcpKeysService.deleteCommunityConnector(name)
        await get().fetchCommunityConnectors()
    },
}))
