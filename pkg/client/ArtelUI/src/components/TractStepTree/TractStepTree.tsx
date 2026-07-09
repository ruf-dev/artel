import {Button, ConfirmDialog} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {cn} from "@/app/utils/cn.ts"
import TemplateInput from "@/components/TemplateInput/TemplateInput.tsx"
import ActionCard from "@/components/TractStepTree/components/ActionCard.tsx"
import CardHeader from "@/components/TractStepTree/components/CardHeader.tsx"
import InsertRow from "@/components/TractStepTree/components/InsertRow.tsx"
import {buildSourcesFor, collectIdsFromRoot} from "@/components/TractStepTree/processes/tractStepTreeHelpers.ts"
import StepPickerDialog, {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"
import {SchemaNode, TractCondition, TractStep, TractTool} from "@/processes/Tracts.ts"
import {appendStep, buildStepFromDraft, hasChildren, insertStepAt, Location, removeStepAt, replaceStep} from "@/processes/tractSteps.ts"
import cls from "@/components/TractStepTree/TractStepTree.module.css"

const CONDITION_OPS = ["==", "!=", ">", "<", ">=", "<=", "contains", "glob", "regex"] as const

interface TractStepTreeProps {
    rootSteps: TractStep[]
    location: Location
    list: TractStep[]
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChange: (newRootSteps: TractStep[]) => void
    direction?: "column" | "row"
}

export default function TractStepTree(props: TractStepTreeProps) {
    const {OpenDialog} = useDialog()
    const direction = props.direction ?? "column"

    function openPicker(index: number) {
        OpenDialog(
            <StepPickerDialog
                onConfirm={(draft: StepDraft) => {
                    const existingIds = collectIdsFromRoot(props.rootSteps)
                    const newStep = buildStepFromDraft(draft, existingIds)
                    props.onChange(insertStepAt(props.rootSteps, props.location, index, newStep))
                }}
            />
        )
    }

    function appendLane() {
        OpenDialog(
            <StepPickerDialog
                onConfirm={(draft: StepDraft) => {
                    const existingIds = collectIdsFromRoot(props.rootSteps)
                    const newStep = buildStepFromDraft(draft, existingIds)
                    props.onChange(appendStep(props.rootSteps, props.location, newStep))
                }}
            />
        )
    }

    if (direction === "row") {
        return (
            <div className={cls.ListRow}>
                {props.list.map(step => (
                    <div className={cls.Lane} key={step.id}>
                        <StepCard
                            rootSteps={props.rootSteps}
                            location={props.location}
                            step={step}
                            tools={props.tools}
                            triggerSchema={props.triggerSchema}
                            onChange={props.onChange}
                        />
                    </div>
                ))}
                <Button variant="ghost" onClick={appendLane}>+ Add lane</Button>
            </div>
        )
    }

    return (
        <div className={cls.List}>
            <InsertRow onClick={() => openPicker(0)}/>
            {props.list.map((step, i) => (
                <div className={cls.ListItem} key={step.id}>
                    <StepCard
                        rootSteps={props.rootSteps}
                        location={props.location}
                        step={step}
                        tools={props.tools}
                        triggerSchema={props.triggerSchema}
                        onChange={props.onChange}
                    />
                    <InsertRow onClick={() => openPicker(i + 1)}/>
                </div>
            ))}
        </div>
    )
}

// Not extracted to components/: StepCard renders TractStepTree recursively (for
// "parallel"/"group" nested lanes), and moving it to its own file would create
// StepCard.tsx -> TractStepTree.tsx -> StepCard.tsx, an import-x/no-cycle error.
// Deliberately colocated with TractStepTree despite the react/no-multi-comp warning.
function StepCard({rootSteps, location, step, tools, triggerSchema, onChange}: {
    rootSteps: TractStep[]
    location: Location
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChange: (newRootSteps: TractStep[]) => void
}) {
    const {OpenDialog, CloseDialog} = useDialog()

    function updateStep(updater: (s: TractStep) => TractStep) {
        onChange(replaceStep(rootSteps, step.id, updater))
    }

    function deleteStep() {
        if (hasChildren(step)) {
            OpenDialog(
                <ConfirmDialog
                    title="Delete step"
                    message={`Delete "${step.name || step.id}"? Its nested branches/lanes will be removed too.`}
                    confirmLabel="Delete"
                    danger
                    onClose={CloseDialog}
                    onConfirm={() => onChange(removeStepAt(rootSteps, location, step.id))}
                />
            )
            return
        }
        onChange(removeStepAt(rootSteps, location, step.id))
    }

    if (step.type === "action") {
        return <ActionCard rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema} onUpdate={updateStep} onDelete={deleteStep}/>
    }
    if (step.type === "condition") {
        return (
            <ConditionCard
                rootSteps={rootSteps}
                step={step}
                tools={tools}
                triggerSchema={triggerSchema}
                onUpdate={updateStep}
                onDelete={deleteStep}
                onChange={onChange}
            />
        )
    }
    if (step.type === "parallel") {
        return (
            <div className={cls.Card}>
                <CardHeader step={step} onUpdate={updateStep} onDelete={deleteStep}/>
                <div className={cls.Container}>
                    <TractStepTree
                        rootSteps={rootSteps}
                        location={{parentId: step.id, branch: "steps"}}
                        list={step.steps ?? []}
                        tools={tools}
                        triggerSchema={triggerSchema}
                        onChange={onChange}
                        direction="row"
                    />
                </div>
            </div>
        )
    }
    // group
    return (
        <div className={cls.Card}>
            <CardHeader step={step} onUpdate={updateStep} onDelete={deleteStep}/>
            <div className={cls.Container}>
                <TractStepTree
                    rootSteps={rootSteps}
                    location={{parentId: step.id, branch: "steps"}}
                    list={step.steps ?? []}
                    tools={tools}
                    triggerSchema={triggerSchema}
                    onChange={onChange}
                />
            </div>
        </div>
    )
}

interface ConditionCardProps {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onUpdate: (updater: (s: TractStep) => TractStep) => void
    onDelete: () => void
    onChange: (newRootSteps: TractStep[]) => void
}

// Same recursion constraint as StepCard above: ConditionCard renders TractStepTree
// for its "then"/"else" branches, so it can't move out without an import cycle.
function ConditionCard(props: ConditionCardProps) {
    const sources = buildSourcesFor(props.rootSteps, props.step.id, props.tools, props.triggerSchema)
    const conditions = props.step.conditions ?? []

    function updateCondition(index: number, patch: Partial<TractCondition>) {
        props.onUpdate(s => {
            const next = [...(s.conditions ?? [])]
            next[index] = {...next[index], ...patch}
            return {...s, conditions: next}
        })
    }

    function addRow() {
        props.onUpdate(s => ({...s, conditions: [...(s.conditions ?? []), {left: "", op: "==", right: ""}]}))
    }

    function removeRow(index: number) {
        props.onUpdate(s => ({...s, conditions: (s.conditions ?? []).filter((_, i) => i !== index)}))
    }

    return (
        <div className={cls.Card}>
            <CardHeader step={props.step} onUpdate={props.onUpdate} onDelete={props.onDelete}/>
            <div className={cls.ConditionRows}>
                {conditions.map((cond, i) => (
                    <div className={cls.ConditionRow} key={i}>
                        <TemplateInput value={cond.left} onChange={v => updateCondition(i, {left: v})} sources={sources} placeholder="left"/>
                        <select
                            className={cls.ConditionOp}
                            value={cond.op}
                            onChange={e => updateCondition(i, {op: e.target.value as TractCondition["op"]})}
                        >
                            {CONDITION_OPS.map(op => <option key={op} value={op}>{op}</option>)}
                        </select>
                        <TemplateInput value={cond.right} onChange={v => updateCondition(i, {right: v})} sources={sources} placeholder="right"/>
                        <Button variant="iconDanger" className={cls.RemoveRowBtn} onClick={() => removeRow(i)} aria-label="Remove condition">✕</Button>
                    </div>
                ))}
                <Button variant="ghost" onClick={addRow}>+ Add condition (AND)</Button>
            </div>
            <div className={cls.Branches}>
                <div className={cls.Branch}>
                    <span className={cn(cls.BranchLabel, cls.BranchLabelTrue)}>✓ true</span>
                    <TractStepTree
                        rootSteps={props.rootSteps}
                        location={{parentId: props.step.id, branch: "then"}}
                        list={props.step.then ?? []}
                        tools={props.tools}
                        triggerSchema={props.triggerSchema}
                        onChange={props.onChange}
                    />
                </div>
                <div className={cls.Branch}>
                    <span className={cn(cls.BranchLabel, cls.BranchLabelFalse)}>✗ false</span>
                    <TractStepTree
                        rootSteps={props.rootSteps}
                        location={{parentId: props.step.id, branch: "else"}}
                        list={props.step.else ?? []}
                        tools={props.tools}
                        triggerSchema={props.triggerSchema}
                        onChange={props.onChange}
                    />
                </div>
            </div>
        </div>
    )
}
