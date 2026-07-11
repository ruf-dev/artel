import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/TractCanvasLogPanel.module.css"
import {TractRun} from "@/processes/Tracts.ts"
import {cn} from "@/app/utils/cn.ts"
import LogPanelBar from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/LogPanelBar/LogPanelBar.tsx"
import RunLog from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/RunLog/RunLog.tsx"

interface Props {
    open: boolean
    runs: TractRun[]
    selectedRunUuid: string | null
    onSelectRun: (uuid: string) => void
    onClose: () => void
}

export default function TractCanvasLogPanel({open, runs, selectedRunUuid, onSelectRun, onClose}: Props) {
    const [enlarged, setEnlarged] = useState(false)

    useEffect(() => {
        if (!open) setEnlarged(false)
    }, [open])

    return (
        <>
            {open && enlarged && (
                <div className={cls.Backdrop} onClick={() => setEnlarged(false)}/>
            )}
            <div className={cn(cls.Panel, open && cls.PanelOpen, open && enlarged && cls.Enlarged)}>
                <LogPanelBar enlarged={enlarged} onToggleEnlarge={() => setEnlarged(e => !e)} onClose={onClose}/>
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
                            ? <RunLog runUuid={selectedRunUuid}/>
                            : <p className={cls.Empty}>Select a run to see its log.</p>}
                    </div>
                </div>
            </div>
        </>
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

