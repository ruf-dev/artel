import {useEffect, useMemo, useState} from "react"

import cls from "@/pages/tract-canvas/components/TractCanvasBuilder/TractCanvasBuilder.module.css"

import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import Button from "@/components/shared/Button/Button.tsx"
import TractBlockPicker from "@/pages/tract-canvas/components/TractBlockPicker/TractBlockPicker.tsx"
import {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"

import TractCanvasArea from "@/pages/tract-canvas/components/TractCanvasArea/TractCanvasArea.tsx"
import TractCanvasInspector from "@/pages/tract-canvas/components/TractCanvasInspector/TractCanvasInspector.tsx"
import TractCanvasLogPanel from "@/pages/tract-canvas/components/TractCanvasLogPanel/TractCanvasLogPanel.tsx"
import {NodeStatus} from "@/pages/tract-canvas/components/TractCanvasNode/TractCanvasNode.tsx"
import RunStatusBadge from "@/pages/tract-canvas/components/RunStatusBadge/RunStatusBadge.tsx"
import RunButton from "@/pages/tract-canvas/components/RunButton/RunButton.tsx"

import {buildStepFromDraft, collectAllStepIds, insertStepAt, Location, ROOT_LOCATION} from "@/processes/tractSteps.ts"
import {layoutTract, TRIGGER_NODE_ID} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {Tract, TractDefinition, TractRun, TractTool, Trigger} from "@/processes/Tracts.ts"

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms))

interface Props {
    tract: Tract
    tools: TractTool[]
    triggers: Trigger[]
    runs: TractRun[]
    onBack: () => void
}

export default function TractCanvasBuilder({tract, tools, triggers, runs, onBack}: Props) {
    const {updateTract, runTract, fetchRuns, currentRun, currentRunSteps} = useTracts()
    const bakeError = useBakeError()
    const {OpenDialog} = useDialog()

    const [seeded, setSeeded] = useState(false)
    const [name, setName] = useState(tract.name)
    const [definition, setDefinition] = useState<TractDefinition>(tract.definition)
    const [saving, setSaving] = useState(false)
    const [warnings, setWarnings] = useState<string[]>([])
    const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null)
    const [logOpen, setLogOpen] = useState(false)
    const [selectedRunUuid, setSelectedRunUuid] = useState<string | null>(null)
    const [running, setRunning] = useState(false)

    useEffect(() => {
        if (!seeded) {
            setName(tract.name)
            setDefinition(tract.definition)
            setSeeded(true)
        }
    }, [seeded, tract])

    const isDirty = seeded && (name !== tract.name || JSON.stringify(definition) !== JSON.stringify(tract.definition))

    const layout = useMemo(() => layoutTract(definition.steps), [definition])
    const selectedNode = layout.nodes.find(n => n.id === selectedNodeId) ?? null

    const linkedSummaries = tract.triggers ?? []
    const firstLinked = linkedSummaries[0]
    const triggerInfo = firstLinked ? {name: firstLinked.name, kind: firstLinked.kind, source: firstLinked.source} : undefined
    const triggerSchema = triggers.find(t => linkedSummaries.some(l => l.uuid === t.uuid) && Object.keys(t.payloadSchema.properties).length > 0)?.payloadSchema

    const activeRunSteps = useMemo(
        () => selectedRunUuid && currentRun?.uuid === selectedRunUuid ? currentRunSteps : [],
        [selectedRunUuid, currentRun, currentRunSteps],
    )
    const lastOutputByStepId = useMemo(() => {
        const m: Record<string, unknown> = {}
        for (const s of activeRunSteps) {
            if (s.status === "done" && s.output && typeof s.output === "object") m[s.stepId] = s.output
        }
        return m
    }, [activeRunSteps])

    function nodeStatus(id: string): NodeStatus {
        if (id === TRIGGER_NODE_ID) return activeRunSteps.length > 0 ? "ok" : "idle"
        const step = activeRunSteps.find(s => s.stepId === id)
        if (!step) return "idle"
        if (step.status === "done") return "ok"
        if (step.status === "failed") return "err"
        return "running"
    }

    const runningEdgeIds = useMemo(() => {
        if (!running) return new Set<string>()
        return new Set(layout.edges.map(e => e.id))
    }, [running, layout])

    function handleSave() {
        setSaving(true)
        updateTract(tract.uuid, name, tract.description, definition)
            .then(({warnings: w}) => setWarnings(w))
            .catch(err => bakeError("Failed to save tract", err))
            .finally(() => setSaving(false))
    }

    function handleChangeSteps(steps: TractDefinition["steps"]) {
        setDefinition({steps})
    }

    function openAddBlock(location: Location, index: number) {
        OpenDialog(
            <TractBlockPicker
                onConfirm={(draft: StepDraft) => {
                    const existingIds = collectAllStepIds(definition.steps)
                    const newStep = buildStepFromDraft(draft, existingIds)
                    setDefinition(d => ({steps: insertStepAt(d.steps, location, index, newStep)}))
                }}
            />
        )
    }

    function handleRun() {
        setRunning(true)
        runTract(tract.uuid, {})
            .then(() => sleep(900))
            .then(() => fetchRuns(tract.uuid))
            .then(() => sleep(1200))
            .then(() => fetchRuns(tract.uuid))
            .then(() => {
                const latest = runsByLatest()
                if (latest) {
                    setSelectedRunUuid(latest.uuid)
                    setLogOpen(true)
                }
            })
            .catch(err => bakeError("Failed to run tract", err))
            .finally(() => setRunning(false))
    }

    function runsByLatest(): TractRun | undefined {
        // Read the live store rather than the `runs` prop — this closure outlives the render
        // that created it (it runs after two awaited fetches), so the prop would be stale.
        return useTracts.getState().runsByTract[tract.uuid]?.[0]
    }

    return (
        <div className={cls.Root}>
            <div className={cls.Bar}>
                <Button variant="ghost" onClick={onBack}>← Tracts</Button>
                <span className={cls.Divider}/>
                <input className={cls.NameInput} value={name} onChange={e => setName(e.target.value)}/>
                {isDirty && <span className={cls.DirtyDot} title="Unsaved changes"/>}
                <div className={cls.BarRight}>
                    <RunStatusBadge running={running} lastRunStatus={runs[0]?.status}/>
                    <Button variant="ghost" onClick={() => setLogOpen(o => !o)}>Logs</Button>
                    <Button variant="ghost" onClick={() => openAddBlock(ROOT_LOCATION, definition.steps.length)}>+ Add block</Button>
                    {isDirty && (
                        <Button onClick={handleSave} disabled={saving}>{saving ? "Saving…" : "Save"}</Button>
                    )}
                    <RunButton running={running} onClick={handleRun}/>
                </div>
            </div>
            {warnings.length > 0 && (
                <div className={cls.Warnings}>
                    {warnings.map((w, i) => <div key={i} className={cls.WarningRow}>{w}</div>)}
                </div>
            )}
            <div className={cls.Main}>
                <TractCanvasArea
                    layout={layout}
                    tools={tools}
                    triggerInfo={triggerInfo}
                    selectedNodeId={selectedNodeId}
                    onSelectNode={id => setSelectedNodeId(id)}
                    onBackgroundClick={() => setSelectedNodeId(null)}
                    nodeStatus={nodeStatus}
                    runningEdgeIds={runningEdgeIds}
                />
                <TractCanvasInspector
                    node={selectedNode}
                    rootSteps={definition.steps}
                    tools={tools}
                    triggerSchema={triggerSchema}
                    lastOutputByStepId={lastOutputByStepId}
                    tractUuid={tract.uuid}
                    linkedTriggerSummaries={linkedSummaries}
                    onChangeSteps={handleChangeSteps}
                    onOpenAddBlock={openAddBlock}
                    onClose={() => setSelectedNodeId(null)}
                />
            </div>
            <TractCanvasLogPanel
                open={logOpen}
                runs={runs}
                selectedRunUuid={selectedRunUuid}
                onSelectRun={setSelectedRunUuid}
                onClose={() => setLogOpen(false)}
            />
        </div>
    )
}
