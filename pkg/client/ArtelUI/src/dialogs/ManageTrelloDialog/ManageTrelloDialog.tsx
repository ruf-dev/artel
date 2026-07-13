import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageTrelloDialog/ManageTrelloDialog.module.css"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageTrelloDialog/components/DialogHead/DialogHead.tsx"
import AccountsSection from "@/dialogs/ManageTrelloDialog/components/AccountsSection/AccountsSection.tsx"
import TrelloAddDialog from "@/dialogs/ManageTrelloDialog/components/TrelloAddDialog/TrelloAddDialog.tsx"

export default function ManageTrelloDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, loading, disconnectConnection} = useExternalConnections()
    const bakeError = useBakeError()

    const trelloConnections = connections.filter(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_TRELLO)

    function handleRemove(conn: ExternalConnectionInfo) {
        const fullName = conn.generic?.fields?.full_name ?? "this account"
        OpenDialog(
            <ConfirmDialog
                title="Disconnect Trello"
                message={`Remove "${fullName}"? You can reconnect at any time.`}
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnectConnection(conn.id ?? "").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead onClose={CloseDialog}/>
            <AccountsSection
                loading={loading}
                connections={trelloConnections}
                onAdd={() => OpenDialog(<TrelloAddDialog/>)}
                onRemove={handleRemove}
            />
        </div>
    )
}
