import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'

import {workbenchService} from "@/processes/Workbench.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"

export function workbenchQueryKey(vaultId?: string) {
    return ['workbench', vaultId] as const
}

export function workbenchTerminalTabsQueryKey(vaultId?: string) {
    return ['workbench-terminal-tabs', vaultId] as const
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

    const removeMutation = useMutation({
        mutationFn: () => workbenchService.remove(vaultId!),
        onSuccess: () => queryClient.invalidateQueries({queryKey: workbenchQueryKey(vaultId)}),
    })

    return {
        create: createMutation.mutateAsync,
        start: startMutation.mutateAsync,
        stop: stopMutation.mutateAsync,
        remove: removeMutation.mutateAsync,
        removing: removeMutation.isPending,
    }
}

export function useWorkbenchTerminalTabs(vaultId: string | undefined, enabled: boolean) {
    const q = useQuery({
        queryKey: workbenchTerminalTabsQueryKey(vaultId),
        queryFn: () => workbenchService.listTerminalTabs(vaultId!),
        enabled: !!vaultId && enabled,
        refetchOnWindowFocus: true,
        // A natural tmux-session death (e.g. double Ctrl+C) never goes through any of
        // useWorkbenchTerminalTabMutations's create/select/close, so nothing invalidates
        // this query on its own — a short poll is the only way the UI notices the tab
        // count dropped to zero without the user refocusing the browser window.
        refetchInterval: 5000,
    })

    return {
        tabs: q.data ?? [],
        isLoading: q.isLoading,
    }
}

type TerminalTabData = {id: string; name: string; active: boolean}[]

export function useWorkbenchTerminalTabMutations(vaultId?: string) {
    const queryClient = useQueryClient()
    const queryKey = workbenchTerminalTabsQueryKey(vaultId)

    const createMutation = useMutation({
        mutationFn: () => workbenchService.createTerminalTab(vaultId!),
        onSuccess: () => queryClient.invalidateQueries({queryKey}),
    })

    // Optimistic: flips the clicked tab active in the cache immediately, before the
    // network round trip resolves, so the tab strip highlights instantly instead of
    // waiting on the select call *and* a follow-up list refetch. The actual terminal
    // content still catches up once the backend's tmux select-window call lands —
    // this just removes the avoidable extra latency from the tab UI itself.
    const selectMutation = useMutation({
        mutationFn: (tabId: string) => workbenchService.selectTerminalTab(vaultId!, tabId),
        onMutate: async (tabId: string) => {
            await queryClient.cancelQueries({queryKey})
            const previous = queryClient.getQueryData<TerminalTabData>(queryKey)
            queryClient.setQueryData<TerminalTabData>(queryKey, old =>
                old?.map(tab => ({...tab, active: tab.id === tabId})))
            return {previous}
        },
        onError: (_err, _tabId, context) => {
            if (context?.previous) queryClient.setQueryData(queryKey, context.previous)
        },
        onSettled: () => queryClient.invalidateQueries({queryKey}),
    })

    const closeMutation = useMutation({
        mutationFn: (tabId: string) => workbenchService.closeTerminalTab(vaultId!, tabId),
        onMutate: async (tabId: string) => {
            await queryClient.cancelQueries({queryKey})
            const previous = queryClient.getQueryData<TerminalTabData>(queryKey)
            queryClient.setQueryData<TerminalTabData>(queryKey, old => old?.filter(tab => tab.id !== tabId))
            return {previous}
        },
        onError: (_err, _tabId, context) => {
            if (context?.previous) queryClient.setQueryData(queryKey, context.previous)
        },
        onSettled: () => queryClient.invalidateQueries({queryKey}),
    })

    return {
        create: createMutation.mutateAsync,
        select: selectMutation.mutateAsync,
        close: closeMutation.mutateAsync,
    }
}
