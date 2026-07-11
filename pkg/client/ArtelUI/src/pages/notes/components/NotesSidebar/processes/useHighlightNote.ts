import { useEffect, useRef, useState } from "react"

const HIGHLIGHT_DURATION_MS = 2700
const SCROLL_LOOKUP_MAX_ATTEMPTS = 20

export function useHighlightNote() {
    const [highlightedPath, setHighlightedPath] = useState<string | null>(null)
    const scrollAreaRef = useRef<HTMLDivElement>(null)

    function highlightNote(path: string) {
        setHighlightedPath(null)
        requestAnimationFrame(() => setHighlightedPath(path))
    }

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

    useEffect(() => {
        if (!highlightedPath) return
        const timer = setTimeout(() => setHighlightedPath(null), HIGHLIGHT_DURATION_MS)
        return () => clearTimeout(timer)
    }, [highlightedPath])

    return { highlightedPath, scrollAreaRef, highlightNote }
}
