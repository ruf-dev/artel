import {useState} from "react"

import cls from "@/pages/tract-canvas/components/NameField/NameField.module.css"

import {TractStep} from "@/processes/Tracts.ts"
import {cn} from "@/app/utils/cn.ts"

interface Props {
    step: TractStep
    onRename: (name: string) => void
}

export default function NameField({step, onRename}: Props) {
    const [name, setName] = useState(step.name || step.id)
    const pattern = /^[a-z][a-z0-9_]*$/
    const invalid = !pattern.test(name) || name === "trigger"

    function commit() {
        if (invalid || name === step.name) return
        onRename(name)
    }

    return (
        <div className={cls.NameFieldContainer}>
            <span className={cls.Key}>name</span>
            <input
                className={cn(cls.Input, invalid && cls.InputInvalid)}
                value={name}
                onChange={e => setName(e.target.value)}
                onBlur={commit}
                data-tooltip-id={invalid ? "root-tooltip" : undefined}
                data-tooltip-content={invalid ? "Must match ^[a-z][a-z0-9_]*$ and not be \"trigger\"" : undefined}
            />
        </div>
    )
}
