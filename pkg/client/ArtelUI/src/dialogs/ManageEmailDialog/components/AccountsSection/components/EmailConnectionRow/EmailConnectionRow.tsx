import {Button} from "@vervstack/chures"

import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import {mailProviderIcon} from "@/app/utils/mailProviderIcon"
import cls from "@/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.module.css"

interface EmailConnectionRowProps {
    connection: ExternalConnectionInfo
    onRemove: (conn: ExternalConnectionInfo) => void
    onEdit: (conn: ExternalConnectionInfo) => void
}

export default function EmailConnectionRow({connection, onRemove, onEdit}: EmailConnectionRowProps) {
    const email = connection.generic?.fields?.username ?? connection.generic?.fields?.email ?? "—"
    const icon = mailProviderIcon(email)
    return (
        <div className={cls.AccountRowContainer} onClick={() => onEdit(connection)} role="button" tabIndex={0}>
            <div className={cls.AccountRowLeft}>
                {icon
                    ? <img src={icon} className={cls.AccountEmailIcon} alt=""/>
                    : <span className={cls.AccountEmailAtBadge}>@</span>
                }
                <span className={cls.AccountEmail}>{email}</span>
            </div>
            <div className={cls.AccountRowActions}>
                <Button variant="iconDanger" onClick={e => { e.stopPropagation(); onRemove(connection) }} aria-label="Remove">×</Button>
            </div>
        </div>
    )
}
