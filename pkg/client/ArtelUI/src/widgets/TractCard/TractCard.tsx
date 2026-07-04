import cls from "@/widgets/TractCard/TractCard.module.css"

import {Tract, TractLastRun, TractTriggerSummary} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import Button from "@/components/shared/Button/Button.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"

interface Props {
    tract: Tract
    onClick: () => void
}

export default function TractCard({tract, onClick}: Props) {
    const {setEnabled, deleteTract} = useTracts()
    const bakeError = useBakeError()
    const {OpenDialog} = useDialog()

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
                onConfirm={() => deleteTract(tract.uuid).catch(err => bakeError("Failed to delete tract", err))}
            />
        )
    }

    return (
        <div
            className={`${cls.Card} ${!tract.enabled ? cls.CardPaused : ""}`}
            onClick={onClick}
            role="button"
            tabIndex={0}
        >
            <TriggerIcon triggers={tract.triggers}/>
            <CardInfo name={tract.name} description={tract.description}/>
            <CardMeta triggers={tract.triggers} lastRun={tract.lastRun} enabled={tract.enabled} onToggle={handleToggle}/>
            <Button className={cls.DeleteBtn} variant="iconDanger" onClick={handleDelete} aria-label="Delete tract">
                <TrashIcon/>
            </Button>
            <ChevronIcon/>
        </div>
    )
}

function TriggerIcon({triggers}: { triggers?: TractTriggerSummary[] }) {
    const kind = triggers?.[0]?.kind
    return (
        <div className={`${cls.IconWrap} ${kind === "webhook" ? cls.IconWebhook : cls.IconManual}`}>
            {kind === "webhook" ? <WebhookIcon/> : <ManualIcon/>}
        </div>
    )
}

function CardInfo({name, description}: { name: string; description: string }) {
    return (
        <div className={cls.Info}>
            <span className={cls.Name}>{name}</span>
            {description && <span className={cls.Desc}>{description}</span>}
        </div>
    )
}

function CardMeta({triggers, lastRun, enabled, onToggle}: {
    triggers?: TractTriggerSummary[]
    lastRun?: TractLastRun
    enabled: boolean
    onToggle: () => void
}) {
    const trig = triggers?.[0]
    const chipLabel = trig
        ? `${trig.kind === "webhook" ? "Webhook" : "Manual"}${trig.source ? ` · ${trig.source}` : ""}`
        : "Manual"

    return (
        <div className={cls.Meta}>
            <span className={cls.Chip}>{chipLabel}</span>
            <RunStatusDot lastRun={lastRun}/>
            <label className={cls.Toggle} onClick={e => e.stopPropagation()}>
                <input type="checkbox" checked={enabled} onChange={onToggle}/>
                <span className={cls.ToggleTrack}/>
            </label>
        </div>
    )
}

function RunStatusDot({lastRun}: { lastRun?: TractLastRun }) {
    if (!lastRun) return <span className={`${cls.Dot} ${cls.DotNever}`}/>
    if (lastRun.status === "done") return <span className={`${cls.Dot} ${cls.DotOk}`}/>
    if (lastRun.status === "failed") return <span className={`${cls.Dot} ${cls.DotFailed}`}/>
    return <span className={`${cls.Dot} ${cls.DotRunning}`}/>
}

function WebhookIcon() {
    return (
        <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <path d="M18 16.98h-5.99c-1.1 0-1.95.94-2.48 1.9A4 4 0 0 1 2 17c0-1.1.94-2 2.06-2H8"/>
            <path d="M17 7.82a2 2 0 0 1 3 1.73V13"/>
            <path d="M13 6.51a2 2 0 0 0-3.46 0L6 12.5"/>
        </svg>
    )
}

function ManualIcon() {
    return (
        <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            <polygon points="6 3 20 12 6 21 6 3"/>
        </svg>
    )
}

function TrashIcon() {
    return (
        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="3 6 5 6 21 6"/>
            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
        </svg>
    )
}

function ChevronIcon() {
    return (
        <svg className={cls.Chevron} viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="9 18 15 12 9 6"/>
        </svg>
    )
}
