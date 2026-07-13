import {Button} from "@vervstack/chures"

import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import ProviderIcon from "@/components/ProviderIcon/ProviderIcon.tsx"
import cls
// eslint-disable-next-line max-len -- path too long to fit even unindented, can't shorten without changing import
from "@/dialogs/ManageTrelloDialog/components/AccountsSection/components/TrelloConnectionRow/TrelloConnectionRow.module.css"

interface TrelloConnectionRowProps {
    connection: ExternalConnectionInfo
    onRemove: (conn: ExternalConnectionInfo) => void
}

export default function TrelloConnectionRow({connection, onRemove}: TrelloConnectionRowProps) {
    const fullName = connection.generic?.fields?.full_name ?? "—"
    return (
        <div className={cls.AccountRowContainer}>
            <div className={cls.AccountRowLeft}>
                <div className={cls.AccountIcon}>
                    <ProviderIcon provider={ExternalProvider.EXTERNAL_PROVIDER_TRELLO}/>
                </div>
                <span className={cls.AccountName}>{fullName}</span>
            </div>
            <div className={cls.AccountRowActions}>
                <Button
                    variant="iconDanger"
                    onClick={() => onRemove(connection)}
                    aria-label="Remove"
                >
                    ×
                </Button>
            </div>
        </div>
    )
}
