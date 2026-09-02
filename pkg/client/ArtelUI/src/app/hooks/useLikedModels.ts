import {useCallback} from "react"
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query"

import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import {UserSettings, userSettingsQueryKey, userSettingsService} from "@/processes/UserSettings.ts"
import useUser from "@/hooks/user/User.ts"

export function useLikedModels() {
    const {auth} = useUser()
    const bakeError = useBakeError()
    const queryClient = useQueryClient()

    const q = useQuery({
        queryKey: userSettingsQueryKey(),
        queryFn: () => userSettingsService.getUserSettings(),
        enabled: auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    const likedIds = q.data?.likedOpenrouterModels ?? []

    const mutation = useMutation({
        mutationFn: (nextLiked: string[]) => userSettingsService.setLikedModels(nextLiked),
    })

    const isLiked = useCallback((id: string) => likedIds.includes(id), [likedIds])

    const toggleLiked = useCallback((id: string) => {
        const previous = likedIds
        const next = previous.includes(id) ? previous.filter(existing => existing !== id) : [...previous, id]

        queryClient.setQueryData<UserSettings>(userSettingsQueryKey(), current =>
            current ? {...current, likedOpenrouterModels: next} : current)

        mutation.mutate(next, {
            onError: err => {
                queryClient.setQueryData<UserSettings>(userSettingsQueryKey(), current =>
                    current ? {...current, likedOpenrouterModels: previous} : current)
                bakeError("Failed to update liked models", err)
            },
        })
    }, [likedIds, mutation.mutate, queryClient, bakeError])

    return {likedIds, isLiked, toggleLiked}
}
