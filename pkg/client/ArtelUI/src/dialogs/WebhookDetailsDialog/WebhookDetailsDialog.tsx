import {useState} from "react"
import {useQueryClient} from "@tanstack/react-query"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/WebhookDetailsDialog/WebhookDetailsDialog.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {triggersQueryKey, useTracts, useTrigger, useTriggerSources} from "@/app/hooks/Tracts.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {providerEnumFor, providerLabel} from "@/components/TriggerPanel/triggerLabels.ts"
import TokenRevealDialog from "@/dialogs/TokenRevealDialog/TokenRevealDialog.tsx"
import ManageGitlabDialog from "@/dialogs/ManageGitlabDialog/ManageGitlabDialog.tsx"

interface WebhookDetailsDialogProps {
    triggerUuid: string
}

export default function WebhookDetailsDialog({triggerUuid}: WebhookDetailsDialogProps) {
    const {OpenDialog, CloseDialog} = useDialog()
    const {rotateTriggerToken} = useTracts()
    const {connections} = useExternalConnections()
    const bakeError = useBakeError()
    const queryClient = useQueryClient()

    const [copied, setCopied] = useState(false)
    const {trigger} = useTrigger(triggerUuid)
    const {triggerSources} = useTriggerSources()

    const preset = triggerSources.find(s => s.key === trigger?.source)
    const sharedProvider = trigger?.kind === "webhook" ? preset?.provider ?? "" : ""
    const sharedConnection = sharedProvider ? connections.find(c => c.provider === providerEnumFor(sharedProvider)) : undefined

    const webhookUrl = sharedProvider
        ? (sharedConnection ? `${window.location.origin}/webhooks/gitlab/${sharedConnection.id}` : "")
        : trigger?.kind === "webhook"
            ? `${window.location.origin}/tract/hook/${trigger.triggerUuid}`
            : ""

    function handleCopy() {
        navigator.clipboard.writeText(webhookUrl).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    function handleRotate() {
        CloseDialog()
        OpenDialog(
            <ConfirmDialog
                title="Rotate webhook token"
                message="The old webhook URL stops working immediately."
                confirmLabel="Rotate"
                danger
                onClose={CloseDialog}
                onConfirm={() => rotateTriggerToken(triggerUuid)
                    .then(result => {
                        queryClient.invalidateQueries({queryKey: triggersQueryKey})
                        setTimeout(() => OpenDialog(<TokenRevealDialog webhookUrl={result.webhookUrl} webhookToken={result.webhookToken}/>), 0)
                    })
                    .catch(err => bakeError("Failed to rotate token", err))}
            />
        )
    }

    function handleManageConnection() {
        CloseDialog()
        OpenDialog(<ManageGitlabDialog/>)
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Hook info</h2>

            <span className={cls.FieldLabel}>URL</span>
            <div className={cls.UrlRow}>
                <div className={cls.UrlBox}>{webhookUrl}</div>
                <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy"}</Button>
            </div>

            {sharedProvider ? (
                <p className={cls.Notice}>
                    This trigger uses your shared {providerLabel(sharedProvider)} connection's webhook —
                    manage its secret from the connection settings.
                </p>
            ) : (
                <>
                    <span className={cls.TokenLabel}>Token</span>
                    <div className={cls.TokenBox}>•••••••• {trigger?.tokenSuffix || "····"}</div>

                    <p className={cls.Notice}>
                        The full token was shown once, right after this webhook was created or last
                        rotated, and can't be displayed again. Rotate to get a new one.
                    </p>
                </>
            )}

            <div className={cls.DialogActions}>
                {sharedProvider
                    ? <Button variant="ghost" onClick={handleManageConnection}>Manage connection</Button>
                    : <Button variant="ghost" onClick={handleRotate}>Rotate token</Button>}
                <Button variant="primary" onClick={CloseDialog}>Done</Button>
            </div>
        </div>
    )
}
