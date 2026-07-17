import {useState} from "react"

import cls from "@/components/TractRunTimeline/TractRunTimeline.module.css"
import {cn} from "@/app/utils/cn.ts"
import {TractRunStep} from "@/processes/Tracts.ts"
import JsonBlock from "@/components/JsonBlock/JsonBlock.tsx"

interface Props {
    step: TractRunStep
}

function statusClass(status: string): string {
    if (status === "done") return cls.StatusDone
    if (status === "failed") return cls.StatusFailed
    return cls.StatusRunning
}

export default function StepRow({step}: Props) {
    const [open, setOpen] = useState(false)
    const isCondition = step.stepType === "condition"
    const conditionResult = isCondition && step.output && typeof step.output === "object"
        ? (step.output as Record<string, unknown>).result
        : undefined

    return (
        <div className={cls.Step}>
            <div className={cls.StepHeader} onClick={() => setOpen(o => !o)}>
                <span className={cn(cls.StatusChip, statusClass(step.status))}>{step.status}</span>
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
