import {useState} from "react"

import {Button} from "@vervstack/chures"
import cls from "@/dialogs/WebhookDetailsDialog/WebhookDetailsDialog.module.css"

import {useDialog} from "@/app/hooks/Dialog"

export default function WebhookDetailsDialog({webhookUrl, onRotate}: {
    webhookUrl: string
    onRotate: () => void
}) {
    const {CloseDialog} = useDialog()
    const [copied, setCopied] = useState(false)

    function handleCopy() {
        navigator.clipboard.writeText(webhookUrl).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    function handleRotate() {
        CloseDialog()
        onRotate()
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Webhook details</h2>

            <span className={cls.FieldLabel}>URL</span>
            <div className={cls.UrlRow}>
                <div className={cls.UrlBox}>{webhookUrl}</div>
                <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy"}</Button>
            </div>

            <p className={cls.Notice}>
                The secret token was shown once, right after this webhook was created or last
                rotated, and can't be displayed again. Rotate the token to get a new one.
            </p>

            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={handleRotate}>Rotate token</Button>
                <Button variant="primary" onClick={CloseDialog}>Done</Button>
            </div>
        </div>
    )
}
