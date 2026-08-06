import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageAnthropicDialog/ManageAnthropicDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import LlmKeyDialogHead from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.tsx"
import LlmKeyConnectForm from "@/components/LlmKeyConnectForm/LlmKeyConnectForm.tsx"
import LlmKeyConnectedContent from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.tsx"

export default function ManageAnthropicDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, addAnthropicConnection, checkAnthropicConnection} = useExternalConnections()
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
            <LlmKeyDialogHead
                provider={ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC}
                title="Claude (Anthropic)"
                onClose={CloseDialog}
            />
            {connection
                ? <LlmKeyConnectedContent
                    fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}
                />
                : <LlmKeyConnectForm
                    providerName="Claude (Anthropic)"
                    bodyCopy={
                        "Connect your Anthropic API key to let Claude access it as a BYOK LLM "
                        + "provider. We'll verify the key against Anthropic before saving it."
                    }
                    apiKeyPlaceholder="sk-ant-…"
                    baseUrlPlaceholder="https://api.anthropic.com"
                    modelPlaceholder="claude-opus-4-8"
                    addConnection={addAnthropicConnection}
                    checkConnection={checkAnthropicConnection}
                />
            }
        </div>
    )
}
