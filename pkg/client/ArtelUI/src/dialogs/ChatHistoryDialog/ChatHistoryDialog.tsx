import {useEffect, useState} from "react"
import {Button} from "@vervstack/chures"

import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatSessionSummary, getChatSession, listChatSessions} from "@/pages/workbench/processes/chatHistoryApi.ts"
import ChatMessageList from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.tsx"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import cls from "@/dialogs/ChatHistoryDialog/ChatHistoryDialog.module.css"

interface Props {
    vaultId: string
}

type Screen = "list" | "detail"

interface DetailState {
    session: ChatSessionSummary
    items: ChatItem[]
}

export default function ChatHistoryDialog({vaultId}: Props) {
    const [screen, setScreen] = useState<Screen>("list")
    const [sessions, setSessions] = useState<ChatSessionSummary[]>([])
    const [loading, setLoading] = useState(false)
    const [detail, setDetail] = useState<DetailState | null>(null)
    const bakeError = useBakeError()

    useEffect(() => {
        setLoading(true)
        listChatSessions(vaultId)
            .then(setSessions)
            .catch(err => bakeError("Failed to load chat history", err))
            .finally(() => setLoading(false))
    }, [vaultId, bakeError])

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
        <div className={cls.ChatHistoryDialogContainer}>
            {screen === "list" && (
                <>
                    <div className={cls.Header}>
                        <h2 className={cls.Title}>Chat History</h2>
                    </div>
                    <div className={cls.ListContainer}>
                        {loading && sessions.length === 0 && (
                            <p className={cls.EmptyState}>Loading chat history…</p>
                        )}
                        {!loading && sessions.length === 0 && (
                            <p className={cls.EmptyState}>No previous chat sessions yet.</p>
                        )}
                        {sessions.map(session => (
                            <Button
                                key={session.id}
                                variant="secondary"
                                className={cls.SessionRow}
                                onClick={() => handleSelectSession(session)}
                                disabled={loading}
                            >
                                <div className={cls.SessionTitle}>
                                    {session.firstUserMessage || "(empty session)"}
                                </div>
                                {session.lastActivityAt && (
                                    <div className={cls.SessionTime}>{formatDate(session.lastActivityAt)}</div>
                                )}
                            </Button>
                        ))}
                    </div>
                </>
            )}
            {screen === "detail" && detail && (
                <>
                    <div className={cls.Header}>
                        <Button
                            variant="secondary"
                            className={cls.BackButton}
                            onClick={handleBackToList}
                            aria-label="Back to list"
                            title="Back to list"
                        >
                            ←
                        </Button>
                        <h2 className={cls.Title}>{detail.session.firstUserMessage || "Empty Session"}</h2>
                    </div>
                    <div className={cls.DetailContainer}>
                        <ChatMessageList items={detail.items} onPermissionDecision={() => {}}/>
                    </div>
                </>
            )}
        </div>
    )
}

function formatDate(isoString: string): string {
    try {
        const date = new Date(isoString)
        return date.toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        })
    } catch {
        return isoString
    }
}
