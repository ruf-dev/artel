import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import {useTracts} from "@/app/hooks/Tracts.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/RunLog.module.css"
import {buildLogLines}
    from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/processes/logLines.ts"
import LogEntryRow
    from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/components/LogEntryRow/LogEntryRow.tsx"

interface RunLogProps {
    runUuid: string
    runStatus: string
    onRetry: () => Promise<unknown>
}

export default function RunLog({runUuid, runStatus, onRetry}: RunLogProps) {
    const {currentRun, currentRunSteps, currentRunLoading, watchRun, fetchRun, clearCurrentRun} = useTracts()
    const bakeError = useBakeError()
    const [retrying, setRetrying] = useState(false)

    useEffect(() => {
        if (runStatus !== "running") {
            fetchRun(runUuid).catch(err => bakeError("Failed to load run", err))
            return () => clearCurrentRun()
        }

        const stopWatching = watchRun(runUuid)
        return () => {
            stopWatching()
            clearCurrentRun()
        }
    }, [runUuid, runStatus])

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
            <div className={cls.Entries}>
                {lines.map(l => <LogEntryRow key={l.key} line={l}/>)}
            </div>
            <div className={cls.RetryBar}>
                <Button variant="ghost" onClick={handleRetry} disabled={retrying}>
                    {retrying ? "Retrying…" : "Retry run"}
                </Button>
            </div>
        </>
    )
}
