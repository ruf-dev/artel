import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageGitlabDialog/ManageGitlabDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import DialogHead from "@/dialogs/ManageGitlabDialog/components/DialogHead/DialogHead.tsx"
import ConnectForm from "@/dialogs/ManageGitlabDialog/components/ConnectForm/ConnectForm.tsx"
import ConnectedContent from "@/dialogs/ManageGitlabDialog/components/ConnectedContent/ConnectedContent.tsx"

export default function ManageGitlabDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_GITLAB)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect GitLab"
                message="Remove the connection to GitLab? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("gitlab").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <DialogHead onClose={CloseDialog}/>
            {connection
                ? <ConnectedContent connectionId={connection.id ?? ""} fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}/>
                : <ConnectForm/>}
        </div>
    )
}
