import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageTrelloDialog/components/DialogHead/DialogHead.tsx"
import TrelloCheckButton from "@/dialogs/ManageTrelloDialog/components/TrelloCheckButton/TrelloCheckButton.tsx"
import cls from "@/dialogs/ManageTrelloDialog/components/TrelloAddDialog/TrelloAddDialog.module.css"

const API_KEY_URL = "https://trello.com/power-ups/admin/"

function tokenAuthorizeUrl(apiKey: string): string | null {
    if (!apiKey.trim()) return null
    const params = new URLSearchParams({
        expiration: "never",
        scope: "read,write",
        response_type: "token",
        name: "Artel",
        key: apiKey.trim(),
    })
    return `https://trello.com/1/authorize?${params.toString()}`
}

export default function TrelloAddDialog() {
    const [connecting, setConnecting] = useState(false)
    const [apiKey, setApiKey] = useState("")
    const [apiToken, setApiToken] = useState("")
    const [verified, setVerified] = useState(false)

    const {addTrelloConnection} = useExternalConnections()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    const generateTokenUrl = tokenAuthorizeUrl(apiKey)

    function handleApiKeyChange(value: string) {
        setApiKey(value)
        setVerified(false)
    }

    function handleApiTokenChange(value: string) {
        setApiToken(value)
        setVerified(false)
    }

    function handleConnect() {
        setConnecting(true)
        addTrelloConnection({apiKey, apiToken})
            .then(CloseDialog)
            .catch(e => bakeError("Failed to connect Trello", e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead title="Add Trello account" onClose={CloseDialog} disabled={connecting}/>
            <p className={cls.ModalSub}>
                Connect a Trello account to read and comment on cards. We&apos;ll verify the credentials
                against Trello before saving them.
            </p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>API key</span>
                <Input type="text" value={apiKey} setValue={handleApiKeyChange} disabled={connecting}
                    autoComplete="off"/>
                <div className={cls.TokenHint}>
                    <a href={API_KEY_URL} target="_blank" rel="noopener noreferrer">Create a Power-Up</a>
                    {" "}to get your API key
                </div>
            </label>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>API token</span>
                <Input type="text" inputClassName={cls.TokenInput} value={apiToken} setValue={handleApiTokenChange}
                    disabled={connecting} autoComplete="off"/>
                {generateTokenUrl && (
                    <div className={cls.TokenHint}>
                        Click{" "}
                        <a href={generateTokenUrl} target="_blank" rel="noopener noreferrer">here</a>
                        {" "}to generate a token
                    </div>
                )}
            </label>
            <div className={cls.ModalActions}>
                <TrelloCheckButton req={{apiKey, apiToken}} disabled={connecting || !apiKey || !apiToken}
                    onResult={setVerified}/>
                <Button variant="primary" onClick={handleConnect}
                    disabled={connecting || !apiKey || !apiToken || !verified}>
                    {connecting ? "Connecting…" : "Connect"}
                </Button>
            </div>
        </div>
    )
}
