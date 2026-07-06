import {useEffect, useState} from "react"

import cls from "@/pages/tract-canvas/components/TractCanvasInspector/TractCanvasInspector.module.css"

import {CanvasNode} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {SchemaNode, TractCondition, TractStep, TractTool, TractTriggerSummary} from "@/processes/Tracts.ts"
import {computeVisibleStepIds, findStepById} from "@/processes/tractTemplate.ts"
import {hasChildren, Location, removeStepAt, replaceStep} from "@/processes/tractSteps.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"

import {Button, ConfirmDialog} from "@vervstack/chures"
import TemplateInput, {TemplateSource} from "@/components/TemplateInput/TemplateInput.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import {ChevronRightIcon, CloseIcon} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"
import TractStepTree from "@/components/TractStepTree/TractStepTree.tsx"
import TriggerPanel from "@/components/TriggerPanel/TriggerPanel.tsx"

const BUILTIN_MCP = "artel"
const CONDITION_OPS = ["==", "!=", ">", "<", ">=", "<=", "contains", "glob", "regex"] as const

function buildSources(rootSteps: TractStep[], targetId: string, tools: TractTool[], triggerSchema?: SchemaNode): TemplateSource[] {
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

interface Props {
    node: CanvasNode | null
    rootSteps: TractStep[]
    tools: TractTool[]
    triggerSchema?: SchemaNode
    lastOutputByStepId: Record<string, unknown>
    tractUuid: string
    linkedTriggerSummaries: TractTriggerSummary[]
    onChangeSteps: (newRootSteps: TractStep[]) => void
    onOpenAddBlock: (location: Location, index: number) => void
    onClose: () => void
}

export default function TractCanvasInspector({node, rootSteps, tools, triggerSchema, lastOutputByStepId, tractUuid, linkedTriggerSummaries, onChangeSteps, onOpenAddBlock, onClose}: Props) {
    const [enlarged, setEnlarged] = useState(false)

    useEffect(() => {
        setEnlarged(false)
    }, [node?.id])

    return (
        <>
            {node && enlarged && (
                <div className={cls.Backdrop} onClick={() => setEnlarged(false)}/>
            )}
            <div className={`${cls.Panel} ${node ? cls.Open : ""} ${node && enlarged ? cls.Enlarged : ""}`}>
                {node && (
                    <Body
                        key={node.id}
                        node={node}
                        rootSteps={rootSteps}
                        tools={tools}
                        triggerSchema={triggerSchema}
                        lastOutputByStepId={lastOutputByStepId}
                        tractUuid={tractUuid}
                        linkedTriggerSummaries={linkedTriggerSummaries}
                        onChangeSteps={onChangeSteps}
                        onOpenAddBlock={onOpenAddBlock}
                        onClose={onClose}
                        enlarged={enlarged}
                        onToggleEnlarge={() => setEnlarged(e => !e)}
                    />
                )}
            </div>
        </>
    )
}

function Body({node, rootSteps, tools, triggerSchema, lastOutputByStepId, tractUuid, linkedTriggerSummaries, onChangeSteps, onOpenAddBlock, onClose, enlarged, onToggleEnlarge}: Omit<Props, "node"> & {
    node: CanvasNode
    enlarged: boolean
    onToggleEnlarge: () => void
}) {
    const step = node.step

    return (
        <>
            <div className={cls.Head}>
                <Button
                    variant="ghost"
                    className={cls.EnlargeBtn}
                    onClick={onToggleEnlarge}
                    aria-label={enlarged ? "Shrink inspector" : "Enlarge inspector"}
                    aria-pressed={enlarged}
                >
                    <ChevronRightIcon className={`${cls.EnlargeChevron} ${enlarged ? cls.EnlargeChevronOpen : ""}`}/>
                </Button>
                <div className={cls.Titles}>
                    <div className={cls.Title}>{node.kind === "trigger" ? "Trigger" : step?.name || step?.id || node.kind}</div>
                    <div className={cls.Sub}>{node.kind === "trigger" ? "trigger" : node.kind}</div>
                </div>
                <Button variant="iconDanger" className={cls.CloseBtn} onClick={onClose} aria-label="Close inspector">
                    <CloseIcon/>
                </Button>
            </div>
            <div className={cls.Body}>
                {node.kind === "trigger" && (
                    <TriggerPanel tractUuid={tractUuid} linkedTriggerSummaries={linkedTriggerSummaries}/>
                )}
                {step && node.kind === "action" && (
                    <ActionBody
                        rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema}
                        lastOutput={lastOutputByStepId[step.id]}
                        onChangeSteps={onChangeSteps}
                    />
                )}
                {step && node.kind === "condition" && (
                    <ConditionBody
                        rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema}
                        onChangeSteps={onChangeSteps} onOpenAddBlock={onOpenAddBlock}
                    />
                )}
                {step && node.kind === "parallel" && (
                    <ParallelBody step={step} onOpenAddBlock={onOpenAddBlock}/>
                )}
                {step && node.kind === "group" && (
                    <GroupBody rootSteps={rootSteps} step={step} tools={tools} triggerSchema={triggerSchema} onChangeSteps={onChangeSteps}/>
                )}
                {step && (
                    <DangerZone rootSteps={rootSteps} step={step} location={node.location} onChangeSteps={onChangeSteps} onClose={onClose}/>
                )}
                <Section title="Flow">
                    <Button variant="ghost" onClick={() => onOpenAddBlock(node.location, node.index + 1)}>
                        {node.kind === "trigger" ? "+ Add first step" : "+ Add step after"}
                    </Button>
                </Section>
            </div>
        </>
    )
}

function Section({title, children}: { title: string; children: React.ReactNode }) {
    return (
        <div className={cls.Section}>
            <div className={cls.SectionTitle}>{title}</div>
            {children}
        </div>
    )
}

function NameField({step, onRename}: { step: TractStep; onRename: (name: string) => void }) {
    const [name, setName] = useState(step.name || step.id)
    const pattern = /^[a-z][a-z0-9_]*$/
    const invalid = !pattern.test(name) || name === "trigger"

    function commit() {
        if (invalid || name === step.name) return
        onRename(name)
    }

    return (
        <div className={cls.Row}>
            <span className={cls.Key}>name</span>
            <input
                className={`${cls.Input} ${invalid ? cls.InputInvalid : ""}`}
                value={name}
                onChange={e => setName(e.target.value)}
                onBlur={commit}
                data-tooltip-id={invalid ? "root-tooltip" : undefined}
                data-tooltip-content={invalid ? "Must match ^[a-z][a-z0-9_]*$ and not be \"trigger\"" : undefined}
            />
        </div>
    )
}

function ActionBody({rootSteps, step, tools, triggerSchema, lastOutput, onChangeSteps}: {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    lastOutput: unknown
    onChangeSteps: (s: TractStep[]) => void
}) {
    const {momCandidates, fetchMomCandidates} = useMcpKeys()
    const tool = tools.find(t => t.mcp === step.mcp && t.tool === step.tool)
    const sources = buildSources(rootSteps, step.id, tools, triggerSchema)

    useEffect(() => {
        void fetchMomCandidates()
    }, [fetchMomCandidates])

    function update(updater: (s: TractStep) => TractStep) {
        onChangeSteps(replaceStep(rootSteps, step.id, updater))
    }

    const candidate = momCandidates.find(c => c.name === step.mcp)
    const stepConnections = candidate?.connections ?? []

    return (
        <>
            <Section title="Config">
                <NameField step={step} onRename={name => update(s => ({...s, name}))}/>
                <div className={cls.Row}>
                    <span className={cls.Key}>tool</span>
                    <span className={cls.Val}>{step.mcp}.{step.tool}</span>
                </div>
                {step.mcp !== BUILTIN_MCP && (
                    <div className={cls.Row}>
                        <span className={cls.Key}>connection</span>
                        {stepConnections.length > 0 ? (
                            <select
                                className={`${cls.Select} ${!step.connection_uuid ? cls.InputInvalid : ""}`}
                                value={step.connection_uuid ?? ""}
                                onChange={e => update(s => ({...s, connection_uuid: e.target.value}))}
                            >
                                <option value="">— select —</option>
                                {stepConnections.map(c => (
                                    <option key={c.id} value={c.id}>{connectionLabel(c)}</option>
                                ))}
                            </select>
                        ) : (
                            <p className={cls.Empty}>No connections available for {step.mcp}.</p>
                        )}
                    </div>
                )}
            </Section>
            {tool && Object.keys(tool.inputSchema.properties).length > 0 && (
                <Section title="Inputs">
                    {Object.entries(tool.inputSchema.properties).map(([name, def]) => (
                        <div className={cls.ParamRow} key={name}>
                            <span className={cls.Key}>{name}{tool.inputSchema.required?.includes(name) ? " *" : ""}</span>
                            {def.enum?.length ? (
                                <select
                                    className={cls.Select}
                                    value={step.params?.[name] ?? ""}
                                    onChange={e => update(s => ({...s, params: {...s.params, [name]: e.target.value}}))}
                                >
                                    <option value="">—</option>
                                    {def.enum.map(v => <option key={v} value={v}>{v}</option>)}
                                </select>
                            ) : (
                                <TemplateInput
                                    value={step.params?.[name] ?? ""}
                                    onChange={v => update(s => ({...s, params: {...s.params, [name]: v}}))}
                                    sources={sources}
                                    placeholder={def.type}
                                />
                            )}
                        </div>
                    ))}
                </Section>
            )}
            {lastOutput !== undefined && lastOutput !== null && typeof lastOutput === "object" && (
                <Section title="Last output">
                    {Object.entries(lastOutput as Record<string, unknown>).map(([k, v]) => (
                        <div className={cls.Row} key={k}>
                            <span className={cls.Key}>{k}</span>
                            <span className={`${cls.Val} ${cls.ValDim}`}>{String(v)}</span>
                        </div>
                    ))}
                </Section>
            )}
        </>
    )
}

function ConditionBody({rootSteps, step, tools, triggerSchema, onChangeSteps, onOpenAddBlock}: {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChangeSteps: (s: TractStep[]) => void
    onOpenAddBlock: (location: Location, index: number) => void
}) {
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
        <>
            <Section title="Config">
                <NameField step={step} onRename={name => update(s => ({...s, name}))}/>
            </Section>
            <Section title="Conditions (AND)">
                {conditions.map((cond, i) => (
                    <div className={cls.ConditionRow} key={i}>
                        <TemplateInput value={cond.left} onChange={v => updateCondition(i, {left: v})} sources={sources} placeholder="left"/>
                        <select className={cls.Select} value={cond.op} onChange={e => updateCondition(i, {op: e.target.value as TractCondition["op"]})}>
                            {CONDITION_OPS.map(op => <option key={op} value={op}>{op}</option>)}
                        </select>
                        <TemplateInput value={cond.right} onChange={v => updateCondition(i, {right: v})} sources={sources} placeholder="right"/>
                        <Button
                            variant="iconDanger"
                            onClick={() => update(s => ({...s, conditions: (s.conditions ?? []).filter((_, xi) => xi !== i)}))}
                            aria-label="Remove condition"
                        >
                            <CloseIcon/>
                        </Button>
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
        </>
    )
}

function ParallelBody({step, onOpenAddBlock}: {
    step: TractStep
    onOpenAddBlock: (location: Location, index: number) => void
}) {
    const lanes = step.steps ?? []
    return (
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
    )
}

function GroupBody({rootSteps, step, tools, triggerSchema, onChangeSteps}: {
    rootSteps: TractStep[]
    step: TractStep
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onChangeSteps: (s: TractStep[]) => void
}) {
    return (
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
    )
}

function DangerZone({rootSteps, step, location, onChangeSteps, onClose}: {
    rootSteps: TractStep[]
    step: TractStep
    location: Location
    onChangeSteps: (s: TractStep[]) => void
    onClose: () => void
}) {
    const {OpenDialog, CloseDialog} = useDialog()

    function doDelete() {
        onChangeSteps(removeStepAt(rootSteps, location, step.id))
        onClose()
    }

    function handleDelete() {
        if (hasChildren(step)) {
            OpenDialog(
                <ConfirmDialog
                    title="Delete step"
                    message={`Delete "${step.name || step.id}"? Its nested branches/lanes will be removed too.`}
                    confirmLabel="Delete"
                    danger
                    onClose={CloseDialog}
                    onConfirm={async () => doDelete()}
                />
            )
            return
        }
        doDelete()
    }

    return (
        <Section title="Danger zone">
            <Button variant="danger" onClick={handleDelete}>Delete step</Button>
        </Section>
    )
}
