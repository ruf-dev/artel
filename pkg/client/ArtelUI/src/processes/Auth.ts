import {UserInfo} from "@/processes/AuthMiddleware.ts";
import {LoginRequest, AuthAPI} from "@/app/api/artel";
import useUser from "@/hooks/user/User.ts"

// Session no longer carries refresh-token/expiry data — the backend manages
// the token lifecycle entirely via httpOnly cookies now. `token` survives
// only because it's still needed in-memory for the unrelated MCP OAuth
// bearer-token handoff (see AuthMiddleware.ts), not for /api/* auth.
export interface Session {
    token: string
}

export interface IAuthService {
    LoginViaTelegram: (idToken: string) => Promise<Session>
    LoginWithPassword: (email: string, password: string) => Promise<Session>
    Register: (email: string, password: string) => Promise<void>
    FetchUserInfo: () => Promise<UserInfo>
}

export class AuthService implements IAuthService {
    async LoginViaTelegram(idToken: string): Promise<Session> {
        const r: LoginRequest = {
            telegram: {idToken},
        } as LoginRequest

        return AuthAPI.Login(r, useUser.getState().auth.getInitReq()).then(res => ({
            token: res.token ?? "",
        }))
    }

    async LoginWithPassword(email: string, password: string): Promise<Session> {
        const r: LoginRequest = {
            password: {email, password},
        } as LoginRequest

        return AuthAPI.Login(r, useUser.getState().auth.getInitReq()).then(res => ({
            token: res.token ?? "",
        }))
    }

    async Register(email: string, password: string): Promise<void> {
        // Registration alone doesn't establish a session on the backend —
        // the caller must follow up with LoginWithPassword to log in.
        return AuthAPI.Register({email, password}, useUser.getState().auth.getInitReq()).then(() => undefined)
    }

    async FetchUserInfo(): Promise<UserInfo> {
        const res = await AuthAPI.GetMe({}, useUser.getState().auth.getInitReq())
        return {
            id: res.id ?? "",
            username: res.username ?? "",
            email: res.email ?? "",
            isAdministrator: res.permissions?.isAdministrator === true,
            hasEmails: res.permissions?.hasEmails === true,
            hasTaskTrackers: res.permissions?.hasTaskTrackers === true,
            photoUrl: res.photoUrl || undefined,
            hasNotes: res.permissions?.hasNotes === true,
        }
    }
}

export const authService = new AuthService()
