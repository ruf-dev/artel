import {McpToolInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import JsonView from "@/components/JsonView/JsonView.tsx"
import {buildTaskTrackerTable}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/taskTrackerAdapters.ts"
import DisplayTaskTrackerTables
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/DisplayTaskTrackerTables/DisplayTaskTrackerTables.tsx"
import cls
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/ResultView.module.css"

function tryParseJson(result: string): { ok: true; value: unknown } | { ok: false } {
    try {
        return {ok: true, value: JSON.parse(result)}
    } catch {
        return {ok: false}
    }
}

export default function ResultView({result, candidate, tool}: {
    result: string; candidate: MomCandidate; tool: McpToolInfo
}) {
    const parsed = tryParseJson(result)
    const table = parsed.ok ? buildTaskTrackerTable(candidate.name ?? "", parsed.value) : null

    return (
        <div className={cls.ResultViewContainer}>
            <span className={cls.SectionTitle}>Result</span>
            {!parsed.ok && <pre className={cls.ResponsePre}>{result}</pre>}
            {parsed.ok && table && <DisplayTaskTrackerTables toolName={tool.name ?? ""} table={table}/>}
            {parsed.ok && !table && <JsonView value={parsed.value}/>}
        </div>
    )
}
