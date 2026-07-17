import {McpToolInfo} from "@/app/api/artel/mcp_keys.pb.ts"
import SmtpIcon from "@/pages/toolbox/components/ToolsDialog/components/ToolRow/icons/SmtpIcon.tsx"
import ImapIcon from "@/pages/toolbox/components/ToolsDialog/components/ToolRow/icons/ImapIcon.tsx"
import GenericToolIcon from "@/pages/toolbox/components/ToolsDialog/components/ToolRow/icons/GenericToolIcon.tsx"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolRow/ToolRow.module.css"

export default function ToolRow({tool, onClick}: { tool: McpToolInfo; onClick: () => void }) {
    return (
        <div className={cls.ToolRowContainer} onClick={onClick} role="button" tabIndex={0}>
            <div className={cls.ToolHeader}>
                <div className={cls.ToolIcon}>
                    {tool.smtp ? <SmtpIcon/> : tool.imap ? <ImapIcon/> : <GenericToolIcon/>}
                </div>
                <span className={cls.ToolName}>{tool.name}</span>
            </div>
            <p className={cls.ToolDesc}>{tool.description}</p>
        </div>
    )
}
