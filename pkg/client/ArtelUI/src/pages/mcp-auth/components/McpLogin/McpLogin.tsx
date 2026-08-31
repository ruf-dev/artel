import {useEffect, useRef, useState} from "react"

import {useDialog} from "@/app/hooks/Dialog.ts"
import VaultSelect from "@/pages/mcp-auth/components/VaultSelect/VaultSelect.tsx"
import cls from "@/pages/mcp-auth/components/McpLogin/McpLogin.module.css"

const SESSION_KEY = "mcpSessionToken"

interface McpLoginProps {
    botId: string
    clientId: string
    redirectUri: string
    codeChallenge: string
    state: string
}

export default function McpLogin({botId, clientId, redirectUri, codeChallenge, state}: McpLoginProps) {
    const {OpenDialog} = useDialog()
    const telegramRef = useRef<HTMLDivElement>(null)
    const [error, setError] = useState("")
    const [loading, setLoading] = useState(false)

    useEffect(function () {
        if (!telegramRef.current) return

        ;(window as unknown as Record<string, unknown>).onMcpTelegramAuth = async function (data: {id_token: string}) {
            setLoading(true)
            setError("")
            try {
                const res = await fetch("/api/oauth/login", {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify({id_token: data.id_token}),
                })
                if (!res.ok) throw new Error("Authentication failed")
                const body = await res.json()
                localStorage.setItem(SESSION_KEY, body.session_token)
                OpenDialog(
                    <VaultSelect
                        botId={botId}
                        vaults={body.vaults ?? []}
                        sessionToken={body.session_token}
                        clientId={clientId}
                        redirectUri={redirectUri}
                        codeChallenge={codeChallenge}
                        state={state}
                    />
                )
            } catch (e) {
                setError(e instanceof Error ? e.message : "Login failed")
            } finally {
                setLoading(false)
            }
        }

        const script = document.createElement("script")
        script.src = "https://oauth.telegram.org/js/telegram-login.js?3"
        script.setAttribute("data-client-id", botId)
        script.setAttribute("data-size", "large")
        script.setAttribute("data-onauth", "onMcpTelegramAuth(data)")
        script.setAttribute("data-request-access", "write phone")
        script.async = true
        telegramRef.current.appendChild(script)

        return function () {
            delete (window as unknown as Record<string, unknown>).onMcpTelegramAuth
        }
    }, [])

    return (
        <div className={cls.McpLoginContainer}>
            <div className={cls.Logo}>artel</div>
            <p className={cls.Sub}>
                Sign in to grant Claude access to your vault
            </p>
            {error && <div className={cls.Error}>{error}</div>}
            {loading && <span className={cls.Checking}>Checking…</span>}
            <div ref={telegramRef} className={cls.TelegramContainer}/>
        </div>
    )
}
