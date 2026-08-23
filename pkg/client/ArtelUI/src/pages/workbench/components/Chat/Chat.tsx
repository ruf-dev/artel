import {useState} from "react"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/pages/workbench/components/Chat/Chat.module.css"
import ChatHeader from "@/pages/workbench/components/Chat/components/ChatHeader/ChatHeader.tsx"
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
}

export default function Chat({items, status, sendMessage, sendPermissionDecision, onNewChat, vaultId}: Props) {
    const [draft, setDraft] = useState("")
    const [historyOpen, setHistoryOpen] = useState(false)

    function handleSend() {
        const text = draft.trim()
        if (!text) return
        sendMessage(text)
        setDraft("")
    }

    function closeHistory() {
        setHistoryOpen(false)
    }

    function toggleHistory() {
        setHistoryOpen(open => !open)
    }

    return (
        <div className={cls.ChatContainer}>
            <ChatHeader
                onNewChat={onNewChat}
                onToggleHistory={toggleHistory}
                disabled={status !== "open"}
                vaultId={vaultId}
            />
            <div className={cn(cls.ChatContent, historyOpen && cls.ChatContentBlurred)}>
                <ChatStatusBanner status={status}/>
                <ChatMessageList items={items} onPermissionDecision={sendPermissionDecision}/>
                <ChatComposer
                    value={draft}
                    onChange={setDraft}
                    onSend={handleSend}
                    disabled={status !== "open"}
                    placeholder="Message the workbench…"
                />
                {historyOpen && <div className={cls.ChatContentOverlay} onClick={closeHistory}/>}
            </div>
            {vaultId && <ChatHistorySidebar vaultId={vaultId} open={historyOpen} onClose={closeHistory}/>}
        </div>
    )
}
