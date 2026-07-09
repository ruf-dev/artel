import {useEffect} from "react"

import cls from "@/components/TractRunTimeline/TractRunTimeline.module.css"
import {useTracts} from "@/app/hooks/Tracts.ts"
import TriggerPayloadSection from "@/components/TractRunTimeline/components/TriggerPayloadSection.tsx"
import StepRow from "@/components/TractRunTimeline/components/StepRow.tsx"

interface Props {
    runUuid: string
}

export default function TractRunTimeline({runUuid}: Props) {
    const {currentRun, currentRunSteps, currentRunLoading, fetchRun, clearCurrentRun} = useTracts()

    useEffect(() => {
        void fetchRun(runUuid)
        return () => clearCurrentRun()
    }, [runUuid])

    if (currentRunLoading || !currentRun) {
        return <p className={cls.Empty}>Loading…</p>
    }

    return (
        <div className={cls.TimelineContainer}>
            <TriggerPayloadSection payload={currentRun.triggerPayload}/>
            {currentRunSteps.map(step => (
                <StepRow key={step.stepId} step={step}/>
            ))}
            {currentRun.error && <p className={cls.RunError}>{currentRun.error}</p>}
        </div>
    )
}
