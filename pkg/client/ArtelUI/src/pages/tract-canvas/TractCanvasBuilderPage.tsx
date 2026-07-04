import {useEffect, useMemo, useState} from "react"
import {useNavigate, useParams} from "react-router-dom"

import cls from "@/pages/tract-canvas/TractCanvasBuilderPage.module.css"

import {Path} from "@/app/routing/Router.tsx"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"

import Button from "@/components/shared/Button/Button.tsx"
import TractBlockPicker from "@/components/TractBlockPicker/TractBlockPicker.tsx"
import {PlayIcon} from "@/components/TractIcons/TractIcons.tsx"
import {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"

import TractCanvasArea from "@/widgets/TractCanvas/TractCanvasArea.tsx"
import TractCanvasInspector from "@/widgets/TractCanvas/TractCanvasInspector.tsx"
import TractCanvasLogPanel from "@/widgets/TractCanvas/TractCanvasLogPanel.tsx"
import {NodeStatus} from "@/widgets/TractCanvas/TractCanvasNode.tsx"

import {buildStepFromDraft, collectAllStepIds, insertStepAt, Location, ROOT_LOCATION} from "@/processes/tractSteps.ts"
import {layoutTract, TRIGGER_NODE_ID} from "@/processes/tractCanvasLayout.ts"
import {Tract, TractDefinition, TractRun, TractTool, Trigger} from "@/processes/Tracts.ts"

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms))

export default function TractCanvasBuilderPage() {
    const {tractUuid} = useParams<{ tractUuid: string }>()
    const navigate = useNavigate()
    const {auth} = useUser()
    const {tracts, fetch, fetchRuns, runsByTract, fetchTools, tools, triggers, fetchTriggers} = useTracts()
    const {fetch: fetchConnections} = useExternalConnections()

    useEffect(() => {
        if (!auth.isAuthenticated()) navigate(Path.InitPage)
    }, [auth, navigate])

    useEffect(() => {
        if (!auth.isAuthenticated() || !tractUuid) return
        if (tracts.length === 0) void fetch()
        void fetchRuns(tractUuid)
        void fetchTools()
        void fetchTriggers()
        void fetchConnections()
    }, [auth, tractUuid])

    const tract = tracts.find(t => t.uuid === tractUuid)
    const runs = tractUuid ? runsByTract[tractUuid] ?? [] : []

    if (!tract) {
        return (
            <div className={cls.Root}>
                <div className={cls.Bar}>
                    <Button variant="ghost" onClick={() => navigate(Path.TractsPage)}>← Tracts</Button>
                    <span className={cls.Divider}/>
                    <span className={cls.Name}>…</span>
                </div>
            </div>
        )
    }

    return (
        <Builder
            tract={tract}
            tools={tools}
            triggers={triggers}
            runs={runs}
            onBack={() => navigate(Path.TractsPage)}
        />
    )
}

function Builder({tract, tools, triggers, runs, onBack}: {
    tract: Tract
    tools: TractTool[]
    triggers: Trigger[]
    runs: TractRun[]
    onBack: () => void
}) {
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

function RunStatusBadge({running, lastRunStatus}: { running: boolean; lastRunStatus?: string }) {
    return (
        <span className={cls.RunBadge}>
            <span className={`${cls.BadgeDot} ${running ? cls.BadgeDotRunning : ""}`}/>
            {running ? "Running…" : lastRunStatus ? `Last run: ${lastRunStatus}` : "No runs yet"}
        </span>
    )
}

function RunButton({running, onClick}: { running: boolean; onClick: () => void }) {
    return (
        <Button variant="primary" onClick={onClick} disabled={running}>
            <PlayIcon/> {running ? "Running…" : "Run"}
        </Button>
    )
}
