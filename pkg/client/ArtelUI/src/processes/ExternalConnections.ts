import {ExternalConnectionInfo, ExternalConnectionsAPI} from "@/app/api/artel/external_connections.pb.ts"
import * as fm from "@/app/api/artel/fetch.pb.ts"
import useUser from "@/hooks/user/User.ts"

export interface IExternalConnectionsService {
    list: () => Promise<ExternalConnectionInfo[]>
    initiateGoogleOAuth: () => Promise<string>
    exchangeGoogleOAuth: (code: string, state: string) => Promise<void>
    disconnect: (provider: string) => Promise<void>
}

export class ExternalConnectionsService implements IExternalConnectionsService {
    async list(): Promise<ExternalConnectionInfo[]> {
        const res = await ExternalConnectionsAPI.ListConnections({}, useUser.getState().auth.getInitReq())
        return res.connections ?? []
    }

    async initiateGoogleOAuth(): Promise<string> {
        const res = await ExternalConnectionsAPI.InitiateGoogleOAuth({origin: window.location.origin}, useUser.getState().auth.getInitReq())
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
}

export const externalConnectionsService = new ExternalConnectionsService()
