import {useState} from "react"
import {motion} from "framer-motion"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/Chat/Chat.module.css"
import ChatHistorySidebar from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/ChatHistorySidebar.tsx"
import ChatStatusBanner from "@/pages/workbench/components/Chat/components/ChatStatusBanner/ChatStatusBanner.tsx"
import ChatMessageList from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.tsx"
import ChatComposer from "@/pages/workbench/components/Chat/components/ChatComposer/ChatComposer.tsx"
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
    const [draft, setDraft] = useState("")

    function handleSend() {
        const text = draft.trim()
        if (!text) return
        props.sendMessage(text)
        setDraft("")
    }

    return (
        <motion.div
            className={cls.ChatContainer}
            initial={{opacity: 0, y: 20}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0, y: -12}}
            transition={{duration: 0.22, ease: "easeOut"}}
        >
            <div className={cn(cls.ChatContent, props.historyOpen && cls.ChatContentBlurred)}>
                <ChatStatusBanner status={props.status}/>
                <ChatMessageList items={props.items} onPermissionDecision={props.sendPermissionDecision}/>
                <ChatComposer
                    value={draft}
                    onChange={setDraft}
                    onSend={handleSend}
                    onNewChat={props.onNewChat}
                    disabled={props.status !== "open"}
                    placeholder="Message the workbench…"
                />
                {props.historyOpen && <div className={cls.ChatContentOverlay} onClick={props.onCloseHistory}/>}
            </div>
            {props.vaultId && (
                <ChatHistorySidebar vaultId={props.vaultId} open={props.historyOpen} onClose={props.onCloseHistory}/>
            )}
        </motion.div>
    )
}
