import {useEffect, useState} from "react"
import {NavigateFunction} from "react-router-dom"
import {useToaster, TelegramAuth, TelegramAuthData, Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog.ts"
import {useAppConfig} from "@/app/hooks/AppConfig.ts"
import {authService, Session} from "@/processes/Auth.ts"
import {UserInfo} from "@/processes/AuthMiddleware.ts"
import {Path, REDIRECT_AFTER_LOGIN_KEY} from "@/app/routing/Router.tsx"
import {AuthAPI} from "@/app/api/artel"
import {apiPrefix} from "@/app/api/api.ts"
import PasswordLoginForm from "@/pages/init/components/LoginContent/components/PasswordLoginForm/PasswordLoginForm.tsx"
import RegisterForm from "@/pages/init/components/LoginContent/components/RegisterForm/RegisterForm.tsx"
import cls from "@/pages/init/components/LoginContent/LoginContent.module.css"

interface LoginContentProps {
    login: (session: Session, userInfo?: UserInfo) => void
    navigate: NavigateFunction
}

type Mode = "login" | "register"

export default function LoginContent({login, navigate}: LoginContentProps) {
    const {CloseDialog, UnlockClosing} = useDialog()
    const {bake} = useToaster()
    const passwordAuthEnabled = useAppConfig((s) => s.passwordAuthEnabled)
    const telegramAuthEnabled = useAppConfig((s) => s.telegramAuthEnabled)
    const selfRegisterEnabled = useAppConfig((s) => s.selfRegisterEnabled)

    const [botId, setBotId] = useState("")
    const [mode, setMode] = useState<Mode>("login")

    useEffect(function () {
        AuthAPI.GetConfig({}, apiPrefix()).then(function (cfg) {
            setBotId(cfg.telegramClientId ?? "")
        }).catch(() => {
        })
    }, [])

    async function handleLoggedIn(session: Session) {
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
    }

    function onSuccessLoginTg(data: TelegramAuthData) {
        authService.LoginViaTelegram(data.id_token)
            .then(handleLoggedIn)
            .catch(function (err: unknown) {
                bake({
                    title: "Telegram login failed",
                    description: err instanceof Error ? err.message : "Telegram login failed",
                    level: "Error",
                })
            })
            .finally(CloseDialog)
    }

    function onPasswordLoggedIn(session: Session) {
        handleLoggedIn(session).finally(CloseDialog)
    }

    return (
        <div className={cls.LoginContentContainer}>
            <div className={cls.Logo}>artel</div>
            {mode === "login" && telegramAuthEnabled && (
                <TelegramAuth
                    title="Sign in with Telegram"
                    botId={botId}
                    className={cls.LoginButton}
                    onSuccess={onSuccessLoginTg}
                />
            )}
            {mode === "login" && passwordAuthEnabled && (
                <PasswordLoginForm onLoggedIn={onPasswordLoggedIn}/>
            )}
            {mode === "register" && passwordAuthEnabled && (
                <RegisterForm onRegistered={onPasswordLoggedIn} onSwitchToLogin={() => setMode("login")}/>
            )}
            {mode === "login" && selfRegisterEnabled && passwordAuthEnabled && (
                <Button
                    type="button"
                    variant="unstyled"
                    className={cls.RegisterLink}
                    onClick={() => setMode("register")}
                >
                    Don&apos;t have an account? Register
                </Button>
            )}
        </div>
    )
}
