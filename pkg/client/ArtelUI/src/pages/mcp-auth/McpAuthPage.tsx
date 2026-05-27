import {useEffect, useRef, useState} from "react"
import {useSearchParams} from "react-router-dom"

import cls from "@/pages/init/InitPage.module.css"

import {useDialog} from "@/app/hooks/Dialog.ts"
import useUser from "@/hooks/user/User.ts"

type Vault = {
    id: string
    name: string
}

const SESSION_KEY = "mcpSessionToken"

export default function McpAuthPage() {
    const [params] = useSearchParams()
    const clientId = params.get("client_id") ?? ""
    const redirectUri = params.get("redirect_uri") ?? ""
    const codeChallenge = params.get("code_challenge") ?? ""
    const state = params.get("state") ?? ""

    const {OpenDialog, LockClosing} = useDialog()
    const {auth} = useUser()

    useEffect(function () {
        LockClosing()

        async function init() {
            // 1. Try main Artel session token
            if (auth.isAuthenticated()) {
                const token = auth.getToken()
                try {
                    const vaults = await fetchVaults(token)
                    OpenDialog(
                        <VaultSelect
                            vaults={vaults}
                            sessionToken={token}
                            clientId={clientId}
                            redirectUri={redirectUri}
                            codeChallenge={codeChallenge}
                            state={state}
                        />
                    )
                    return
                } catch {
                    // main token rejected by oauth endpoint — fall through
                }
            }

            // 2. Try stored MCP session token
            const stored = localStorage.getItem(SESSION_KEY)
            if (stored) {
                try {
                    const vaults = await fetchVaults(stored)
                    OpenDialog(
                        <VaultSelect
                            vaults={vaults}
                            sessionToken={stored}
                            clientId={clientId}
                            redirectUri={redirectUri}
                            codeChallenge={codeChallenge}
                            state={state}
                        />
                    )
                    return
                } catch {
                    localStorage.removeItem(SESSION_KEY)
                }
            }

            // 3. Show login dialog
            OpenDialog(
                <McpLogin
                    clientId={clientId}
                    redirectUri={redirectUri}
                    codeChallenge={codeChallenge}
                    state={state}
                />
            )
        }

        void init()
    }, [])

    return <div className={cls.Root}/>
}

async function fetchVaults(sessionToken: string): Promise<Vault[]> {
    const res = await fetch("/oauth/vaults", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({session_token: sessionToken}),
    })
    if (!res.ok) throw new Error("failed")
    const body = await res.json()
    return body.vaults ?? []
}

interface McpLoginProps {
    clientId: string
    redirectUri: string
    codeChallenge: string
    state: string
}

function McpLogin({clientId, redirectUri, codeChallenge, state}: McpLoginProps) {
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
                const res = await fetch("/oauth/login", {
                    method: "POST",
                    headers: {"Content-Type": "application/json"},
                    body: JSON.stringify({id_token: data.id_token}),
                })
                if (!res.ok) throw new Error("Authentication failed")
                const body = await res.json()
                localStorage.setItem(SESSION_KEY, body.session_token)
                OpenDialog(
                    <VaultSelect
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
        script.setAttribute("data-client-id", import.meta.env.VITE_TELEGRAM_BOT_ID ?? "")
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
        <div className={cls.Card}>
            <div className={cls.Logo}>artel</div>
            <p style={{color: "var(--secondary-fg-color)", fontSize: "var(--font-size-sm)", textAlign: "center", margin: 0}}>
                Sign in to grant Claude access to your vault
            </p>
            {error && <div className={cls.Error}>{error}</div>}
            {loading && <span style={{color: "var(--secondary-fg-color)"}}>Checking…</span>}
            <div ref={telegramRef} className={cls.TelegramContainer}/>
        </div>
    )
}

interface VaultSelectProps {
    vaults: Vault[]
    sessionToken: string
    clientId: string
    redirectUri: string
    codeChallenge: string
    state: string
}

function VaultSelect({vaults, sessionToken, clientId, redirectUri, codeChallenge, state}: VaultSelectProps) {
    const {OpenDialog} = useDialog()
    const [error, setError] = useState("")
    const [loading, setLoading] = useState(false)

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
            window.location.href = body.redirect_url
        } catch (e) {
            setError(e instanceof Error ? e.message : "Failed to grant access")
        } finally {
            setLoading(false)
        }
    }

    function switchAccount() {
        localStorage.removeItem(SESSION_KEY)
        OpenDialog(
            <McpLogin
                clientId={clientId}
                redirectUri={redirectUri}
                codeChallenge={codeChallenge}
                state={state}
            />
        )
    }

    return (
        <div className={cls.Card}>
            <div className={cls.Logo}>artel</div>
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
