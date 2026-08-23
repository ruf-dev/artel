import {Button} from "@vervstack/chures"

import CloseIcon from "@/icons/common/CloseIcon.tsx"
import ChatMessageList from "@/pages/workbench/components/Chat/components/ChatMessageList/ChatMessageList.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/workbench/components/Chat/components/ChatHistorySidebar/components/ChatHistoryDetailScreen/ChatHistoryDetailScreen.module.css"
import {ChatSessionSummary} from "@/pages/workbench/processes/chatHistoryApi.ts"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

interface Props {
    session: ChatSessionSummary
    items: ChatItem[]
    onBack: () => void
    onClose: () => void
}

export default function ChatHistoryDetailScreen({session, items, onBack, onClose}: Props) {
    return (
        <div className={cls.ChatHistoryDetailScreenContainer}>
            <div className={cls.Header}>
                <Button
                    variant="secondary"
                    className={cls.BackButton}
                    onClick={onBack}
                    aria-label="Back to list"
                    title="Back to list"
                >
                    ←
                </Button>
                <h2 className={cls.Title}>{session.firstUserMessage || "Empty Session"}</h2>
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
            <div className={cls.DetailContainer}>
                <ChatMessageList items={items} onPermissionDecision={() => {}}/>
            </div>
        </div>
    )
}
