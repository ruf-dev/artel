import {useEffect, useState} from "react"
import {NavigateFunction, useNavigate} from "react-router-dom"
import {TelegramAuth} from "@vervstack/telegram-auth"

import cls from "@/pages/init/InitPage.module.css"

import useUser from "@/hooks/user/User.ts"
import {AuthService, Session} from "@/processes/Auth.ts"
import {UserInfo} from "@/processes/AuthMiddleware.ts"
import {Path} from "@/app/routing/Router.tsx"
import {AuthAPI} from "@/app/api/artel"
import {apiPrefix} from "@/app/api/api.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"

export default function InitPage() {
    const navigate = useNavigate()
    const {auth, login} = useUser()
    const {OpenDialog, LockClosing} = useDialog()

    useEffect(function () {
        if (auth.isAuthenticated()) {
            navigate(Path.HomePage)
            return
        }
        LockClosing()
        OpenDialog(<LoginContent login={login} navigate={navigate}/>)
    }, [])

    return <div className={cls.Root}/>
}

interface LoginContentProps {
    login: (session: Session, userInfo?: UserInfo) => void
    navigate: NavigateFunction
}

function LoginContent({login, navigate}: LoginContentProps) {
    const svc = new AuthService()
    const {CloseDialog, UnlockClosing} = useDialog()

    const [botId, setBotId] = useState("")

    useEffect(function () {
        AuthAPI.GetConfig({}, apiPrefix()).then(function (cfg) {
            setBotId(cfg.telegramClientId ?? "")
        }).catch(() => {
        })
    }, [])

    return (
        <div className={cls.Card}>
            <div className={cls.Logo}>artel</div>
            <TelegramAuth
                botId={botId}
                onSuccess={function (data) {
                    svc.LoginViaTelegram(data.id_token)
                        .then(async function (session) {
                            login(session)
                            try {
                                svc.login(session)
                                const userInfo = await svc.FetchUserInfo()
                                login(session, userInfo)
                                UnlockClosing()
                            } catch {
                                // permissions unavailable — proceed without admin flag
                            }
                            navigate(Path.HomePage)
                        })
                        .catch(function (err: unknown) {
                            alert(err instanceof Error ? err.message : "Telegram login failed")
                        })
                        .finally(CloseDialog)
                }}
            >
                {TelegramAuthButton}
            </TelegramAuth>
        </div>
    )
}

interface TgAuthProps {
    login: () => void
    isReady: boolean
}

function TelegramAuthButton({login: telegramLogin, isReady}: TgAuthProps) {
    return (
        <button
            className={cls.SubmitBtn}
            onClick={telegramLogin}
            disabled={!isReady}
        >
            Sign in with Telegram
        </button>
    )
}
