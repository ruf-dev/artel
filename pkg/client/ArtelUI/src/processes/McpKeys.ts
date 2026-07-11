import useUser from "@/hooks/user/User.ts"
import {
    McpKeyInfo,
    McpKeysAPI,
    CreateMcpKeyResponse,
    McpConnectorInfo,
    MomCandidate,
} from "@/app/api/artel/mcp_keys.pb.ts"

export interface IMcpKeysService {
    list: () => Promise<McpKeyInfo[]>
    create: (name: string, vaultId: string) => Promise<CreateMcpKeyResponse>
    revoke: (keyId: string, vaultId: string) => Promise<void>
    setAccess: (keyId: string, vaultId: string) => Promise<void>
    listConnectors: (keyId: string) => Promise<McpConnectorInfo[]>
    addConnector: (keyId: string, mcpName: string, externalConnectionId: string) => Promise<McpConnectorInfo>
    removeConnector: (keyId: string, mcpName: string) => Promise<void>
    listMomCandidates: () => Promise<MomCandidate[]>
    executeMomTool: (
        mcpName: string, toolName: string, externalConnectionId: string, params: Record<string, unknown>,
    ) => Promise<string>
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

    async setAccess(keyId: string, vaultId: string): Promise<void> {
        await McpKeysAPI.SetMcpKeyAccess({keyId, vaultId}, useUser.getState().auth.getInitReq())
    }

    async listConnectors(keyId: string): Promise<McpConnectorInfo[]> {
        const res = await McpKeysAPI.ListMcpConnectors({keyId}, useUser.getState().auth.getInitReq())
        return res.connectors ?? []
    }

    async addConnector(keyId: string, mcpName: string, externalConnectionId: string): Promise<McpConnectorInfo> {
        const res = await McpKeysAPI.AddMcpConnector(
            {keyId, mcpName, externalConnectionId}, useUser.getState().auth.getInitReq(),
        )
        return res.connector ?? {}
    }

    async removeConnector(keyId: string, mcpName: string): Promise<void> {
        await McpKeysAPI.RemoveMcpConnector({keyId, mcpName}, useUser.getState().auth.getInitReq())
    }

    async listMomCandidates(): Promise<MomCandidate[]> {
        const res = await McpKeysAPI.ListMomCandidates({}, useUser.getState().auth.getInitReq())
        return res.candidates ?? []
    }

    async executeMomTool(
        mcpName: string, toolName: string, externalConnectionId: string, params: Record<string, unknown>,
    ): Promise<string> {
        const req = {mcpName, toolName, externalConnectionId, paramsJson: JSON.stringify(params)}
        const res = await McpKeysAPI.ExecuteMomTool(req, useUser.getState().auth.getInitReq())
        return res.result ?? ""
    }
}

export const mcpKeysService = new McpKeysService()
