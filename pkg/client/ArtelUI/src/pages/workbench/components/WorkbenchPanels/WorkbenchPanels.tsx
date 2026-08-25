import {AnimatePresence} from "framer-motion"

import Chat from "@/pages/workbench/components/Chat/Chat.tsx"
import SimpleChat from "@/pages/workbench/components/SimpleChat/SimpleChat.tsx"
import AnimatedTerminalView from "@/pages/workbench/components/TerminalView/AnimatedTerminalView.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import type {WorkbenchView} from "@/pages/workbench/components/WorkbenchToolbar/WorkbenchToolbar.tsx"
import type {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import type {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import type {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface ChatSession {
    items: ChatItem[]
    status: ChatConnectionStatus
    sendMessage: (text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    startNewChat: () => void
}

interface SimpleChatSessionBundle {
    items: ChatItem[]
    status: ChatConnectionStatus
    sendMessage: (text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    onNewChat: () => void
}

interface Props {
    effectiveMode: WorkbenchMode | "picking"
    exists: boolean
    status: string
    vaultId: string
    view: WorkbenchView
    awaitingAuth: boolean
    chatSession: ChatSession
    historyOpen: boolean
    onCloseHistory: () => void
    tabs: {id: string; name: string; active: boolean}[]
    pendingTerminalAuthLink?: string
    onSelectTab: (tabId: string) => void
    onCreateTab: () => void
    onCloseTab: (tabId: string) => void
    simpleChatId?: string
    onSelectSimpleChat: (chatId: string | undefined) => void
    simpleHistoryOpen: boolean
    onCloseSimpleHistory: () => void
    simpleChatSession: SimpleChatSessionBundle
}

// The three mutually-exclusive workbench panels (Docker chat, Docker terminal,
// Simple Chat), split out purely to keep WorkbenchPage.tsx's render function
// under the max-lines-per-function lint limit.
export default function WorkbenchPanels(props: Props) {
    const isDockerRunning = props.effectiveMode === "docker" && props.exists && props.status === "running"

    return (
        <AnimatePresence mode="wait">
            {isDockerRunning && !props.awaitingAuth && props.view === "chat" && (
                <Chat key="chat" items={props.chatSession.items} status={props.chatSession.status}
                      sendMessage={props.chatSession.sendMessage}
                      sendPermissionDecision={props.chatSession.sendPermissionDecision}
                      onNewChat={props.chatSession.startNewChat} vaultId={props.vaultId}
                      historyOpen={props.historyOpen} onCloseHistory={props.onCloseHistory}/>
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
            {props.effectiveMode === "simple-chat" && (
                <SimpleChat
                    key="simple-chat"
                    chatId={props.simpleChatId}
                    vaultId={props.vaultId}
                    onSelectChat={props.onSelectSimpleChat}
                    session={props.simpleChatSession}
                    historyOpen={props.simpleHistoryOpen || !props.simpleChatId}
                    onCloseHistory={props.onCloseSimpleHistory}
                />
            )}
        </AnimatePresence>
    )
}
