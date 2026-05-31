import {create} from 'zustand'

import {VaultItem, VaultMemberInfo, VaultInviteItem, VaultsAPI} from "@/app/api/artel/vaults.pb.ts"
import useUser from "@/hooks/user/User.ts"
import {Errors, GrpcError} from "@/processes/Errors.ts"

export type {VaultMemberInfo, VaultInviteItem}

interface VaultsState {
    vaults: VaultItem[]
    loading: boolean
    forbidden: boolean
    fetch: () => Promise<void>
    create: (name: string) => Promise<void>
    remove: (id: string) => Promise<void>
    listMembers: (vaultId: string) => Promise<VaultMemberInfo[]>
    removeMember: (vaultId: string, userId: string) => Promise<void>
    listInvites: (vaultId: string) => Promise<VaultInviteItem[]>
    createInvite: (vaultId: string, role: string) => Promise<VaultInviteItem>
    revokeInvite: (vaultId: string, inviteId: string) => Promise<void>
    acceptInvite: (token: string) => Promise<string>
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

    listMembers: async (vaultId: string) => {
        const {auth} = useUser.getState()
        const res = await VaultsAPI.ListMembers({vaultId}, auth.getInitReq())
        return res.members ?? []
    },

    removeMember: async (vaultId: string, userId: string) => {
        const {auth} = useUser.getState()
        await VaultsAPI.RemoveMember({vaultId, userId}, auth.getInitReq())
    },

    listInvites: async (vaultId: string) => {
        const {auth} = useUser.getState()
        const res = await VaultsAPI.ListInviteLinks({vaultId}, auth.getInitReq())
        return res.invites ?? []
    },

    createInvite: async (vaultId: string, role: string) => {
        const {auth} = useUser.getState()
        const res = await VaultsAPI.CreateInviteLink({vaultId, role}, auth.getInitReq())
        return res.invite!
    },

    revokeInvite: async (vaultId: string, inviteId: string) => {
        const {auth} = useUser.getState()
        await VaultsAPI.RevokeInviteLink({vaultId, inviteId}, auth.getInitReq())
    },

    acceptInvite: async (token: string) => {
        const {auth} = useUser.getState()
        const res = await VaultsAPI.AcceptInvite({token}, auth.getInitReq())
        return res.vaultId ?? ""
    },
}))