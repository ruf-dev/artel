import {toGenericTable}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/taskTrackerAdapters.ts"
import DisplayTaskTrackerTables
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/DisplayTaskTrackerTables.tsx"
import {ResultViewWidgetProps}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/resultViewWidgetTypes.ts"

export default function TrelloTableWidget({value, toolName}: ResultViewWidgetProps) {
    return <DisplayTaskTrackerTables toolName={toolName} table={toGenericTable(value)}/>
}
