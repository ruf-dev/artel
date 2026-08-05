export interface TelegramLoginResponse {
    token: string
    expiresAt: string
}

const CSRF_COOKIE_NAME = "csrf_token"
const CSRF_HEADER_NAME = "Grpc-Metadata-X-Csrf-Token"

export interface InitReq {
    pathPrefix: string,
    credentials: RequestCredentials,
    headers: {
        "Content-Type"?: string,
        "Grpc-Metadata-X-Csrf-Token"?: string,
    },
}

// getCsrfToken reads the non-httpOnly, root-scoped csrf_token cookie set by
// the backend on Login/Refresh (and cleared on Logout). It's the only
// auth-related cookie readable from JS by design — access/refresh tokens are
// httpOnly and never touch document.cookie.
export function getCsrfToken(): string | undefined {
    const match = document.cookie.match(new RegExp(`(?:^|; )${CSRF_COOKIE_NAME}=([^;]*)`))
    return match ? decodeURIComponent(match[1]) : undefined
}

// clearCsrfCookie expires the csrf_token cookie immediately, client-side. Needed because
// forceLogout() (AuthFetchInterceptor.ts) does a synchronous hard redirect right after firing
// AuthAPI.Logout() — that call clears the cookie server-side via Set-Cookie, but is async and can
// lose the race against the redirect, leaving the cookie behind. isAuthenticated() then still
// reads true on the next page load, the next authenticated call 401s again, and forceLogout()
// fires again — an infinite reload loop. Clearing the cookie here removes that race entirely.
export function clearCsrfCookie(): void {
    document.cookie = `${CSRF_COOKIE_NAME}=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`
}

function csrfHeader(): { "Grpc-Metadata-X-Csrf-Token"?: string } {
    const token = getCsrfToken()
    return token ? {[CSRF_HEADER_NAME]: token} : {}
}

export function apiPrefix(): InitReq {
    return {
        pathPrefix: "",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            ...csrfHeader(),
        }
    }
}
