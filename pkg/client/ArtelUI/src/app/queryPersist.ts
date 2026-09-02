import type {Query} from "@tanstack/react-query"

import {getStoredUserId, queryCacheKeyPrefix} from "@/processes/AuthMiddleware.ts"

// react-query cache persistence (main.tsx). Only the query keys listed here are
// written to localStorage and rehydrated on the next page load, so a refresh
// paints the last-known data immediately and refetches in the background.
// First segment of each queryKey — keep in sync with the hooks that own them:
//   "vaults" -> vaultsQueryKey in app/hooks/Vaults.ts
const persistedFirstSegments = new Set<string>(["vaults"])

// Bump when the persisted shape of any whitelisted query changes in a way an
// old cached blob can't be trusted for; maxAge still bounds staleness and the
// background refetch always corrects the data, this is just belt-and-braces.
export const queryCacheBuster = "v1"

export const queryCacheMaxAge = 1000 * 60 * 60 * 24

export function queryCacheStorageKey(): string {
    return `${queryCacheKeyPrefix}_${getStoredUserId() ?? "anon"}`
}

export function shouldPersistQuery(query: Query): boolean {
    const first = query.queryKey[0]
    return typeof first === "string" && persistedFirstSegments.has(first)
}
