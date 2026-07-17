import {useState} from "react"

import {McpToolInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {cn} from "@/app/utils/cn.ts"
import JsonView from "@/components/JsonView/JsonView.tsx"
import {getResultViewWidget}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/resultViewWidgets.ts"
import ViewModeToggle, {ResultViewMode}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/ViewModeToggle/ViewModeToggle.tsx"
import cls
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/ResultView.module.css"

function tryParseJson(result: string): { ok: true; value: unknown } | { ok: false } {
    try {
        return {ok: true, value: JSON.parse(result)}
    } catch {
        return {ok: false}
    }
}

export default function ResultView({result, candidate, tool, className}: {
    result: string; candidate: MomCandidate; tool: McpToolInfo; className?: string
}) {
    const parsed = tryParseJson(result)
    const widget = parsed.ok ? getResultViewWidget(candidate.name ?? "", tool.name ?? "") : null
    const [viewMode, setViewMode] = useState<ResultViewMode>("widget")

    return (
        <div className={cn(cls.ResultViewContainer, className)}>
            <div className={cls.ResultViewHeader}>
                <span className={cls.SectionTitle}>Result</span>
                {parsed.ok && widget && (
                    <ViewModeToggle mode={viewMode} widgetLabel={widget.label} onChange={setViewMode}/>
                )}
            </div>
            {!parsed.ok && <pre className={cls.ResponsePre}>{result}</pre>}
            {parsed.ok && widget && viewMode === "widget" && (
                <widget.Component value={parsed.value} toolName={tool.name ?? ""}/>
            )}
            {parsed.ok && (!widget || viewMode === "json") && <JsonView value={parsed.value}/>}
        </div>
    )
}
