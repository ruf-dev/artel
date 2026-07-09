import cls from "@/pages/tract-canvas/components/GroupBody/GroupBody.module.css"
import {SchemaNode, TractStep, TractTool} from "@/processes/Tracts.ts"
import TractStepTree from "@/components/TractStepTree/TractStepTree.tsx"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"

interface Props {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChangeSteps: (s: TractStep[]) => void
}

export default function GroupBody({rootSteps, step, tools, triggerSchema, onChangeSteps}: Props) {
    return (
        <div className={cls.GroupBodyContainer}>
            <Section title="Steps in this group">
                <TractStepTree
                    rootSteps={rootSteps}
                    location={{parentId: step.id, branch: "steps"}}
                    list={step.steps ?? []}
                    tools={tools}
                    triggerSchema={triggerSchema}
                    onChange={onChangeSteps}
                />
            </Section>
        </div>
    )
}
