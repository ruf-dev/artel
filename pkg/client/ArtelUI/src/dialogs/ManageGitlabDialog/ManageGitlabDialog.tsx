import {useState} from "react"

import cls from "@/dialogs/ManageGitlabDialog/ManageGitlabDialog.module.css"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import Button from "@/components/shared/Button/Button.tsx"
import Input from "@/components/shared/Input/Input.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import ModalClose from "@/components/ModalClose/ModalClose.tsx"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"

export default function ManageGitlabDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_GITLAB)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect GitLab"
                message="Remove the connection to GitLab? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onConfirm={() =>
                    disconnect("gitlab").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead onClose={CloseDialog}/>
            {connection
                ? <ConnectedContent connectionId={connection.id ?? ""} fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}/>
                : <ConnectForm/>}
        </div>
    )
}

function DialogHead({onClose}: { onClose: () => void }) {
    return (
        <div className={cls.ModalHead}>
            <div className={cls.ModalHeadLeft}>
                <div className={cls.ModalIcon}>
                    <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_GITLAB}/>
                </div>
                <span className={cls.ModalTitle}>GitLab</span>
            </div>
            <ModalClose onClick={onClose} className={cls.ModalClose}/>
        </div>
    )
}

function ConnectForm() {
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
        <>
            <p className={cls.ModalSub}>
                Connect a GitLab account to create and list merge requests and issues, and comment on them.
                We&apos;ll verify the token against the instance before saving it.
            </p>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Personal access token</span>
                <Input type="password" placeholder="glpat-…" value={personalAccessToken}
                    onChange={e => setPersonalAccessToken(e.target.value)} disabled={connecting} autoComplete="off"/>
            </label>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Instance URL (optional)</span>
                <Input placeholder="https://gitlab.com" value={instanceUrl}
                    onChange={e => setInstanceUrl(e.target.value)} disabled={connecting} autoComplete="off"/>
            </label>
            <div className={cls.ModalActions}>
                <Button variant="primary" onClick={handleConnect} disabled={connecting || !personalAccessToken}>
                    {connecting ? "Verifying…" : "Connect"}
                </Button>
            </div>
        </>
    )
}

function ConnectedContent({connectionId, fields, onDisconnect}: {
    connectionId: string
    fields: Record<string, string>
    onDisconnect: () => void
}) {
    const [copied, setCopied] = useState(false)
    const webhookUrl = `${window.location.origin}/webhooks/gitlab/${connectionId}`

    function handleCopy() {
        navigator.clipboard.writeText(webhookUrl).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <>
            <p className={cls.ModalSub}>
                Connected as <b>@{fields.username}</b> on <b>{fields.instance_url}</b>.
            </p>
            <WebhookSecretSection webhookSecretSet={fields.webhook_secret_set === "true"}/>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Webhook URL</span>
                <Input value={webhookUrl} readOnly/>
            </label>
            <p className={cls.ModalSub}>
                Add this URL as a webhook in your GitLab project settings, with the secret above, to receive merge
                request and issue events.
            </p>
            <div className={cls.ModalActions}>
                <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy URL"}</Button>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </>
    )
}

function WebhookSecretSection({webhookSecretSet}: { webhookSecretSet: boolean }) {
    const [editing, setEditing] = useState(!webhookSecretSet)
    const [saving, setSaving] = useState(false)
    const [webhookSecret, setWebhookSecret] = useState("")

    const {setGitlabWebhookSecret} = useExternalConnections()
    const bakeError = useBakeError()

    function handleSave() {
        setSaving(true)
        setGitlabWebhookSecret({webhookSecret})
            .then(() => setEditing(false))
            .catch(e => bakeError("Failed to save webhook secret", e))
            .finally(() => setSaving(false))
    }

    if (!editing) {
        return (
            <div className={cls.Field}>
                <span className={cls.FieldLabel}>Webhook secret</span>
                <div className={cls.ModalActions}>
                    <span>Configured</span>
                    <Button variant="ghost" onClick={() => setEditing(true)}>Rotate</Button>
                </div>
            </div>
        )
    }

    return (
        <label className={cls.Field}>
            <span className={cls.FieldLabel}>Webhook secret</span>
            <Input type="password" placeholder="Verifies incoming GitLab webhooks" value={webhookSecret}
                onChange={e => setWebhookSecret(e.target.value)} disabled={saving} autoComplete="off"/>
            <div className={cls.ModalActions}>
                <Button variant="secondary" onClick={handleSave} disabled={saving || !webhookSecret}>
                    {saving ? "Saving…" : "Save webhook secret"}
                </Button>
            </div>
        </label>
    )
}
