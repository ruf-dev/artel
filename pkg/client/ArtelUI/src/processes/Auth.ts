import {AuthMiddleware, UserInfo} from "@/processes/AuthMiddleware.ts";
import {
    LoginRequest,
    AuthAPI
} from "@/app/api/artel";

export interface Session {
    token: string
    expiresAt: string
}

export interface IAuthService {
    LoginViaTelegram: (idToken: string) => Promise<Session>
    FetchUserInfo: () => Promise<UserInfo>
}

export class AuthService extends AuthMiddleware implements IAuthService {
    constructor() {
        super();
    }

    async LoginViaTelegram(idToken: string): Promise<Session> {
        const initR = this.getInitReq()

        const r: LoginRequest = {
            telegram: {
                idToken: idToken,
            },
        } as LoginRequest

        return AuthAPI.Login(r, initR).then(res => {
            return {
                token: res.token,
                expiresAt: res.expiresAt,
            } as Session
        })
    }

    async FetchUserInfo(): Promise<UserInfo> {
        const initR = this.getInitReq()
        const res = await AuthAPI.GetMe({}, initR)
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
