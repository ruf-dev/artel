import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManagePostgresDialog/ManagePostgresDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import LlmKeyDialogHead from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.tsx"
import LlmKeyConnectedContent from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.tsx"
import PostgresConnectForm from "@/components/PostgresConnectForm/PostgresConnectForm.tsx"

export default function ManagePostgresDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, addPostgresConnection, checkPostgresConnection} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_POSTGRES)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect Postgres"
                message="Remove this Postgres connection? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("postgres").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <LlmKeyDialogHead
                provider={ExternalProvider.EXTERNAL_PROVIDER_POSTGRES}
                title="PostgreSQL"
                onClose={CloseDialog}
            />
            {connection
                ? <LlmKeyConnectedContent
                    fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}
                />
                : <PostgresConnectForm
                    addConnection={addPostgresConnection}
                    checkConnection={checkPostgresConnection}
                />
            }
        </div>
    )
}
