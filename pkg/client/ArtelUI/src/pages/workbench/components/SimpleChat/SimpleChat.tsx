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
}

interface Props {
    chatId: string | undefined
    vaultId: string
    onSelectChat: (chatId: string | undefined) => void
    session: SimpleChatSessionBundle
    historyOpen: boolean
    onCloseHistory: () => void
}

// Thin wrapper around ChatPanelShell — Simple Chat's own session/model state now
// lives in useSimpleChatController, called from WorkbenchPage.tsx (mirroring
// docker's useChatSession), so this component only shapes that bundle into
// ChatPanelShell's generic props and supplies its own history sidebar. The
// model-switcher/new-chat/history-toggle row that used to live in this component's
// own header now lives in the page-level SimpleChatTopBar instead — hence
// hideNewChatButton is always true here (the top bar owns the only affordance).
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function SimpleChat(props: Props) {
    return (
        <ChatPanelShell
            items={props.session.items}
            bannerStatus={props.chatId ? props.session.status : "closed"}
            sendMessage={props.session.sendMessage}
            sendPermissionDecision={props.session.sendPermissionDecision}
            onNewChat={props.session.onNewChat}
            composerDisabled={props.session.status !== "open" || !props.chatId}
            composerPlaceholder="Message the agent…"
            hideNewChatButton
            historyOpen={props.historyOpen}
            onCloseHistory={props.onCloseHistory}
            historySidebar={(
                <SimpleChatHistorySidebar
                    vaultId={props.vaultId}
                    activeChatId={props.chatId}
                    open={props.historyOpen}
                    onClose={props.onCloseHistory}
                    onSelectChat={props.onSelectChat}
                    onNewChat={props.session.onNewChat}
                />
            )}
        />
    )
}
