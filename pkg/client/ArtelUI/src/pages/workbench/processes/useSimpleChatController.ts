import {useMemo} from "react"

import {useSimpleChat, useSimpleChatMutations} from "@/app/hooks/SimpleChat.ts"
import {simpleChatMessagesToItems} from "@/pages/workbench/processes/simpleChatMessages.ts"
import {useOpenRouterModels} from "@/pages/workbench/processes/useOpenRouterModels.ts"
import {useSimpleChatSession} from "@/pages/workbench/processes/useSimpleChatSession.ts"

interface Params {
    chatId: string | undefined
    vaultId: string
    active: boolean
    onChatCreated: (chatId: string) => void
    bakeError: (title: string, err: unknown) => void
}

// Lifts Simple Chat's session/model state out of SimpleChat.tsx up to
// WorkbenchPage.tsx, mirroring how docker's chatSession (useChatSession) already
// lives there — needed so the page-level SimpleChatTopBar (a sibling of the chat
// panel, not nested inside it) can render the model switcher/new-chat button.
// `active` gates every underlying hook exactly like docker's chatSession is only
// constructed while Docker mode is the effective mode, so switching to Simple
// Chat mode doesn't leave two sets of network calls/sockets running at once.
export function useSimpleChatController({chatId, vaultId, active, onChatCreated, bakeError}: Params) {
    const effectiveChatId = active ? chatId : undefined

    const {chat, messages} = useSimpleChat(effectiveChatId)
    const {models, recommendedDefaultModel, isLoading: modelsLoading} = useOpenRouterModels(active)
    const {create} = useSimpleChatMutations(vaultId)

    const initialItems = useMemo(() => simpleChatMessagesToItems(messages), [messages])
    const initialModel = chat?.model || recommendedDefaultModel || models[0] || ""

    const {items, status, sendMessage, sendPermissionDecision, currentModel, setModel} =
        useSimpleChatSession(effectiveChatId, initialModel, initialItems)

    function handleNewChat() {
        const model = currentModel || recommendedDefaultModel || models[0] || ""
        create({model, vaultAccess: true})
            .then(newChat => onChatCreated(newChat.id))
            .catch(e => bakeError("Failed to create chat", e))
    }

    return {
        items,
        status,
        sendMessage,
        sendPermissionDecision,
        currentModel,
        setModel,
        models,
        modelsLoading,
        handleNewChat,
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
        sendPermissionDecision: controller.sendPermissionDecision,
        onNewChat: controller.handleNewChat,
    }
}

// Same rationale as toSimpleChatSessionBundle, for WorkbenchHeader.tsx's
// SimpleChatTopBar-facing prop bundle. `onToggleHistory` is page-owned
// (simpleHistoryOpen toggling lives in useWorkbenchModeControls), not part of
// the controller itself, so it's threaded through as a parameter.
export function toSimpleChatTopBarProps(controller: SimpleChatController, onToggleHistory: () => void) {
    return {
        models: controller.models,
        currentModel: controller.currentModel,
        modelsLoading: controller.modelsLoading,
        onChangeModel: controller.setModel,
        onNewChat: controller.handleNewChat,
        onToggleHistory,
    }
}
