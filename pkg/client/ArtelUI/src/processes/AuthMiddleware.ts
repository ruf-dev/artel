import {apiPrefix, clearCsrfCookie, getCsrfToken} from "@/app/api/api.ts"
import {Session} from "@/processes/Auth.ts";
import {useAppConfig} from "@/app/hooks/AppConfig.ts";
import {AuthAPI} from "@/app/api/artel"

export type UserInfo = {
    id: string
    username: string
    email: string
    isAdministrator: boolean
    hasEmails: boolean
    hasTaskTrackers: boolean
    hasNotes: boolean
    photoUrl?: string
}

type StoredAuth = {
    userInfo?: UserInfo
}

export class AuthMiddleware {
    // In-memory only, never persisted: the access token itself now travels
    // as an httpOnly cookie the browser attaches automatically, so JS no
    // longer needs to hold it for authenticating /api/* calls. It's kept
    // here only for the unrelated MCP OAuth bearer-token handoff
    // (McpAuthPage.tsx passes it to the separate /oauth/vaults endpoint),
    // which is out of scope for this migration — it now only works for the
    // lifetime of the current page load rather than surviving a reload.
    session?: Session
    userInfo?: UserInfo

    constructor() {
        const stored = fromLocalStorage()
        this.userInfo = stored?.userInfo
    }

    // isAuthenticated is a synchronous, UX-only "was I logged in" signal —
    // access/refresh tokens are httpOnly and unreadable from JS, so this
    // checks for the non-httpOnly, root-scoped csrf_token cookie instead,
    // which the backend sets alongside the auth cookies on Login/Refresh and
    // clears on Logout. Real authorization is still enforced server-side on
    // every request regardless of what this returns; it only avoids a
    // login-page flash before an inevitable 401 redirect.
    isAuthenticated(): boolean {
        if (useAppConfig.getState().noAuthEnabled) return true
        return Boolean(getCsrfToken())
    }

    login(s: Session, info?: UserInfo) {
        this.session = s
        this.userInfo = info

        saveToLocalStorage({userInfo: info})
    }

    setUserInfo(info: UserInfo) {
        this.userInfo = info
        saveToLocalStorage({userInfo: info})
    }

    logout() {
        this.session = undefined
        this.userInfo = undefined
        clearLocalStorage()
        // Drop the readable csrf_token cookie now, synchronously — the server
        // clears it too, but AuthAPI.Logout() below is an async round-trip that
        // loses the race against the post-logout navigation, leaving
        // isAuthenticated() reading true (so InitPage bounces straight back to
        // the app). Same rationale as forceLogout()'s clearCsrfCookie() call.
        clearCsrfCookie()
        // Best-effort: clears the httpOnly auth cookies + csrf_token server
        // side. Local state is cleared regardless of whether this succeeds.
        AuthAPI.Logout({}, apiPrefix()).catch(() => {})
    }

    getToken(): string {
        if (!this.session) return ""
        return this.session.token
    }

    isAdmin(): boolean {
        return this.userInfo?.isAdministrator === true
    }

    hasEmailsPermission(): boolean {
        return this.userInfo?.hasEmails === true
    }

    hasTaskTrackersPermission(): boolean {
        return this.userInfo?.hasTaskTrackers === true
    }

    hasNotesPermission(): boolean {
        return this.userInfo?.hasNotes === true
    }

    getPhotoUrl(): string | undefined {
        return this.userInfo?.photoUrl || undefined
    }

    getInitReq() {
        return apiPrefix()
    }
}

function saveToLocalStorage(stored: StoredAuth) {
    localStorage.setItem("artel_session", JSON.stringify(stored))
}

function fromLocalStorage(): StoredAuth | undefined {
    const raw = localStorage.getItem("artel_session")
    if (!raw) return undefined
    try {
        const parsed = JSON.parse(raw)
        // Pre-migration browsers may still have the old {session, userInfo}
        // shape cached; only the userInfo part survives, tokens are dropped.
        if (parsed && typeof parsed === "object" && "userInfo" in parsed) {
            return {userInfo: parsed.userInfo as UserInfo | undefined}
        }
        return undefined
    } catch {
        return undefined
    }
}

function clearLocalStorage() {
    localStorage.removeItem("artel_session")
}
