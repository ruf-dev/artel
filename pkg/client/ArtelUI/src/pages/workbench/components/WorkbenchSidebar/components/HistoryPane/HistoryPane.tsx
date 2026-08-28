import {WorkbenchHistory, WorkbenchHistoryRow} from "@/pages/workbench/processes/useWorkbenchHistory.ts"
import HistoryRow from
    "@/pages/workbench/components/WorkbenchSidebar/components/HistoryPane/components/HistoryRow/HistoryRow.tsx"
import cls from "@/pages/workbench/components/WorkbenchSidebar/components/HistoryPane/HistoryPane.module.css"

interface Props {
    history: WorkbenchHistory
}

const DAY_MS = 86_400_000

// The History tab body: one flat list of WorkbenchHistoryRow grouped into
// Today / Yesterday / Earlier day buckets (rows with no timestamp fall into
// Earlier). Selection/delete/loading all come from the WorkbenchHistory bundle.
export default function HistoryPane({history}: Props) {
    const {onDelete} = history
    const groups = groupRowsByDay(history.rows)

    return (
        <div className={cls.HistoryPaneContainer}>
            {history.loading && history.rows.length === 0 && (
                <p className={cls.State}>Loading history…</p>
            )}
            {!history.loading && history.rows.length === 0 && (
                <p className={cls.State}>No conversations yet.</p>
            )}
            {groups.map(group => (
                <div key={group.label} className={cls.Group}>
                    <div className={cls.GroupLabel}>{group.label}</div>
                    {group.rows.map(row => (
                        <HistoryRow
                            key={row.id}
                            row={row}
                            active={row.id === history.activeId}
                            onSelect={() => history.onSelect(row.id)}
                            onDelete={onDelete ? () => onDelete(row.id) : undefined}
                        />
                    ))}
                </div>
            ))}
        </div>
    )
}

interface RowGroup {
    label: string
    rows: WorkbenchHistoryRow[]
}

function groupRowsByDay(rows: WorkbenchHistoryRow[]): RowGroup[] {
    const now = new Date()
    const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
    const startOfYesterday = startOfToday - DAY_MS

    const today: WorkbenchHistoryRow[] = []
    const yesterday: WorkbenchHistoryRow[] = []
    const earlier: WorkbenchHistoryRow[] = []

    rows.forEach(row => {
        const t = row.timestamp ? new Date(row.timestamp).getTime() : NaN
        if (Number.isNaN(t) || t < startOfYesterday) earlier.push(row)
        else if (t >= startOfToday) today.push(row)
        else yesterday.push(row)
    })

    return [
        {label: "Today", rows: today},
        {label: "Yesterday", rows: yesterday},
        {label: "Earlier", rows: earlier},
    ].filter(group => group.rows.length > 0)
}
