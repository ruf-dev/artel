import useUser from "@/hooks/user/User.ts"
import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"

export interface IVaultService {
    ListVaults: () => Promise<VaultItem[]>
}

export class VaultService implements IVaultService {
    async ListVaults(): Promise<VaultItem[]> {
        const initR = useUser.getState().auth.getInitReq()
        return VaultsAPI.ListVaults({}, initR).then(res => res.vaults ?? [])
    }
}

export const vaultService = new VaultService()
