import {create} from 'zustand'

import {VaultItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import useUser from "@/hooks/user/User.ts"
import {Errors, GrpcError} from "@/processes/Errors.ts"

interface VaultsState {
    vaults: VaultItem[]
    loading: boolean
    forbidden: boolean
    fetch: () => Promise<void>
    create: (name: string) => Promise<void>
    remove: (id: string) => Promise<void>
}

export const useVaults = create<VaultsState>((set, get) => ({
    vaults: [],
    loading: false,
    forbidden: false,

    fetch: async () => {
        set({loading: true})
        try {
            const {auth} = useUser.getState()
            const res = await VaultsAPI.ListVaults({}, auth.getInitReq())
            set({vaults: res.vaults ?? []})
        } catch (err) {
            if ((err as GrpcError)?.code === Errors.PERMISSION_DENIED) {
                set({forbidden: true})
            }
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