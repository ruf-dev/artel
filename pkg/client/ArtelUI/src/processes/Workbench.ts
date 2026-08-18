import useUser from "@/hooks/user/User.ts"
import {VaultsAPI} from "@/app/api/artel/vaults.pb.ts"

export interface IWorkbenchService {
    getStatus: (vaultId: string) => Promise<{ exists: boolean; status: string }>
    create: (vaultId: string) => Promise<{ status: string }>
    start: (vaultId: string, authMode: string) => Promise<{ status: string; authMode: string }>
    stop: (vaultId: string) => Promise<{ status: string }>
    remove: (vaultId: string) => Promise<{ status: string }>
}

export class WorkbenchService implements IWorkbenchService {
    async getStatus(vaultId: string): Promise<{ exists: boolean; status: string }> {
        const res = await VaultsAPI.GetVault({id: vaultId}, useUser.getState().auth.getInitReq())
        return {exists: res.workbenchExists ?? false, status: res.workbenchStatus ?? ""}
    }

    async create(vaultId: string): Promise<{ status: string }> {
        const res = await VaultsAPI.CreateWorkbench({vaultId}, useUser.getState().auth.getInitReq())
        return {status: res.status ?? ""}
    }

    async start(vaultId: string, authMode: string): Promise<{ status: string; authMode: string }> {
        const res = await VaultsAPI.StartWorkbench({vaultId, authMode}, useUser.getState().auth.getInitReq())
        return {status: res.status ?? "", authMode: res.authMode ?? ""}
    }

    async stop(vaultId: string): Promise<{ status: string }> {
        const res = await VaultsAPI.StopWorkbench({vaultId}, useUser.getState().auth.getInitReq())
        return {status: res.status ?? ""}
    }

    async remove(vaultId: string): Promise<{ status: string }> {
        const res = await VaultsAPI.DeleteWorkbench({vaultId}, useUser.getState().auth.getInitReq())
        return {status: res.status ?? ""}
    }
}

export const workbenchService = new WorkbenchService()
