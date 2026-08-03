import {apiPrefix} from "@/app/api/api.ts"
import {AuthAPI} from "@/app/api/artel"
import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"

const SKIP_REFRESH_PATHS = [
    "/api/auth/login", "/api/auth/register", "/api/auth/refresh", "/api/auth/config", "/api/auth/logout",
]

const originalFetch = window.fetch.bind(window)

let inFlightRefresh: Promise<boolean> | null = null

function requestUrl(input: RequestInfo | URL): string {
    if (typeof input === "string") return input
    if (input instanceof URL) return input.toString()
    return input.url
}

function refreshTokens(): Promise<boolean> {
    if (inFlightRefresh) return inFlightRefresh

    // No body needed: the refresh token now travels via its own httpOnly
    // cookie, read server-side. A successful response updates the auth
    // cookies in place via Set-Cookie — nothing for the client to store.
    inFlightRefresh = AuthAPI.Refresh({}, apiPrefix())
        .then(() => true)
        .catch(() => false)
        .finally(() => {
            inFlightRefresh = null
        })

    return inFlightRefresh
}

function forceLogout() {
    const hadSession = useUser.getState().auth.isAuthenticated()
    useUser.getState().logout()

    // Nothing to force: there was no session to lose, or we're already on
    // the login page — a hard redirect here would just reload in place and,
    // repeated on every subsequent 401, loop forever.
    if (!hadSession || window.location.pathname === Path.InitPage) {
        return
    }

    window.location.href = Path.InitPage
}

window.fetch = function authAwareFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
    const url = requestUrl(input)

    return originalFetch(input, init).then(response => {
        if (response.status !== 401 || SKIP_REFRESH_PATHS.includes(url)) {
            return response
        }

        return refreshTokens().then(refreshed => {
            if (!refreshed) {
                forceLogout()
                return response
            }

            // Cookies are already updated in place by Set-Cookie on the
            // refresh response — the retried request needs nothing more
            // than what it already sent (credentials: 'include' flows
            // through apiPrefix() on every request).
            return originalFetch(input, init)
        })
    })
}
