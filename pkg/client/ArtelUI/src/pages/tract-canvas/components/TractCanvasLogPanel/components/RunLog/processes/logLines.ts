import {TractRun, TractRunStep} from "@/processes/Tracts.ts"

export type LogLineKind = "info" | "ok" | "err"

export interface LogLine {
    key: string
    t: string
    title: string
    meta?: string
    kind: LogLineKind
    // Marks a run-level event (trigger fired / run completed / run failed) rather than a step
    // outcome — rendered with a distinct accent so it reads as separate from the step-by-step log.
    terminal?: boolean
    error?: string
    input?: unknown
}

export function formatTime(iso: string | undefined): string {
    if (!iso) return ""
    const d = new Date(iso)
    if (isNaN(d.getTime())) return ""
    return d.toLocaleTimeString(undefined, {hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false})
}

function formatDuration(startedAt: string | undefined, finishedAt: string | undefined): string | undefined {
    if (!startedAt || !finishedAt) return undefined
    const ms = new Date(finishedAt).getTime() - new Date(startedAt).getTime()
    if (isNaN(ms) || ms < 0) return undefined
    return ms < 1000 ? `${ms}ms` : `${(ms / 1000).toFixed(1)}s`
}

function stepMeta(s: TractRunStep): string | undefined {
    return [s.stepType, formatDuration(s.startedAt, s.finishedAt)].filter(Boolean).join(" · ") || undefined
}

export function buildLogLines(run: TractRun, steps: TractRunStep[]): LogLine[] {
    const lines: LogLine[] = [{
        key: "trigger",
        t: formatTime(run.createdAt),
        title: "Trigger fired",
        meta: run.startedBy,
        kind: "info",
        terminal: true,
        input: run.triggerPayload,
    }]

    for (const s of steps) {
        const name = s.stepName || s.stepId
        if (s.status === "done") {
            lines.push({
                key: `${s.stepId}-done`,
                t: formatTime(s.finishedAt) || formatTime(s.startedAt),
                title: name,
                meta: stepMeta(s),
                kind: "ok",
                input: s.input,
            })
        } else if (s.status === "failed") {
            lines.push({
                key: `${s.stepId}-fail`,
                t: formatTime(s.finishedAt) || formatTime(s.startedAt),
                title: name,
                meta: stepMeta(s),
                kind: "err",
                error: s.error,
                input: s.input,
            })
        } else {
            lines.push({
                key: `${s.stepId}-running`,
                t: formatTime(s.startedAt),
                title: name,
                meta: s.stepType ? `${s.stepType} · running…` : "running…",
                kind: "info",
                input: s.input,
            })
        }
    }

    if (run.status === "done") {
        lines.push({
            key: "done", t: formatTime(run.updatedAt), title: "Tract run completed",
            kind: "ok", terminal: true,
        })
    } else if (run.status === "failed") {
        lines.push({
            key: "failed", t: formatTime(run.updatedAt), title: "Tract run failed",
            kind: "err", terminal: true, error: run.error,
        })
    }

    return lines
}
