import { useEffect, useRef } from "react"

const SCROLL_LOOKUP_MAX_ATTEMPTS = 20

export function useHighlightNote(highlightedPath: string | null) {
    const scrollAreaRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        if (!highlightedPath) return
        let cancelled = false
        let frame: number

        function tryScroll(attempt: number) {
            if (cancelled) return
            const selector = `[data-path="${CSS.escape(highlightedPath ?? "")}"]`
            const el = scrollAreaRef.current?.querySelector<HTMLElement>(selector)
            if (el) {
                el.scrollIntoView({ block: "center", behavior: "smooth" })
                return
            }
            if (attempt < SCROLL_LOOKUP_MAX_ATTEMPTS) {
                frame = requestAnimationFrame(() => tryScroll(attempt + 1))
            }
        }

        frame = requestAnimationFrame(() => tryScroll(0))
        return () => {
            cancelled = true
            cancelAnimationFrame(frame)
        }
    }, [highlightedPath])

    return { scrollAreaRef }
}
