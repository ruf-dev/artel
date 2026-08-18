import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/Terminal/components/TerminalTabBar/components/TerminalTabChip/TerminalTabChip.module.css" // eslint-disable-line max-len

interface Props {
    tab: {id: string; name: string; active: boolean}
    canClose: boolean
    onSelect: () => void
    onClose: () => void
}

export default function TerminalTabChip({tab, canClose, onSelect, onClose}: Props) {
    function handleCloseClick(e: React.MouseEvent) {
        e.stopPropagation()
        onClose()
    }

    return (
        <div
            className={cn(cls.TerminalTabChipContainer, tab.active && cls.Active)}
            onClick={onSelect}
        >
            <span className={cls.TabName}>{tab.name}</span>
            <Button
                variant="secondary"
                className={cls.CloseButton}
                onClick={handleCloseClick}
                aria-label="Close tab"
                disabled={!canClose}
            >
                <svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor">
                    <line x1="6" y1="6" x2="18" y2="18" stroke="currentColor" strokeWidth="2"/>
                    <line x1="18" y1="6" x2="6" y2="18" stroke="currentColor" strokeWidth="2"/>
                </svg>
            </Button>
        </div>
    )
}
