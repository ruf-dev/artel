import cls from "@/pages/workbench/components/SimpleChat/SimpleChat.module.css"
import ChatPanelShell from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"
import type {WorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"

interface SimpleChatSessionBundle {
    items: ChatItem[]
    status: ChatConnectionStatus
    sendMessage: (text: string) => void
    resendMessage: (id: string, text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    onNewChat: () => void
    models: string[]
    currentModel: string
    modelsLoading?: boolean
    onChangeModel: (model: string) => void
    pendingTurn: boolean
}

interface Props {
    chatId: string | undefined
    vaultId: string
    session: SimpleChatSessionBundle
    ctx: WorkbenchContext
    attachedPaths: string[]
    onRemoveAttachment: (path: string) => void
    onClearAttachments: () => void
}

// Display name for the assistant label: the last path segment of the model id
// (e.g. "anthropic/claude-sonnet-4" -> "claude-sonnet-4"), falling back to "Model".
function prettyModel(id: string): string {
    return (id.split("/").pop() || "") || "Model"
}

// The api chat body — a thin shaper over ChatPanelShell. Since Stage 2 the model
// switcher, new-chat button, history list and close-to-picker control all live in
// the page-level WorkbenchSidebar, so this component is just the message
// panel now. hideNewChatButton stays true: the sidebar owns "new chat".
export default function SimpleChat(props: Props) {
    const noModelSelected = !props.session.currentModel

    return (
        <div className={cls.SimpleChatContainer}>
            <ChatPanelShell
                items={props.session.items}
                vaultId={props.vaultId}
                bannerStatus={props.chatId ? props.session.status : "closed"}
                sendMessage={props.session.sendMessage}
                onResendMessage={props.session.resendMessage}
                sendPermissionDecision={props.session.sendPermissionDecision}
                onNewChat={props.session.onNewChat}
                pendingTurn={props.session.pendingTurn}
                composerDisabled={props.session.status !== "open" || !props.chatId || noModelSelected}
                composerPlaceholder={noModelSelected ? "Select a model to start chatting…" : "Message the agent…"}
                assistantLabel={prettyModel(props.session.currentModel)}
                hideNewChatButton
                ctx={props.ctx}
                attachedPaths={props.attachedPaths}
                onRemoveAttachment={props.onRemoveAttachment}
                onClearAttachments={props.onClearAttachments}
            />
        </div>
    )
}
