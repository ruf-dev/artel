import {McpToolInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolRow/ToolRow.module.css"

export default function ToolRow({tool, onClick}: { tool: McpToolInfo; onClick: () => void }) {
    return (
        <div className={cls.ToolRowContainer} onClick={onClick} role="button" tabIndex={0}>
            <span className={cls.ToolName}>{tool.name}</span>
            <p className={cls.ToolDesc}>{tool.description}</p>
        </div>
    )
}
