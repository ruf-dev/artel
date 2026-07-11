import {useEffect} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/components/TriggerPanel/TriggerPanel.module.css"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import AddTriggerDialog from "@/dialogs/AddTriggerDialog/AddTriggerDialog.tsx"
import TriggerRow from "@/components/TriggerPanel/TriggerRow.tsx"

interface Props {
    tractUuid: string
    linkedTriggerSummaries: { uuid: string; name: string; kind: string; source: string }[]
}

export default function TriggerPanel({tractUuid, linkedTriggerSummaries}: Props) {
    const {triggers, fetchTriggers, unlinkTrigger} = useTracts()
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()

    useEffect(() => {
        void fetchTriggers()
    }, [fetchTriggers])

    const linkedUuids = new Set(linkedTriggerSummaries.map(t => t.uuid))
    const linked = triggers.filter(t => linkedUuids.has(t.uuid))

    function handleUnlink(triggerUuid: string) {
        OpenDialog(
            <ConfirmDialog
                title="Unlink trigger"
                message="This tract will no longer start on this trigger's events."
                confirmLabel="Unlink"
                danger
                onClose={CloseDialog}
                onConfirm={() => unlinkTrigger(triggerUuid, tractUuid)
                    .catch(err => bakeError("Failed to unlink trigger", err))}
            />
        )
    }

    return (
        <div className={cls.Panel}>
            <div className={cls.PanelHeader}>
                <span className={cls.PanelTitle}>Trigger</span>
                {linked.length === 0 && (
                    <Button variant="ghost" onClick={() => OpenDialog(<AddTriggerDialog tractUuid={tractUuid}
                                                                                        linkedUuids={linkedUuids}/>)}>
                        + Add trigger
                    </Button>
                )}
            </div>
            {linked.length === 0 && (
                <p className={cls.Empty}>
                    No triggers linked — use "Run" in the Runs panel to fire manually, or add a trigger.
                </p>
            )}
            {linked.map(t => (
                <TriggerRow
                    key={t.uuid}
                    trigger={t}
                    onUnlink={() => handleUnlink(t.uuid)}
                />
            ))}
        </div>
    )
}
