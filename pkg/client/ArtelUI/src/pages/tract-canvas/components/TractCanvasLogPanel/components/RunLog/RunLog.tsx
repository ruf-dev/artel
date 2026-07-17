import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import {TractRun, TractRunStep} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {cn} from "@/app/utils/cn.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import JsonBlock from "@/components/JsonBlock/JsonBlock.tsx"
import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/RunLog.module.css"

interface LogLine {
    key: string
    t: string
    msg: string
    kind: "info" | "ok" | "err"
    input?: unknown
}

function formatTime(iso: string | undefined): string {
    if (!iso) return ""
    const d = new Date(iso)
    if (isNaN(d.getTime())) return ""
    return d.toLocaleTimeString(undefined, {hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false})
}

function buildLogLines(run: TractRun, steps: TractRunStep[]): LogLine[] {
    const lines: LogLine[] = []
    const triggerT = formatTime(run.createdAt)
    lines.push({
        key: "trigger",
        t: triggerT,
        msg: `Trigger fired (${run.startedBy})`,
        kind: "info",
        input: run.triggerPayload,
    })

    for (const s of steps) {
        const t = formatTime(s.startedAt)
        lines.push({
            key: `${s.stepId}-start`,
            t,
            msg: `${s.stepType}.${s.stepName || s.stepId} →`,
            kind: "info",
            input: s.input,
        })
        if (s.status === "done") {
            lines.push({
                key: `${s.stepId}-done`,
                t: formatTime(s.finishedAt) || t,
                msg: `✓ ${s.stepName || s.stepId} → ${JSON.stringify(s.output)}`,
                kind: "ok",
            })
        } else if (s.status === "failed") {
            lines.push({
                key: `${s.stepId}-fail`,
                t: formatTime(s.finishedAt) || t,
                msg: `✗ ${s.stepName || s.stepId} — ${s.error}`,
                kind: "err",
            })
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

function cap(s: string): string {
    return s.charAt(0).toUpperCase() + s.slice(1)
}

interface RunLogProps {
    runUuid: string
    onRetry: () => Promise<unknown>
}

export default function RunLog({runUuid, onRetry}: RunLogProps) {
    const {currentRun, currentRunSteps, currentRunLoading, fetchRun, clearCurrentRun} = useTracts()
    const bakeError = useBakeError()
    const [retrying, setRetrying] = useState(false)

    useEffect(() => {
        void fetchRun(runUuid)
        return () => clearCurrentRun()
    }, [runUuid])

    function handleRetry() {
        setRetrying(true)
        onRetry()
            .catch(err => bakeError("Failed to retry run", err))
            .finally(() => setRetrying(false))
    }

    if (currentRunLoading || !currentRun) {
        return <p className={cls.Empty}>Loading…</p>
    }

    const lines = buildLogLines(currentRun, currentRunSteps)

    return (
        <>
            {lines.map(l => (
                <div className={cls.Entry} key={l.key}>
                    <div className={cls.Line}>
                        <span className={cls.LineT}>{l.t}</span>
                        <span className={cn(cls.LineM, cls[`LineM${cap(l.kind)}`])}>{l.msg}</span>
                    </div>
                    {l.input !== undefined && l.input !== null && (
                        <div className={cls.EntryInput}>
                            <JsonBlock label="Input" value={l.input}/>
                        </div>
                    )}
                </div>
            ))}
            <div className={cls.RetryBar}>
                <Button variant="ghost" onClick={handleRetry} disabled={retrying}>
                    {retrying ? "Retrying…" : "Retry run"}
                </Button>
            </div>
        </>
    )
}
