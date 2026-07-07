import cls from "@/pages/tract-canvas/components/ConditionBody/ConditionBody.module.css"

import {SchemaNode, TractCondition, TractStep, TractTool} from "@/processes/Tracts.ts"
import {Location, replaceStep} from "@/processes/tractSteps.ts"
import {buildSources} from "@/pages/tract-canvas/processes/inspectorSources.ts"
import {cn} from "@/app/utils/cn.ts"

import {Button} from "@vervstack/chures"
import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"
import {CloseIcon} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"
import NameField from "@/pages/tract-canvas/components/NameField/NameField.tsx"

const CONDITION_OPS = ["==", "!=", ">", "<", ">=", "<=", "contains", "glob", "regex"] as const

interface Props {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChangeSteps: (s: TractStep[]) => void
    onOpenAddBlock: (location: Location, index: number) => void
    enlarged: boolean
}

export default function ConditionBody({rootSteps, step, tools, triggerSchema, onChangeSteps, onOpenAddBlock, enlarged}: Props) {
    const sources = buildSources(rootSteps, step.id, tools, triggerSchema)
    const conditions = step.conditions ?? []

    function update(updater: (s: TractStep) => TractStep) {
        onChangeSteps(replaceStep(rootSteps, step.id, updater))
    }

    function updateCondition(i: number, patch: Partial<TractCondition>) {
        update(s => {
            const next = [...(s.conditions ?? [])]
            next[i] = {...next[i], ...patch}
            return {...s, conditions: next}
        })
    }

    return (
        <div className={cls.ConditionBodyContainer}>
            <Section title="Config">
                <NameField step={step} onRename={name => update(s => ({...s, name}))}/>
            </Section>
            <Section title="Conditions (AND)">
                {conditions.map((cond, i) => (
                    <div className={cn(cls.ConditionRow, enlarged && cls.ConditionRowEnlarged)} key={i}>
                        <TemplateInput value={cond.left} onChange={v => updateCondition(i, {left: v})} sources={sources} placeholder="left"/>
                        <div className={cls.ConditionOpRow}>
                            <select className={cls.Select} value={cond.op} onChange={e => updateCondition(i, {op: e.target.value as TractCondition["op"]})}>
                                {CONDITION_OPS.map(op => <option key={op} value={op}>{op}</option>)}
                            </select>
                            <Button
                                variant="iconDanger"
                                onClick={() => update(s => ({...s, conditions: (s.conditions ?? []).filter((_, xi) => xi !== i)}))}
                                aria-label="Remove condition"
                            >
                                <CloseIcon/>
                            </Button>
                        </div>
                        <TemplateInput value={cond.right} onChange={v => updateCondition(i, {right: v})} sources={sources} placeholder="right"/>
                    </div>
                ))}
                <Button variant="ghost" onClick={() => update(s => ({...s, conditions: [...(s.conditions ?? []), {left: "", op: "==", right: ""}]}))}>
                    + Add condition
                </Button>
            </Section>
            <Section title="Branches">
                <Button variant="ghost" onClick={() => onOpenAddBlock({parentId: step.id, branch: "then"}, step.then?.length ?? 0)}>
                    ✓ Add to true branch
                </Button>
                <Button variant="ghost" onClick={() => onOpenAddBlock({parentId: step.id, branch: "else"}, step.else?.length ?? 0)}>
                    ✗ Add to false branch
                </Button>
            </Section>
        </div>
    )
}
