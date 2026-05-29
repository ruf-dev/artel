import {useEffect, useState} from "react"
import {NavigateFunction, useNavigate} from "react-router-dom"
import {TelegramAuth, TelegramAuthData} from "@vervstack/telegram-auth"

import cls from "@/pages/init/InitPage.module.css"

import useUser from "@/hooks/user/User.ts"
import {AuthService, Session} from "@/processes/Auth.ts"
import {UserInfo} from "@/processes/AuthMiddleware.ts"
import {Path} from "@/app/routing/Router.tsx"
import {AuthAPI} from "@/app/api/artel"
import {apiPrefix} from "@/app/api/api.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {useToaster} from "@vervstack/chures"

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
    const {bake} = useToaster()

    const [botId, setBotId] = useState("")

    useEffect(function () {
        AuthAPI.GetConfig({}, apiPrefix()).then(function (cfg) {
            setBotId(cfg.telegramClientId ?? "")
        }).catch(() => {
        })
    }, [])

    function onSuccessLoginTg(data: TelegramAuthData) {
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
                bake({title: "Telegram login failed", description: err instanceof Error ? err.message : "Telegram login failed", level: "Error"})
            })
            .finally(CloseDialog)
    }


    return (
        <div className={cls.CardContainer}>
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
