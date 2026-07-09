import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/TractCanvasLogPanel.module.css"
import {TractRun, TractRunStep} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {ChevronRightIcon, CloseIcon} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    open: boolean
    runs: TractRun[]
    selectedRunUuid: string | null
    onSelectRun: (uuid: string) => void
    onClose: () => void
}

export default function TractCanvasLogPanel({open, runs, selectedRunUuid, onSelectRun, onClose}: Props) {
    const [enlarged, setEnlarged] = useState(false)

    useEffect(() => {
        if (!open) setEnlarged(false)
    }, [open])

    return (
        <>
            {open && enlarged && (
                <div className={cls.Backdrop} onClick={() => setEnlarged(false)}/>
            )}
            <div className={cn(cls.Panel, open && cls.PanelOpen, open && enlarged && cls.Enlarged)}>
                <LogPanelBar enlarged={enlarged} onToggleEnlarge={() => setEnlarged(e => !e)} onClose={onClose}/>
                <div className={cls.Split}>
                    <div className={cls.RunList}>
                        {runs.length === 0 && <p className={cls.Empty}>No runs yet.</p>}
                        {runs.map(r => (
                            <Button
                                key={r.uuid}
                                variant="ghost"
                                className={cn(cls.RunRow, r.uuid === selectedRunUuid && cls.RunRowSelected)}
                                onClick={() => onSelectRun(r.uuid)}
                            >
                                <span className={cn(cls.Dot, dotClass(r.status, cls))}/>
                                <span className={cls.RunMeta}>{r.startedBy}</span>
                                <span className={cls.RunMeta}>{formatDate(r.createdAt)}</span>
                            </Button>
                        ))}
                    </div>
                    <div className={cls.LogScroll}>
                        {selectedRunUuid ? <RunLog runUuid={selectedRunUuid}/> : <p className={cls.Empty}>Select a run to see its log.</p>}
                    </div>
                </div>
            </div>
        </>
    )
}

function LogPanelBar({enlarged, onToggleEnlarge, onClose}: {
    enlarged: boolean
    onToggleEnlarge: () => void
    onClose: () => void
}) {
    return (
        <div className={cls.Bar}>
            <Button
                variant="ghost"
                className={cls.EnlargeBtn}
                onClick={onToggleEnlarge}
                aria-label={enlarged ? "Shrink log panel" : "Enlarge log panel"}
                aria-pressed={enlarged}
            >
                <ChevronRightIcon className={cn(cls.EnlargeChevron, enlarged && cls.EnlargeChevronOpen)}/>
            </Button>
            <span className={cls.BarLabel}>Runs</span>
            <Button variant="iconDanger" onClick={onClose} aria-label="Close log">
                <CloseIcon/>
            </Button>
        </div>
    )
}

function dotClass(status: string, cls: Record<string, string>): string {
    if (status === "done") return cls.DotOk
    if (status === "failed") return cls.DotErr
    return cls.DotRunning
}

function formatDate(iso: string | undefined): string {
    if (!iso) return ""
    const d = new Date(iso)
    if (isNaN(d.getTime())) return ""
    return d.toLocaleString(undefined, {month: "short", day: "numeric", hour: "2-digit", minute: "2-digit"})
}

interface LogLine {
    key: string
    t: string
    msg: string
    kind: "info" | "ok" | "err"
}

function buildLogLines(run: TractRun, steps: TractRunStep[]): LogLine[] {
    const lines: LogLine[] = []
    const triggerT = formatTime(run.createdAt)
    lines.push({key: "trigger", t: triggerT, msg: `Trigger fired (${run.startedBy}) — payload ${JSON.stringify(run.triggerPayload)}`, kind: "info"})

    for (const s of steps) {
        const t = formatTime(s.startedAt)
        lines.push({key: `${s.stepId}-start`, t, msg: `${s.stepType}.${s.stepName || s.stepId} →`, kind: "info"})
        if (s.status === "done") {
            lines.push({key: `${s.stepId}-done`, t: formatTime(s.finishedAt) || t, msg: `✓ ${s.stepName || s.stepId} → ${JSON.stringify(s.output)}`, kind: "ok"})
        } else if (s.status === "failed") {
            lines.push({key: `${s.stepId}-fail`, t: formatTime(s.finishedAt) || t, msg: `✗ ${s.stepName || s.stepId} — ${s.error}`, kind: "err"})
        } else {
            lines.push({key: `${s.stepId}-running`, t, msg: `${s.stepName || s.stepId} running…`, kind: "info"})
        }
    }

    if (run.status === "done") {
        lines.push({key: "done", t: formatTime(run.updatedAt), msg: "✓ Tract run completed", kind: "ok"})
    } else if (run.status === "failed") {
        lines.push({key: "failed", t: formatTime(run.updatedAt), msg: `✗ Tract run failed — ${run.error}`, kind: "err"})
    }

    return lines
}

function formatTime(iso: string | undefined): string {
    if (!iso) return ""
    const d = new Date(iso)
    if (isNaN(d.getTime())) return ""
    return d.toLocaleTimeString(undefined, {hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false})
}

function RunLog({runUuid}: { runUuid: string }) {
    const {currentRun, currentRunSteps, currentRunLoading, fetchRun, clearCurrentRun} = useTracts()

    useEffect(() => {
        void fetchRun(runUuid)
        return () => clearCurrentRun()
    }, [runUuid])

    if (currentRunLoading || !currentRun) {
        return <p className={cls.Empty}>Loading…</p>
    }

    const lines = buildLogLines(currentRun, currentRunSteps)

    return (
        <>
            {lines.map(l => (
                <div className={cls.Line} key={l.key}>
                    <span className={cls.LineT}>{l.t}</span>
                    <span className={cn(cls.LineM, cls[`LineM${cap(l.kind)}`])}>{l.msg}</span>
                </div>
            ))}
        </>
    )
}

function cap(s: string): string {
    return s.charAt(0).toUpperCase() + s.slice(1)
}
