import {useState} from "react"
import {Button} from "@vervstack/chures"

import {getCsrfToken} from "@/app/api/api.ts"
import {useDialog} from "@/app/hooks/Dialog.ts"
import {Vault} from "@/pages/mcp-auth/processes/vault.ts"
import VaultList from "@/pages/mcp-auth/components/VaultList/VaultList.tsx"
import McpLogin from "@/pages/mcp-auth/components/McpLogin/McpLogin.tsx"
import cls from "@/pages/mcp-auth/components/VaultSelect/VaultSelect.module.css"

const SESSION_KEY = "mcpSessionToken"

interface VaultSelectProps {
    botId: string
    vaults: Vault[]
    sessionToken: string
    clientId: string
    redirectUri: string
    codeChallenge: string
    state: string
}

export default function VaultSelect(props: VaultSelectProps) {
    const {OpenDialog} = useDialog()
    const [error, setError] = useState("")
    const [loading, setLoading] = useState(false)

    async function selectVault(vaultId: string) {
        setLoading(true)
        setError("")
        try {
            const res = await fetch("/api/oauth/vault", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-Csrf-Token": getCsrfToken() ?? "",
                },
                body: JSON.stringify({
                    session_token: props.sessionToken,
                    vault_id: vaultId,
                    client_id: props.clientId,
                    redirect_uri: props.redirectUri,
                    code_challenge: props.codeChallenge,
                    state: props.state,
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
                botId={props.botId}
                clientId={props.clientId}
                redirectUri={props.redirectUri}
                codeChallenge={props.codeChallenge}
                state={props.state}
            />
        )
    }

    return (
        <div className={cls.VaultSelectContainer}>
            <div className={cls.Logo}>artel</div>
            <p className={cls.Sub}>
                Select a vault for Claude to access
            </p>
            {error && <div className={cls.Error}>{error}</div>}
            <VaultList vaults={props.vaults} loading={loading} onSelect={selectVault}/>
            <Button variant="ghost" onClick={switchAccount} className={cls.SwitchAccountBtn}>
                Use a different account
            </Button>
        </div>
    )
}
