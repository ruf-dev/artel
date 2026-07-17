import {TaskTrackerTable}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/taskTrackerAdapters.ts"
import TaskTrackerTableHead
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableHead/TaskTrackerTableHead.tsx"
import TaskTrackerTableBody
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableBody/TaskTrackerTableBody.tsx"
import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/DisplayTaskTrackerTables.module.css"

export default function DisplayTaskTrackerTables({toolName, table}: { toolName: string; table: TaskTrackerTable }) {
    return (
        <div className={cls.DisplayTaskTrackerTablesContainer}>
            {table.columns.length === 0 ? (
                <p className={cls.EmptyState}>{toolName} returned no results.</p>
            ) : (
                <table className={cls.Table}>
                    <TaskTrackerTableHead columns={table.columns}/>
                    <TaskTrackerTableBody columns={table.columns} rows={table.rows}/>
                </table>
            )}
        </div>
    )
}
