import TaskTrackerCell
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableBody/components/TaskTrackerCell/TaskTrackerCell.tsx"
import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableBody/TaskTrackerTableBody.module.css"

export default function TaskTrackerTableBody(
    {columns, rows}: { columns: string[]; rows: Record<string, unknown>[] }
) {
    return (
        <tbody className={cls.TaskTrackerTableBodyContainer}>
            {rows.map((row, i) => (
                <tr key={i}>
                    {columns.map(c => <td key={c}><TaskTrackerCell value={row[c]}/></td>)}
                </tr>
            ))}
        </tbody>
    )
}
