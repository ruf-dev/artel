import {useEffect, useState} from "react"
import {useNavigate, useParams} from "react-router-dom"

import cls from "@/pages/tracts/TractEditorPage.module.css"

import {Path} from "@/app/routing/Router.tsx"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import useUser from "@/hooks/user/User.ts"

import {Button} from "@vervstack/chures"
import TractRunTimeline from "@/components/TractRunTimeline/TractRunTimeline.tsx"
import RunTractDialog from "@/components/RunTractDialog/RunTractDialog.tsx"
import TriggerPanel from "@/components/TriggerPanel/TriggerPanel.tsx"
import TractStepTree from "@/components/TractStepTree/TractStepTree.tsx"
import {ROOT_LOCATION} from "@/processes/tractSteps.ts"
import {SchemaNode, Tract, TractDefinition, TractRun, TractTool} from "@/processes/Tracts.ts"

export default function TractEditorPage() {
    const {tractUuid} = useParams<{ tractUuid: string }>()
    const navigate = useNavigate()
    const {auth} = useUser()
    const {tracts, fetch, fetchRuns, runsByTract, fetchTools, tools, triggers, fetchTriggers} = useTracts()
    const {fetch: fetchConnections} = useExternalConnections()

    useEffect(() => {
        if (!auth.isAuthenticated()) {
            navigate(Path.InitPage)
        }
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

    const linkedTrigger = triggers.find(t => (tract?.triggers ?? []).some(l => l.uuid === t.uuid) && Object.keys(t.payloadSchema.properties).length > 0)

    return (
        <div className={cls.Root}>
            {tract ? (
                <Builder
                    tract={tract}
                    tools={tools}
                    triggerSchema={linkedTrigger?.payloadSchema}
                    onBack={() => navigate(Path.TractsPage)}
                />
            ) : (
                <EditorHeader tract={undefined} onBack={() => navigate(Path.TractsPage)}/>
            )}
            <RunsPanel tractUuid={tractUuid ?? ""} runs={runs} tract={tract}/>
        </div>
    )
}

function Builder({tract, tools, triggerSchema, onBack}: {
    tract: Tract
    tools: TractTool[]
    triggerSchema?: SchemaNode
    onBack: () => void
}) {
    const {updateTract} = useTracts()
    const bakeError = useBakeError()

    const [seeded, setSeeded] = useState(false)
    const [name, setName] = useState(tract.name)
    const [definition, setDefinition] = useState<TractDefinition>(tract.definition)
    const [saving, setSaving] = useState(false)
    const [warnings, setWarnings] = useState<string[]>([])

    useEffect(() => {
        if (!seeded) {
            setName(tract.name)
            setDefinition(tract.definition)
            setSeeded(true)
        }
    }, [seeded, tract])

    const isDirty = seeded && (name !== tract.name || JSON.stringify(definition) !== JSON.stringify(tract.definition))

    function handleSave() {
        setSaving(true)
        updateTract(tract.uuid, name, tract.description, definition)
            .then(({warnings: w}) => setWarnings(w))
            .catch(err => bakeError("Failed to save tract", err))
            .finally(() => setSaving(false))
    }

    return (
        <>
            <div className={cls.Header}>
                <Button variant="ghost" onClick={onBack}>← Back</Button>
                <span className={cls.Divider}/>
                <input
                    className={cls.NameInput}
                    value={name}
                    onChange={e => setName(e.target.value)}
                />
                {isDirty && <span className={cls.DirtyDot} title="Unsaved changes"/>}
                <Button className={cls.SaveBtn} variant="primary" disabled={!isDirty || saving} onClick={handleSave}>
                    {saving ? "Saving…" : "Save"}
                </Button>
            </div>
            {warnings.length > 0 && (
                <div className={cls.Warnings}>
                    {warnings.map((w, i) => <div key={i} className={cls.WarningRow}>{w}</div>)}
                </div>
            )}
            <div className={cls.BuilderContainer}>
                <TriggerPanel tractUuid={tract.uuid} linkedTriggerSummaries={tract.triggers ?? []}/>
                <div className={cls.StepsSection}>
                    <h2 className={cls.StepsTitle}>Steps</h2>
                    <TractStepTree
                        rootSteps={definition.steps}
                        location={ROOT_LOCATION}
                        list={definition.steps}
                        tools={tools}
                        triggerSchema={triggerSchema}
                        onChange={newSteps => setDefinition({steps: newSteps})}
                    />
                </div>
            </div>
        </>
    )
}

function EditorHeader({tract, onBack}: { tract?: Tract; onBack: () => void }) {
    return (
        <div className={cls.Header}>
            <Button variant="ghost" onClick={onBack}>← Back</Button>
            <span className={cls.Divider}/>
            <span className={cls.Name}>{tract?.name ?? "…"}</span>
        </div>
    )
}

function RunsPanel({tractUuid, runs, tract}: { tractUuid: string; runs: TractRun[]; tract?: Tract }) {
    const {OpenDialog} = useDialog()
    const {runTract} = useTracts()
    const [selectedRunUuid, setSelectedRunUuid] = useState<string | null>(null)

    function startRun(params: unknown): Promise<TractRun | undefined> {
        return runTract(tractUuid, params).then(() => useTracts.getState().runsByTract[tractUuid]?.[0])
    }

    return (
        <div className={cls.RunsContainer}>
            <div className={cls.RunsToolbar}>
                <h2 className={cls.RunsTitle}>Runs</h2>
                <Button variant="primary" onClick={() => OpenDialog(<RunTractDialog tract={tract} onRun={startRun}/>)}>
                    Run
                </Button>
            </div>
            <div className={cls.RunsList}>
                {runs.map(r => (
                    <RunRow key={r.uuid} run={r} selected={r.uuid === selectedRunUuid} onClick={() => setSelectedRunUuid(r.uuid)}/>
                ))}
                {runs.length === 0 && <p className={cls.Empty}>No runs yet.</p>}
            </div>
            {selectedRunUuid && <TractRunTimeline runUuid={selectedRunUuid}/>}
        </div>
    )
}

function statusClass(status: string): string {
    if (status === "done") return cls.StatusDone
    if (status === "failed") return cls.StatusFailed
    return cls.StatusRunning
}

function RunRow({run, selected, onClick}: { run: TractRun; selected: boolean; onClick: () => void }) {
    return (
        <div className={`${cls.RunRow} ${selected ? cls.RunRowSelected : ""}`} onClick={onClick} role="button" tabIndex={0}>
            <span className={`${cls.StatusChip} ${statusClass(run.status)}`}>{run.status}</span>
            <span className={cls.RunMeta}>{run.startedBy}</span>
            <span className={cls.RunMeta}>{formatDate(run.createdAt)}</span>
            {run.error && <span className={cls.RunError}>{run.error}</span>}
        </div>
    )
}

function formatDate(iso: string | undefined): string {
    if (!iso) return ""
    const d = new Date(iso)
    if (isNaN(d.getTime())) return ""
    return d.toLocaleString(undefined, {year: "numeric", month: "short", day: "numeric", hour: "2-digit", minute: "2-digit"})
}
