import {Button} from "@vervstack/chures"

import {CloseIcon} from "@/pages/tract-canvas/icons/CloseIcon/CloseIcon.tsx"
import cls from "@/pages/tract-canvas/components/TractCanvasLogPanel/components/LogPanelBar/LogPanelBar.module.css"

interface LogPanelBarProps {
    onClose: () => void
}

export default function LogPanelBar({onClose}: LogPanelBarProps) {
    return (
        <div className={cls.LogPanelBarContainer}>
            <span className={cls.BarLabel}>Runs</span>
            <Button variant="iconDanger" onClick={onClose} aria-label="Close log">
                <CloseIcon/>
            </Button>
        </div>
    )
}
