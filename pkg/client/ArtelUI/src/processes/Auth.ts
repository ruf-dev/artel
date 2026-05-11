import {AuthMiddleware} from "@/processes/AuthMiddleware.ts";
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

        console.debug(r, initR)

        return AuthAPI.Login(r, initR).then(res => {
            return {
                token: res.token,
                expiresAt: res.expiresAt,
            } as Session
        })
    }
}

