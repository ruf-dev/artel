import {apiPrefix} from "@/app/api/api.ts"
import {Session} from "@/processes/Auth.ts";

export type UserInfo = {
    id: string
    username: string
    email: string
    isAdministrator: boolean
    hasEmails: boolean
    hasTaskTrackers: boolean
    photoUrl?: string
}

type StoredAuth = {
    session: Session
    userInfo?: UserInfo
}

export class AuthMiddleware {
    session?: Session
    userInfo?: UserInfo

    constructor() {
        const stored = fromLocalStorage()
        this.session = stored?.session
        this.userInfo = stored?.userInfo
    }

    isAuthenticated(): boolean {
        if (!this.session) return false
        return new Date(this.session.expiresAt) > new Date()
    }

    login(s: Session, info?: UserInfo) {
        this.session = s
        this.userInfo = info

        saveToLocalStorage({session: s, userInfo: info})
    }

    setUserInfo(info: UserInfo) {
        this.userInfo = info
        if (this.session) {
            saveToLocalStorage({session: this.session, userInfo: info})
        }
    }

    logout() {
        this.session = undefined
        this.userInfo = undefined
        clearLocalStorage()
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

    getPhotoUrl(): string | undefined {
        return this.userInfo?.photoUrl || undefined
    }

    getInitReq() {
        return apiPrefix({accessToken: this.getToken()})
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
        if (parsed.session) return parsed as StoredAuth
        // legacy: old format stored session directly
        return {session: parsed}
    } catch {
        return undefined
    }
}

function clearLocalStorage() {
    localStorage.removeItem("artel_session")
}
