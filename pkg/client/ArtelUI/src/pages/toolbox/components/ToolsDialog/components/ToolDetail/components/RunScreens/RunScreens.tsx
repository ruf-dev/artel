import {ReactNode} from "react"

import {McpToolInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {cn} from "@/app/utils/cn.ts"
import ResultView
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/ResultView.tsx"
import cls
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/RunScreens/RunScreens.module.css"

export default function RunScreens({showResult, result, candidate, tool, children}: {
    showResult: boolean
    result: string | null
    candidate: MomCandidate
    tool: McpToolInfo
    children: ReactNode
}) {
    return (
        <div className={cls.RunScreensContainer}>
            <div className={cn(cls.Track, showResult && cls.TrackShowResult)}>
                <div className={cls.FormPane}>{children}</div>
                {result !== null
                    ? <ResultView className={cls.ResultPane} result={result} candidate={candidate} tool={tool}/>
                    : <div className={cls.ResultPane}/>}
            </div>
        </div>
    )
}
