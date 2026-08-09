import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'

import {
    VaultItem,
    VaultMemberInfo,
    VaultInviteItem,
    VaultsAPI,
    EnablePostgresDatabaseResponse,
    GetPostgresDatabaseResponse,
} from "@/app/api/artel/vaults.pb.ts"
import {vaultService} from "@/processes/Vaults.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"

export type {VaultItem, VaultMemberInfo, VaultInviteItem}

export const vaultsQueryKey = ['vaults'] as const


export function useVaults() {
    const {auth} = useUser()

    const q = useQuery({
        queryKey: vaultsQueryKey,
        queryFn: () => vaultService.ListVaults(),
        enabled: auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    return {
        vaults: q.data || [],
        isLoading: q.isLoading,
        error: q.error
    }
}

export function useVaultMutations() {
    const queryClient = useQueryClient()
    const {auth} = useUser()

    const createMutation = useMutation({
        mutationFn: (name: string) => VaultsAPI.CreateVault({name}, auth.getInitReq()),
        onSuccess: () => queryClient.invalidateQueries({queryKey: vaultsQueryKey}),
    })

    const removeMutation = useMutation({
        mutationFn: (id: string) => VaultsAPI.DeleteVault({id}, auth.getInitReq()),
        onSuccess: () => queryClient.invalidateQueries({queryKey: vaultsQueryKey}),
    })

    return {
        create: createMutation.mutateAsync,
        remove: removeMutation.mutateAsync,
        listMembers: async (vaultId: string): Promise<VaultMemberInfo[]> => {
            const res = await VaultsAPI.ListMembers({vaultId}, auth.getInitReq())
            return res.members ?? []
        },
        removeMember: async (vaultId: string, userId: string): Promise<void> => {
            await VaultsAPI.RemoveMember({vaultId, userId}, auth.getInitReq())
        },
        listInvites: async (vaultId: string): Promise<VaultInviteItem[]> => {
            const res = await VaultsAPI.ListInviteLinks({vaultId}, auth.getInitReq())
            return res.invites ?? []
        },
        createInvite: async (vaultId: string, role: string): Promise<VaultInviteItem> => {
            const res = await VaultsAPI.CreateInviteLink({vaultId, role}, auth.getInitReq())
            return res.invite!
        },
        revokeInvite: async (vaultId: string, inviteId: string): Promise<void> => {
            await VaultsAPI.RevokeInviteLink({vaultId, inviteId}, auth.getInitReq())
        },
        setBinaryStorage: async (vaultId: string, useCouchDb: boolean): Promise<void> => {
            await VaultsAPI.SetVaultBinaryStorage({vaultId, useCouchdb: useCouchDb}, auth.getInitReq())
        },
        publish: async (vaultId: string, slug: string): Promise<VaultItem> => {
            const res = await VaultsAPI.PublishVault({vaultId, slug}, auth.getInitReq())
            return res.vault!
        },
        unpublish: async (vaultId: string): Promise<VaultItem> => {
            const res = await VaultsAPI.UnpublishVault({vaultId}, auth.getInitReq())
            return res.vault!
        },
        acceptInvite: async (token: string): Promise<string> => {
            const res = await VaultsAPI.AcceptInvite({token}, auth.getInitReq())
            return res.vaultId ?? ""
        },
        enablePostgresDatabase: async (vaultId: string): Promise<EnablePostgresDatabaseResponse> => {
            return VaultsAPI.EnablePostgresDatabase({vaultId}, auth.getInitReq())
        },
        getPostgresDatabase: async (vaultId: string): Promise<GetPostgresDatabaseResponse> => {
            return VaultsAPI.GetPostgresDatabase({vaultId}, auth.getInitReq())
        },
        disablePostgresDatabase: async (vaultId: string): Promise<void> => {
            await VaultsAPI.DisablePostgresDatabase({vaultId}, auth.getInitReq())
        },
    }
}
