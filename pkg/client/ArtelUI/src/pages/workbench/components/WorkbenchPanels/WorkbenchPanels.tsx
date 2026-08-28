import {AnimatePresence} from "framer-motion"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import SimpleChat from "@/pages/workbench/components/SimpleChat/SimpleChat.tsx"
import AnimatedTerminalView from "@/pages/workbench/components/TerminalView/AnimatedTerminalView.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {WorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"
import type {WorkbenchAttachments} from "@/pages/workbench/processes/useWorkbenchAttachments.ts"
import type {WorkbenchView} from "@/pages/workbench/processes/workbenchView.ts"
import type {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import type {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import type {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface ChatSession {
    items: ChatItem[]
    status: ChatConnectionStatus
    sendMessage: (text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    startNewChat: () => void
    pendingTurn: boolean
}

interface SimpleChatSessionBundle {
    items: ChatItem[]
    status: ChatConnectionStatus
    sendMessage: (text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    onNewChat: () => void
    models: string[]
    currentModel: string
    modelsLoading?: boolean
    onChangeModel: (model: string) => void
    pendingTurn: boolean
}

interface Props {
    effectiveMode: WorkbenchMode | "picking"
    exists: boolean
    status: string
    vaultId: string
    view: WorkbenchView
    awaitingAuth: boolean
    chatSession: ChatSession
    tabs: {id: string; name: string; active: boolean}[]
    pendingTerminalAuthLink?: string
    onSelectTab: (tabId: string) => void
    onCreateTab: () => void
    onCloseTab: (tabId: string) => void
    simpleChatId?: string
    simpleChatSession: SimpleChatSessionBundle
    ctx: WorkbenchContext
    attachments: WorkbenchAttachments
}

// The three mutually-exclusive workbench panels (Docker chat, Docker terminal,
// api chat), split out purely to keep WorkbenchPage.tsx's render function
// under the max-lines-per-function lint limit.
export default function WorkbenchPanels(props: Props) {
    const isDockerRunning = props.effectiveMode === "docker" && props.exists && props.status === "running"

    return (
        <AnimatePresence mode="wait">
            {isDockerRunning && !props.awaitingAuth && props.view === "chat" && (
                <Chat key="chat" items={props.chatSession.items} status={props.chatSession.status}
                      sendMessage={props.chatSession.sendMessage}
                      sendPermissionDecision={props.chatSession.sendPermissionDecision}
                      onNewChat={props.chatSession.startNewChat}
                      pendingTurn={props.chatSession.pendingTurn}
                      ctx={props.ctx}
                      attachedPaths={props.attachments.paths}
                      onRemoveAttachment={props.attachments.remove}
                      onClearAttachments={props.attachments.clear}/>
            )}
            {isDockerRunning && props.view === "terminal" && (
                <AnimatedTerminalView
                    key="terminal"
                    vaultId={props.vaultId}
                    tabs={props.tabs}
                    pendingTerminalAuthLink={props.pendingTerminalAuthLink}
                    onSelectTab={props.onSelectTab}
                    onCreateTab={props.onCreateTab}
                    onCloseTab={props.onCloseTab}
                />
            )}
            {props.effectiveMode === "api" && (
                <SimpleChat
                    key="api"
                    chatId={props.simpleChatId}
                    session={props.simpleChatSession}
                    ctx={props.ctx}
                    attachedPaths={props.attachments.paths}
                    onRemoveAttachment={props.attachments.remove}
                    onClearAttachments={props.attachments.clear}
                />
            )}
        </AnimatePresence>
    )
}
