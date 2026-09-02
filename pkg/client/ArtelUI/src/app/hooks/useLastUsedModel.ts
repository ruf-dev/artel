import {useCallback} from "react"
import {useMutation, useQuery, useQueryClient} from "@tanstack/react-query"

import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {retryOnStatus} from "@/processes/grpcErrors.ts"
import {UserSettings, userSettingsQueryKey, userSettingsService} from "@/processes/UserSettings.ts"
import useUser from "@/hooks/user/User.ts"

// Remembers the last OpenRouter model the user picked in Simple Chat, so a new
// chat (or a chat with no model of its own yet) defaults to it instead of
// always falling back to the connection's recommendedDefaultModel. Reads the
// same GetUserSettings row as useLikedModels.ts under the same query key, so
// the two hooks share one cached fetch instead of double-fetching.
export function useLastUsedModel() {
    const {auth} = useUser()
    const bakeError = useBakeError()
    const queryClient = useQueryClient()

    const q = useQuery({
        queryKey: userSettingsQueryKey(),
        queryFn: () => userSettingsService.getUserSettings(),
        enabled: auth.isAuthenticated(),
        retry: retryOnStatus(),
    })

    const mutation = useMutation({
        mutationFn: (model: string) => userSettingsService.setLastUsedModel(model),
        onError: err => bakeError("Failed to save last used model", err),
    })

    const setLastUsedModel = useCallback((model: string) => {
        queryClient.setQueryData<UserSettings>(userSettingsQueryKey(), current =>
            current ? {...current, lastUsedModel: model} : current)
        mutation.mutate(model)
    }, [queryClient, mutation.mutate])

    return {lastUsedModel: q.data?.lastUsedModel, setLastUsedModel}
}
