import {Button} from "@vervstack/chures"
import cls from "@/components/TriggerPanel/TriggerRow.module.css"

import {Trigger} from "@/processes/Tracts.ts"
import {providerLabel, triggerChipLabel} from "@/components/TriggerPanel/triggerLabels.ts"
import {useTriggerSources} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"

import WebhookDetailsDialog from "@/dialogs/WebhookDetailsDialog/WebhookDetailsDialog.tsx"

export default function TriggerRow({trigger, onUnlink}: {
    trigger: Trigger
    onUnlink: () => void
}) {
    const {OpenDialog} = useDialog()
    const {triggerSources} = useTriggerSources()
    const preset = triggerSources.find(s => s.key === trigger.source)
    const sharedProvider = trigger.kind === "webhook" ? preset?.provider ?? "" : ""

    function handleHookInfo() {
        OpenDialog(<WebhookDetailsDialog triggerUuid={trigger.uuid}/>)
    }

    return (
        <div className={cls.Row}>
            <span className={cls.Name}>{trigger.name}</span>
            <span className={cls.Chip}>{triggerChipLabel(trigger, preset)}</span>
            <div className={cls.RowActions}>
                {trigger.kind === "webhook" &&
                    sharedProvider &&
                    <span className={cls.Chip}>Shared · {providerLabel(sharedProvider)} connection</span>}

                <Button variant="ghost" onClick={handleHookInfo}>Hook info</Button>

                <Button variant="iconDanger" onClick={onUnlink} aria-label="Unlink trigger">✕</Button>
            </div>
        </div>
    )
}
