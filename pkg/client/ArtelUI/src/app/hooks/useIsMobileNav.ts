import { useEffect, useState } from "react"

const MOBILE_NAV_QUERY = "(orientation: portrait), (width <= 767px)"

export function useIsMobileNav(): boolean {
    const [isMobile, setIsMobile] = useState(
        () => window.matchMedia(MOBILE_NAV_QUERY).matches
    )

    useEffect(() => {
        const mq = window.matchMedia(MOBILE_NAV_QUERY)
        function handler(e: MediaQueryListEvent) {
            setIsMobile(e.matches)
        }
        mq.addEventListener("change", handler)
        return () => mq.removeEventListener("change", handler)
    }, [])

    return isMobile
}
