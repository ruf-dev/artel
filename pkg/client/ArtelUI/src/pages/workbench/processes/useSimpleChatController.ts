import {useCallback, useMemo} from "react"

import {useSimpleChat, useSimpleChatMutations} from "@/app/hooks/SimpleChat.ts"
import {useLastUsedModel} from "@/app/hooks/useLastUsedModel.ts"
import {simpleChatMessagesToItems} from "@/pages/workbench/processes/simpleChatMessages.ts"
import {useOpenRouterModels} from "@/pages/workbench/processes/useOpenRouterModels.ts"
import {useSimpleChatSession} from "@/pages/workbench/processes/useSimpleChatSession.ts"

interface Params {
    chatId: string | undefined
    vaultId: string
    active: boolean
    onChatIdChange: (chatId: string | undefined) => void
    bakeError: (title: string, err: unknown) => void
}

// Lifts Simple Chat's session/model state out of SimpleChat.tsx up to
// WorkbenchPage.tsx, mirroring how docker's chatSession (useChatSession) already
// lives there. `active` gates every underlying hook exactly like docker's
// chatSession is only constructed while Docker mode is the effective mode, so
// switching to Simple Chat mode doesn't leave two sets of network calls/sockets
// running at once.
export function useSimpleChatController({chatId, vaultId, active, onChatIdChange, bakeError}: Params) {
    const effectiveChatId = active ? chatId : undefined

    const {chat, messages} = useSimpleChat(effectiveChatId)
    const {models, recommendedDefaultModel, isLoading: modelsLoading} = useOpenRouterModels(active)
    const {create, creating} = useSimpleChatMutations(vaultId)
    const {lastUsedModel, setLastUsedModel} = useLastUsedModel()

    const initialItems = useMemo(() => simpleChatMessagesToItems(messages), [messages])
    const initialModel = chat?.model || lastUsedModel || recommendedDefaultModel || models[0] || ""

    const session = useSimpleChatSession(effectiveChatId, initialModel, initialItems)
    const {setModel} = session

    // Remembers whatever model ends up selected (explicit switcher change, or the
    // one a new chat gets created with) so the next chat defaults to it — see
    // useLastUsedModel.ts.
    const handleSetModel = useCallback((model: string) => {
        setModel(model)
        setLastUsedModel(model)
    }, [setModel, setLastUsedModel])

    // Clears chatId immediately, before the CreateChat round-trip even starts: effectiveChatId
    // dropping to undefined closes the current session's websocket right away (see
    // useSimpleChatSession's status/composerDisabled), so there is no window where the old
    // chat's live socket — and the composer sending into it — is still reachable while the UI's
    // intent has already moved to "starting a new chat". Without this, a message sent during the
    // create() round-trip lands on the old thread, with its full history, instead of the new one.
    function handleNewChat() {
        const model = session.currentModel || lastUsedModel || recommendedDefaultModel || models[0] || ""
        if (model) {
            setLastUsedModel(model)
        }
        onChatIdChange(undefined)
        create({model, vaultAccess: true})
            .then(newChat => onChatIdChange(newChat.id))
            .catch(e => bakeError("Failed to create chat", e))
    }

    return {
        items: session.items,
        status: session.status,
        sendMessage: session.sendMessage,
        resendMessage: session.resendMessage,
        sendPermissionDecision: session.sendPermissionDecision,
        currentModel: session.currentModel,
        setModel: handleSetModel,
        models,
        modelsLoading,
        handleNewChat,
        creatingChat: creating,
        pendingTurn: session.pendingTurn,
    }
}

export type SimpleChatController = ReturnType<typeof useSimpleChatController>

// Reshapes the controller's return value into the bundle SimpleChat.tsx/
// WorkbenchPanels.tsx expect (session prop) — pulled out to a plain function so
// the mapping lines don't count against WorkbenchPage.tsx's function, which is
// already at the max-lines-per-function lint budget.
export function toSimpleChatSessionBundle(controller: SimpleChatController) {
    return {
        items: controller.items,
        status: controller.status,
        sendMessage: controller.sendMessage,
        resendMessage: controller.resendMessage,
        sendPermissionDecision: controller.sendPermissionDecision,
        onNewChat: controller.handleNewChat,
        models: controller.models,
        currentModel: controller.currentModel,
        modelsLoading: controller.modelsLoading,
        onChangeModel: controller.setModel,
        pendingTurn: controller.pendingTurn,
    }
}
