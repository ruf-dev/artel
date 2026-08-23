import {useEffect, useState} from "react"

import {cn} from "@/app/utils/cn.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import cls from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/ChatHistorySidebar.module.css"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import ChatHistoryDetailScreen from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/components/ChatHistoryDetailScreen/ChatHistoryDetailScreen.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import ChatHistoryListScreen from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/components/ChatHistoryListScreen/ChatHistoryListScreen.tsx"
import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatSessionSummary, getChatSession, listChatSessions} from "@/pages/workbench/processes/chatHistoryApi.ts"

interface Props {
    vaultId: string
    open: boolean
    onClose: () => void
}

type Screen = "list" | "detail"

interface DetailState {
    session: ChatSessionSummary
    items: ChatItem[]
}

export default function ChatHistorySidebar({vaultId, open, onClose}: Props) {
    const [screen, setScreen] = useState<Screen>("list")
    const [sessions, setSessions] = useState<ChatSessionSummary[]>([])
    const [loading, setLoading] = useState(false)
    const [detail, setDetail] = useState<DetailState | null>(null)
    const bakeError = useBakeError()

    // Re-fetch and reset to the list screen every time the sidebar is opened — mirrors
    // the old dialog's mount-time fetch, since the sidebar (unlike the dialog) stays
    // mounted between opens to support the slide transition.
    useEffect(() => {
        if (!open) return
        setScreen("list")
        setDetail(null)
        setLoading(true)
        listChatSessions(vaultId)
            .then(setSessions)
            .catch(err => bakeError("Failed to load chat history", err))
            .finally(() => setLoading(false))
    }, [open, vaultId, bakeError])

    function handleSelectSession(session: ChatSessionSummary) {
        setLoading(true)
        getChatSession(vaultId, session.id)
            .then(events => {
                const items = events.reduce(applyEvent, [] as ChatItem[])
                setDetail({session, items})
                setScreen("detail")
            })
            .catch(err => bakeError(`Failed to load session "${session.firstUserMessage || "untitled"}"`, err))
            .finally(() => setLoading(false))
    }

    function handleBackToList() {
        setDetail(null)
        setScreen("list")
    }

    return (
        <div className={cn(cls.ChatHistorySidebarContainer, open && cls.ChatHistorySidebarOpen)}>
            {screen === "list" && (
                <ChatHistoryListScreen
                    sessions={sessions}
                    loading={loading}
                    onSelectSession={handleSelectSession}
                    onClose={onClose}
                />
            )}
            {screen === "detail" && detail && (
                <ChatHistoryDetailScreen
                    session={detail.session}
                    items={detail.items}
                    onBack={handleBackToList}
                    onClose={onClose}
                />
            )}
        </div>
    )
}
