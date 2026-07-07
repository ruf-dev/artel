import {useState} from "react"

import {Button} from "@vervstack/chures"
import cls from "@/dialogs/TokenRevealDialog/TokenRevealDialog.module.css"

import {useDialog} from "@/app/hooks/Dialog"

export default function TokenRevealDialog({webhookUrl, webhookToken}: { webhookUrl: string; webhookToken: string }) {
    const {CloseDialog} = useDialog()
    const [copied, setCopied] = useState(false)

    function handleCopy() {
        navigator.clipboard.writeText(webhookToken).then(() => {
            setCopied(true)
            setTimeout(() => setCopied(false), 2000)
        })
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Webhook token</h2>
            <p className={cls.Notice}>
                This token is shown once. Store it now — send it as <code>X-Tract-Token</code> (or{" "}
                <code>X-Gitlab-Token</code>) on webhook deliveries to <code>{webhookUrl}</code>.
            </p>
            <div className={cls.TokenBox}>{webhookToken}</div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={handleCopy}>{copied ? "Copied!" : "Copy token"}</Button>
                <Button variant="primary" onClick={CloseDialog}>Done</Button>
            </div>
        </div>
    )
}
