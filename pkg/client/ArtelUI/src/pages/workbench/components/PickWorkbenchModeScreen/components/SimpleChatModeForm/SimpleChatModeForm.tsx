import {useState} from "react"
import {Button, Toggle} from "@vervstack/chures"

import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useSimpleChatMutations} from "@/app/hooks/SimpleChat.ts"
import cls from
    "@/pages/workbench/components/PickWorkbenchModeScreen/components/SimpleChatModeForm/SimpleChatModeForm.module.css"

interface Props {
    vaultId: string
    onCreated: (chatId: string) => void
}

// Inline form shown once "Simple Chat" is picked in PickWorkbenchModeScreen. Model
// selection lives on the OpenRouter BYOK connection card (LlmKeyConnectedContent)
// as a persisted default now, not here — this form only needs a vault-access
// toggle before creating the chat with an empty model string; SimpleChat.tsx
// resolves a real model (chat.model || recommendedDefaultModel || models[0]) once
// the chat is opened, and the in-chat ModelSwitcher lets the user change it.
export default function SimpleChatModeForm({vaultId, onCreated}: Props) {
    const {create, creating} = useSimpleChatMutations(vaultId)
    const bakeError = useBakeError()

    const [vaultAccess, setVaultAccess] = useState(true)

    function handleStart() {
        create({model: "", vaultAccess})
            .then(chat => onCreated(chat.id))
            .catch(e => bakeError("Failed to create chat", e))
    }

    return (
        <div className={cls.SimpleChatModeFormContainer}>
            <label className={cls.ToggleField}>
                <span className={cls.FieldLabel}>Give this chat vault access</span>
                <Toggle checked={vaultAccess} onChange={setVaultAccess}/>
            </label>
            <div className={cls.ModalFooter}>
                <Button variant="primary" onClick={handleStart} disabled={creating}>
                    {creating ? "Starting…" : "Start Chat"}
                </Button>
            </div>
        </div>
    )
}
