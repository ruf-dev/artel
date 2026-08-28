import {createElement, useEffect, useState} from "react"
import {ConfirmDialog} from "@vervstack/chures"

import {useDialog} from "@/app/hooks/Dialog"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {useSimpleChatMutations, useSimpleChats} from "@/app/hooks/SimpleChat.ts"
import {SimpleChat} from "@/processes/SimpleChat.ts"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import WorkbenchHistoryTranscriptDialog from "@/pages/workbench/components/WorkbenchHistoryTranscriptDialog/WorkbenchHistoryTranscriptDialog.tsx"
import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatSessionSummary, getChatSession, listChatSessions} from "@/pages/workbench/processes/chatHistoryApi.ts"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"

export interface WorkbenchHistoryRow {
    id: string
    title: string
    source: "api" | "docker"
    subtitle?: string
    timestamp?: string
}

export interface WorkbenchHistory {
    rows: WorkbenchHistoryRow[]
    loading: boolean
    activeId?: string
    onSelect: (id: string) => void
    onDelete?: (id: string) => void
    onNewChat: () => void
}

interface Params {
    mode: WorkbenchMode
    vaultId: string
    activeApiChatId?: string
    onSelectApiChat: (id: string | undefined) => void
    onNewChat: () => void
}

// Single flat-list interface behind the unified WorkbenchSidebar's History tab.
// Both modes' history sources (api: useSimpleChats saved threads; docker:
// chatHistoryApi past bridge-session transcripts) are mapped to WorkbenchHistoryRow
// here so the sidebar renders one list regardless of mode. All hooks run
// unconditionally; only the docker fetch effect is gated on mode.
export function useWorkbenchHistory(p: Params): WorkbenchHistory {
    const {chats, isLoading: apiLoading} = useSimpleChats(p.vaultId)
    const {delete: deleteApiChat} = useSimpleChatMutations(p.vaultId)
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()

    const [dockerSessions, setDockerSessions] = useState<ChatSessionSummary[]>([])
    const [dockerLoading, setDockerLoading] = useState(false)

    useEffect(() => {
        if (p.mode !== "docker" || !p.vaultId) return
        let cancelled = false
        setDockerLoading(true)
        listChatSessions(p.vaultId)
            .then(sessions => {
                if (!cancelled) setDockerSessions(sessions)
            })
            .catch(e => {
                if (!cancelled) bakeError("Failed to load chat history", e)
            })
            .finally(() => {
                if (!cancelled) setDockerLoading(false)
            })
        return () => {
            cancelled = true
        }
    }, [p.mode, p.vaultId, bakeError])

    function handleDeleteApi(id: string) {
        const chat = chats.find(c => c.id === id)
        OpenDialog(createElement(ConfirmDialog, {
            title: "Delete Chat",
            message: `Delete "${chat?.title || "Untitled chat"}"? This cannot be undone.`,
            confirmLabel: "Delete",
            cancelLabel: "Cancel",
            danger: true,
            onClose: CloseDialog,
            onConfirm: () => deleteApiChat(id)
                .then(() => {
                    if (id === p.activeApiChatId) p.onSelectApiChat(undefined)
                })
                .catch(e => bakeError("Failed to delete chat", e)),
        }))
    }

    function handleSelectDocker(id: string) {
        const session = dockerSessions.find(s => s.id === id)
        getChatSession(p.vaultId, id)
            .then(events => events.reduce(applyEvent, [] as ChatItem[]))
            .then(items => OpenDialog(createElement(WorkbenchHistoryTranscriptDialog, {
                title: session?.firstUserMessage || "Untitled session",
                items,
            })))
            .catch(e => bakeError("Failed to load session transcript", e))
    }

    if (p.mode === "docker") {
        return {
            rows: mapDockerRows(dockerSessions),
            loading: dockerLoading,
            activeId: undefined,
            onSelect: handleSelectDocker,
            onDelete: undefined,
            onNewChat: p.onNewChat,
        }
    }

    return {
        rows: mapApiRows(chats),
        loading: apiLoading,
        activeId: p.activeApiChatId,
        onSelect: p.onSelectApiChat,
        onDelete: handleDeleteApi,
        onNewChat: p.onNewChat,
    }
}

function mapApiRows(chats: SimpleChat[]): WorkbenchHistoryRow[] {
    return chats.map(c => ({
        id: c.id,
        title: c.title || "Untitled chat",
        source: "api" as const,
        subtitle: [c.model, formatTimestamp(c.lastActivityAt)].filter(Boolean).join(" · ") || undefined,
        timestamp: c.lastActivityAt || undefined,
    }))
}

function mapDockerRows(sessions: ChatSessionSummary[]): WorkbenchHistoryRow[] {
    return sessions.map(s => ({
        id: s.id,
        title: s.firstUserMessage || "Untitled session",
        source: "docker" as const,
        subtitle: formatTimestamp(s.lastActivityAt),
        timestamp: s.lastActivityAt,
    }))
}

function formatTimestamp(iso?: string): string | undefined {
    if (!iso) return undefined
    const d = new Date(iso)
    if (Number.isNaN(d.getTime())) return undefined
    return d.toLocaleDateString(undefined, {month: "short", day: "numeric", hour: "2-digit", minute: "2-digit"})
}
