import {Button} from "@vervstack/chures"

import cls from "@/pages/tract-canvas/components/TractCanvasTopBar/TractCanvasTopBar.module.css"
import RunStatusBadge from "@/pages/tract-canvas/components/RunStatusBadge/RunStatusBadge.tsx"
import RunButton from "@/pages/tract-canvas/components/RunButton/RunButton.tsx"

interface Props {
    name: string
    onNameChange: (name: string) => void
    isDirty: boolean
    saving: boolean
    running: boolean
    lastRunStatus?: string
    onBack: () => void
    onSave: () => void
    onRun: () => void
    onToggleLog: () => void
}

export default function TractCanvasTopBar({
    name, onNameChange, isDirty, saving, running, lastRunStatus, onBack, onSave, onRun, onToggleLog,
}: Props) {
    return (
        <div className={cls.BarContainer}>
            <Button variant="ghost" onClick={onBack}>← Tracts</Button>
            <span className={cls.Divider}/>
            <input className={cls.NameInput} value={name} onChange={e => onNameChange(e.target.value)}/>
            {isDirty && <span className={cls.DirtyDot} title="Unsaved changes"/>}
            <div className={cls.BarRight}>
                <RunStatusBadge running={running} lastRunStatus={lastRunStatus}/>
                <Button variant="ghost" onClick={onToggleLog}>Logs</Button>
                {isDirty && (
                    <Button onClick={onSave} disabled={saving}>{saving ? "Saving…" : "Save"}</Button>
                )}
                <RunButton running={running} onClick={onRun}/>
            </div>
        </div>
    )
}
