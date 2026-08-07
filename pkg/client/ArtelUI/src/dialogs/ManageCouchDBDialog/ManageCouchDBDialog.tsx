import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageCouchDBDialog/ManageCouchDBDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import LlmKeyDialogHead from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.tsx"
import LlmKeyConnectedContent from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.tsx"
import CouchDBConnectForm from "@/components/CouchDBConnectForm/CouchDBConnectForm.tsx"

export default function ManageCouchDBDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, addCouchDBConnection, checkCouchDBConnection} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_COUCHDB)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect CouchDB"
                message="Remove this CouchDB connection? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("couchdb").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <LlmKeyDialogHead
                provider={ExternalProvider.EXTERNAL_PROVIDER_COUCHDB}
                title="CouchDB"
                onClose={CloseDialog}
            />
            {connection
                ? <LlmKeyConnectedContent
                    fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}
                />
                : <CouchDBConnectForm
                    addConnection={addCouchDBConnection}
                    checkConnection={checkCouchDBConnection}
                />
            }
        </div>
    )
}
