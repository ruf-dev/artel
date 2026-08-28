import cls from "@/pages/workbench/components/SimpleChat/SimpleChat.module.css"
import ChatPanelShell from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx"
import SimpleChatHistorySidebar from
    "@/pages/workbench/components/SimpleChatHistorySidebar/SimpleChatHistorySidebar.tsx"
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
    vaultId: string
    onSelectChat: (chatId: string | undefined) => void
    session: SimpleChatSessionBundle
    onClose: () => void
}

// Simple Chat's own session/model state lives in useSimpleChatController, called
// from WorkbenchPage.tsx (mirroring docker's useChatSession), so this component
// only shapes that bundle into ChatPanelShell's generic props. The history
// sidebar is a permanently-visible pinned left pane (not a toggleable overlay),
// rendered directly here as a flex sibling of ChatPanelShell — it owns the
// model-switcher/new-chat header row itself, since there's no page-level top bar
// anymore for Simple Chat mode. Hence hideNewChatButton is always true here.
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function SimpleChat(props: Props) {
    return (
        <div className={cls.SimpleChatContainer}>
            <SimpleChatHistorySidebar
                vaultId={props.vaultId}
                activeChatId={props.chatId}
                onSelectChat={props.onSelectChat}
                onNewChat={props.session.onNewChat}
                models={props.session.models}
                currentModel={props.session.currentModel}
                modelsLoading={props.session.modelsLoading}
                onChangeModel={props.session.onChangeModel}
                onClose={props.onClose}
            />
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
