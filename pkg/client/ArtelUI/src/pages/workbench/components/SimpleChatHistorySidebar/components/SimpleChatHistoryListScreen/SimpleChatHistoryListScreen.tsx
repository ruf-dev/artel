import {Button, ChevronDownIcon} from "@vervstack/chures"

import PlusIcon from "@/icons/common/PlusIcon.tsx"
import {SimpleChat} from "@/processes/SimpleChat.ts"
import CloseIcon from "@/icons/common/CloseIcon.tsx"
import ModelSwitcher from "@/pages/workbench/components/SimpleChat/components/ModelSwitcher/ModelSwitcher.tsx"
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
    onCollapse: () => void
    models: string[]
    currentModel: string
    modelsLoading?: boolean
    onChangeModel: (model: string) => void
    onClose: () => void
}

export default function SimpleChatHistoryListScreen(props: Props) {
    return (
        <div className={cls.SimpleChatHistoryListScreenContainer}>
            <div className={cls.Header}>
                <Button
                    variant="secondary"
                    className={cls.CollapseButton}
                    onClick={props.onCollapse}
                    aria-label="Collapse sidebar"
                    title="Collapse sidebar"
                >
                    <ChevronDownIcon className={cls.CollapseIcon}/>
                </Button>
                <h2 className={cls.Title}>History</h2>
                <Button variant="secondary" className={cls.CloseButton} onClick={props.onClose}
                        aria-label="Close chat" title="Back to workbench mode selection">
                    <CloseIcon/>
                </Button>
                <Button
                    variant="secondary"
                    className={cls.NewChatButton}
                    onClick={props.onNewChat}
                    aria-label="New chat"
                    title="New chat"
                >
                    <PlusIcon/>
                </Button>
            </div>
            <div className={cls.ModelRow}>
                <ModelSwitcher
                    models={props.models}
                    value={props.currentModel}
                    isLoading={props.modelsLoading}
                    onChange={props.onChangeModel}
                />
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
