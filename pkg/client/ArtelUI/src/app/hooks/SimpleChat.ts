import {useMutation, useQuery, useQueryClient} from '@tanstack/react-query'

import {simpleChatService} from "@/processes/SimpleChat.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import useUser from "@/hooks/user/User.ts"

export function simpleChatsQueryKey(vaultId?: string) {
    return ['simple-chats', vaultId] as const
}

export function simpleChatQueryKey(chatId?: string) {
    return ['simple-chat', chatId] as const
}

export function useSimpleChats(vaultId?: string) {
    const {auth} = useUser()

    const q = useQuery({
        queryKey: simpleChatsQueryKey(vaultId),
        queryFn: () => simpleChatService.listChats(vaultId!),
        enabled: !!vaultId && auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    return {
        chats: q.data ?? [],
        isLoading: q.isLoading,
        refetch: q.refetch,
    }
}

export function useSimpleChat(chatId?: string) {
    const {auth} = useUser()

    const q = useQuery({
        queryKey: simpleChatQueryKey(chatId),
        queryFn: () => simpleChatService.getChat(chatId!),
        enabled: !!chatId && auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    return {
        chat: q.data?.chat,
        messages: q.data?.messages ?? [],
        isLoading: q.isLoading,
        refetch: q.refetch,
    }
}

export function useSimpleChatMutations(vaultId?: string) {
    const queryClient = useQueryClient()

    const createMutation = useMutation({
        mutationFn: ({model, vaultAccess}: { model: string; vaultAccess: boolean }) =>
            simpleChatService.createChat(vaultId!, model, vaultAccess),
        onSuccess: () => queryClient.invalidateQueries({queryKey: simpleChatsQueryKey(vaultId)}),
    })

    const deleteMutation = useMutation({
        mutationFn: (chatId: string) => simpleChatService.deleteChat(chatId),
        onSuccess: () => queryClient.invalidateQueries({queryKey: simpleChatsQueryKey(vaultId)}),
    })

    return {
        create: createMutation.mutateAsync,
        delete: deleteMutation.mutateAsync,
        creating: createMutation.isPending,
        deleting: deleteMutation.isPending,
    }
}
