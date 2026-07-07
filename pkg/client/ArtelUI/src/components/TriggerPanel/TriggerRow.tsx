import {Button} from "@vervstack/chures"
import cls from "@/components/TriggerPanel/TriggerRow.module.css"

import {Trigger} from "@/processes/Tracts.ts"
import {providerLabel, triggerChipLabel} from "@/components/TriggerPanel/triggerLabels.ts"
import {useTriggerSources} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"

import WebhookDetailsDialog from "@/dialogs/WebhookDetailsDialog/WebhookDetailsDialog.tsx"

export default function TriggerRow({trigger, onUnlink, onToggle, onRotate}: {
    trigger: Trigger
    onUnlink: () => void
    onToggle: (enabled: boolean) => void
    onRotate: () => void
}) {
    const {OpenDialog} = useDialog()
    const {triggerSources} = useTriggerSources()
    const preset = triggerSources.find(s => s.key === trigger.source)
    const sharedProvider = trigger.kind === "webhook" ? preset?.provider ?? "" : ""
    const webhookUrl = trigger.kind === "webhook" && !sharedProvider ? `${window.location.origin}/tract/hook/${trigger.triggerUuid}` : ""

    function handleViewWebhook() {
        OpenDialog(<WebhookDetailsDialog webhookUrl={webhookUrl} onRotate={onRotate}/>)
    }

    return (
        <div className={cls.Row}>
            <span className={cls.Name}>{trigger.name}</span>
            <span className={cls.Chip}>{triggerChipLabel(trigger, preset)}</span>
            <div className={cls.RowActions}>
                {trigger.kind === "webhook" && (
                    sharedProvider ? (
                        <span className={cls.Chip}>Shared · {providerLabel(sharedProvider)} connection</span>
                    ) : (
                        <>
                            <Button variant="ghost" onClick={handleViewWebhook}>View webhook</Button>
                            <Button variant="ghost" onClick={onRotate}>Rotate token</Button>
                        </>
                    )
                )}
                <Button variant="ghost" onClick={() => onToggle(!trigger.enabled)}>{trigger.enabled ? "Enabled" : "Disabled"}</Button>
                <Button variant="iconDanger" onClick={onUnlink} aria-label="Unlink trigger">✕</Button>
            </div>
        </div>
    )
}
