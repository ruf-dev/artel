import {ImapOperation, ImapToolAction} from "@/app/api/artel/mcp_keys.pb.ts"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ImapActionView/ImapActionView.module.css"

export default function ImapActionView({action}: { action: ImapToolAction }) {
    return (
        <div className={cls.ImapActionViewContainer}>
            <label className={cls.ActionLabel}>Operation</label>
            <select className={cls.ActionSelect} disabled value={action.operation ?? ImapOperation.IMAP_OP_UNSPECIFIED}>
                <option value={ImapOperation.IMAP_OP_LIST_FOLDERS}>list folders</option>
                <option value={ImapOperation.IMAP_OP_LIST_MESSAGES}>list messages</option>
                <option value={ImapOperation.IMAP_OP_FETCH_MESSAGE}>fetch message</option>
            </select>
        </div>
    )
}
