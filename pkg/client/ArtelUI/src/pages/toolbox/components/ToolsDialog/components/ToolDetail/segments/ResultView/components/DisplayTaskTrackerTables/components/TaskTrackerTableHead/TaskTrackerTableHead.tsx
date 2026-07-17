import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableHead/TaskTrackerTableHead.module.css"

export default function TaskTrackerTableHead({columns}: { columns: string[] }) {
    return (
        <thead className={cls.TaskTrackerTableHeadContainer}>
            <tr>
                {columns.map(c => <th key={c}>{c}</th>)}
            </tr>
        </thead>
    )
}
