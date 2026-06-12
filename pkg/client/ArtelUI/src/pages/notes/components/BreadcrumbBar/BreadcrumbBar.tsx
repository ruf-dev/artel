import cls from "./BreadcrumbBar.module.css"
import ModeBar from "@/pages/notes/components/ModeBar/ModeBar.tsx"
import SaveStatusIndicator, { SaveStatus } from "@/pages/notes/components/SaveStatusIndicator/SaveStatusIndicator.tsx"

type Mode = 'edit' | 'preview' | 'read'

interface BreadcrumbBarProps {
    path: string | null
    mode: Mode
    onModeChange: (mode: Mode) => void
    saveStatus: SaveStatus
    saveError?: string
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

export default function BreadcrumbBar({ path, mode, onModeChange, saveStatus, saveError }: BreadcrumbBarProps) {
    return (
        <div className={cls.BreadcrumbBarContainer}>
            <div className={cls.LeftSlot}>
                {path && <BreadcrumbPath path={path} />}
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
