import {apiPrefix} from "@/app/api/api.ts"
import {Session} from "@/processes/Auth.ts";

export class AuthMiddleware {
    session?: Session

    constructor(session?: Session) {
        this.session = session ?? fromLocalStorage()
    }

    isAuthenticated(): boolean {
        if (!this.session) return false
        return new Date(this.session.expiresAt) > new Date()
    }

    login(s: Session) {
        this.session = s
        saveToLocalStorage(s)
    }

    logout() {
        this.session = undefined
        clearLocalStorage()
    }

    getToken(): string {
        if (!this.session) return ""
        return this.session.token
    }

    getInitReq() {
        return apiPrefix({accessToken: this.getToken()})
    }
}

function saveToLocalStorage(session: Session) {
    localStorage.setItem("artel_session", JSON.stringify(session))
}

function fromLocalStorage(): Session | undefined {
    const raw = localStorage.getItem("artel_session")
    if (!raw) return undefined
    return JSON.parse(raw)
}

function clearLocalStorage() {
    localStorage.removeItem("artel_session")
}
