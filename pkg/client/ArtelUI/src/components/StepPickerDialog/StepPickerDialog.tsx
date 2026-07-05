import {useEffect, useState} from "react"

import cls from "@/components/StepPickerDialog/StepPickerDialog.module.css"

import {TractTool} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useDialog} from "@/app/hooks/Dialog"

import {Button} from "@vervstack/chures"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"

const BUILTIN_MCP = "artel"

export interface StepDraft {
    type: "action" | "condition" | "parallel" | "group"
    mcp?: string
    tool?: string
    connectionUuid?: string
}

type Step = "kind" | "tool" | "connection"

interface Props {
    onConfirm: (draft: StepDraft) => void
}

export default function StepPickerDialog({onConfirm}: Props) {
    const {CloseDialog} = useDialog()
    const {tools, fetchTools} = useTracts()
    const [step, setStep] = useState<Step>("kind")
    const [selectedTool, setSelectedTool] = useState<TractTool | null>(null)
    const [selectedConnectionId, setSelectedConnectionId] = useState("")

    useEffect(() => {
        void fetchTools()
    }, [fetchTools])

    function handleSelectTool(tool: TractTool) {
        setSelectedTool(tool)
        if (tool.mcp === BUILTIN_MCP) {
            onConfirm({type: "action", mcp: tool.mcp, tool: tool.tool})
            CloseDialog()
            return
        }
        setStep("connection")
    }

    function handleConfirmConnection() {
        if (!selectedTool || !selectedConnectionId) return
        onConfirm({type: "action", mcp: selectedTool.mcp, tool: selectedTool.tool, connectionUuid: selectedConnectionId})
        CloseDialog()
    }

    function handleSelectSimpleKind(type: "condition" | "parallel" | "group") {
        onConfirm({type})
        CloseDialog()
    }

    if (step === "connection" && selectedTool) {
        return (
            <ConnectionStep
                mcp={selectedTool.mcp}
                selectedConnectionId={selectedConnectionId}
                onSelect={setSelectedConnectionId}
                onBack={() => setStep("tool")}
                onConfirm={handleConfirmConnection}
            />
        )
    }

    if (step === "tool") {
        return (
            <ToolStep
                tools={tools}
                onSelect={handleSelectTool}
                onBack={() => setStep("kind")}
            />
        )
    }

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Add step</h2>
            <div className={cls.List}>
                <button type="button" className={cls.KindOption} onClick={() => setStep("tool")}>
                    <span className={cls.KindTitle}>Action</span>
                    <span className={cls.KindDesc}>Executes an MoM tool (or a builtin) with rendered params.</span>
                </button>
                <button type="button" className={cls.KindOption} onClick={() => handleSelectSimpleKind("condition")}>
                    <span className={cls.KindTitle}>Condition</span>
                    <span className={cls.KindDesc}>Evaluates rules and routes to a true/false branch.</span>
                </button>
                <button type="button" className={cls.KindOption} onClick={() => handleSelectSimpleKind("parallel")}>
                    <span className={cls.KindTitle}>Parallel</span>
                    <span className={cls.KindDesc}>Runs each lane concurrently. Add lanes after creating it.</span>
                </button>
                <button type="button" className={cls.KindOption} onClick={() => handleSelectSimpleKind("group")}>
                    <span className={cls.KindTitle}>Group</span>
                    <span className={cls.KindDesc}>A sequential chain nested as one visual unit.</span>
                </button>
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
            </div>
        </div>
    )
}

function ToolStep({tools, onSelect, onBack}: {
    tools: TractTool[]
    onSelect: (tool: TractTool) => void
    onBack: () => void
}) {
    const {CloseDialog} = useDialog()

    const grouped = tools.reduce<Record<string, TractTool[]>>((acc, t) => {
        (acc[t.mcp] ??= []).push(t)
        return acc
    }, {})

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Choose an action</h2>
            <div className={cls.List}>
                {Object.entries(grouped).map(([mcp, mcpTools]) => (
                    <div className={cls.McpGroup} key={mcp}>
                        <span className={cls.McpGroupLabel}>{mcp}</span>
                        {mcpTools.map(t => (
                            <button
                                key={t.tool}
                                type="button"
                                className={cls.ToolRow}
                                onClick={() => onSelect(t)}
                            >
                                <span className={cls.ToolName}>{t.tool}</span>
                                {t.description && <span className={cls.ToolDesc}>{t.description}</span>}
                            </button>
                        ))}
                    </div>
                ))}
                {tools.length === 0 && <p className={cls.Empty}>No tools available.</p>}
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
            </div>
        </div>
    )
}

function ConnectionStep({mcp, selectedConnectionId, onSelect, onBack, onConfirm}: {
    mcp: string
    selectedConnectionId: string
    onSelect: (id: string) => void
    onBack: () => void
    onConfirm: () => void
}) {
    const {CloseDialog} = useDialog()
    const {momCandidates, fetchMomCandidates} = useMcpKeys()

    useEffect(() => {
        void fetchMomCandidates()
    }, [fetchMomCandidates])

    const candidate = momCandidates.find(c => c.name === mcp)
    const connections = candidate?.connections ?? []

    return (
        <div className={cls.DialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>Choose a connection</h2>
            <div className={cls.List}>
                {connections.map(c => (
                    <SelectOption
                        key={c.id}
                        label={connectionLabel(c)}
                        selected={c.id === selectedConnectionId}
                        onSelect={() => onSelect(c.id ?? "")}
                    />
                ))}
                {connections.length === 0 && <p className={cls.Empty}>No connections available for {mcp}.</p>}
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <div className={cls.ActionsRight}>
                    <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
                    <Button variant="primary" disabled={!selectedConnectionId} onClick={onConfirm}>Add</Button>
                </div>
            </div>
        </div>
    )
}
