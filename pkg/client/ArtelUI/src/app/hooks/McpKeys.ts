import {create} from 'zustand'

import {McpKeyInfo, CreateMcpKeyResponse} from "@/app/api/artel/mcp_keys.pb.ts"
import {mcpKeysService} from "@/processes/McpKeys.ts"

interface McpKeysState {
    keys: McpKeyInfo[]
    loading: boolean
    fetch: () => Promise<void>
    create: (name: string, vaultId: string) => Promise<CreateMcpKeyResponse>
    revoke: (keyId: string, vaultId: string) => Promise<void>
    setAccess: (keyId: string, vaultId: string, emailAccountId: string) => Promise<void>
}

export const useMcpKeys = create<McpKeysState>((set, get) => ({
    keys: [],
    loading: false,

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

    setAccess: async (keyId: string, vaultId: string, emailAccountId: string) => {
        await mcpKeysService.setAccess(keyId, vaultId, emailAccountId)
        await get().fetch()
    },
}))
