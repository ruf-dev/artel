import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageOpenAIDialog/ManageOpenAIDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import LlmKeyDialogHead from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.tsx"
import LlmKeyConnectForm from "@/components/LlmKeyConnectForm/LlmKeyConnectForm.tsx"
import LlmKeyConnectedContent from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.tsx"

export default function ManageOpenAIDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, addOpenAIConnection, checkOpenAIConnection} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_OPENAI)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect GPT (OpenAI)"
                message="Remove this API key? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("openai").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <LlmKeyDialogHead
                provider={ExternalProvider.EXTERNAL_PROVIDER_OPENAI}
                title="GPT (OpenAI)"
                onClose={CloseDialog}
            />
            {connection
                ? <LlmKeyConnectedContent
                    fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}
                />
                : <LlmKeyConnectForm
                    providerName="GPT (OpenAI)"
                    bodyCopy={
                        "Connect your OpenAI API key to let Artel access it as a BYOK LLM "
                        + "provider. We'll verify the key against OpenAI before saving it."
                    }
                    apiKeyPlaceholder="sk-…"
                    baseUrlPlaceholder="https://api.openai.com/v1"
                    modelPlaceholder="gpt-4o"
                    addConnection={addOpenAIConnection}
                    checkConnection={checkOpenAIConnection}
                />
            }
        </div>
    )
}
