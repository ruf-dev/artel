import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageS3Dialog/ManageS3Dialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import LlmKeyDialogHead from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.tsx"
import LlmKeyConnectedContent from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.tsx"
import S3ConnectForm from "@/components/S3ConnectForm/S3ConnectForm.tsx"

export default function ManageS3Dialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, addS3Connection, checkS3Connection} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_S3)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect S3 / MinIO"
                message="Remove this S3 connection? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("s3").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <LlmKeyDialogHead
                provider={ExternalProvider.EXTERNAL_PROVIDER_S3}
                title="S3 / MinIO"
                onClose={CloseDialog}
            />
            {connection
                ? <LlmKeyConnectedContent
                    fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}
                />
                : <S3ConnectForm
                    addConnection={addS3Connection}
                    checkConnection={checkS3Connection}
                />
            }
        </div>
    )
}
