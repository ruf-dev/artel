import {useEffect, useState} from "react"
import {NavigateFunction} from "react-router-dom"
import {useToaster, TelegramAuth, TelegramAuthData} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog.ts"
import {authService, Session} from "@/processes/Auth.ts"
import {UserInfo} from "@/processes/AuthMiddleware.ts"
import {Path, REDIRECT_AFTER_LOGIN_KEY} from "@/app/routing/Router.tsx"
import {AuthAPI} from "@/app/api/artel"
import {apiPrefix} from "@/app/api/api.ts"
import cls from "@/pages/init/components/LoginContent/LoginContent.module.css"

interface LoginContentProps {
    login: (session: Session, userInfo?: UserInfo) => void
    navigate: NavigateFunction
}

export default function LoginContent({login, navigate}: LoginContentProps) {
    const {CloseDialog, UnlockClosing} = useDialog()
    const {bake} = useToaster()

    const [botId, setBotId] = useState("")

    useEffect(function () {
        AuthAPI.GetConfig({}, apiPrefix()).then(function (cfg) {
            setBotId(cfg.telegramClientId ?? "")
        }).catch(() => {
        })
    }, [])

    function onSuccessLoginTg(data: TelegramAuthData) {
        authService.LoginViaTelegram(data.id_token)
            .then(async function (session) {
                login(session)
                try {
                    const userInfo = await authService.FetchUserInfo()
                    login(session, userInfo)
                    UnlockClosing()
                } catch {
                    // permissions unavailable — proceed without admin flag
                }
                const redirect = localStorage.getItem(REDIRECT_AFTER_LOGIN_KEY)
                if (redirect) {
                    localStorage.removeItem(REDIRECT_AFTER_LOGIN_KEY)
                    navigate(redirect)
                } else {
                    navigate(Path.HomePage)
                }
            })
            .catch(function (err: unknown) {
                bake({title: "Telegram login failed", description: err instanceof Error ? err.message : "Telegram login failed", level: "Error"})
            })
            .finally(CloseDialog)
    }


    return (
        <div className={cls.LoginContentContainer}>
            <div className={cls.Logo}>artel</div>
            <TelegramAuth
                title="Sign in with Telegram"
                botId={botId}
                className={cls.LoginButton}
                onSuccess={onSuccessLoginTg}
            />
        </div>
    )
}
