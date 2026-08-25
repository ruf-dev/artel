import {ReactNode, useState} from "react"
import {motion} from "framer-motion"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/ChatPanelShell/ChatPanelShell.module.css"
import ChatStatusBanner from "@/pages/workbench/components/Chat/components/ChatStatusBanner/ChatStatusBanner.tsx"
import ChatMessageList from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.tsx"
import ChatComposer from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    items: ChatItem[]
    bannerStatus: ChatConnectionStatus
    sendMessage: (text: string) => void
    sendPermissionDecision: (id: string, decision: PermissionDecision) => void
    onNewChat: () => void
    composerDisabled: boolean
    composerPlaceholder: string
    hideNewChatButton?: boolean
    historyOpen: boolean
    onCloseHistory: () => void
    historySidebar?: ReactNode
}

// Shared chat-panel body (status banner, message list, composer, blur overlay,
// history-sidebar slot) reused by the Docker workbench's Chat and Simple Chat's
// SimpleChat — each mode owns its own session data/handlers and history sidebar.
//
// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function ChatPanelShell(props: Props) {
    const [draft, setDraft] = useState("")

    function handleSend() {
        const text = draft.trim()
        if (!text) return
        props.sendMessage(text)
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
            <div className={cn(cls.ChatContent, props.historyOpen && cls.ChatContentBlurred)}>
                <ChatStatusBanner status={props.bannerStatus}/>
                <ChatMessageList items={props.items} onPermissionDecision={props.sendPermissionDecision}/>
                <ChatComposer
                    value={draft}
                    onChange={setDraft}
                    onSend={handleSend}
                    onNewChat={props.onNewChat}
                    disabled={props.composerDisabled}
                    placeholder={props.composerPlaceholder}
                    hideNewChatButton={props.hideNewChatButton}
                />
                {props.historyOpen && <div className={cls.ChatContentOverlay} onClick={props.onCloseHistory}/>}
            </div>
            {props.historySidebar}
        </motion.div>
    )
}
