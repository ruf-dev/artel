import {GrpcAuthHeader} from "@/app/api/api.ts"
import {Path} from "@/app/routing/Router.tsx"
import useUser from "@/hooks/user/User.ts"

const REFRESH_PATH = "/api/auth/refresh"
const SKIP_REFRESH_PATHS = ["/api/auth/login", "/api/auth/register", "/api/auth/refresh", "/api/auth/config"]

type RefreshResponseBody = {
    token?: string
    expiresAt?: string
    refreshToken?: string
    refreshExpiresAt?: string
}

const originalFetch = window.fetch.bind(window)

let inFlightRefresh: Promise<boolean> | null = null

function requestUrl(input: RequestInfo | URL): string {
    if (typeof input === "string") return input
    if (input instanceof URL) return input.toString()
    return input.url
}

function refreshTokens(): Promise<boolean> {
    if (inFlightRefresh) return inFlightRefresh

    const auth = useUser.getState().auth
    const refreshToken = auth.session?.refreshToken

    if (!refreshToken) {
        return Promise.resolve(false)
    }

    const refreshInit: RequestInit = {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({refreshToken}),
    }

    inFlightRefresh = originalFetch(REFRESH_PATH, refreshInit)
        .then(response => response.json().then((body: RefreshResponseBody) => ({response, body})))
        .then(({response, body}) => {
            if (!response.ok || !body.token || !body.refreshToken) return false

            auth.updateTokens(body.token, body.expiresAt ?? "", body.refreshToken, body.refreshExpiresAt ?? "")
            return true
        })
        .catch(() => false)
        .finally(() => {
            inFlightRefresh = null
        })

    return inFlightRefresh
}

function forceLogout() {
    const hadSession = Boolean(useUser.getState().auth.session)
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

            const retryHeaders = {
                ...(init?.headers as Record<string, string> | undefined),
                [GrpcAuthHeader]: useUser.getState().auth.getToken(),
            }
            const retryInit: RequestInit = {...init, headers: retryHeaders}

            return originalFetch(input, retryInit)
        })
    })
}
