import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/components/StepPickerDialog/StepPickerDialog.module.css"
import {TractTool} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useDialog} from "@/app/hooks/Dialog"
import ToolStep from "@/components/StepPickerDialog/components/ToolStep.tsx"
import ConnectionStep from "@/components/StepPickerDialog/components/ConnectionStep.tsx"

const BUILTIN_MCP = "artel"

export interface StepDraft {
    type: "action" | "condition" | "parallel"
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

    function handleSelectSimpleKind(type: "condition" | "parallel") {
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
                <Button variant="ghost" className={cls.KindOption} onClick={() => setStep("tool")}>
                    <span className={cls.KindTitle}>Action</span>
                    <span className={cls.KindDesc}>Executes an MoM tool (or a builtin) with rendered params.</span>
                </Button>
                <Button variant="ghost" className={cls.KindOption} onClick={() => handleSelectSimpleKind("condition")}>
                    <span className={cls.KindTitle}>Condition</span>
                    <span className={cls.KindDesc}>Evaluates rules and routes to a true/false branch.</span>
                </Button>
                <Button variant="ghost" className={cls.KindOption} onClick={() => handleSelectSimpleKind("parallel")}>
                    <span className={cls.KindTitle}>Parallel</span>
                    <span className={cls.KindDesc}>Runs each lane concurrently. Add lanes after creating it.</span>
                </Button>
            </div>
            <div className={cls.DialogActions}>
                <Button variant="ghost" onClick={CloseDialog}>Cancel</Button>
            </div>
        </div>
    )
}