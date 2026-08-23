import {useEffect, useRef, useState} from "react"
import {AnimatePresence} from "framer-motion"

import cls from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.module.css"
import ChatMessageItem from "@/pages/workbench/components/Chat/components/ChatMessageItem/ChatMessageItem.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"

interface Props {
    items: ChatItem[]
    onPermissionDecision: (id: string, decision: PermissionDecision) => void
}

// Auto-scrolls to the newest content unless the user has scrolled up to read history —
// standard chat UX: a user actively reviewing earlier messages shouldn't get yanked to
// the bottom by a streaming assistant reply.
const BOTTOM_THRESHOLD_PX = 48

export default function ChatMessageList({items, onPermissionDecision}: Props) {
    const containerRef = useRef<HTMLDivElement>(null)
    const [autoScroll, setAutoScroll] = useState(true)

    useEffect(() => {
        if (!autoScroll) return
        const el = containerRef.current
        if (!el) return
        el.scrollTop = el.scrollHeight
    }, [items, autoScroll])

    function handleScroll() {
        const el = containerRef.current
        if (!el) return
        const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight
        setAutoScroll(distanceFromBottom < BOTTOM_THRESHOLD_PX)
    }

    return (
        <div className={cls.ChatMessageListContainer} ref={containerRef} onScroll={handleScroll}>
            {items.length === 0 && <p className={cls.EmptyState}>Send a message to start the conversation.</p>}
            <AnimatePresence initial={false}>
                {items.map(item => (
                    <ChatMessageItem key={item.key} item={item} onPermissionDecision={onPermissionDecision}/>
                ))}
            </AnimatePresence>
        </div>
    )
}
