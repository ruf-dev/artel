import {useCallback, useEffect, useState} from "react"

const STORAGE_KEY = "artel.likedOpenRouterModels"

function readLikedIds(): string[] {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        if (!raw) {
            return []
        }
        const parsed: unknown = JSON.parse(raw)
        if (!Array.isArray(parsed)) {
            return []
        }
        return parsed.filter((id): id is string => typeof id === "string")
    } catch {
        // Garbage/corrupt localStorage content: default to no liked models
        // rather than throwing.
        return []
    }
}

export function useLikedModels() {
    const [likedIds, setLikedIds] = useState<string[]>(() => readLikedIds())

    useEffect(() => {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(likedIds))
    }, [likedIds])

    const isLiked = useCallback((id: string) => likedIds.includes(id), [likedIds])

    const toggleLiked = useCallback((id: string) => {
        setLikedIds(prev => (prev.includes(id) ? prev.filter(existing => existing !== id) : [...prev, id]))
    }, [])

    return {likedIds, isLiked, toggleLiked}
}
