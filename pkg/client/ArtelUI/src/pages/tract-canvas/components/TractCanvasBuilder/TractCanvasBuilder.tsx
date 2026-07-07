import {useEffect, useMemo, useState} from "react"

import cls from "@/pages/tract-canvas/components/TractCanvasBuilder/TractCanvasBuilder.module.css"

import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

import TractBlockPicker from "@/pages/tract-canvas/components/TractBlockPicker/TractBlockPicker.tsx"
import {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"

import TractCanvasTopBar from "@/pages/tract-canvas/components/TractCanvasTopBar/TractCanvasTopBar.tsx"
import TractCanvasArea from "@/pages/tract-canvas/components/TractCanvasArea/TractCanvasArea.tsx"
import TractCanvasInspector from "@/pages/tract-canvas/widgets/TractCanvasInspector/TractCanvasInspector.tsx"
import TractCanvasLogPanel from "@/pages/tract-canvas/components/TractCanvasLogPanel/TractCanvasLogPanel.tsx"
import {useTractRunTracking} from "@/pages/tract-canvas/components/TractCanvasBuilder/useTractRunTracking.ts"
import RunTractDialog from "@/components/RunTractDialog/RunTractDialog.tsx"

import {buildStepFromDraft, collectAllStepIds, insertBlockAfter, Location} from "@/processes/tractSteps.ts"
import {layoutTract} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"
import {Tract, TractDefinition, TractRun, TractTool, Trigger} from "@/processes/Tracts.ts"
import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"

interface Props {
    tract: Tract
    tools: TractTool[]
    triggers: Trigger[]
    runs: TractRun[]
    momCandidates: MomCandidate[]
    onBack: () => void
}

export default function TractCanvasBuilder({tract, tools, triggers, runs, momCandidates, onBack}: Props) {
    const {updateTract} = useTracts()
    const bakeError = useBakeError()
    const {OpenDialog} = useDialog()

    const [seeded, setSeeded] = useState(false)
    const [name, setName] = useState(tract.name)
    const [definition, setDefinition] = useState<TractDefinition>(tract.definition)
    const [savedName, setSavedName] = useState(tract.name)
    const [savedDefinition, setSavedDefinition] = useState<TractDefinition>(tract.definition)
    const [saving, setSaving] = useState(false)
    const [warnings, setWarnings] = useState<string[]>([])
    const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null)

    useEffect(() => {
        if (!seeded) {
            setName(tract.name)
            setDefinition(tract.definition)
            setSavedName(tract.name)
            setSavedDefinition(tract.definition)
            setSeeded(true)
        }
    }, [seeded, tract])

    // Compared against the locally-built snapshot from the last successful save, not the
    // `tract` prop — the prop is re-derived from the server response via definitionFromProto,
    // whose proto-default filling (e.g. an empty params map coming back as undefined instead
    // of {}) doesn't byte-for-byte match the locally-built object, which left isDirty stuck
    // true forever after a real save.
    const isDirty = seeded && (name !== savedName || JSON.stringify(definition) !== JSON.stringify(savedDefinition))

    const layout = useMemo(() => layoutTract(definition.steps), [definition])
    const selectedNode = layout.nodes.find(n => n.id === selectedNodeId) ?? null

    const linkedSummaries = tract.triggers ?? []
    const firstLinked = linkedSummaries[0]
    const triggerInfo = firstLinked ? {name: firstLinked.name, kind: firstLinked.kind, source: firstLinked.source} : undefined
    const triggerSchema = triggers.find(t => linkedSummaries.some(l => l.uuid === t.uuid) && Object.keys(t.payloadSchema.properties).length > 0)?.payloadSchema

    const {
        running, logOpen, toggleLog, closeLog, selectedRunUuid, setSelectedRunUuid,
        lastOutputByStepId, nodeStatus, runningEdgeIds, startRun,
    } = useTractRunTracking(tract.uuid, layout)

    function openRunDialog() {
        OpenDialog(<RunTractDialog tract={tract} onRun={startRun}/>)
    }

    function handleSave() {
        setSaving(true)
        updateTract(tract.uuid, name, tract.description, definition)
            .then(({warnings: w}) => {
                setWarnings(w)
                setSavedName(name)
                setSavedDefinition(definition)
            })
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
                    existingIds.add(newStep.id)
                    setDefinition(d => ({steps: insertBlockAfter(d.steps, location, index, newStep, existingIds)}))
                }}
            />
        )
    }

    return (
        <div className={cls.Root}>
            <TractCanvasTopBar
                name={name}
                onNameChange={setName}
                isDirty={isDirty}
                saving={saving}
                running={running}
                lastRunStatus={runs[0]?.status}
                onBack={onBack}
                onSave={handleSave}
                onRun={openRunDialog}
                onToggleLog={toggleLog}
            />
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
                    momCandidates={momCandidates}
                    selectedNodeId={selectedNodeId}
                    onSelectNode={id => setSelectedNodeId(id)}
                    onBackgroundClick={() => setSelectedNodeId(null)}
                    nodeStatus={nodeStatus}
                    runningEdgeIds={runningEdgeIds}
                    onAddBlock={openAddBlock}
                />
                <TractCanvasInspector
                    node={selectedNode}
                    rootSteps={definition.steps}
                    tools={tools}
                    triggerSchema={triggerSchema}
                    momCandidates={momCandidates}
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
                onClose={closeLog}
            />
        </div>
    )
}
