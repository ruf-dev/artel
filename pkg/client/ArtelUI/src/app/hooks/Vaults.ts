import {create} from 'zustand'

import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import {VaultService} from "@/processes/Vaults.ts"
import useUser from "@/hooks/user/User.ts"

interface VaultsState {
    vaults: VaultItem[]
    loading: boolean
    fetch: () => Promise<void>
    create: (name: string) => Promise<void>
    remove: (id: string) => Promise<void>
}

const vaultService = new VaultService()

export const useVaults = create<VaultsState>((set, get) => ({
    vaults: [],
    loading: false,

    fetch: async () => {
        set({loading: true})
        try {
            const vaults = await vaultService.ListVaults()
            set({vaults})
        } finally {
            set({loading: false})
        }
    },

    create: async (name: string) => {
        const {auth} = useUser.getState()
        await VaultsAPI.CreateVault({name}, auth.getInitReq())
        await get().fetch()
    },

    remove: async (id: string) => {
        const {auth} = useUser.getState()
        await VaultsAPI.DeleteVault({id}, auth.getInitReq())
        await get().fetch()
    },
}))