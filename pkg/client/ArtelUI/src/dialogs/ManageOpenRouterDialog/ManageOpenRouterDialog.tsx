import {ConfirmDialog} from "@vervstack/chures"

import cls from "@/dialogs/ManageOpenRouterDialog/ManageOpenRouterDialog.module.css"
import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import LlmKeyDialogHead from "@/components/LlmKeyDialogHead/LlmKeyDialogHead.tsx"
import LlmKeyConnectForm from "@/components/LlmKeyConnectForm/LlmKeyConnectForm.tsx"
import LlmKeyConnectedContent from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.tsx"
import ViewOpenRouterUsageDialog from
    "@/dialogs/ManageOpenRouterDialog/components/ViewOpenRouterUsageDialog/ViewOpenRouterUsageDialog.tsx"

export default function ManageOpenRouterDialog() {
    const {CloseDialog, OpenDialog} = useDialog()
    const {connections, disconnect, addOpenAIConnection, checkOpenAIConnection} = useExternalConnections()
    const bakeError = useBakeError()

    const connection = connections.find(c => c.provider === ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER)

    function handleDisconnect() {
        OpenDialog(
            <ConfirmDialog
                title="Disconnect OpenRouter"
                message="Remove this API key? You can reconnect at any time."
                confirmLabel="Disconnect"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    disconnect("openrouter").catch(e => bakeError("Failed to disconnect", e))}
            />
        )
    }

    return (
        <div className={cls.ModalContainer} onClick={e => e.stopPropagation()} role="dialog" aria-modal="true">
            <LlmKeyDialogHead
                provider={ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER}
                title="OpenRouter"
                onClose={CloseDialog}
            />
            {connection
                ? <LlmKeyConnectedContent
                    fields={connection.generic?.fields ?? {}}
                    onDisconnect={handleDisconnect}
                    onViewUsage={() => OpenDialog(<ViewOpenRouterUsageDialog/>)}
                />
                : <LlmKeyConnectForm
                    providerName="OpenRouter"
                    bodyCopy={
                        "Connect your OpenRouter API key to let Artel access it as a BYOK LLM "
                        + "provider. We'll verify the key against OpenRouter before saving it."
                    }
                    apiKeyPlaceholder="sk-or-…"
                    baseUrlPlaceholder="https://openrouter.ai/api/v1"
                    modelPlaceholder="anthropic/claude-sonnet-4.5"
                    addConnection={(req) => addOpenAIConnection({
                        ...req,
                        provider: ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER,
                    })}
                    checkConnection={(req) => checkOpenAIConnection({
                        ...req,
                        provider: ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER,
                    })}
                />
            }
        </div>
    )
}
