import {useEffect, useState} from "react"

import {AuthAPI} from "@/app/api/artel"
import useUser from "@/hooks/user/User.ts"

const PING_INTERVAL_MS = 10_000
const PING_TIMEOUT_MS = 5_000

async function pingServer(): Promise<boolean> {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), PING_TIMEOUT_MS)
    try {
        await AuthAPI.GetMe({}, {
            ...useUser.getState().auth.getInitReq(),
            signal: controller.signal,
        })
        return true
    } catch (err) {
        // TypeError = network unreachable (ECONNREFUSED)
        // AbortError = our timeout fired
        // Any other throw = HTTP error body (e.g. 401) → server is up
        if (err instanceof TypeError) return false
        if (err instanceof DOMException && err.name === "AbortError") return false
        return true
    } finally {
        clearTimeout(timeoutId)
    }
}

export function useServerStatus(): boolean {
    // Start false: show the unavailable page until first successful ping.
    // This prevents routes from rendering (and showing a black screen) before
    // we know the server is reachable.
    const [isAvailable, setIsAvailable] = useState(false)

    useEffect(() => {
        let cancelled = false

        async function check() {
            const ok = await pingServer()
            if (!cancelled) setIsAvailable(ok)
        }

        check()
        const id = setInterval(check, PING_INTERVAL_MS)
        return () => {
            cancelled = true
            clearInterval(id)
        }
    }, [])

    return isAvailable
}
