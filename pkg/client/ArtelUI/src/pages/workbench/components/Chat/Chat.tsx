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
    sendAuthCode: (code: string) => void
}

export default function Chat({items, status, sendMessage, sendPermissionDecision, sendAuthCode}: Props) {
    const [draft, setDraft] = useState("")

    const pendingAuthCode = items.some(item => item.kind === "auth_code_needed" && !item.resolved)

    function handleSend() {
        const text = draft.trim()
        if (!text) return
        if (pendingAuthCode) {
            sendAuthCode(text)
        } else {
            sendMessage(text)
        }
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
                placeholder={pendingAuthCode ? "Enter the authentication code…" : "Message the workbench…"}
            />
        </div>
    )
}
