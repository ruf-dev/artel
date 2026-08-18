import {Button} from "@vervstack/chures"

import TerminalTabChip
    from "@/pages/workbench/components/Terminal/components/TerminalTabBar/components/TerminalTabChip/TerminalTabChip"
import cls
    from "@/pages/workbench/components/Terminal/components/TerminalTabBar/TerminalTabBar.module.css"

interface Props {
    tabs: {id: string; name: string; active: boolean}[]
    onSelect: (id: string) => void
    onCreate: () => void
    onClose: (id: string) => void
}

export default function TerminalTabBar({tabs, onSelect, onCreate, onClose}: Props) {
    const canClose = tabs.length > 1

    return (
        <div className={cls.TerminalTabBarContainer}>
            <div className={cls.TabsWrapper}>
                {tabs.map(tab => (
                    <TerminalTabChip
                        key={tab.id}
                        tab={tab}
                        canClose={canClose}
                        onSelect={() => onSelect(tab.id)}
                        onClose={() => onClose(tab.id)}
                    />
                ))}
            </div>
            <Button
                variant="secondary"
                className={cls.NewTabButton}
                onClick={onCreate}
                aria-label="New terminal tab"
            >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <line x1="12" y1="5" x2="12" y2="19" stroke="currentColor" strokeWidth="2"/>
                    <line x1="5" y1="12" x2="19" y2="12" stroke="currentColor" strokeWidth="2"/>
                </svg>
            </Button>
        </div>
    )
}
