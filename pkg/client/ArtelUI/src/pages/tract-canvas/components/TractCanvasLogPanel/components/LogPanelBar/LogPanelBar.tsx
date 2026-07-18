import {Button} from "@vervstack/chures"

import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"
import {ExpandIcon} from "@/pages/tract-canvas/icons/ExpandIcon/ExpandIcon.tsx"
import {CollapseIcon} from "@/pages/tract-canvas/icons/CollapseIcon/CollapseIcon.tsx"
import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/LogPanelBar/LogPanelBar.module.css"

interface LogPanelBarProps {
    expanded: boolean
    onTogglePreset: () => void
    onClose: () => void
}

export default function LogPanelBar({expanded, onTogglePreset, onClose}: LogPanelBarProps) {
    return (
        <div className={cls.LogPanelBarContainer}>
            <span className={cls.BarLabel}>Runs</span>
            <Button
                variant="ghost"
                className={cls.SizeButton}
                onClick={onTogglePreset}
                aria-label={expanded ? "Shrink log panel" : "Expand log panel"}
                title={expanded ? "Shrink to 20% of screen" : "Expand to 85% of available height"}
            >
                {expanded ? <CollapseIcon/> : <ExpandIcon/>}
            </Button>
            <Button variant="iconDanger" onClick={onClose} aria-label="Close log">
                <CloseIcon/>
            </Button>
        </div>
    )
}
