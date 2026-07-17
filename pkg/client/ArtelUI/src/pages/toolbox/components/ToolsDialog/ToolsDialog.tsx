import {useState} from "react"

import {McpToolInfo, MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import SearchInput from "@/components/SearchInput/SearchInput.tsx"
import ToolRow from "@/pages/toolbox/components/ToolsDialog/components/ToolRow/ToolRow.tsx"
import ToolDetail from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/ToolDetail.tsx"
import cls from "@/pages/toolbox/components/ToolsDialog/ToolsDialog.module.css"

export default function ToolsDialog({candidate}: { candidate: MomCandidate }) {
    const tools = candidate.tools ?? []
    const [activeTool, setActiveTool] = useState<McpToolInfo | null>(null)
    const [query, setQuery] = useState("")

    if (activeTool) {
        return <ToolDetail candidate={candidate} tool={activeTool} onBack={() => setActiveTool(null)}/>
    }

    const q = query.trim().toLowerCase()
    const filtered = q
        ? tools.filter(t => t.name?.toLowerCase().includes(q) || t.description?.toLowerCase().includes(q))
        : tools

    return (
        <div className={cls.ToolsDialogContainer} role="dialog" aria-modal="true">
            <h2 className={cls.DialogTitle}>{candidate.name}</h2>
            <p className={cls.DialogDesc}>{candidate.description}</p>
            {tools.length > 0 && (
                <SearchInput
                    className={cls.SearchWrapper}
                    placeholder="Search tools…"
                    value={query}
                    setValue={setQuery}
                />
            )}
            <div className={cls.ToolList}>
                {filtered.map(t => (
                    <ToolRow key={t.name} tool={t} onClick={() => setActiveTool(t)}/>
                ))}
                {tools.length === 0 && <p className={cls.ToolEmpty}>No tools defined.</p>}
                {tools.length > 0 && filtered.length === 0 && (
                    <p className={cls.ToolEmpty}>No tools match your search.</p>
                )}
            </div>
        </div>
    )
}
