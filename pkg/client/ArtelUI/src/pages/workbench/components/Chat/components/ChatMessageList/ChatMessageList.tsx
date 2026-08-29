import {useEffect, useRef, useState} from "react"

import cls from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.module.css"
import ChatMessageColumn from "@/pages/workbench/components/Chat/components/ChatMessageColumn/ChatMessageColumn.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    items: ChatItem[]
    vaultId: string
    assistantLabel?: string
    retryDisabled?: boolean
    onRetryMessage?: (text: string) => void
    onResendMessage?: (id: string, text: string) => void
    onPermissionDecision: (id: string, decision: PermissionDecision) => void
    pending?: {bucket: "normal" | "slow" | "stuck"; label?: string}
}

// Auto-scrolls to the newest content unless the user has scrolled up to read history —
// standard chat UX: a user actively reviewing earlier messages shouldn't get yanked to
// the bottom by a streaming assistant reply.
const BOTTOM_THRESHOLD_PX = 48

export default function ChatMessageList(props: Props) {
    const containerRef = useRef<HTMLDivElement>(null)
    const [autoScroll, setAutoScroll] = useState(true)

    useEffect(() => {
        if (!autoScroll) return
        const el = containerRef.current
        if (!el) return
        el.scrollTop = el.scrollHeight
    }, [props.items, autoScroll])

    function handleScroll() {
        const el = containerRef.current
        if (!el) return
        const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight
        setAutoScroll(distanceFromBottom < BOTTOM_THRESHOLD_PX)
    }

    return (
        <div className={cls.ChatMessageListContainer} ref={containerRef} onScroll={handleScroll}>
            {props.items.length === 0 && <p className={cls.EmptyState}>Send a message to start the conversation.</p>}
            <ChatMessageColumn
                items={props.items}
                vaultId={props.vaultId}
                assistantLabel={props.assistantLabel}
                retryDisabled={props.retryDisabled}
                onRetryMessage={props.onRetryMessage}
                onResendMessage={props.onResendMessage}
                onPermissionDecision={props.onPermissionDecision}
                pending={props.pending}
            />
        </div>
    )
}
