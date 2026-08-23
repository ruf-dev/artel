import useUser from "@/hooks/user/User.ts"
import {VaultsAPI} from "@/app/api/artel/vaults.pb.ts"

export interface IWorkbenchService {
    getStatus: (vaultId: string) => Promise<{ exists: boolean; status: string; pendingTerminalAuthLink: string }>
    create: (vaultId: string) => Promise<{ status: string }>
    start: (vaultId: string, authMode: string) => Promise<{ status: string; authMode: string }>
    stop: (vaultId: string) => Promise<{ status: string }>
    remove: (vaultId: string) => Promise<{ status: string }>
    listTerminalTabs: (vaultId: string) => Promise<{id: string; name: string; active: boolean}[]>
    createTerminalTab: (vaultId: string) => Promise<{id: string; name: string; active: boolean}>
    selectTerminalTab: (vaultId: string, tabId: string) => Promise<void>
    closeTerminalTab: (vaultId: string, tabId: string) => Promise<void>
}

export class WorkbenchService implements IWorkbenchService {
    async getStatus(vaultId: string): Promise<{ exists: boolean; status: string; pendingTerminalAuthLink: string }> {
        const res = await VaultsAPI.GetVault({id: vaultId}, useUser.getState().auth.getInitReq())
        return {
            exists: res.workbenchExists ?? false,
            status: res.workbenchStatus ?? "",
            pendingTerminalAuthLink: res.pendingTerminalAuthLink ?? "",
        }
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

    async listTerminalTabs(vaultId: string): Promise<{id: string; name: string; active: boolean}[]> {
        const res = await VaultsAPI.ListWorkbenchTerminalTabs({vaultId}, useUser.getState().auth.getInitReq())
        return (res.tabs ?? []).map(t => ({id: t.id ?? "", name: t.name ?? "", active: t.active ?? false}))
    }

    async createTerminalTab(vaultId: string): Promise<{id: string; name: string; active: boolean}> {
        const res = await VaultsAPI.CreateWorkbenchTerminalTab({vaultId}, useUser.getState().auth.getInitReq())
        const tab = res.tab
        return {id: tab?.id ?? "", name: tab?.name ?? "", active: tab?.active ?? false}
    }

    async selectTerminalTab(vaultId: string, tabId: string): Promise<void> {
        await VaultsAPI.SelectWorkbenchTerminalTab({vaultId, tabId}, useUser.getState().auth.getInitReq())
    }

    async closeTerminalTab(vaultId: string, tabId: string): Promise<void> {
        await VaultsAPI.CloseWorkbenchTerminalTab({vaultId, tabId}, useUser.getState().auth.getInitReq())
    }
}

export const workbenchService = new WorkbenchService()
