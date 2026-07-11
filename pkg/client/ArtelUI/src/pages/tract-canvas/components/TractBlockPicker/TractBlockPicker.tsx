import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import SearchInput from "@/components/SearchInput/SearchInput.tsx"
import cls from "@/pages/tract-canvas/components/TractBlockPicker/TractBlockPicker.module.css"
import {TractTool} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useDialog} from "@/app/hooks/Dialog"
import type {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"
import {useTractBlockPickerData} from "@/pages/tract-canvas/components/TractBlockPicker/useTractBlockPickerData.ts"
import ConnectionFilterRow
    from "@/pages/tract-canvas/components/TractBlockPicker/components/ConnectionFilterRow/ConnectionFilterRow.tsx"
import LogicCell from "@/pages/tract-canvas/components/TractBlockPicker/components/LogicCell/LogicCell.tsx"
import ToolCell from "@/pages/tract-canvas/components/TractBlockPicker/components/ToolCell/ToolCell.tsx"
import ConnectionStep
    from "@/pages/tract-canvas/components/TractBlockPicker/components/ConnectionStep/ConnectionStep.tsx"

interface Props {
    onConfirm: (draft: StepDraft) => void
}

export default function TractBlockPicker({onConfirm}: Props) {
    const {CloseDialog} = useDialog()
    const {tools, fetchTools} = useTracts()
    const {momCandidates, fetchMomCandidates} = useMcpKeys()
    const [query, setQuery] = useState("")
    const [connFilter, setConnFilter] = useState<string | null>(null)
    const [selectedTool, setSelectedTool] = useState<TractTool | null>(null)
    const [selectedConnectionId, setSelectedConnectionId] = useState("")

    useEffect(() => {
        void fetchTools()
    }, [fetchTools])

    useEffect(() => {
        void fetchMomCandidates()
    }, [fetchMomCandidates])

    const {connectedCandidates, grouped, orderedGroups, logicOptions} = useTractBlockPickerData({
        tools, momCandidates, query, connFilter})

    function handleConfirmConnection() {
        if (!selectedTool || !selectedConnectionId) return
        onConfirm(
            {type: "action", mcp: selectedTool.mcp, tool: selectedTool.tool, connectionUuid: selectedConnectionId})
        CloseDialog()
    }

    if (selectedTool) {
        return (
            <ConnectionStep
                mcp={selectedTool.mcp}
                selectedConnectionId={selectedConnectionId}
                onSelect={setSelectedConnectionId}
                onBack={() => setSelectedTool(null)} onConfirm={handleConfirmConnection}
            />
        )
    }

    return (
        <div className={cls.Modal} role="dialog" aria-modal="true" aria-label="Add block">
            <div className={cls.Head}>
                <div className={cls.Title}>Add block</div>
                <Button variant="ghost" onClick={CloseDialog} aria-label="Close">✕</Button>
            </div>
            <SearchInput
                className={cls.SearchWrapper}
                placeholder="Search actions…"
                value={query}
                setValue={setQuery}
                autoFocus
            />
            {connectedCandidates.length > 0 && (
                <ConnectionFilterRow
                    candidates={connectedCandidates}
                    connFilter={connFilter}
                    onSelect={setConnFilter}
                />
            )}
            <div className={cls.Scroll}>
                {logicOptions.length > 0 && (
                    <>
                        <div className={cls.CatTitle}>Logic</div>
                        <div className={cls.Grid}>
                            {logicOptions.map(o => (
                                <LogicCell key={o.type} option={o} onSelect={() => {
                                    onConfirm({type: o.type})
                                    CloseDialog()
                                }}/>
                            ))}
                        </div>
                    </>
                )}
                {orderedGroups.map(([mcp, mcpTools]) => (
                    <div className={cls.McpGroup} key={mcp}>
                        <div className={cls.CatTitle}>{mcp}</div>
                        <div className={cls.Grid}>
                            {mcpTools.map(t => (
                                <ToolCell key={t.tool} tool={t} onSelect={() => setSelectedTool(t)}/>
                            ))}
                        </div>
                    </div>
                ))}
                {logicOptions.length === 0 && Object.keys(grouped).length === 0 && (
                    <p className={cls.Empty}>No matching blocks.</p>
                )}
            </div>
        </div>
    )
}
