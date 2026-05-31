import {create} from 'zustand'

import {McpKeyInfo, McpKeysAPI, CreateMcpKeyResponse} from "@/app/api/artel/mcp_keys.pb.ts"
import useUser from "@/hooks/user/User.ts"

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
            const {auth} = useUser.getState()
            const resp = await McpKeysAPI.ListUserMcpKeys({}, auth.getInitReq())
            set({keys: resp.keys ?? []})
        } finally {
            set({loading: false})
        }
    },

    create: async (name: string, vaultId: string) => {
        const {auth} = useUser.getState()
        const resp = await McpKeysAPI.CreateMcpKey({name, vaultId}, auth.getInitReq())
        await get().fetch()
        return resp
    },

    revoke: async (keyId: string, vaultId: string) => {
        const {auth} = useUser.getState()
        await McpKeysAPI.RevokeMcpKey({keyId, vaultId}, auth.getInitReq())
        await get().fetch()
    },

    setAccess: async (keyId: string, vaultId: string, emailAccountId: string) => {
        const {auth} = useUser.getState()
        await McpKeysAPI.SetMcpKeyAccess({keyId, vaultId, emailAccountId}, auth.getInitReq())
        await get().fetch()
    },
}))
