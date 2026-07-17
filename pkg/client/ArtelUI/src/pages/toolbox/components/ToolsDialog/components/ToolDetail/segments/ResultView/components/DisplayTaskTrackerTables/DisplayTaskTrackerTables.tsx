import {useState} from "react"

import {TaskTrackerTable, transposeTable}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/taskTrackerAdapters.ts"
import TaskTrackerTableHead
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableHead/TaskTrackerTableHead.tsx"
import TaskTrackerTableBody
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableBody/TaskTrackerTableBody.tsx"
import TransposeToggle
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TransposeToggle/TransposeToggle.tsx"
import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/DisplayTaskTrackerTables.module.css"

export default function DisplayTaskTrackerTables({toolName, table}: { toolName: string; table: TaskTrackerTable }) {
    const [transposed, setTransposed] = useState(false)
    const displayTable = transposed ? transposeTable(table) : table

    return (
        <div className={cls.DisplayTaskTrackerTablesContainer}>
            {table.columns.length === 0 ? (
                <p className={cls.EmptyState}>{toolName} returned no results.</p>
            ) : (
                <>
                    <div className={cls.TableToolbar}>
                        <TransposeToggle transposed={transposed} onChange={setTransposed}/>
                    </div>
                    <table className={cls.Table}>
                        <TaskTrackerTableHead columns={displayTable.columns}/>
                        <TaskTrackerTableBody columns={displayTable.columns} rows={displayTable.rows}/>
                    </table>
                </>
            )}
        </div>
    )
}
