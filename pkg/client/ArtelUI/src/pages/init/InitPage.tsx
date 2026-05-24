import {useEffect, useRef} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/init/InitPage.module.css"

import useUser from "@/hooks/user/User.ts"
import {AuthService} from "@/processes/Auth.ts"
import {Path} from "@/app/routing/Router.tsx"

type TelegramLoginWindow = Window & {
    Telegram?: {
        Login: {
            init: (botId: string, callback: (data: { id_token: string }) => void) => void
            auth: (options: { request_access: string }, callback: (data: { id_token: string }) => void) => void
        }
    }
}

const BOT_ID = import.meta.env.VITE_TELEGRAM_BOT_ID ?? ""

export default function InitPage() {
    const navigate = useNavigate()
    const {auth, login} = useUser()

    const svc = new AuthService()
    const scriptReadyRef = useRef(false)
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
        script.onload = function () {
            const tg = (window as TelegramLoginWindow).Telegram
            if (tg) {
                tg.Login.init(BOT_ID, function (data) { authCallbackRef.current(data) })
                scriptReadyRef.current = true
            }
        }
        document.head.appendChild(script)

        return function () {
            document.head.removeChild(script)
            scriptReadyRef.current = false
        }
    }, [])

    function handleLoginClick() {
        const tg = (window as TelegramLoginWindow).Telegram
        if (!tg || !scriptReadyRef.current) return
        tg.Login.auth({request_access: "write"}, function (data) { authCallbackRef.current(data) })
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