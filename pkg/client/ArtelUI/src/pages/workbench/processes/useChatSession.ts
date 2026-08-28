import {useCallback, useEffect, useRef, useState} from "react"

import {ChatEvent, PermissionDecision} from "@/pages/workbench/processes/chatProtocol.ts"
import {applyEvent, ChatItem} from "@/pages/workbench/processes/chatReducer.ts"

export type ChatConnectionStatus = "connecting" | "open" | "reconnecting" | "closed"

const INITIAL_BACKOFF_MS = 1000
const MAX_BACKOFF_MS = 15000

function buildWsUrl(vaultId: string): string {
    const proto = window.location.protocol === "https:" ? "wss:" : "ws:"
    return `${proto}//${window.location.host}/api/vaults/workbench/${vaultId}/terminal/`
}

// Connects to the workbench chat WebSocket (same route the old ttyd terminal iframe used —
// the backend track keeps the route/port, only the wire protocol changes from raw terminal
// bytes to ChatEvent JSON), reconnecting with exponential backoff on drop. Exposes the
// folded item list (see chatReducer.ts) plus helpers for the consumer -> bridge event
// types; each helper both sends over the wire and applies the event to local state
// optimistically. The bridge does echo user_message back to every consumer (including
// the sender, since consumers are treated as equals — see hub.go), so sendMessage tags
// it with a client-generated id: chatReducer's applyUserMessage recognizes the echo by
// that id and skips re-appending it.
export function useChatSession(vaultId: string | undefined) {
    const [items, setItems] = useState<ChatItem[]>([])
    const [status, setStatus] = useState<ChatConnectionStatus>("connecting")
    const [authComplete, setAuthComplete] = useState(false)
    const [pendingTurn, setPendingTurn] = useState(false)
    const wsRef = useRef<WebSocket | null>(null)
    const backoffRef = useRef(INITIAL_BACKOFF_MS)
    const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
    const closedByUserRef = useRef(false)
    const lastSeqRef = useRef(0)

    useEffect(() => {
        if (!vaultId) return

        closedByUserRef.current = false
        setItems([])
        setStatus("connecting")
        setAuthComplete(false)
        setPendingTurn(false)
        lastSeqRef.current = 0

        function connect() {
            if (closedByUserRef.current) return

            const ws = new WebSocket(buildWsUrl(vaultId!))
            wsRef.current = ws

            ws.onopen = () => {
                backoffRef.current = INITIAL_BACKOFF_MS
                setStatus("open")
            }
            ws.onmessage = event => {
                let parsed: ChatEvent
                try {
                    parsed = JSON.parse(event.data as string) as ChatEvent
                } catch {
                    return
                }
                // Deduplicate already-seen events based on seq. The backend stamps a
                // monotonically increasing seq once per event for the lifetime of the
                // bridge process; on reconnect, the backlog replays all events with their
                // original seqs, so we skip any where seq <= lastSeq (and seq !== 0, which
                // means "not set" with omitempty on the backend).
                if (typeof parsed.seq === "number" && parsed.seq > 0 && parsed.seq <= lastSeqRef.current) {
                    return
                }
                if (typeof parsed.seq === "number" && parsed.seq > 0) {
                    lastSeqRef.current = parsed.seq
                }
                if (parsed.type === "auth_complete") setAuthComplete(true)
                if (parsed.type === "turn_done" || parsed.type === "error" || parsed.type === "new_chat") {
                    setPendingTurn(false)
                }
                setItems(prev => applyEvent(prev, parsed))
            }
            ws.onclose = () => {
                if (closedByUserRef.current) return
                setStatus("reconnecting")
                const delay = backoffRef.current
                backoffRef.current = Math.min(backoffRef.current * 2, MAX_BACKOFF_MS)
                reconnectTimerRef.current = setTimeout(connect, delay)
            }
            ws.onerror = () => {
                ws.close()
            }
        }

        connect()

        return () => {
            closedByUserRef.current = true
            clearTimeout(reconnectTimerRef.current)
            wsRef.current?.close()
            wsRef.current = null
            setStatus("closed")
        }
    }, [vaultId])

    const dispatch = useCallback((event: ChatEvent) => {
        const ws = wsRef.current
        if (!ws || ws.readyState !== WebSocket.OPEN) return
        ws.send(JSON.stringify(event))
        setItems(prev => applyEvent(prev, event))
    }, [])

    const sendMessage = useCallback((text: string) => {
        const id = crypto.randomUUID()
        dispatch({type: "user_message", text, id})
        setPendingTurn(true)
    }, [dispatch])

    const sendPermissionDecision = useCallback((id: string, decision: PermissionDecision) => {
        dispatch({type: "permission_decision", id, decision})
        setPendingTurn(true)
    }, [dispatch])

    const startNewChat = useCallback(() => {
        dispatch({type: "new_chat"})
        setPendingTurn(false)
    }, [dispatch])

    return {items, status, authComplete, sendMessage, sendPermissionDecision, startNewChat, pendingTurn}
}
