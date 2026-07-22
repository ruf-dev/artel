import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageAnthropicDialog/ManageAnthropicDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageAnthropicDialog/components/DialogHead/DialogHead.tsx"
import ConnectForm from "@/dialogs/ManageAnthropicDialog/components/ConnectForm/ConnectForm.tsx"
import ConnectedContent from "@/dialogs/ManageAnthropicDialog/components/ConnectedContent/ConnectedContent.tsx"

export default function ManageAnthropicDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect Claude (Anthropic)"
                message="Remove this API key? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("anthropic").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead onClose={CloseDialog}/>
            {connection
                ? <ConnectedContent fields={connection.generic?.fields ?? {}} onDisconnect={handleDisconnect}/>
                : <ConnectForm/>}
        </div>
    )
}
