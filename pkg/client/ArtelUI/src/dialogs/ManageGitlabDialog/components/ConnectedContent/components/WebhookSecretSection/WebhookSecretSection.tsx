import {useState} from "react"
import {Button} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import TokenRevealDialog from "@/dialogs/TokenRevealDialog/TokenRevealDialog.tsx"
import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/dialogs/ManageGitlabDialog/components/ConnectedContent/components/WebhookSecretSection/WebhookSecretSection.module.css"

export default function WebhookSecretSection(
    {webhookSecretSet, webhookUrl}: { webhookSecretSet: boolean; webhookUrl: string },
) {
    const [generating, setGenerating] = useState(false)

    const {generateGitlabWebhookSecret} = useExternalConnections()
    const {OpenDialog} = useDialog()
    const bakeError = useBakeError()

    function handleGenerate() {
        setGenerating(true)
        generateGitlabWebhookSecret()
            .then(webhookSecret => {
                OpenDialog(<TokenRevealDialog webhookUrl={webhookUrl} webhookToken={webhookSecret}/>)
            })
            .catch(e => bakeError("Failed to generate webhook secret", e))
            .finally(() => setGenerating(false))
    }

    return (
        <div className={cls.Field}>
            <span className={cls.FieldLabel}>Webhook secret</span>
            <div className={cls.ModalActions}>
                <span>{webhookSecretSet ? "Configured" : "Not configured"}</span>
                <Button variant="ghost" onClick={handleGenerate} disabled={generating}>
                    {generating ? "Generating…" : webhookSecretSet ? "Rotate" : "Generate"}
                </Button>
            </div>
        </div>
    )
}
