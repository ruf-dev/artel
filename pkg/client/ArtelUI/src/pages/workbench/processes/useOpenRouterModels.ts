import {useCallback, useMemo} from "react"
import {useQuery} from "@tanstack/react-query"

import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useLikedModels} from "@/app/hooks/useLikedModels.ts"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {
    buildLatestModelFlags,
    filterModelTree,
    groupModelsByFamily,
    sortLikedToTop,
} from "@/processes/groupModelsByFamily.ts"

// Simple Chat's model-switcher and its "pick initial model" create-chat form both
// need the user's OpenRouter BYOK connection's available model list. There is no
// RPC that returns this for an already-saved connection by id —
// CheckOpenAIConnectionRequest's apiKey/baseUrl fields exist for pre-save
// verification of a not-yet-connected key (see LlmKeyCheckButton.tsx /
// ManageOpenRouterDialog.tsx's "Test" flow) — so this calls it with only
// `provider` set and no apiKey/baseUrl, relying on the backend falling back to
// the already-saved connection's stored credentials for that provider when no key
// is supplied in the request. Flagging this as an assumption to cross-check
// against the backend track's implementation of the same RPC.
export function useOpenRouterModels(enabled: boolean) {
    const {checkOpenAIConnection} = useExternalConnections()
    const {likedIds, isLiked, toggleLiked} = useLikedModels()

    const q = useQuery({
        queryKey: ["openrouter-models"] as const,
        queryFn: () => checkOpenAIConnection({provider: ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER}),
        enabled,
        staleTime: 5 * 60 * 1000,
    })

    const models = q.data?.availableModels ?? []

    const groupedModels = useMemo(
        () => sortLikedToTop(groupModelsByFamily(models), likedIds),
        [models, likedIds],
    )

    const latestFlags = useMemo(() => buildLatestModelFlags(models), [models])

    const searchModels = useCallback(
        (query: string) => Promise.resolve(filterModelTree(groupedModels, query)),
        [groupedModels],
    )

    return {
        models,
        groupedModels,
        searchModels,
        recommendedDefaultModel: q.data?.recommendedDefaultModel,
        isLoading: q.isLoading,
        error: q.error,
        latestFlags,
        isLiked,
        toggleLiked,
    }
}
