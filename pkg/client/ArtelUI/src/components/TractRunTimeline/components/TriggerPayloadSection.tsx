import {useState} from "react"

import cls from "@/components/TractRunTimeline/TractRunTimeline.module.css"
import JsonBlock from "@/components/JsonBlock/JsonBlock.tsx"

interface Props {
    payload: unknown
}

export default function TriggerPayloadSection({payload}: Props) {
    const [open, setOpen] = useState(false)
    return (
        <div className={cls.Step}>
            <div className={cls.StepHeader} onClick={() => setOpen(o => !o)}>
                <span className={cls.StepName}>trigger</span>
                <span className={cls.StepType}>payload</span>
            </div>
            {open && (
                <div className={cls.StepBody}>
                    <JsonBlock label="Payload" value={payload}/>
                </div>
            )}
        </div>
    )
}
