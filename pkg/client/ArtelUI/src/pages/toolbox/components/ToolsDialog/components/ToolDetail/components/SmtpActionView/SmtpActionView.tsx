import {SmtpOperation, SmtpToolAction} from "@/app/api/artel/mcp_keys.pb.ts"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/SmtpActionView/SmtpActionView.module.css"

export default function SmtpActionView({action}: { action: SmtpToolAction }) {
    return (
        <div className={cls.SmtpActionViewContainer}>
            <label className={cls.ActionLabel}>Operation</label>
            <select className={cls.ActionSelect} disabled value={action.operation ?? SmtpOperation.SMTP_OP_UNSPECIFIED}>
                <option value={SmtpOperation.SMTP_OP_SEND}>send</option>
            </select>
        </div>
    )
}
