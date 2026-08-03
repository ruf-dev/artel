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
