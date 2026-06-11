import useUser from "@/hooks/user/User.ts"
import {McpKeyInfo, McpKeysAPI, CreateMcpKeyResponse} from "@/app/api/artel/mcp_keys.pb.ts"

export interface IMcpKeysService {
    list: () => Promise<McpKeyInfo[]>
    create: (name: string, vaultId: string) => Promise<CreateMcpKeyResponse>
    revoke: (keyId: string, vaultId: string) => Promise<void>
    setAccess: (keyId: string, vaultId: string, emailAccountId: string) => Promise<void>
}

export class McpKeysService implements IMcpKeysService {
    async list(): Promise<McpKeyInfo[]> {
        const res = await McpKeysAPI.ListUserMcpKeys({}, useUser.getState().auth.getInitReq())
        return res.keys ?? []
    }

    async create(name: string, vaultId: string): Promise<CreateMcpKeyResponse> {
        return McpKeysAPI.CreateMcpKey({name, vaultId}, useUser.getState().auth.getInitReq())
    }

    async revoke(keyId: string, vaultId: string): Promise<void> {
        await McpKeysAPI.RevokeMcpKey({keyId, vaultId}, useUser.getState().auth.getInitReq())
    }

    async setAccess(keyId: string, vaultId: string, emailAccountId: string): Promise<void> {
        await McpKeysAPI.SetMcpKeyAccess({keyId, vaultId, emailAccountId}, useUser.getState().auth.getInitReq())
    }
}

export const mcpKeysService = new McpKeysService()
