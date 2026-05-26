import {useEffect, useState} from "react"
import {useNavigate} from "react-router-dom"
import {TelegramAuth} from "telegram-auth"

import cls from "@/pages/init/InitPage.module.css"

import useUser from "@/hooks/user/User.ts"
import {AuthService} from "@/processes/Auth.ts"
import {Path} from "@/app/routing/Router.tsx"
import {AuthAPI} from "@/app/api/artel"
import {apiPrefix} from "@/app/api/api.ts"

export default function InitPage() {
    const navigate = useNavigate()
    const {auth, login} = useUser()
    const svc = new AuthService()
    const [botId, setBotId] = useState("")

    useEffect(function () {
        AuthAPI.GetConfig({}, apiPrefix()).then(function (cfg) {
            setBotId(cfg.telegramClientId ?? "")
        })
    }, [])

    if (auth.isAuthenticated()) {
        navigate(Path.HomePage)
        return null
    }

    return (
        <div className={cls.Root}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>

                <TelegramAuth
                    botId={botId}
                    onSuccess={function (data) {
                        svc.LoginViaTelegram(data.id_token)
                            .then(function (session) {
                                login(session)
                                navigate(Path.HomePage)
                            })
                            .catch(function (err: unknown) {
                                alert(err instanceof Error ? err.message : "Telegram login failed")
                            })
                    }}
                >
                    {function ({login: telegramLogin, isReady}) {
                        return (
                            <button className={cls.SubmitBtn} onClick={telegramLogin} disabled={!isReady}>
                                Sign in with Telegram
                            </button>
                        )
                    }}
                </TelegramAuth>
            </div>
        </div>
    )
}