import {Button} from "@vervstack/chures"

import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import EmailConnectionRow from "@/dialogs/ManageEmailDialog/components/AccountsSection/components/EmailConnectionRow/EmailConnectionRow.tsx"
import cls from "@/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.module.css"

interface AccountsSectionProps {
    loading: boolean
    connections: ExternalConnectionInfo[]
    onAdd: () => void
    onRemove: (conn: ExternalConnectionInfo) => void
    onEdit: (conn: ExternalConnectionInfo) => void
}

export default function AccountsSection({loading, connections, onAdd, onRemove, onEdit}: AccountsSectionProps) {
    return (
        <div className={cls.AccountSectionContainer}>
            <div className={cls.AccountSectionHeader}>
                <span className={cls.AccountSectionLabel}>Accounts</span>
                <Button variant="ghost" onClick={onAdd}>+ Add</Button>
            </div>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : connections.length === 0 ? (
                <p className={cls.Empty}>No email accounts added yet.</p>
            ) : (
                <div className={cls.AccountList}>
                    {connections.map(conn => (
                        <EmailConnectionRow key={conn.id} connection={conn} onRemove={onRemove} onEdit={onEdit}/>
                    ))}
                </div>
            )}
        </div>
    )
}
