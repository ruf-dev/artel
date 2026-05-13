import {useEffect, useRef, useState} from "react"
import {useSearchParams} from "react-router-dom"

import cls from "@/pages/init/InitPage.module.css"

type Vault = {
    id: string
    name: string
}

type Step = "login" | "vaults" | "done"

const SESSION_KEY = "mcpSessionToken"

export default function McpAuthPage() {
    const [params] = useSearchParams()
    const clientId = params.get("client_id") ?? ""
    const redirectUri = params.get("redirect_uri") ?? ""
    const codeChallenge = params.get("code_challenge") ?? ""
    const state = params.get("state") ?? ""

    const [step, setStep] = useState<Step>("login")
    const [sessionToken, setSessionToken] = useState("")
    const [vaults, setVaults] = useState<Vault[]>([])
    const [error, setError] = useState("")
    const [loading, setLoading] = useState(false)

    const telegramRef = useRef<HTMLDivElement>(null)

    useEffect(function () {
        const stored = localStorage.getItem(SESSION_KEY)
        if (!stored) return

        setLoading(true)
        fetch("/oauth/vaults", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({session_token: stored}),
        })
            .then(function (res) {
                if (!res.ok) throw new Error()
                return res.json()
            })
            .then(function (body) {
                setSessionToken(stored)
                setVaults(body.vaults ?? [])
                setStep("vaults")
            })
            .catch(function () {
                localStorage.removeItem(SESSION_KEY)
            })
            .finally(function () {
                setLoading(false)
            })
    }, [])

    useEffect(function () {
        if (step !== "login" || !telegramRef.current) return

        ;(window as unknown as Record<string, unknown>).onMcpTelegramAuth = async function (data: {id_token: string}) {
            setLoading(true)
            setError("")
            try {
                const res = await fetch("/oauth/login", {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify({id_token: data.id_token}),
                })
                if (!res.ok) throw new Error("Authentication failed")
                const body = await res.json()
                localStorage.setItem(SESSION_KEY, body.session_token)
                setSessionToken(body.session_token)
                setVaults(body.vaults ?? [])
                setStep("vaults")
            } catch (e) {
                setError(e instanceof Error ? e.message : "Login failed")
            } finally {
                setLoading(false)
            }
        }

        const script = document.createElement("script")
        script.src = "https://oauth.telegram.org/js/telegram-login.js?3"
        script.setAttribute("data-client-id", import.meta.env.VITE_TELEGRAM_BOT_ID ?? "")
        script.setAttribute("data-size", "large")
        script.setAttribute("data-onauth", "onMcpTelegramAuth(data)")
        script.setAttribute("data-request-access", "write phone")
        script.async = true
        telegramRef.current.appendChild(script)

        return function () {
            delete (window as unknown as Record<string, unknown>).onMcpTelegramAuth
        }
    }, [step])

    async function selectVault(vaultId: string) {
        setLoading(true)
        setError("")
        try {
            const res = await fetch("/oauth/vault", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({
                    session_token: sessionToken,
                    vault_id: vaultId,
                    client_id: clientId,
                    redirect_uri: redirectUri,
                    code_challenge: codeChallenge,
                    state,
                }),
            })
            if (!res.ok) throw new Error("Failed to grant access")
            const body = await res.json()
            setStep("done")
            window.location.href = body.redirect_url
        } catch (e) {
            setError(e instanceof Error ? e.message : "Failed to grant access")
        } finally {
            setLoading(false)
        }
    }

    function switchAccount() {
        localStorage.removeItem(SESSION_KEY)
        setSessionToken("")
        setVaults([])
        setError("")
        setStep("login")
    }

    return (
        <div className={cls.Root}>
            <div className={cls.Card}>
                <div className={cls.Logo}>artel</div>

                {loading && step === "login" && (
                    <span style={{color: "var(--secondary-fg-color)"}}>Checking session…</span>
                )}

                {step === "login" && !loading && (
                    <>
                        <p style={{color: "var(--secondary-fg-color)", fontSize: "var(--font-size-sm)", textAlign: "center", margin: 0}}>
                            Sign in to grant Claude access to your vault
                        </p>
                        {error && <div className={cls.Error}>{error}</div>}
                        <div>
                            <div ref={telegramRef} className={cls.TelegramContainer}/>
                            <button className="tg-auth-button" data-style="shine">Sign In with Telegram</button>
                        </div>
                    </>
                )}

                {step === "vaults" && (
                    <>
                        <p style={{color: "var(--secondary-fg-color)", fontSize: "var(--font-size-sm)", textAlign: "center", margin: 0}}>
                            Select a vault for Claude to access
                        </p>
                        {error && <div className={cls.Error}>{error}</div>}
                        <VaultList vaults={vaults} loading={loading} onSelect={selectVault}/>
                        <button
                            onClick={switchAccount}
                            style={{background: "none", border: "none", color: "var(--thirdy-fg-color)", fontSize: "var(--font-size-sm)", cursor: "pointer", padding: 0}}
                        >
                            Use a different account
                        </button>
                    </>
                )}

                {step === "done" && (
                    <p style={{color: "var(--secondary-fg-color)", fontSize: "var(--font-size-sm)", textAlign: "center", margin: 0}}>
                        Access granted. Redirecting…
                    </p>
                )}
            </div>
        </div>
    )
}

function VaultList({vaults, loading, onSelect}: {vaults: Vault[]; loading: boolean; onSelect: (id: string) => void}) {
    if (vaults.length === 0) {
        return (
            <p style={{color: "var(--thirdy-fg-color)", fontSize: "var(--font-size-sm)", textAlign: "center"}}>
                No vaults found. Create one in Artel first.
            </p>
        )
    }
    return (
        <div style={{display: "flex", flexDirection: "column", gap: "var(--space-2)", width: "100%"}}>
            {vaults.map(function (v) {
                return (
                    <button
                        key={v.id}
                        className={cls.SubmitBtn}
                        disabled={loading}
                        onClick={function () { onSelect(v.id) }}
                    >
                        {v.name}
                    </button>
                )
            })}
        </div>
    )
}
