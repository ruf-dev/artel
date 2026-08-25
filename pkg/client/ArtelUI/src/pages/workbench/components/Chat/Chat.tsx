import ChatHistorySidebar from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/ChatHistorySidebar.tsx"
import ChatPanelShell from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    items: ChatItem[]
    status: ChatConnectionStatus
    sendMessage: (text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    onNewChat: () => void
    vaultId?: string
    historyOpen: boolean
    onCloseHistory: () => void
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function Chat(props: Props) {
    return (
        <ChatPanelShell
            items={props.items}
            bannerStatus={props.status}
            sendMessage={props.sendMessage}
            sendPermissionDecision={props.sendPermissionDecision}
            onNewChat={props.onNewChat}
            composerDisabled={props.status !== "open"}
            composerPlaceholder="Message the workbench…"
            historyOpen={props.historyOpen}
            onCloseHistory={props.onCloseHistory}
            historySidebar={props.vaultId && (
                <ChatHistorySidebar vaultId={props.vaultId} open={props.historyOpen} onClose={props.onCloseHistory}/>
            )}
        />
    )
}
