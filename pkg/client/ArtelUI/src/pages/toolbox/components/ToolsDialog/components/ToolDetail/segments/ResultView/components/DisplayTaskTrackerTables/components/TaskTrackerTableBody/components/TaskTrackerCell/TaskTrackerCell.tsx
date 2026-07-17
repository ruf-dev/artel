import {useState} from "react"
import {Button, ChevronDownIcon} from "@vervstack/chures"

import {stringifyCell}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/taskTrackerAdapters.ts"
import JsonView from "@/components/JsonView/JsonView.tsx"
import {cn} from "@/app/utils/cn.ts"
import cls
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/components/TaskTrackerTableBody/components/TaskTrackerCell/TaskTrackerCell.module.css"

function isJsonValue(value: unknown): boolean {
    return value !== null && typeof value === "object"
}

export default function TaskTrackerCell({value}: { value: unknown }) {
    const [expanded, setExpanded] = useState(false)

    if (!isJsonValue(value)) {
        return <span>{stringifyCell(value)}</span>
    }

    return (
        <div className={cls.TaskTrackerCellContainer}>
            <Button variant="unstyled" className={cls.ExpandBtn} onClick={() => setExpanded(prev => !prev)}>
                <ChevronDownIcon className={cn(cls.Chevron, expanded && cls.ChevronOpen)}/>
                <span className={cls.Preview}>{stringifyCell(value)}</span>
            </Button>
            {expanded && <JsonView value={value}/>}
        </div>
    )
}
