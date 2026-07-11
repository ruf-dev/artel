import {useState} from "react"
import {Button, Input} from "@vervstack/chures"

import WebhookSecretSection
from "@/dialogs/ManageGitlabDialog/components/ConnectedContent/components/WebhookSecretSection/WebhookSecretSection.tsx"
import cls from "@/dialogs/ManageGitlabDialog/components/ConnectedContent/ConnectedContent.module.css"

interface ConnectedContentProps {
    connectionId: string
    fields: Record<string, string>
    onDisconnect: () => void
}

export default function ConnectedContent({connectionId, fields, onDisconnect}: ConnectedContentProps) {
    const [copied, setCopied] = useState(false)
    const webhookUrl = `${window.location.origin}/webhooks/gitlab/${connectionId}`

    function handleCopy() {
        navigator.clipboard.writeText(webhookUrl).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <div className={cls.ConnectedContentContainer}>
            <p className={cls.ModalSub}>
                Connected as <b>@{fields.username}</b> on <b>{fields.instance_url}</b>.
            </p>
            <WebhookSecretSection webhookSecretSet={fields.webhook_secret_set === "true"} webhookUrl={webhookUrl}/>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Webhook URL</span>
                <Input value={webhookUrl} setValue={() => {}} readOnly/>
            </label>
            <p className={cls.ModalSub}>
                Add this URL as a webhook in your GitLab project settings, with the secret above, to receive merge
                request and issue events.
            </p>
            <div className={cls.ModalActions}>
                <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy URL"}</Button>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </div>
    )
}
