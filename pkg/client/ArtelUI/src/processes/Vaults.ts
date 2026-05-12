import {AuthMiddleware} from "@/processes/AuthMiddleware.ts"
import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"

export interface IVaultService {
    ListVaults: () => Promise<VaultItem[]>
}

export class VaultService extends AuthMiddleware implements IVaultService {
    constructor() {
        super()
    }

    async ListVaults(): Promise<VaultItem[]> {
        const initR = this.getInitReq()
        return VaultsAPI.ListVaults({}, initR).then(res => res.vaults ?? [])
    }
}
