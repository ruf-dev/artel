import { useEffect, useState } from "react"

export function usePortrait(): boolean {
    const [portrait, setPortrait] = useState(
        () => window.matchMedia("(orientation: portrait)").matches
    )

    useEffect(() => {
        const mq = window.matchMedia("(orientation: portrait)")
        function handler(e: MediaQueryListEvent) {
            setPortrait(e.matches)
        }
        mq.addEventListener("change", handler)
        return () => mq.removeEventListener("change", handler)
    }, [])

    return portrait
}
