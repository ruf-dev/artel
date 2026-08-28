import {MouseEvent} from "react"
import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import {WorkbenchHistoryRow} from "@/pages/workbench/processes/useWorkbenchHistory.ts"
import cls from
    "@/pages/workbench/components/WorkbenchSidebar/components/HistoryPane/components/HistoryRow/HistoryRow.module.css"

interface Props {
    row: WorkbenchHistoryRow
    active: boolean
    onSelect: () => void
    onDelete?: () => void
}

export default function HistoryRow({row, active, onSelect, onDelete}: Props) {
    function handleDelete(e: MouseEvent) {
        e.stopPropagation()
        onDelete?.()
    }

    return (
        <div className={cn(cls.HistoryRowContainer, active && cls.HistoryRowActive)}>
            <Button variant="unstyled" className={cls.SelectButton} onClick={onSelect}>
                <span className={cn(cls.Badge, row.source === "docker" ? cls.BadgeDocker : cls.BadgeApi)}/>
                <span className={cls.Texts}>
                    <span className={cls.Title}>{row.title}</span>
                    {row.subtitle && <span className={cls.Subtitle}>{row.subtitle}</span>}
                </span>
            </Button>
            {onDelete && (
                <Button
                    variant="unstyled"
                    className={cls.DeleteButton}
                    onClick={handleDelete}
                    aria-label="Delete chat"
                    title="Delete chat"
                >
                    <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor"
                         strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                </Button>
            )}
        </div>
    )
}
