import {useState} from "react"
import {ConfirmDialog} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useSimpleChatMutations, useSimpleChats} from "@/app/hooks/SimpleChat.ts"
import {cn} from "@/app/utils/cn.ts"
import {SimpleChat} from "@/processes/SimpleChat.ts"
import cls from "@/pages/workbench/components/SimpleChatHistorySidebar/SimpleChatHistorySidebar.module.css"
import SimpleChatHistoryListScreen from
    // eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistoryListScreen/SimpleChatHistoryListScreen.tsx"
import SimpleChatHistorySidebarCollapsedRail from
    // eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistorySidebarCollapsedRail/SimpleChatHistorySidebarCollapsedRail.tsx"

interface Props {
    vaultId: string
    activeChatId?: string
    onSelectChat: (chatId: string | undefined) => void
    onNewChat: () => void
    models: string[]
    currentModel: string
    modelsLoading?: boolean
    onChangeModel: (model: string) => void
    onClose: () => void
}

// New sibling to Chat/components/ChatHistorySidebar (not a modification of it) —
// backed by useSimpleChats (a saved-chat-thread list) instead of the docker chat's
// chatHistoryApi.ts (past bridge-session transcripts). A permanently-visible
// pinned left pane rather than a toggleable overlay, so it also owns the model
// switcher/new-chat header row (there's no page-level top bar for Simple Chat
// mode anymore). Only one screen (the list): unlike docker's history, a Simple
// Chat thread is resumable rather than a closed past transcript, so "select"
// means "switch the active chatId and go live" — there's no separate read-only
// detail screen to reuse ChatHistoryDetailScreen for. ChatHistoryListScreen also
// isn't reused as-is: it has no delete affordance and renders a docker-shaped
// ChatSessionSummary (id/lastActivityAt/firstUserMessage) rather than a
// SimpleChat (title/model/vaultAccess), so this renders its own
// SimpleChatHistoryListScreen/SimpleChatHistoryRow variant instead.
export default function SimpleChatHistorySidebar(props: Props) {
    const [collapsed, setCollapsed] = useState(false)
    const {chats, isLoading} = useSimpleChats(props.vaultId)
    const {delete: deleteChat} = useSimpleChatMutations(props.vaultId)
    const {CloseDialog, OpenDialog} = useDialog()
    const bakeError = useBakeError()

    function handleDeleteChat(chat: SimpleChat) {
        OpenDialog(
            <ConfirmDialog
                title="Delete Chat"
                message={`Delete "${chat.title || "Untitled chat"}"? This cannot be undone.`}
                confirmLabel="Delete"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() => deleteChat(chat.id)
                    .then(() => {
                        if (chat.id === props.activeChatId) props.onSelectChat(undefined)
                    })
                    .catch(e => bakeError("Failed to delete chat", e))}
            />
        )
    }

    return (
        <div className={cn(cls.SimpleChatHistorySidebarContainer, collapsed && cls.Collapsed)}>
            {collapsed ? (
                <SimpleChatHistorySidebarCollapsedRail
                    onExpand={() => setCollapsed(false)}
                    onNewChat={props.onNewChat}
                />
            ) : (
                <SimpleChatHistoryListScreen
                    chats={chats}
                    loading={isLoading}
                    activeChatId={props.activeChatId}
                    onSelectChat={props.onSelectChat}
                    onDeleteChat={handleDeleteChat}
                    onNewChat={props.onNewChat}
                    onCollapse={() => setCollapsed(true)}
                    models={props.models}
                    currentModel={props.currentModel}
                    modelsLoading={props.modelsLoading}
                    onChangeModel={props.onChangeModel}
                    onClose={props.onClose}
                />
            )}
        </div>
    )
}
