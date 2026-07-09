import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/ParallelBody/ParallelBody.module.css"
import {TractStep} from "@/processes/Tracts.ts"
import {Location} from "@/processes/tractSteps.ts"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"

interface Props {
    step: TractStep
    onOpenAddBlock: (location: Location, index: number) => void
}

export default function ParallelBody({step, onOpenAddBlock}: Props) {
    const lanes = step.steps ?? []

    return (
        <div className={cls.ParallelBodyContainer}>
            <Section title="Lanes">
                {lanes.length === 0 && <p className={cls.Empty}>No lanes yet.</p>}
                {lanes.map((lane, i) => (
                    <div className={cls.Row} key={lane.id}>
                        <span className={cls.Key}>lane {i + 1}</span>
                        <span className={cls.Val}>{lane.name || lane.id} ({lane.type})</span>
                    </div>
                ))}
                <Button variant="ghost" onClick={() => onOpenAddBlock({parentId: step.id, branch: "steps"}, lanes.length)}>
                    + Add lane
                </Button>
            </Section>
        </div>
    )
}
