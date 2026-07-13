import {Button} from "@vervstack/chures"

import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import TrelloConnectionRow
// eslint-disable-next-line max-len -- path too long to fit even unindented, can't shorten without changing import
    from "@/dialogs/ManageTrelloDialog/components/AccountsSection/components/TrelloConnectionRow/TrelloConnectionRow.tsx"
import cls from "@/dialogs/ManageTrelloDialog/components/AccountsSection/AccountsSection.module.css"

interface AccountsSectionProps {
    loading: boolean
    connections: ExternalConnectionInfo[]
    onAdd: () => void
    onRemove: (conn: ExternalConnectionInfo) => void
}

export default function AccountsSection({loading, connections, onAdd, onRemove}: AccountsSectionProps) {
    return (
        <div className={cls.AccountSectionContainer}>
            <div className={cls.AccountSectionHeader}>
                <span className={cls.AccountSectionLabel}>Accounts</span>
                <Button variant="ghost" onClick={onAdd}>+ Add</Button>
            </div>
            {loading ? (
                <p className={cls.Empty}>Loading…</p>
            ) : connections.length === 0 ? (
                <p className={cls.Empty}>No Trello accounts added yet.</p>
            ) : (
                <div className={cls.AccountList}>
                    {connections.map(conn => (
                        <TrelloConnectionRow key={conn.id} connection={conn} onRemove={onRemove}/>
                    ))}
                </div>
            )}
        </div>
    )
}
