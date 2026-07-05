import {useMemo, useState} from "react"

import {useTracts} from "@/app/hooks/Tracts.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import {NodeStatus} from "@/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx"
import {CanvasLayout, TRIGGER_NODE_ID} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {TractRun} from "@/processes/Tracts.ts"
import {sleep} from "@/app/utils/sleep.ts"

const RUN_POLL_INTERVAL_MS = 700
const RUN_POLL_MAX_ATTEMPTS = 10

export function useTractRunTracking(tractUuid: string, layout: CanvasLayout) {
    const {runTract, fetchRuns, currentRun, currentRunSteps} = useTracts()
    const bakeError = useBakeError()

    const [running, setRunning] = useState(false)
    const [logOpen, setLogOpen] = useState(false)
    const [selectedRunUuid, setSelectedRunUuid] = useState<string | null>(null)

    const activeRunSteps = useMemo(
        () => selectedRunUuid && currentRun?.uuid === selectedRunUuid ? currentRunSteps : [],
        [selectedRunUuid, currentRun, currentRunSteps],
    )

    const lastOutputByStepId = useMemo(() => {
        const m: Record<string, unknown> = {}
        for (const s of activeRunSteps) {
            if (s.status === "done" && s.output && typeof s.output === "object") m[s.stepId] = s.output
        }
        return m
    }, [activeRunSteps])

    function nodeStatus(id: string): NodeStatus {
        if (id === TRIGGER_NODE_ID) return activeRunSteps.length > 0 ? "ok" : "idle"
        const step = activeRunSteps.find(s => s.stepId === id)
        if (!step) return "idle"
        if (step.status === "done") return "ok"
        if (step.status === "failed") return "err"
        return "running"
    }

    const runningEdgeIds = useMemo(() => {
        if (!running) return new Set<string>()
        return new Set(layout.edges.map(e => e.id))
    }, [running, layout])

    function pollUntilRunFinished(attempt = 0): Promise<void> {
        return fetchRuns(tractUuid).then(() => {
            const latest = useTracts.getState().runsByTract[tractUuid]?.[0]
            if (latest && latest.status !== "running") return
            if (attempt >= RUN_POLL_MAX_ATTEMPTS - 1) return
            return sleep(RUN_POLL_INTERVAL_MS).then(() => pollUntilRunFinished(attempt + 1))
        })
    }

    function latestRun(): TractRun | undefined {
        // Read the live store rather than a `runs` prop — this closure outlives the render
        // that created it (it runs after two awaited fetches), so a prop would be stale.
        return useTracts.getState().runsByTract[tractUuid]?.[0]
    }

    function handleRun() {
        setRunning(true)
        runTract(tractUuid, {})
            .then(() => pollUntilRunFinished())
            .then(() => {
                const latest = latestRun()
                if (latest) {
                    setSelectedRunUuid(latest.uuid)
                    setLogOpen(true)
                }
            })
            .catch(err => bakeError("Failed to run tract", err))
            .finally(() => setRunning(false))
    }

    return {
        running,
        logOpen, setLogOpen,
        selectedRunUuid, setSelectedRunUuid,
        activeRunSteps, lastOutputByStepId, nodeStatus, runningEdgeIds,
        handleRun,
    }
}
