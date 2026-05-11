import {useState, useRef, useEffect} from "react"
import {useNavigate} from "react-router-dom"

import cls from "@/pages/init/InitPage.module.css"

import {AuthMiddleware, AuthService, Session} from "@/processes/Auth.ts"
import {Path} from "@/app/routing/Router.tsx"

interface Props {
    auth: AuthMiddleware
    onLogin: (session: Session) => void
}

export default function InitPage({auth, onLogin}: Props) {
    const navigate = useNavigate()
    const [mode, setMode] = useState<"login" | "register">("login")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")
    const [loading, setLoading] = useState(false)

    const svc = new AuthService()
    const telegramContainerRef = useRef<HTMLDivElement>(null)

    useEffect(function() {
        if (!telegramContainerRef.current) return

        ;(window as unknown as Record<string, unknown>).onTelegramAuth = function(data: {id_token: string}) {
            setLoading(true)
            setError("")
            svc.LoginViaTelegram(data.id_token)
                .then(function(session) {
                    onLogin(session)
                    navigate(Path.HomePage)
                })
                .catch(function(err: unknown) {
                    setError(err instanceof Error ? err.message : "Telegram login failed")
                })
                .finally(function() {
                    setLoading(false)
                })
        }

        const script = document.createElement("script")
        script.src = "https://oauth.telegram.org/telegram-login.js"
        script.setAttribute("data-telegram-login", import.meta.env.VITE_TELEGRAM_BOT_ID ?? "")
        script.setAttribute("data-size", "large")
        script.setAttribute("data-onauth", "onTelegramAuth(user)")
        script.setAttribute("data-request-access", "write")
        script.async = true
        telegramContainerRef.current.appendChild(script)

        return function() {
            delete (window as unknown as Record<string, unknown>).onTelegramAuth
        }
    }, [])

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault()
        setError("")
        setLoading(true)

        try {
            if (mode === "login") {
                const session = await svc.LoginViaPass(email, password)
                onLogin(session)
                navigate(Path.HomePage)
            } else {
                await svc.Register(email, password)
                setMode("login")
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : "Unknown error")
        } finally {
            setLoading(false)
        }
    }

    if (auth.isAuthenticated()) {
        navigate(Path.HomePage)
        return null
    }

    return (
        <div className={cls.Root}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>

                <form className={cls.Form} onSubmit={handleSubmit}>
                    <input
                        className={cls.Input}
                        type="email"
                        placeholder="email"
                        value={email}
                        onChange={e => setEmail(e.target.value)}
                        required
                    />
                    <input
                        className={cls.Input}
                        type="password"
                        placeholder="password"
                        value={password}
                        onChange={e => setPassword(e.target.value)}
                        required
                    />
                    {error && <span className={cls.Error}>{error}</span>}
                    <button className={cls.SubmitBtn} type="submit" disabled={loading}>
                        {loading ? "..." : mode === "login" ? "Log in" : "Create account"}
                    </button>
                </form>

                <div className={cls.Divider}>or</div>

                <div ref={telegramContainerRef} className={cls.TelegramContainer} />

                <div className={cls.Toggle}>
                    {mode === "login" ? (
                        <>No account?{" "}
                            <button type="button" onClick={() => setMode("register")}>Sign up</button>
                        </>
                    ) : (
                        <>Have an account?{" "}
                            <button type="button" onClick={() => setMode("login")}>Log in</button>
                        </>
                    )}
                </div>
            </div>
        </div>
    )
}
