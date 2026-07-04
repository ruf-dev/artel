import {useEffect, useState} from "react"

import cls from "@/components/TractRunTimeline/TractRunTimeline.module.css"

import {TractRunStep} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"

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

function TriggerPayloadSection({payload}: { payload: unknown }) {
    const [open, setOpen] = useState(false)
    return (
        <div className={cls.Step}>
            <div className={cls.StepHeader} onClick={() => setOpen(o => !o)}>
                <span className={cls.StepName}>trigger</span>
                <span className={cls.StepType}>payload</span>
            </div>
            {open && (
                <div className={cls.StepBody}>
                    <pre className={cls.Json}>{JSON.stringify(payload, null, 2)}</pre>
                </div>
            )}
        </div>
    )
}

function statusClass(status: string): string {
    if (status === "done") return cls.StatusDone
    if (status === "failed") return cls.StatusFailed
    return cls.StatusRunning
}

function StepRow({step}: { step: TractRunStep }) {
    const [open, setOpen] = useState(false)
    const isCondition = step.stepType === "condition"
    const conditionResult = isCondition && step.output && typeof step.output === "object"
        ? (step.output as Record<string, unknown>).result
        : undefined

    return (
        <div className={cls.Step}>
            <div className={cls.StepHeader} onClick={() => setOpen(o => !o)}>
                <span className={`${cls.StatusChip} ${statusClass(step.status)}`}>{step.status}</span>
                <span className={cls.StepName}>{step.stepName || step.stepId}</span>
                <span className={cls.StepType}>{step.stepType}</span>
                {conditionResult !== undefined && (
                    <span className={cls.ConditionResult}>result: {String(conditionResult)}</span>
                )}
            </div>
            {open && (
                <div className={cls.StepBody}>
                    {step.input !== undefined && <JsonBlock label="Input" value={step.input}/>}
                    {step.output !== undefined && <JsonBlock label="Output" value={step.output}/>}
                    {step.error && <p className={cls.RunError}>{step.error}</p>}
                </div>
            )}
        </div>
    )
}

function JsonBlock({label, value}: { label: string; value: unknown }) {
    return (
        <div className={cls.JsonBlock}>
            <span className={cls.JsonLabel}>{label}</span>
            <pre className={cls.Json}>{JSON.stringify(value, null, 2)}</pre>
        </div>
    )
}
