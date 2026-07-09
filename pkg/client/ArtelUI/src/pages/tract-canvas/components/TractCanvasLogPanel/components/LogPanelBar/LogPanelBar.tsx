import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import {ChevronRightIcon} from "@/pages/tract-canvas/icons/ChevronRightIcon/ChevronRightIcon.tsx"
import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"
import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/LogPanelBar/LogPanelBar.module.css"

interface LogPanelBarProps {
    enlarged: boolean
    onToggleEnlarge: () => void
    onClose: () => void
}

export default function LogPanelBar({enlarged, onToggleEnlarge, onClose}: LogPanelBarProps) {
    return (
        <div className={cls.LogPanelBarContainer}>
            <Button
                variant="ghost"
                className={cls.EnlargeBtn}
                onClick={onToggleEnlarge}
                aria-label={enlarged ? "Shrink log panel" : "Enlarge log panel"}
                aria-pressed={enlarged}
            >
                <ChevronRightIcon className={cn(cls.EnlargeChevron, enlarged && cls.EnlargeChevronOpen)}/>
            </Button>
            <span className={cls.BarLabel}>Runs</span>
            <Button variant="iconDanger" onClick={onClose} aria-label="Close log">
                <CloseIcon/>
            </Button>
        </div>
    )
}
