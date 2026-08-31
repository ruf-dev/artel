import {useEffect} from "react"
import {useSearchParams} from "react-router-dom"

import cls from "@/pages/mcp-auth/McpAuthPage.module.css"
import {useDialog} from "@/app/hooks/Dialog.ts"
import useUser from "@/hooks/user/User.ts"
import {AuthAPI} from "@/app/api/artel"
import {apiPrefix, getCsrfToken} from "@/app/api/api.ts"
import {Vault} from "@/pages/mcp-auth/processes/vault.ts"
import McpLogin from "@/pages/mcp-auth/components/McpLogin/McpLogin.tsx"
import VaultSelect from "@/pages/mcp-auth/components/VaultSelect/VaultSelect.tsx"

const SESSION_KEY = "mcpSessionToken"

async function fetchVaults(sessionToken?: string): Promise<Vault[]> {
    const res = await fetch("/api/oauth/vaults", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "X-Csrf-Token": getCsrfToken() ?? "",
        },
        body: JSON.stringify(sessionToken ? {session_token: sessionToken} : {}),
    })
    if (!res.ok) throw new Error("failed")
    const body = await res.json()
    return body.vaults ?? []
}

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
            const cfg = await AuthAPI.GetConfig({}, apiPrefix()).catch(() => null)
            const botId = cfg?.telegramClientId ?? ""

            // 1. Try the main Artel browser session (httpOnly access_token cookie + CSRF header)
            if (auth.isAuthenticated()) {
                try {
                    const vaults = await fetchVaults()
                    OpenDialog(
                        <VaultSelect
                            botId={botId}
                            vaults={vaults}
                            sessionToken=""
                            clientId={clientId}
                            redirectUri={redirectUri}
                            codeChallenge={codeChallenge}
                            state={state}
                        />
                    )
                    return
                } catch {
                    // cookie session rejected by oauth endpoint — fall through
                }
            }

            // 2. Try stored MCP session token
            const stored = localStorage.getItem(SESSION_KEY)
            if (stored) {
                try {
                    const vaults = await fetchVaults(stored)
                    OpenDialog(
                        <VaultSelect
                            botId={botId}
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
                    botId={botId}
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
