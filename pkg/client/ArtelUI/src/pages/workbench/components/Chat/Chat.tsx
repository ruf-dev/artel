import ChatPanelShell from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"
import type {WorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"

interface Props {
    items: ChatItem[]
    status: ChatConnectionStatus
    sendMessage: (text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    onNewChat: () => void
    ctx: WorkbenchContext
    attachedPaths: string[]
    onRemoveAttachment: (path: string) => void
    onClearAttachments: () => void
}

// Docker workbench chat — a thin shaper over ChatPanelShell. History is now the
// page-level WorkbenchSidebar, so this no longer owns a history sidebar slot.
//
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
            assistantLabel="Claude Code"
            ctx={props.ctx}
            attachedPaths={props.attachedPaths}
            onRemoveAttachment={props.onRemoveAttachment}
            onClearAttachments={props.onClearAttachments}
        />
    )
}
