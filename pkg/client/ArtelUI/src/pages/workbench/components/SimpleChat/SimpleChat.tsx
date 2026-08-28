import cls from "@/pages/workbench/components/SimpleChat/SimpleChat.module.css"
import ChatPanelShell from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

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
}

interface Props {
    chatId: string | undefined
    session: SimpleChatSessionBundle
}

// The api chat body — a thin shaper over ChatPanelShell. Since Stage 2 the model
// switcher, new-chat button, history list and close-to-picker control all live in
// the page-level WorkbenchSidebar, so this component is just the message
// panel now. hideNewChatButton stays true: the sidebar owns "new chat".
export default function SimpleChat(props: Props) {
    return (
        <div className={cls.SimpleChatContainer}>
            <ChatPanelShell
                items={props.session.items}
                bannerStatus={props.chatId ? props.session.status : "closed"}
                sendMessage={props.session.sendMessage}
                sendPermissionDecision={props.session.sendPermissionDecision}
                onNewChat={props.session.onNewChat}
                composerDisabled={props.session.status !== "open" || !props.chatId}
                composerPlaceholder="Message the agent…"
                hideNewChatButton
            />
        </div>
    )
}
