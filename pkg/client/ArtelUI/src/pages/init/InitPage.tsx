import {useEffect, useRef} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/init/InitPage.module.css"

import useUser from "@/hooks/user/User.ts"
import {AuthService} from "@/processes/Auth.ts"
import {Path} from "@/app/routing/Router.tsx"
import {AuthAPI} from "@/app/api/artel"
import {apiPrefix} from "@/app/api/api.ts"

type TelegramLoginWindow = Window & {
    Telegram?: {
        Login: {
            init: (options: { client_id: string; redirect_uri?: string }, callback: (data: { id_token: string }) => void) => void
            auth: (options: { client_id: string; redirect_uri?: string; request_access: string }, callback: (data: { id_token: string }) => void) => void
        }
    }
}

export default function InitPage() {
    const navigate = useNavigate()
    const {auth, login} = useUser()

    const svc = new AuthService()
    const scriptReadyRef = useRef(false)
    const botIdRef = useRef("")
    const authCallbackRef = useRef<(data: { id_token: string }) => void>(() => {})

    authCallbackRef.current = function (data: { id_token: string }) {
        svc.LoginViaTelegram(data.id_token)
            .then(function (session) {
                login(session)
                navigate(Path.HomePage)
            })
            .catch(function (err: unknown) {
                alert(err instanceof Error ? err.message : "Telegram login failed")
            })
    }

    useEffect(function () {
        const script = document.createElement("script")
        script.src = "https://oauth.telegram.org/js/telegram-login.js?3"
        script.async = true

        AuthAPI.GetConfig({}, apiPrefix()).then(function (cfg) {
            botIdRef.current = cfg.telegramClientId ?? ""
            document.head.appendChild(script)
        })

        script.onload = function () {
            const tg = (window as TelegramLoginWindow).Telegram
            if (tg) {
                tg.Login.init(
                    {client_id: botIdRef.current, redirect_uri: window.location.origin},
                    function (data) { authCallbackRef.current(data) },
                )
                scriptReadyRef.current = true
            }
        }

        return function () {
            if (document.head.contains(script)) document.head.removeChild(script)
            scriptReadyRef.current = false
        }
    }, [])

    function handleLoginClick() {
        const tg = (window as TelegramLoginWindow).Telegram
        if (!tg || !scriptReadyRef.current) return
        tg.Login.auth(
            {client_id: botIdRef.current, redirect_uri: window.location.origin, request_access: "write"},
            function (data) { authCallbackRef.current(data) },
        )
    }


    if (auth.isAuthenticated()) {
        navigate(Path.HomePage)
        return null
    }

    return (
        <div className={cls.Root}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>

                <button className={cls.SubmitBtn} onClick={handleLoginClick}>
                    Sign in with Telegram
                </button>
            </div>
        </div>
    )
}