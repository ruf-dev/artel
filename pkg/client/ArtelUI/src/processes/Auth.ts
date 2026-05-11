import {apiPrefix, loginTelegram} from "@/app/api/api.ts"

// Session shape derived from auth.proto Login.Response
export interface Session {
    token: string
    expiresAt: string
}

export interface IAuthService {
    LoginViaPass: (email: string, password: string) => Promise<Session>
    Register: (email: string, password: string) => Promise<void>
    LoginViaTelegram: (idToken: string) => Promise<Session>
}

// TODO: replace with generated grpc-web client after running `moti g`
export class AuthService implements IAuthService {
    async LoginViaPass(_email: string, _password: string): Promise<Session> {
        throw new Error("not implemented — run `moti g` to generate gRPC client")
    }

    async Register(_email: string, _password: string): Promise<void> {
        throw new Error("not implemented — run `moti g` to generate gRPC client")
    }

    async LoginViaTelegram(idToken: string): Promise<Session> {
        const resp = await loginTelegram(idToken)
        return {token: resp.token, expiresAt: resp.expiresAt}
    }
}

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
