import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/TractCanvasLogPanel.module.css"
import {TractRun} from "@/processes/Tracts.ts"
import {cn} from "@/app/utils/cn.ts"
import LogPanelBar from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/LogPanelBar/LogPanelBar.tsx"
import RunLog from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/RunLog.tsx"
import ResizeHandle from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/ResizeHandle/ResizeHandle.tsx"

const HEIGHT_STORAGE_KEY = "artel.tract.logPanelHeight"
const DEFAULT_HEIGHT = 220
const MIN_HEIGHT = 160
const COMPACT_HEIGHT_RATIO = 0.2
const EXPANDED_HEIGHT_RATIO = 0.85

function clampHeight(h: number): number {
    return Math.min(Math.max(h, MIN_HEIGHT), window.innerHeight * EXPANDED_HEIGHT_RATIO)
}

function loadStoredHeight(): number {
    const stored = Number(localStorage.getItem(HEIGHT_STORAGE_KEY))
    return stored > 0 ? clampHeight(stored) : DEFAULT_HEIGHT
}

interface Props {
    open: boolean
    runs: TractRun[]
    selectedRunUuid: string | null
    onSelectRun: (uuid: string) => void
    onRetryRun: (runUuid: string) => Promise<unknown>
    onClose: () => void
}

export default function TractCanvasLogPanel({open, runs, selectedRunUuid, onSelectRun, onRetryRun, onClose}: Props) {
    const [height, setHeight] = useState(loadStoredHeight)
    const [isDragging, setIsDragging] = useState(false)

    useEffect(() => {
        localStorage.setItem(HEIGHT_STORAGE_KEY, String(height))
    }, [height])

    function handleResize(newHeight: number) {
        setHeight(clampHeight(newHeight))
    }

    const expandedHeight = window.innerHeight * EXPANDED_HEIGHT_RATIO
    const compactHeight = window.innerHeight * COMPACT_HEIGHT_RATIO
    const isExpanded = height >= (compactHeight + expandedHeight) / 2

    function handleTogglePreset() {
        setHeight(isExpanded ? compactHeight : expandedHeight)
    }

    return (
        <div
            className={cn(cls.Panel, open && cls.PanelOpen, isDragging && cls.Dragging)}
            style={{height: open ? height : 0}}
        >
            <ResizeHandle
                onResize={handleResize}
                onDragStart={() => setIsDragging(true)}
                onDragEnd={() => setIsDragging(false)}
            />
            <LogPanelBar expanded={isExpanded} onTogglePreset={handleTogglePreset} onClose={onClose}/>
            <div className={cls.Split}>
                <div className={cls.RunList}>
                    {runs.length === 0 && <p className={cls.Empty}>No runs yet.</p>}
                    {runs.map(r => (
                        <Button
                            key={r.uuid}
                            variant="ghost"
                            className={cn(cls.RunRow, r.uuid === selectedRunUuid && cls.RunRowSelected)}
                            onClick={() => onSelectRun(r.uuid)}
                        >
                            <span className={cn(cls.Dot, dotClass(r.status))}/>
                            <span className={cls.RunMeta}>{r.startedBy}</span>
                            <span className={cls.RunMeta}>{formatDate(r.createdAt)}</span>
                        </Button>
                    ))}
                </div>
                <div className={cls.LogScroll}>
                    {selectedRunUuid
                        ? <RunLog
                            runUuid={selectedRunUuid}
                            runStatus={runs.find(r => r.uuid === selectedRunUuid)?.status ?? ""}
                            onRetry={() => onRetryRun(selectedRunUuid)}
                        />
                        : <p className={cls.Empty}>Select a run to see its log.</p>}
                </div>
            </div>
        </div>
    )
}

function dotClass(status: string): string {
    if (status === "done") return cls.DotOk
    if (status === "failed") return cls.DotErr
    return cls.DotRunning
}

function formatDate(iso: string | undefined): string {
    if (!iso) return ""
    const d = new Date(iso)
    if (isNaN(d.getTime())) return ""
    return d.toLocaleString(undefined, {month: "short", day: "numeric", hour: "2-digit", minute: "2-digit"})
}

