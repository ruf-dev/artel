import {useEffect, useState} from "react"

import {AuthAPI} from "@/app/api/artel"
import {apiPrefix} from "@/app/api/api.ts"
import {useAppConfig} from "@/app/hooks/AppConfig.ts"

const PING_TIMEOUT_MS = 5_000

// Uses the public /api/auth/config endpoint rather than an authenticated one:
// this is a pure liveness check and must never depend on login state, or a
// 401 here gets treated as a session event by the fetch interceptor and
// forces a reload loop for logged-out visitors (see AuthFetchInterceptor.ts).
export async function pingServer(): Promise<boolean> {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), PING_TIMEOUT_MS)
    try {
        const cfg = await AuthAPI.GetConfig({}, {
            ...apiPrefix(),
            signal: controller.signal,
        })
        useAppConfig.getState().setNoAuthEnabled(cfg.noAuthEnabled === true)
        useAppConfig.getState().setCredsEncrypted(cfg.credsEncrypted === true)
        useAppConfig.getState().setSetupCompleted(cfg.setupCompleted === true)
        useAppConfig.getState().setPasswordAuthEnabled(cfg.passwordAuthEnabled === true)
        useAppConfig.getState().setTelegramAuthEnabled(cfg.telegramAuthEnabled === true)
        useAppConfig.getState().setSelfRegisterEnabled(cfg.selfRegisterEnabled === true)
        return true
    } catch (err) {
        // Three cases mean "down":
        // - TypeError = network unreachable (ECONNREFUSED)
        // - AbortError = our timeout fired
        // - a thrown Response = fetch.pb.ts's fetchRequest() throws the raw
        //   Response when r.json() fails to parse the body at all. The real
        //   backend (grpc-gateway generated) always returns JSON, even for
        //   error statuses, so a body that isn't JSON means something other
        //   than the real backend answered — e.g. in local dev, when the Go
        //   backend is down and there's no Vite proxy for /api, the request
        //   resolves against Vite's own dev server, which serves its SPA
        //   index.html fallback (200 OK, text/html).
        // Any other throw = a parsed JSON error body → real backend, just a
        // non-OK HTTP status → server is up.
        if (err instanceof TypeError) return false
        if (err instanceof DOMException && err.name === "AbortError") return false
        if (err instanceof Response) return false
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
        return () => {
            cancelled = true
        }
    }, [])

    return isAvailable
}
