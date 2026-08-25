import {Button} from "@vervstack/chures"
import {MouseEvent} from "react"

import {cn} from "@/app/utils/cn.ts"
import {SimpleChat} from "@/processes/SimpleChat.ts"
import cls from
    // eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistoryRow/SimpleChatHistoryRow.module.css"

interface Props {
    chat: SimpleChat
    active: boolean
    onSelect: () => void
    onDelete: () => void
}

export default function SimpleChatHistoryRow({chat, active, onSelect, onDelete}: Props) {
    function handleDelete(e: MouseEvent) {
        e.stopPropagation()
        onDelete()
    }

    return (
        <div className={cn(cls.SimpleChatHistoryRowContainer, active && cls.SimpleChatHistoryRowActive)}>
            <Button variant="secondary" className={cls.SelectButton} onClick={onSelect}>
                <span className={cls.RowTitle}>{chat.title || "Untitled chat"}</span>
                <span className={cls.RowMeta}>
                    {chat.model}
                    {chat.lastActivityAt && ` · ${formatDate(chat.lastActivityAt)}`}
                </span>
            </Button>
            <Button
                variant="secondary"
                className={cls.DeleteButton}
                onClick={handleDelete}
                aria-label="Delete chat"
                title="Delete chat"
            >
                <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor"
                     strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    <line x1="10" y1="11" x2="10" y2="17"/>
                    <line x1="14" y1="11" x2="14" y2="17"/>
                </svg>
            </Button>
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
