import {UserInfo} from "@/processes/AuthMiddleware.ts";
import {LoginRequest, AuthAPI} from "@/app/api/artel";
import useUser from "@/hooks/user/User.ts"

export interface Session {
    token: string
    expiresAt: string
    refreshToken: string
    refreshExpiresAt: string
}

export interface IAuthService {
    LoginViaTelegram: (idToken: string) => Promise<Session>
    FetchUserInfo: () => Promise<UserInfo>
}

export class AuthService implements IAuthService {
    async LoginViaTelegram(idToken: string): Promise<Session> {
        const r: LoginRequest = {
            telegram: {idToken},
        } as LoginRequest

        return AuthAPI.Login(r, useUser.getState().auth.getInitReq()).then(res => ({
            token: res.token,
            expiresAt: res.expiresAt,
            refreshToken: res.refreshToken,
            refreshExpiresAt: res.refreshExpiresAt,
        } as Session))
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
