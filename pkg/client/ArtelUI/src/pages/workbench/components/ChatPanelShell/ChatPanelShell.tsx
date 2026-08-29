import {useState} from "react"
import {motion} from "framer-motion"

import cls from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.module.css"
import ChatStatusBanner from "@/pages/workbench/components/Chat/components/ChatStatusBanner/ChatStatusBanner.tsx"
import ChatMessageList from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.tsx"
import ChatComposer from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"
import {deriveTurnState} from "@/pages/workbench/processes/turnState.ts"
import {withAttachmentsPreamble} from "@/pages/workbench/processes/workbenchAttachments.ts"
import type {WorkbenchContext} from "@/pages/workbench/processes/workbenchContext.ts"
import {usePendingElapsed} from "@/pages/workbench/processes/usePendingElapsed.ts"

interface Props {
    items: ChatItem[]
    vaultId: string
    bannerStatus: ChatConnectionStatus
    sendMessage: (text: string, attachments?: string[]) => void
    onResendMessage?: (id: string, text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    onNewChat: () => void
    pendingTurn: boolean
    composerDisabled: boolean
    composerPlaceholder: string
    hideNewChatButton?: boolean
    assistantLabel?: string
    ctx: WorkbenchContext
    attachedPaths: string[]
    onRemoveAttachment: (path: string) => void
    onClearAttachments: () => void
}

// Shared chat-panel body (status banner, message list, composer) reused by the
// Docker workbench's Chat and the api chat's SimpleChat — each mode owns its own
// session data/handlers. History now lives in the page-level WorkbenchSidebar,
// not inside this shell.
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function ChatPanelShell(props: Props) {
    const [draft, setDraft] = useState("")
    const turnState = deriveTurnState(props.items, props.pendingTurn)
    const pendingActive = turnState === "working"
    const lastKey = props.items[props.items.length - 1]?.key
    const pendingBucket = usePendingElapsed(pendingActive, lastKey)

    function handleSend() {
        const text = draft.trim()
        if (!text) return
        props.sendMessage(withAttachmentsPreamble(text, props.attachedPaths), props.attachedPaths)
        props.onClearAttachments()
        setDraft("")
    }

    return (
        <motion.div
            className={cls.ChatPanelShellContainer}
            initial={{opacity: 0, y: 20}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0, y: -12}}
            transition={{duration: 0.22, ease: "easeOut"}}
        >
            <div className={cls.ChatContent}>
                <ChatStatusBanner status={props.bannerStatus}/>
                <ChatMessageList
                    items={props.items}
                    vaultId={props.vaultId}
                    assistantLabel={props.assistantLabel}
                    onRetryMessage={props.sendMessage}
                    onResendMessage={props.onResendMessage}
                    retryDisabled={props.composerDisabled}
                    onPermissionDecision={props.sendPermissionDecision}
                    pending={pendingActive ? {bucket: pendingBucket, label: props.assistantLabel} : undefined}
                />
                <ChatComposer
                    value={draft}
                    onChange={setDraft}
                    onSend={handleSend}
                    onNewChat={props.onNewChat}
                    disabled={props.composerDisabled}
                    placeholder={props.composerPlaceholder}
                    hideNewChatButton={props.hideNewChatButton}
                    onOpenTweaks={props.ctx.openTweaks}
                    attachedPaths={props.attachedPaths}
                    onRemoveAttachment={props.onRemoveAttachment}
                />
            </div>
        </motion.div>
    )
}
