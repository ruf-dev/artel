import {useRef, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/init/InitPage.module.css"

import useUser from "@/hooks/user/User.ts"
import {AuthService} from "@/processes/Auth.ts"
import {Path} from "@/app/routing/Router.tsx"

export default function InitPage() {
    const navigate = useNavigate()
    const {auth, login} = useUser()

    const svc = new AuthService()
    const telegramContainerRef = useRef<HTMLDivElement>(null)

    useEffect(function () {
        if (!telegramContainerRef.current) return

        (window as unknown as Record<string, unknown>).onTelegramAuth = function (data: { id_token: string }) {

            svc.LoginViaTelegram(data.id_token)
                .then(function (session) {
                    login(session)
                    navigate(Path.HomePage)
                })
                .catch(function (err: unknown) {
                    alert(err instanceof Error ? err.message : "Telegram login failed")
                })
        }

        const script = document.createElement("script")

        script.src = "https://oauth.telegram.org/js/telegram-login.js?3"
        script.setAttribute("data-client-id", import.meta.env.VITE_TELEGRAM_BOT_ID ?? "")
        script.setAttribute("data-size", "large")
        script.setAttribute("data-onauth", "onTelegramAuth(data)")
        script.setAttribute("data-request-access", "write phone")
        script.async = true
        telegramContainerRef.current.appendChild(script)


        return function () {
            delete (window as unknown as Record<string, unknown>).onTelegramAuth
        }
    }, [])


    if (auth.isAuthenticated()) {
        navigate(Path.HomePage)
        return null
    }

    return (
        <div className={cls.Root}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>

                <div ref={telegramContainerRef} className={cls.TelegramContainer}/>
            </div>
        </div>
    )
}