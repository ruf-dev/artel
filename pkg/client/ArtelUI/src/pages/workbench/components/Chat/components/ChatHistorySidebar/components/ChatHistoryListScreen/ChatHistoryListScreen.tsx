import {Button} from "@vervstack/chures"

import CloseIcon from "@/icons/common/CloseIcon.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/components/ChatHistoryListScreen/ChatHistoryListScreen.module.css"
import {ChatSessionSummary} from "@/pages/workbench/processes/chatHistoryApi.ts"

interface Props {
    sessions: ChatSessionSummary[]
    loading: boolean
    onSelectSession: (session: ChatSessionSummary) => void
    onClose: () => void
}

export default function ChatHistoryListScreen({sessions, loading, onSelectSession, onClose}: Props) {
    return (
        <div className={cls.ChatHistoryListScreenContainer}>
            <div className={cls.Header}>
                <h2 className={cls.Title}>Chat History</h2>
                <Button
                    variant="secondary"
                    className={cls.CloseButton}
                    onClick={onClose}
                    aria-label="Close chat history"
                    title="Close chat history"
                >
                    <CloseIcon className={cls.CloseIcon}/>
                </Button>
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
                        onClick={() => onSelectSession(session)}
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
