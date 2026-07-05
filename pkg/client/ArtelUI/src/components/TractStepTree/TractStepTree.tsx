import {useState} from "react"

import cls from "@/components/TractStepTree/TractStepTree.module.css"

import {SchemaNode, SchemaProperty, TractCondition, TractStep, TractTool} from "@/processes/Tracts.ts"
import {computeVisibleStepIds, findStepById} from "@/processes/tractTemplate.ts"
import {appendStep, buildStepFromDraft, hasChildren, insertStepAt, Location, removeStepAt, replaceStep} from "@/processes/tractSteps.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"

import Button from "@/components/shared/Button/Button.tsx"
import ConfirmDialog from "@/components/ConfirmDialog/ConfirmDialog.tsx"
import StepPickerDialog, {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"
import TemplateInput, {TemplateSource} from "@/components/TemplateInput/TemplateInput.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"

const CONDITION_OPS = ["==", "!=", ">", "<", ">=", "<=", "contains", "glob", "regex"] as const

function buildSourcesFor(rootSteps: TractStep[], targetId: string, tools: TractTool[], triggerSchema?: SchemaNode): TemplateSource[] {
    const visibleIds = computeVisibleStepIds(rootSteps, targetId)
    const sources: TemplateSource[] = [{id: "trigger", label: "trigger", schema: triggerSchema}]

    for (const id of visibleIds) {
        const step = findStepById(rootSteps, id)
        if (!step) continue
        if (step.type === "action") {
            const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
            sources.push({id, label: step.name || id, schema: tool?.outputSchema})
        } else if (step.type === "condition") {
            sources.push({id, label: step.name || id, schema: {properties: {result: {type: "boolean", description: "condition result"}}}})
        }
    }

    return sources
}

interface TractStepTreeProps {
    rootSteps: TractStep[]
    location: Location
    list: TractStep[]
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChange: (newRootSteps: TractStep[]) => void
    direction?: "column" | "row"
}

export default function TractStepTree({rootSteps, location, list, tools, triggerSchema, onChange, direction = "column"}: TractStepTreeProps) {
    const {OpenDialog} = useDialog()

    function openPicker(index: number) {
        OpenDialog(
            <StepPickerDialog
                onConfirm={(draft: StepDraft) => {
                    const existingIds = collectIdsFromRoot(rootSteps)
                    const newStep = buildStepFromDraft(draft, existingIds)
                    onChange(insertStepAt(rootSteps, location, index, newStep))
                }}
            />
        )
    }

    function appendLane() {
        OpenDialog(
            <StepPickerDialog
                onConfirm={(draft: StepDraft) => {
                    const existingIds = collectIdsFromRoot(rootSteps)
                    const newStep = buildStepFromDraft(draft, existingIds)
                    onChange(appendStep(rootSteps, location, newStep))
                }}
            />
        )
    }

    if (direction === "row") {
        return (
            <div className={cls.ListRow}>
                {list.map(step => (
                    <div className={cls.Lane} key={step.id}>
                        <StepCard
                            rootSteps={rootSteps}
                            location={location}
                            step={step}
                            tools={tools}
                            triggerSchema={triggerSchema}
                            onChange={onChange}
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
            {list.map((step, i) => (
                <div className={cls.ListItem} key={step.id}>
                    <StepCard
                        rootSteps={rootSteps}
                        location={location}
                        step={step}
                        tools={tools}
                        triggerSchema={triggerSchema}
                        onChange={onChange}
                    />
                    <InsertRow onClick={() => openPicker(i + 1)}/>
                </div>
            ))}
        </div>
    )
}

function collectIdsFromRoot(rootSteps: TractStep[]): Set<string> {
    const ids = new Set<string>()
    function walk(list: TractStep[]) {
        for (const s of list) {
            ids.add(s.id)
            if (s.then) walk(s.then)
            if (s.else) walk(s.else)
            if (s.steps) walk(s.steps)
        }
    }
    walk(rootSteps)
    return ids
}

function InsertRow({onClick}: { onClick: () => void }) {
    return (
        <div className={cls.InsertRow}>
            <span className={cls.InsertLine}/>
            <Button variant="ghost" className={cls.InsertBtn} onClick={onClick} aria-label="Add step">+</Button>
            <span className={cls.InsertLine}/>
        </div>
    )
}

function StepCard({rootSteps, location, step, tools, triggerSchema, onChange}: {
    rootSteps: TractStep[]
    location: Location
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChange: (newRootSteps: TractStep[]) => void
}) {
    const {OpenDialog} = useDialog()

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

const STEP_ID_PATTERN = /^[a-z][a-z0-9_]*$/

function CardHeader({step, onUpdate, onDelete, right}: {
    step: TractStep
    onUpdate: (updater: (s: TractStep) => TractStep) => void
    onDelete: () => void
    right?: React.ReactNode
}) {
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
            <span className={`${cls.CardIcon} ${iconClass}`}>{step.type.slice(0, 1).toUpperCase()}</span>
            <input
                className={`${cls.NameInput} ${nameInvalid ? cls.NameInputInvalid : ""}`}
                value={name}
                onChange={e => setName(e.target.value)}
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

function ActionCard({rootSteps, step, tools, triggerSchema, onUpdate, onDelete}: {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onUpdate: (updater: (s: TractStep) => TractStep) => void
    onDelete: () => void
}) {
    const {connections} = useExternalConnections()
    const [outputOpen, setOutputOpen] = useState(false)

    const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
    const conn = connections.find(c => c.id === step.connection_uuid)
    const sources = buildSourcesFor(rootSteps, step.id, tools, triggerSchema)

    function setParam(name: string, value: string) {
        onUpdate(s => ({...s, params: {...s.params, [name]: value}}))
    }

    return (
        <div className={cls.Card}>
            <CardHeader
                step={step}
                onUpdate={onUpdate}
                onDelete={onDelete}
                right={
                    <>
                        <span className={cls.ConnChip}>{step.mcp}.{step.tool}</span>
                        {conn && <span className={cls.ConnChip}>{connectionLabel(conn)}</span>}
                    </>
                }
            />
            {tool && Object.keys(tool.inputSchema.properties).length > 0 && (
                <div className={cls.ParamsSection}>
                    {Object.entries(tool.inputSchema.properties).map(([name, def]) => (
                        <ParamRow
                            key={name}
                            name={name}
                            def={def}
                            required={tool.inputSchema.required?.includes(name) ?? false}
                            value={step.params?.[name] ?? ""}
                            sources={sources}
                            onChange={setParam}
                        />
                    ))}
                </div>
            )}
            {tool && Object.keys(tool.outputSchema.properties).length > 0 && (
                <div className={cls.Collapsible}>
                    <div className={cls.CollapsibleHeader} onClick={() => setOutputOpen(o => !o)}>
                        <span>{outputOpen ? "▾" : "▸"} Output</span>
                    </div>
                    {outputOpen && <SchemaTree schema={tool.outputSchema}/>}
                </div>
            )}
        </div>
    )
}

function ParamRow({name, def, required, value, sources, onChange}: {
    name: string
    def: SchemaProperty
    required: boolean
    value: string
    sources: TemplateSource[]
    onChange: (name: string, value: string) => void
}) {
    return (
        <div className={cls.ParamRow}>
            <span className={cls.ParamName}>{name}{required ? " *" : ""}</span>
            {def.description && <span className={cls.ParamDesc}>{def.description}</span>}
            {def.enum?.length ? (
                <select className={cls.PlainSelect} value={value} onChange={e => onChange(name, e.target.value)}>
                    <option value="">—</option>
                    {def.enum.map(v => <option key={v} value={v}>{v}</option>)}
                </select>
            ) : (
                <TemplateInput value={value} onChange={v => onChange(name, v)} sources={sources} placeholder={def.type}/>
            )}
        </div>
    )
}

function SchemaTree({schema}: { schema: SchemaNode }) {
    return (
        <div className={cls.SchemaTree}>
            {Object.entries(schema.properties).map(([name, prop]) => (
                <SchemaFieldRow key={name} name={name} prop={prop} depth={0}/>
            ))}
        </div>
    )
}

function SchemaFieldRow({name, prop, depth}: { name: string; prop: SchemaProperty; depth: number }) {
    return (
        <>
            <div className={cls.SchemaField} style={{paddingLeft: `${depth * 0.75}rem`}}>
                <span className={cls.SchemaFieldName}>{name}</span>
                <span className={cls.SchemaFieldType}>{prop.type}</span>
            </div>
            {prop.type === "object" && prop.properties && Object.entries(prop.properties).map(([n, p]) => (
                <SchemaFieldRow key={n} name={n} prop={p} depth={depth + 1}/>
            ))}
            {prop.type === "array" && prop.items && (
                <SchemaFieldRow name="[0]" prop={prop.items} depth={depth + 1}/>
            )}
        </>
    )
}

function ConditionCard({rootSteps, step, tools, triggerSchema, onUpdate, onDelete, onChange}: {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onUpdate: (updater: (s: TractStep) => TractStep) => void
    onDelete: () => void
    onChange: (newRootSteps: TractStep[]) => void
}) {
    const sources = buildSourcesFor(rootSteps, step.id, tools, triggerSchema)
    const conditions = step.conditions ?? []

    function updateCondition(index: number, patch: Partial<TractCondition>) {
        onUpdate(s => {
            const next = [...(s.conditions ?? [])]
            next[index] = {...next[index], ...patch}
            return {...s, conditions: next}
        })
    }

    function addRow() {
        onUpdate(s => ({...s, conditions: [...(s.conditions ?? []), {left: "", op: "==", right: ""}]}))
    }

    function removeRow(index: number) {
        onUpdate(s => ({...s, conditions: (s.conditions ?? []).filter((_, i) => i !== index)}))
    }

    return (
        <div className={cls.Card}>
            <CardHeader step={step} onUpdate={onUpdate} onDelete={onDelete}/>
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
                    <span className={`${cls.BranchLabel} ${cls.BranchLabelTrue}`}>✓ true</span>
                    <TractStepTree
                        rootSteps={rootSteps}
                        location={{parentId: step.id, branch: "then"}}
                        list={step.then ?? []}
                        tools={tools}
                        triggerSchema={triggerSchema}
                        onChange={onChange}
                    />
                </div>
                <div className={cls.Branch}>
                    <span className={`${cls.BranchLabel} ${cls.BranchLabelFalse}`}>✗ false</span>
                    <TractStepTree
                        rootSteps={rootSteps}
                        location={{parentId: step.id, branch: "else"}}
                        list={step.else ?? []}
                        tools={tools}
                        triggerSchema={triggerSchema}
                        onChange={onChange}
                    />
                </div>
            </div>
        </div>
    )
}
