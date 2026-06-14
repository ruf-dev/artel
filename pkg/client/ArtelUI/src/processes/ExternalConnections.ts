import {ExternalConnectionInfo, ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import useUser from "@/hooks/user/User.ts"

export interface IExternalConnectionsService {
    list: () => Promise<ExternalConnectionInfo[]>
    initiateGoogleOAuth: () => Promise<string>
    disconnect: (provider: string) => Promise<void>
}

export class ExternalConnectionsService implements IExternalConnectionsService {
    async list(): Promise<ExternalConnectionInfo[]> {
        const res = await ExternalConnectionsAPI.ListConnections({}, useUser.getState().auth.getInitReq())
        return res.connections ?? []
    }

    async initiateGoogleOAuth(): Promise<string> {
        const res = await ExternalConnectionsAPI.InitiateGoogleOAuth({}, useUser.getState().auth.getInitReq())
        return res.authUrl ?? ""
    }

    async disconnect(provider: string): Promise<void> {
        await ExternalConnectionsAPI.DisconnectProvider({provider}, useUser.getState().auth.getInitReq())
    }
}

export const externalConnectionsService = new ExternalConnectionsService()
