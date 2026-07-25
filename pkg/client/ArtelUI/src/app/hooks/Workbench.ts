import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'

import {workbenchService} from "@/processes/Workbench.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"

export function workbenchQueryKey(vaultId?: string) {
    return ['workbench', vaultId] as const
}

export function useWorkbench(vaultId?: string) {
    const {auth} = useUser()

    const q = useQuery({
        queryKey: workbenchQueryKey(vaultId),
        queryFn: () => workbenchService.getStatus(vaultId!),
        enabled: !!vaultId && auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    return {
        exists: q.data?.exists ?? false,
        status: q.data?.status ?? "",
        isLoading: q.isLoading,
        refetch: q.refetch,
    }
}

export function useWorkbenchMutations(vaultId?: string) {
    const queryClient = useQueryClient()

    const createMutation = useMutation({
        mutationFn: () => workbenchService.create(vaultId!),
        onSuccess: () => queryClient.invalidateQueries({queryKey: workbenchQueryKey(vaultId)}),
    })

    const startMutation = useMutation({
        mutationFn: (authMode: string) => workbenchService.start(vaultId!, authMode),
        onSuccess: () => queryClient.invalidateQueries({queryKey: workbenchQueryKey(vaultId)}),
    })

    const stopMutation = useMutation({
        mutationFn: () => workbenchService.stop(vaultId!),
        onSuccess: () => queryClient.invalidateQueries({queryKey: workbenchQueryKey(vaultId)}),
    })

    return {
        create: createMutation.mutateAsync,
        start: startMutation.mutateAsync,
        stop: stopMutation.mutateAsync,
    }
}
