import {useState} from "react"

import cls from "@/components/TractRunTimeline/TractRunTimeline.module.css"

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
                    <pre className={cls.Json}>{JSON.stringify(payload, null, 2)}</pre>
                </div>
            )}
        </div>
    )
}
