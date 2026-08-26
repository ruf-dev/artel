import {useCallback, useState} from "react"

const STORAGE_KEY = "artel.lastUsedSimpleChatModel"

function readLastUsedModel(): string | undefined {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        return raw || undefined
    } catch {
        // Garbage/inaccessible localStorage: default to no remembered model
        // rather than throwing.
        return undefined
    }
}

// Remembers the last OpenRouter model the user picked in Simple Chat, so a new
// chat (or a chat with no model of its own yet) defaults to it instead of
// always falling back to the connection's recommendedDefaultModel. See
// useLikedModels.ts for the identical localStorage-backed-hook shape.
export function useLastUsedModel() {
    const [lastUsedModel, setLastUsedModelState] = useState<string | undefined>(() => readLastUsedModel())

    const setLastUsedModel = useCallback((model: string) => {
        setLastUsedModelState(model)
        try {
            localStorage.setItem(STORAGE_KEY, model)
        } catch {
            // Best-effort persistence only.
        }
    }, [])

    return {lastUsedModel, setLastUsedModel}
}
