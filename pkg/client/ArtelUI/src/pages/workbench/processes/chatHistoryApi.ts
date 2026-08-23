import {ChatEvent} from "@/pages/workbench/processes/chatProtocol.ts"

export interface ChatSessionSummary {
    id: string
    lastActivityAt?: string
    firstUserMessage?: string
}

function buildHistoryUrl(vaultId: string): string {
    const proto = window.location.protocol === "https:" ? "https:" : "http:"
    return `${proto}//${window.location.host}/api/vaults/workbench/${vaultId}/terminal/history`
}

/**
 * Fetches a list of past chat sessions for a vault, with summaries
 * (id, last activity time, first user message snippet).
 */
export function listChatSessions(vaultId: string): Promise<ChatSessionSummary[]> {
    const url = buildHistoryUrl(vaultId)
    return fetch(url, {credentials: "include"})
        .then(response => {
            if (!response.ok) throw new Error(`HTTP ${response.status} fetching chat history list`)
            return response.json() as Promise<ChatSessionSummary[]>
        })
}

/**
 * Fetches the full ordered event list for a past chat session.
 */
export function getChatSession(vaultId: string, id: string): Promise<ChatEvent[]> {
    const url = `${buildHistoryUrl(vaultId)}/${encodeURIComponent(id)}`
    return fetch(url, {credentials: "include"})
        .then(response => {
            if (!response.ok) throw new Error(`HTTP ${response.status} fetching chat session ${id}`)
            return response.json() as Promise<ChatEvent[]>
        })
}
