import {useState} from "react"

import {McpToolInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import ToolRow from "@/pages/toolbox/components/ToolsDialog/components/ToolRow/ToolRow.tsx"
import ToolDetail from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/ToolDetail.tsx"
import cls from "@/pages/toolbox/components/ToolsDialog/ToolsDialog.module.css"

export default function ToolsDialog({candidate}: { candidate: MomCandidate }) {
    const tools = candidate.tools ?? []
    const [activeTool, setActiveTool] = useState<McpToolInfo | null>(null)

    if (activeTool) {
        return <ToolDetail candidate={candidate} tool={activeTool} onBack={() => setActiveTool(null)}/>
    }

    return (
        <div className={cls.ToolsDialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>{candidate.name}</h2>
            <p className={cls.DialogDesc}>{candidate.description}</p>
            <div className={cls.ToolList}>
                {tools.map(t => (
                    <ToolRow key={t.name} tool={t} onClick={() => setActiveTool(t)}/>
                ))}
                {tools.length === 0 && <p className={cls.ToolEmpty}>No tools defined.</p>}
            </div>
        </div>
    )
}
