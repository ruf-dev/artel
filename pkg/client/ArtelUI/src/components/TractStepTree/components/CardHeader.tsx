import {useState} from "react"
import {Button, Input} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/components/TractStepTree/TractStepTree.module.css"
import {TractStep} from "@/processes/Tracts.ts"

const STEP_ID_PATTERN = /^[a-z][a-z0-9_]*$/

interface Props {
    step: TractStep
    onUpdate: (updater: (s: TractStep) => TractStep) => void
    onDelete: () => void
    right?: React.ReactNode
}

export default function CardHeader({step, onUpdate, onDelete, right}: Props) {
    const [name, setName] = useState(step.name || step.id)
    const nameInvalid = !STEP_ID_PATTERN.test(name) || name === "trigger"

    function commitName() {
        if (nameInvalid || name === step.name) return
        onUpdate(s => ({...s, name}))
    }

    const iconClass = step.type === "action" ? cls.CardIconAction
        : step.type === "condition" ? cls.CardIconCondition
        : step.type === "parallel" ? cls.CardIconParallel
        : cls.CardIconGroup

    return (
        <div className={cls.CardHeader}>
            <span className={cn(cls.CardIcon, iconClass)}>{step.type.slice(0, 1).toUpperCase()}</span>
            <Input
                className={cls.NameInputWrapper}
                inputClassName={cn(cls.NameInput, nameInvalid && cls.NameInputInvalid)}
                value={name}
                setValue={setName}
                onBlur={commitName}
                data-tooltip-id={nameInvalid ? "root-tooltip" : undefined}
                data-tooltip-content={nameInvalid ? "Must match ^[a-z][a-z0-9_]*$ and not be \"trigger\"" : undefined}
            />
            <span className={cls.TypeLabel}>{step.type}</span>
            {right}
            <div className={cls.CardActions}>
                <Button variant="iconDanger" onClick={onDelete} aria-label="Delete step">✕</Button>
            </div>
        </div>
    )
}
