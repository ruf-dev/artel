import {useEffect, useMemo, useState} from "react"

import cls from "@/pages/tract-canvas/components/TractBlockPicker/TractBlockPicker.module.css"

import {TractTool} from "@/processes/Tracts.ts"
import {useTracts} from "@/app/hooks/Tracts.ts"
import {useMcpKeys} from "@/app/hooks/McpKeys.ts"
import {useDialog} from "@/app/hooks/Dialog"

import Button from "@/components/shared/Button/Button.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/ConnectorChip.tsx"
import {BranchIcon, ForkIcon, LayersIcon} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"
import {colorForKind, iconForTool} from "@/pages/tract-canvas/components/TractIcons/tractIconHelpers.ts"
import type {StepDraft} from "@/components/StepPickerDialog/StepPickerDialog.tsx"

const BUILTIN_MCP = "artel"

interface Props {
    onConfirm: (draft: StepDraft) => void
}

interface LogicOption {
    type: "condition" | "parallel" | "group"
    name: string
    desc: string
    Icon: (p: { className?: string }) => React.JSX.Element
}

const LOGIC_OPTIONS: LogicOption[] = [
    {type: "condition", name: "Condition", desc: "Route to a true/false branch", Icon: BranchIcon},
    {type: "parallel", name: "Parallel", desc: "Run lanes concurrently", Icon: ForkIcon},
    {type: "group", name: "Group", desc: "Nest a sequential chain", Icon: LayersIcon},
]

export default function TractBlockPicker({onConfirm}: Props) {
    const {CloseDialog} = useDialog()
    const {tools, fetchTools} = useTracts()
    const [query, setQuery] = useState("")
    const [selectedTool, setSelectedTool] = useState<TractTool | null>(null)
    const [selectedConnectionId, setSelectedConnectionId] = useState("")

    useEffect(() => {
        void fetchTools()
    }, [fetchTools])

    const q = query.trim().toLowerCase()

    const grouped = useMemo(() => {
        const acc: Record<string, TractTool[]> = {}
        for (const t of tools) {
            if (q && !`${t.mcp}.${t.tool}`.toLowerCase().includes(q) && !t.description?.toLowerCase().includes(q)) continue
            ;(acc[t.mcp] ??= []).push(t)
        }
        return acc
    }, [tools, q])

    const logicOptions = LOGIC_OPTIONS.filter(o => !q || o.name.toLowerCase().includes(q) || o.desc.toLowerCase().includes(q))

    function handleSelectTool(tool: TractTool) {
        if (tool.mcp === BUILTIN_MCP) {
            onConfirm({type: "action", mcp: tool.mcp, tool: tool.tool})
            CloseDialog()
            return
        }
        setSelectedTool(tool)
    }

    function handleConfirmConnection() {
        if (!selectedTool || !selectedConnectionId) return
        onConfirm({type: "action", mcp: selectedTool.mcp, tool: selectedTool.tool, connectionUuid: selectedConnectionId})
        CloseDialog()
    }

    if (selectedTool) {
        return (
            <ConnectionStep
                mcp={selectedTool.mcp}
                selectedConnectionId={selectedConnectionId}
                onSelect={setSelectedConnectionId}
                onBack={() => setSelectedTool(null)}
                onConfirm={handleConfirmConnection}
            />
        )
    }

    return (
        <div className={cls.Modal} role="dialog" aria-modal="true" aria-label="Add block">
            <div className={cls.Head}>
                <div className={cls.Title}>Add block</div>
                <Button variant="ghost" onClick={CloseDialog} aria-label="Close">✕</Button>
            </div>
            <input
                className={cls.Search}
                type="text"
                placeholder="Search actions…"
                value={query}
                onChange={e => setQuery(e.target.value)}
                autoFocus
            />
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
                {Object.entries(grouped).map(([mcp, mcpTools]) => (
                    <div className={cls.McpGroup} key={mcp}>
                        <div className={cls.CatTitle}>{mcp}</div>
                        <div className={cls.Grid}>
                            {mcpTools.map(t => (
                                <ToolCell key={t.tool} tool={t} onSelect={() => handleSelectTool(t)}/>
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

function LogicCell({option, onSelect}: { option: LogicOption; onSelect: () => void }) {
    const color = colorForKind(option.type)
    return (
        <button type="button" className={cls.Opt} onClick={onSelect}>
            <div className={cls.OptIcon} style={{background: color.bg, borderColor: color.border, color: color.fg}}>
                <option.Icon/>
            </div>
            <div className={cls.OptText}>
                <div className={cls.OptName}>{option.name}</div>
                <div className={cls.OptFn}>{option.desc}</div>
            </div>
        </button>
    )
}

function ToolCell({tool, onSelect}: { tool: TractTool; onSelect: () => void }) {
    const color = colorForKind("action", tool.mcp)
    const Icon = iconForTool(tool.mcp, tool.tool)
    return (
        <button type="button" className={cls.Opt} onClick={onSelect}>
            <div className={cls.OptIcon} style={{background: color.bg, borderColor: color.border, color: color.fg}}>
                <Icon/>
            </div>
            <div className={cls.OptText}>
                <div className={cls.OptName}>{tool.tool}</div>
                <div className={cls.OptFn}>{tool.mcp}.{tool.tool}</div>
            </div>
        </button>
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
        <div className={cls.Modal} role="dialog" aria-modal="true">
            <div className={cls.Head}>
                <div className={cls.Title}>Choose a connection</div>
                <Button variant="ghost" onClick={CloseDialog} aria-label="Close">✕</Button>
            </div>
            <div className={cls.Scroll}>
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
            </div>
            <div className={cls.Actions}>
                <Button variant="ghost" onClick={onBack}>Back</Button>
                <Button variant="primary" disabled={!selectedConnectionId} onClick={onConfirm}>Add</Button>
            </div>
        </div>
    )
}
