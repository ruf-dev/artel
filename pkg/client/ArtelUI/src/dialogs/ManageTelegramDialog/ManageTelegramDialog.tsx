import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageTelegramDialog/ManageTelegramDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageTelegramDialog/components/DialogHead/DialogHead.tsx"
import ConnectForm from "@/dialogs/ManageTelegramDialog/components/ConnectForm/ConnectForm.tsx"
import ConnectedContent from "@/dialogs/ManageTelegramDialog/components/ConnectedContent/ConnectedContent.tsx"

export default function ManageTelegramDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_TELEGRAM)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect Telegram"
                message="Remove the connection to Telegram? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("telegram").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead onClose={CloseDialog}/>
            {connection
                ? <ConnectedContent botUsername={connection.generic?.fields?.bot_username ?? ""}
                                    onDisconnect={handleDisconnect}/>
                : <ConnectForm/>}
        </div>
    )
}
