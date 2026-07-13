import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import UseDefaultButton
    from "@/dialogs/ManageGitlabDialog/components/ConnectForm/components/UseDefaultButton/UseDefaultButton.tsx"
import GitlabCheckButton from "@/dialogs/ManageGitlabDialog/components/GitlabCheckButton/GitlabCheckButton.tsx"
import cls from "@/dialogs/ManageGitlabDialog/components/ConnectForm/ConnectForm.module.css"

const DEFAULT_INSTANCE_URL = "https://gitlab.com"

function tokenSettingsUrl(instanceUrl: string): string | null {
    let parsed: URL
    try {
        parsed = new URL(instanceUrl.trim())
    } catch {
        return null
    }
    return `${parsed.origin}/-/user_settings/personal_access_tokens/legacy/new`
}

export default function ConnectForm() {
    const [connecting, setConnecting] = useState(false)
    const [personalAccessToken, setPersonalAccessToken] = useState("")
    const [instanceUrl, setInstanceUrl] = useState("")
    const [verified, setVerified] = useState(false)

    const {addGitlabConnection} = useExternalConnections()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

    const generateTokenUrl = tokenSettingsUrl(instanceUrl)

    function handleTokenChange(value: string) {
        setPersonalAccessToken(value)
        setVerified(false)
    }

    function handleUrlChange(value: string) {
        setInstanceUrl(value)
        setVerified(false)
    }

    function handleConnect() {
        setConnecting(true)
        addGitlabConnection({personalAccessToken, instanceUrl})
            .then(CloseDialog)
            .catch(e => bakeError("Failed to connect GitLab", e))
            .finally(() => setConnecting(false))
    }

    return (
        <div className={cls.ConnectFormContainer}>
            <p className={cls.ModalSub}>
                Connect a GitLab account to create and list merge requests and issues, and comment on them.
                We&apos;ll verify the token against the instance before saving it.
            </p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Instance URL (optional)</span>
                <span className={cls.InputRow}>
                    <Input className={cls.InstanceUrlInput} placeholder="https://gitlab.com" value={instanceUrl}
                        setValue={handleUrlChange} disabled={connecting} autoComplete="off"/>
                    <UseDefaultButton onClick={() => handleUrlChange(DEFAULT_INSTANCE_URL)} disabled={connecting}/>
                </span>
            </label>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Personal access token</span>
                <Input type="text" inputClassName={cls.TokenInput} placeholder="glpat-…" value={personalAccessToken}
                    setValue={handleTokenChange} disabled={connecting} autoComplete="off"/>
                {generateTokenUrl && (
                    <div className={cls.TokenHint}>
                        <a href={generateTokenUrl} target="_blank" rel="noopener noreferrer">
                            Click here to generate token
                        </a>
                    </div>
                )}
            </label>
            <div className={cls.ModalActions}>
                <GitlabCheckButton req={{personalAccessToken, instanceUrl}}
                    disabled={connecting || !personalAccessToken} onResult={setVerified}/>
                <Button variant="primary" onClick={handleConnect}
                    disabled={connecting || !personalAccessToken || !verified}>
                    {connecting ? "Verifying…" : "Connect"}
                </Button>
            </div>
        </div>
    )
}
