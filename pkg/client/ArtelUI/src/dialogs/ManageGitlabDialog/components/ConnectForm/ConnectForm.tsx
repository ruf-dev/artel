import {useState} from "react"
import {Button, Input} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import cls from "@/dialogs/ManageGitlabDialog/components/ConnectForm/ConnectForm.module.css"

export default function ConnectForm() {
    const [connecting, setConnecting] = useState(false)
    const [personalAccessToken, setPersonalAccessToken] = useState("")
    const [instanceUrl, setInstanceUrl] = useState("")

    const {addGitlabConnection} = useExternalConnections()
    const {CloseDialog} = useDialog()
    const bakeError = useBakeError()

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
                <span className={cls.FieldLabel}>Personal access token</span>
                <Input type="password" placeholder="glpat-…" value={personalAccessToken}
                    setValue={setPersonalAccessToken} disabled={connecting} autoComplete="off"/>
            </label>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Instance URL (optional)</span>
                <Input placeholder="https://gitlab.com" value={instanceUrl}
                    setValue={setInstanceUrl} disabled={connecting} autoComplete="off"/>
            </label>
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleConnect} disabled={connecting || !personalAccessToken}>
                    {connecting ? "Verifying…" : "Connect"}
                </Button>
            </div>
        </div>
    )
}
