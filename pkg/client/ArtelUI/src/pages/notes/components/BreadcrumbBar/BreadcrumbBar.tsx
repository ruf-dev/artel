import {Button} from "@vervstack/chures"

import ModeBar from "@/pages/notes/components/ModeBar/ModeBar.tsx"
import SaveStatusIndicator, { SaveStatus } from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.tsx"

import cls from "./BreadcrumbBar.module.css"


type Mode = 'edit' | 'preview' | 'read'

interface BreadcrumbBarProps {
    path: string | null
    mode: Mode
    onModeChange: (mode: Mode) => void
    saveStatus: SaveStatus
    saveError?: string
    onRename?: () => void
}

function BreadcrumbPath({ path }: { path: string }) {
    const segments = path.split('/').filter(Boolean)
    return (
        <div className={cls.BreadcrumbPath}>
            {segments.map((seg, i) => (
                <span key={i} className={cls.SegmentWrapper}>
                    {i > 0 && <span className={cls.Separator}>›</span>}
                    <span className={i === segments.length - 1 ? cls.SegmentActive : cls.Segment}>
                        {seg}
                    </span>
                </span>
            ))}
        </div>
    )
}

function PencilIcon() {
    return (
        <svg viewBox="0 0 14 14" width={12} height={12} fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
            <path d="M9.5 2L12 4.5l-7 7H2.5V9l7-7z" />
        </svg>
    )
}

export default function BreadcrumbBar({ path, mode, onModeChange, saveStatus, saveError, onRename }: BreadcrumbBarProps) {
    return (
        <div className={cls.BreadcrumbBarContainer}>
            <div className={cls.LeftSlot}>
                {path && <BreadcrumbPath path={path} />}
                {path && onRename && (
                    <Button variant="ghost" className={cls.RenameBtn} onClick={onRename} title="Rename / Move">
                        <PencilIcon />
                    </Button>
                )}
            </div>
            <div className={cls.CenterSlot}>
                <SaveStatusIndicator status={saveStatus} errorMessage={saveError} />
            </div>
            <div className={cls.RightSlot}>
                <ModeBar active={mode} onModeChange={onModeChange} />
            </div>
        </div>
    )
}
