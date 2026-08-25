import {Button} from "@vervstack/chures"

import CloseIcon from "@/icons/common/CloseIcon.tsx"
import PlusIcon from "@/icons/common/PlusIcon.tsx"
import {SimpleChat} from "@/processes/SimpleChat.ts"
import SimpleChatHistoryRow from
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistoryRow/SimpleChatHistoryRow.tsx"
import cls from
    // eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistoryListScreen/SimpleChatHistoryListScreen.module.css"

interface Props {
    chats: SimpleChat[]
    loading: boolean
    activeChatId?: string
    onSelectChat: (chatId: string) => void
    onDeleteChat: (chat: SimpleChat) => void
    onNewChat: () => void
    onClose: () => void
}

export default function SimpleChatHistoryListScreen(props: Props) {
    return (
        <div className={cls.SimpleChatHistoryListScreenContainer}>
            <div className={cls.Header}>
                <h2 className={cls.Title}>Simple Chats</h2>
                <Button
                    variant="secondary"
                    className={cls.NewChatButton}
                    onClick={props.onNewChat}
                    aria-label="New chat"
                    title="New chat"
                >
                    <PlusIcon/>
                </Button>
                <Button
                    variant="secondary"
                    className={cls.CloseButton}
                    onClick={props.onClose}
                    aria-label="Close chat history"
                    title="Close chat history"
                >
                    <CloseIcon className={cls.CloseIcon}/>
                </Button>
            </div>
            <div className={cls.ListContainer}>
                {props.loading && props.chats.length === 0 && (
                    <p className={cls.EmptyState}>Loading chats…</p>
                )}
                {!props.loading && props.chats.length === 0 && (
                    <p className={cls.EmptyState}>No chats yet — start one!</p>
                )}
                {props.chats.map(chat => (
                    <SimpleChatHistoryRow
                        key={chat.id}
                        chat={chat}
                        active={chat.id === props.activeChatId}
                        onSelect={() => props.onSelectChat(chat.id)}
                        onDelete={() => props.onDeleteChat(chat)}
                    />
                ))}
            </div>
        </div>
    )
}
