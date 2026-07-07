import cls from "@/pages/tract-canvas/widgets/TractCanvasCard/TractCanvasCard.module.css"

import {Tract} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import {ChevronRightIcon, ManualTriggerIcon, TrashIcon, WebhookIcon} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"
import RunStatusDot from "@/pages/tract-canvas/components/RunStatusDot/RunStatusDot.tsx"

import {Button, ConfirmDialog} from "@vervstack/chures"
interface Props {
    tract: Tract
    onClick: () => void
}

export default function TractCanvasCard({tract, onClick}: Props) {
    const {setEnabled, deleteTract} = useTracts()
    const bakeError = useBakeError()
    const {OpenDialog, CloseDialog} = useDialog()

    function handleToggle() {
        setEnabled(tract.uuid, !tract.enabled).catch(err => bakeError("Failed to update tract", err))
    }

    function handleDelete(e: React.MouseEvent) {
        e.stopPropagation()
        OpenDialog(
            <ConfirmDialog
                title="Delete tract"
                message={`Delete "${tract.name}"? This cannot be undone.`}
                confirmLabel="Delete"
                danger
                onClose={CloseDialog}
                onConfirm={() => deleteTract(tract.uuid).catch(err => bakeError("Failed to delete tract", err))}
            />
        )
    }

    const trig = tract.triggers?.[0]
    const isWebhook = trig?.kind === "webhook"

    return (
        <div
            className={`${cls.Card} ${!tract.enabled ? cls.CardPaused : ""}`}
            onClick={onClick}
            role="button"
            tabIndex={0}
        >
            <div className={`${cls.IconWrap} ${isWebhook ? cls.IconWebhook : cls.IconManual}`}>
                {isWebhook ? <WebhookIcon/> : <ManualTriggerIcon/>}
            </div>
            <div className={cls.Info}>
                <span className={cls.Name}>{tract.name}</span>
                {tract.description && <span className={cls.Desc}>{tract.description}</span>}
            </div>
            <div className={cls.Meta}>
                <span className={cls.Chip}>
                    {trig ? `${isWebhook ? "Webhook" : "Manual"}${trig.source ? ` · ${trig.source}` : ""}` : "Manual"}
                </span>
                <RunStatusDot lastRun={tract.lastRun}/>
                <label className={cls.Toggle} onClick={e => e.stopPropagation()}>
                    <input type="checkbox" checked={tract.enabled} onChange={handleToggle}/>
                    <span className={cls.ToggleTrack}/>
                </label>
            </div>
            <Button className={cls.DeleteBtn} variant="iconDanger" onClick={handleDelete} aria-label="Delete tract">
                <TrashIcon className={cls.DeleteIcon}/>
            </Button>
            <ChevronRightIcon className={cls.Chevron}/>
        </div>
    )
}
