import {useState} from "react"

import cls from "@/pages/workbench/components/Chat/Chat.module.css"
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
}

export default function Chat({items, status, sendMessage, sendPermissionDecision}: Props) {
    const [draft, setDraft] = useState("")

    function handleSend() {
        const text = draft.trim()
        if (!text) return
        sendMessage(text)
        setDraft("")
    }

    return (
        <div className={cls.ChatContainer}>
            <ChatStatusBanner status={status}/>
            <ChatMessageList items={items} onPermissionDecision={sendPermissionDecision}/>
            <ChatComposer
                value={draft}
                onChange={setDraft}
                onSend={handleSend}
                disabled={status !== "open"}
                placeholder="Message the workbench…"
            />
        </div>
    )
}
