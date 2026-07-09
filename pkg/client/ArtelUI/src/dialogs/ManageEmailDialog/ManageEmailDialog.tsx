import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageEmailDialog/ManageEmailDialog.module.css"
import {ExternalConnectionInfo, ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageEmailDialog/components/DialogHead/DialogHead.tsx"
import AccountsSection from "@/dialogs/ManageEmailDialog/components/AccountsSection/AccountsSection.tsx"
import EmailAddDialog from "@/dialogs/ManageEmailDialog/components/EmailAddDialog/EmailAddDialog.tsx"
import EmailEditDialog from "@/dialogs/ManageEmailDialog/components/EmailEditDialog/EmailEditDialog.tsx"

export default function ManageEmailDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, loading, disconnect} = useExternalConnections()
    const bakeError = useBakeError()

    const emailConnections = connections.filter(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_EMAIL)

    function handleEdit(conn: ExternalConnectionInfo) {
        OpenDialog(<EmailEditDialog connection={conn}/>)
    }

    function handleRemove(conn: ExternalConnectionInfo) {
        const email = conn.generic?.fields?.username ?? conn.generic?.fields?.email ?? ""
        OpenDialog(
            <ConfirmDialog
                title="Remove email account"
                message={`Remove "${email}"? You can re-add it at any time.`}
                confirmLabel="Remove"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("email")
                        .catch(e => bakeError("Failed to remove account", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer}
             onClick={e => e.stopPropagation()}
             role="dialog"
             aria-modal="true">
            <DialogHead title="Email" onClose={CloseDialog}/>
            <AccountsSection
                loading={loading}
                connections={emailConnections}
                onAdd={() => OpenDialog(<EmailAddDialog/>)}
                onRemove={handleRemove}
                onEdit={handleEdit}
            />
        </div>
    )
}
