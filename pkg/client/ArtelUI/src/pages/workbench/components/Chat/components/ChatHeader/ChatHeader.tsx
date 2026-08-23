import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/Chat/components/ChatHeader/ChatHeader.module.css"

interface Props {
    onNewChat: () => void
    onToggleHistory: () => void
    disabled?: boolean
    vaultId?: string
}

export default function ChatHeader({onNewChat, onToggleHistory, disabled, vaultId}: Props) {
    return (
        <div className={cls.ChatHeaderContainer}>
            {vaultId && (
                <Button
                    variant="secondary"
                    className={cls.HistoryButton}
                    onClick={onToggleHistory}
                    aria-label="View chat history"
                    title="View chat history"
                >
                    <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                         strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                        <circle cx="12" cy="12" r="10"/>
                        <polyline points="12 6 12 12 16 14"/>
                    </svg>
                </Button>
            )}
            <Button variant="secondary" className={cls.NewChatButton} onClick={onNewChat} disabled={disabled}>
                <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor"
                     strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                New Chat
            </Button>
        </div>
    )
}
